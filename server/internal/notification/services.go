package notification

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService struct {
	notificationRepo *NotificationRepo
}

func NewService(notificationRepo *NotificationRepo) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) Subscribe(ctx context.Context, userId uuid.UUID, endpoint string, p256dh string, auth string) (NotificationSubscription, error) {
	return s.notificationRepo.CreateSubscription(ctx, userId, endpoint, p256dh, auth)
}

func (s *NotificationService) Unsubscribe(ctx context.Context, userId uuid.UUID, endpoint string) error {
	return s.notificationRepo.DeleteSubscriptionByEndpoint(ctx, userId, endpoint)
}

func (s *NotificationService) GetSubscriptionsOfUser(ctx context.Context, userId uuid.UUID) ([]NotificationSubscription, error) {
	return s.notificationRepo.GetSubscriptionsOfUser(ctx, userId)
}
