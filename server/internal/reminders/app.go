package reminders

import (
	"context"
	"errors"
	"log/slog"
	"server/internal/notification"
	"time"

	"github.com/jackc/pgx/v5"
)

// pollInterval is how long the loop waits when there is nothing due right now
// (or after an error) before checking again.
const pollInterval = 5 * time.Second

type ReminderApp struct {
	remindersRepo *RemindersRepo
}

func NewReminderApp(remindersRepo *RemindersRepo) *ReminderApp {
	return &ReminderApp{
		remindersRepo: remindersRepo,
	}
}

// updateReminderAfterNotify advances a reminder once it has fired: one-time
// reminders are deactivated; recurring ones move to their next occurrence.
// AddDate is used (not a fixed number of seconds) so monthly/yearly land on the
// correct calendar date regardless of month length or leap years.
func (a *ReminderApp) updateReminderAfterNotify(ctx context.Context, reminder Reminder) error {
	switch reminder.ReminderType {
	case ONE_TIME_REMINDER_TYPE:
		_, err := a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, reminder.RemindTimestamp, reminder.ReminderType, false, reminder.Name, reminder.Description)
		return err
	case DAILY_REMINDER_TYPE:
		next := reminder.RemindTimestamp.AddDate(0, 0, 1)
		_, err := a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, next, reminder.ReminderType, true, reminder.Name, reminder.Description)
		return err
	case MONTHLY_REMINDER_TYPE:
		next := reminder.RemindTimestamp.AddDate(0, 1, 0)
		_, err := a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, next, reminder.ReminderType, true, reminder.Name, reminder.Description)
		return err
	case YEARLY_REMINDER_TYPE:
		next := reminder.RemindTimestamp.AddDate(1, 0, 0)
		_, err := a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, next, reminder.ReminderType, true, reminder.Name, reminder.Description)
		return err
	}
	return nil
}

func (a *ReminderApp) runReminderLoop(ctx context.Context, notifyService notification.Notification) {
	for {
		// Stop promptly when the app is shutting down.
		select {
		case <-ctx.Done():
			slog.Info("reminder loop stopped", "reason", ctx.Err())
			return
		default:
		}

		// GetNextReminder only returns reminders whose time has already arrived,
		// so anything it returns should fire now.
		nextReminder, err := a.remindersRepo.GetNextReminder(ctx)
		if err != nil {
			// ErrNoRows just means nothing is due right now — only real errors
			// are worth logging.
			if !errors.Is(err, pgx.ErrNoRows) {
				slog.Error("failed to get next reminder", "error", err.Error())
			}
			if !wait(ctx, pollInterval) {
				return
			}
			continue
		}

		notifyService.Notify(ctx, nextReminder.UserId, nextReminder.Name, nextReminder.Description)

		if err := a.updateReminderAfterNotify(ctx, nextReminder); err != nil {
			// If we can't advance/deactivate it, the same reminder stays "due"
			// and we'd resend it in a tight loop — back off instead.
			slog.Error("failed to update reminder after notify", "error", err.Error(), "reminderId", nextReminder.Id.String())
			if !wait(ctx, pollInterval) {
				return
			}
		}
	}
}

// wait sleeps for d, but returns false immediately if ctx is cancelled — so a
// shutdown doesn't have to wait out a full poll interval.
func wait(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (a *ReminderApp) Run(ctx context.Context, notifyService notification.Notification) {
	go a.runReminderLoop(ctx, notifyService)
}
