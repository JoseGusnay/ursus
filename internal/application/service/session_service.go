package service

import (
	"context"
	"fmt"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
	"github.com/google/uuid"
	"time"
)

type SessionService struct {
	repo repository.SessionRepository
}

func NewSessionService(repo repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Repository() repository.SessionRepository {
	return s.repo
}

func (s *SessionService) Start(ctx context.Context, title string) (*entity.Session, error) {
	// Deactivate any currently active session
	if err := s.repo.DeactivateAll(ctx); err != nil {
		return nil, fmt.Errorf("failed to deactivate previous sessions: %w", err)
	}

	session := entity.NewSession(uuid.New().String(), title)
	if err := s.repo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save new session: %w", err)
	}

	return session, nil
}

func (s *SessionService) End(ctx context.Context) error {
	active, err := s.repo.GetActive(ctx)
	if err != nil {
		return err
	}
	if active == nil {
		return nil // No active session to end
	}

	active.IsActive = false
	active.EndTime = time.Now()
	return s.repo.Save(ctx, active)
}

func (s *SessionService) GetActive(ctx context.Context) (*entity.Session, error) {
	return s.repo.GetActive(ctx)
}

func (s *SessionService) List(ctx context.Context) ([]*entity.Session, error) {
	return s.repo.List(ctx)
}
