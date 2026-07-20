package things

import (
	"context"
	"server/internal/reminders"
	"time"

	"github.com/google/uuid"
)

type ThingsService struct {
	thingsRepo   *ThingsRepo
	reminderRepo *reminders.RemindersRepo
}

func NewService(thingsRepo *ThingsRepo, reminderRepo *reminders.RemindersRepo) *ThingsService {
	return &ThingsService{
		thingsRepo:   thingsRepo,
		reminderRepo: reminderRepo,
	}
}

func (s *ThingsService) GetThingsOfUser(ctx context.Context, userId uuid.UUID) ([]Things, error) {
	return s.thingsRepo.GetThingsOfUser(ctx, userId)
}

func (s *ThingsService) GetThingsById(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID) (Things, error) {
	return s.thingsRepo.GetThingsById(ctx, userId, thingsId)
}

func (s *ThingsService) CreateThings(ctx context.Context, userId uuid.UUID, name string, description string, quantity int, expiredAt *time.Time, nextRemindTimestamp *float64) (Things, error) {
	tx, err := s.thingsRepo.db.Begin(ctx)
	if err != nil {
		return Things{}, err
	}
	// Rollback is a no-op once the tx has been committed, so this safely undoes
	// everything if we return early with an error.
	defer tx.Rollback(ctx)

	thing, err := s.thingsRepo.CreateThingsTx(ctx, tx, userId, name, description, quantity, expiredAt)
	if err != nil {
		return Things{}, err
	}

	// Only create a reminder when the caller asked for one.
	if nextRemindTimestamp != nil {
		if _, err := s.reminderRepo.CreateReminderTx(ctx, tx, userId, remindTimeFromOffset(*nextRemindTimestamp), reminders.ONE_TIME_REMINDER_TYPE, true, name, description); err != nil {
			return Things{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Things{}, err
	}

	return thing, nil
}

func (s *ThingsService) UpdateThings(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID, name string, description string, quantity int, expiredAt *time.Time, nextRemindTimestamp *float64) (Things, error) {
	tx, err := s.thingsRepo.db.Begin(ctx)
	if err != nil {
		return Things{}, err
	}
	defer tx.Rollback(ctx)

	// If the thing doesn't exist this returns pgx.ErrNoRows and the tx rolls
	// back, so no reminder is left behind.
	thing, err := s.thingsRepo.UpdateThingsTx(ctx, tx, userId, thingsId, name, description, quantity, expiredAt)
	if err != nil {
		return Things{}, err
	}

	if nextRemindTimestamp != nil {
		if _, err := s.reminderRepo.CreateReminderTx(ctx, tx, userId, remindTimeFromOffset(*nextRemindTimestamp), reminders.ONE_TIME_REMINDER_TYPE, true, name, description); err != nil {
			return Things{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Things{}, err
	}

	return thing, nil
}

// remindTimeFromOffset turns next_remind_timestamp (seconds from now) into an
// absolute time. Fractional seconds are preserved.
func remindTimeFromOffset(offsetSeconds float64) time.Time {
	return time.Now().Add(time.Duration(offsetSeconds * float64(time.Second)))
}

func (s *ThingsService) DeleteThings(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID) error {
	return s.thingsRepo.DeleteThings(ctx, userId, thingsId)
}
