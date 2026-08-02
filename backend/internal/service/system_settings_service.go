package service

import (
	"context"

	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type SystemSettingsService struct {
	repo *repository.SystemSettingsRepository
}

func NewSystemSettingsService(repo *repository.SystemSettingsRepository) *SystemSettingsService {
	return &SystemSettingsService{repo: repo}
}

func (s *SystemSettingsService) GetSettings(ctx context.Context) (map[string]interface{}, error) {
	settings, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.NewInternal("failed to retrieve system settings")
	}
	return settings, nil
}

func (s *SystemSettingsService) UpdateSettings(ctx context.Context, settings map[string]interface{}) error {
	if len(settings) == 0 {
		return errors.NewBadRequest("settings object cannot be empty")
	}
	for key := range settings {
		if key == "gemini_api_key" || key == "gemini_api_key_encrypted" {
			return errors.NewBadRequest("AI secrets must be updated through the AI settings endpoint")
		}
	}
	if err := s.repo.UpdateSettings(ctx, settings); err != nil {
		return errors.NewInternal("failed to update system settings")
	}
	return nil
}
