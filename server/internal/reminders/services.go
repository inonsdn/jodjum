package reminders

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RemindersService struct {
	remindersRepo *RemindersRepo
}

func NewService(remindersRepo *RemindersRepo) *RemindersService {
	return &RemindersService{
		remindersRepo: remindersRepo,
	}
}

func (s *RemindersService) GetRemindersOfUser(ctx context.Context, userId uuid.UUID) ([]Reminder, error) {
	return s.remindersRepo.GetRemindersOfUser(ctx, userId)
}

func (s *RemindersService) GetReminderById(ctx context.Context, userId uuid.UUID, reminderId uuid.UUID) (Reminder, error) {
	return s.remindersRepo.GetReminderById(ctx, userId, reminderId)
}

func (s *RemindersService) CreateReminder(ctx context.Context, userId uuid.UUID, remindTimestamp time.Time, reminderType int, isActive bool, name string, description string) (Reminder, error) {
	return s.remindersRepo.CreateReminder(ctx, userId, remindTimestamp, reminderType, isActive, name, description)
}

func (s *RemindersService) UpdateReminder(ctx context.Context, userId uuid.UUID, reminderId uuid.UUID, remindTimestamp time.Time, reminderType int, isActive bool, name string, description string) (Reminder, error) {
	return s.remindersRepo.UpdateReminder(ctx, userId, reminderId, remindTimestamp, reminderType, isActive, name, description)
}
