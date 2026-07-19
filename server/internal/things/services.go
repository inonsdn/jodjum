package things

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ThingsService struct {
	thingsRepo *ThingsRepo
}

func NewService(thingsRepo *ThingsRepo) *ThingsService {
	return &ThingsService{
		thingsRepo: thingsRepo,
	}
}

func (s *ThingsService) GetThingsOfUser(ctx context.Context, userId uuid.UUID) ([]Things, error) {
	return s.thingsRepo.GetThingsOfUser(ctx, userId)
}

func (s *ThingsService) GetThingsById(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID) (Things, error) {
	return s.thingsRepo.GetThingsById(ctx, userId, thingsId)
}

func (s *ThingsService) CreateThings(ctx context.Context, userId uuid.UUID, name string, description string, quantity int, expiredAt time.Time) (Things, error) {
	return s.thingsRepo.CreateThings(ctx, userId, name, description, quantity, expiredAt)
}

func (s *ThingsService) UpdateThings(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID, name string, description string, quantity int, expiredAt time.Time) (Things, error) {
	return s.thingsRepo.UpdateThings(ctx, userId, thingsId, name, description, quantity, expiredAt)
}

func (s *ThingsService) DeleteThings(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID) error {
	return s.thingsRepo.DeleteThings(ctx, userId, thingsId)
}
