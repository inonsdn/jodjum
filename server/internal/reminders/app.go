package reminders

import (
	"context"
	"server/internal/constants"
	"server/internal/notification"
	"time"
)

type ReminderApp struct {
	remindersRepo *RemindersRepo
}

func NewReminderApp(remindersRepo *RemindersRepo) *ReminderApp {
	return &ReminderApp{
		remindersRepo: remindersRepo,
	}
}

func (a *ReminderApp) updateReminderAfterNotify(ctx context.Context, reminder Reminder) {
	if reminder.ReminderType == ONE_TIME_REMINDER_TYPE {
		a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, reminder.RemindTimestamp, reminder.ReminderType, false, reminder.Name, reminder.Description)
	} else if reminder.ReminderType == DAILY_REMINDER_TYPE {
		reminderTimestamp := reminder.RemindTimestamp.Add(constants.ONE_DAY_TIMESTAMP_SECS)
		a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, reminderTimestamp, reminder.ReminderType, true, reminder.Name, reminder.Description)
	} else if reminder.ReminderType == MONTHLY_REMINDER_TYPE {
		reminderTimestamp := reminder.RemindTimestamp.Add(constants.ONE_MONTH_TIMESTAMP_SECS)
		a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, reminderTimestamp, reminder.ReminderType, true, reminder.Name, reminder.Description)
	} else if reminder.ReminderType == YEARLY_REMINDER_TYPE {
		reminderTimestamp := reminder.RemindTimestamp.Add(constants.ONE_MONTH_TIMESTAMP_SECS * 12)
		a.remindersRepo.UpdateReminder(ctx, reminder.UserId, reminder.Id, reminderTimestamp, reminder.ReminderType, true, reminder.Name, reminder.Description)
	} else {
		return
	}
}

func (a *ReminderApp) runReminderLoop(notifyService notification.Notification) {
	context := context.Background()
	for {
		nextReminder, err := a.remindersRepo.GetNextReminder(context)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		notifyService.Notify(nextReminder.Name, nextReminder.Description)
		a.updateReminderAfterNotify(context, nextReminder)
	}
}

func (a *ReminderApp) Run(notifyService notification.Notification) {
	go a.runReminderLoop(notifyService)
}
