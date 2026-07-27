package service

import (
	"context"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) ListByUser(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error) {
	items, err := s.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, errors.NewInternal("failed to list notifications")
	}
	if items == nil {
		items = []models.Notification{}
	}
	return items, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id string) error {
	err := s.repo.MarkAsRead(ctx, id)
	if err != nil {
		return errors.NewInternal("failed to mark notification as read")
	}
	return nil
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	err := s.repo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return errors.NewInternal("failed to mark all notifications as read")
	}
	return nil
}
