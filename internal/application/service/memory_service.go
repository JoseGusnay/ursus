package service

import (
	"context"
	"time"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
	"github.com/google/uuid"
)

// MemoryService provides high-level operations for managed memories.
type MemoryService struct {
	repo repository.UrsusRepository
}

// NewMemoryService creates a new MemoryService instance.
func NewMemoryService(repo repository.UrsusRepository) *MemoryService {
	return &MemoryService{repo: repo}
}

// Store creates and persists a new memory.
func (s *MemoryService) Store(ctx context.Context, content, metadata string) (*entity.Ursus, error) {
	u := &entity.Ursus{
		ID:        uuid.New().String(),
		Content:   content,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// Search searches for memories matching a query.
func (s *MemoryService) Search(ctx context.Context, query string) ([]*entity.Ursus, error) {
	return s.repo.Search(ctx, query)
}

// List returns all memories ordered by creation date.
func (s *MemoryService) List(ctx context.Context) ([]*entity.Ursus, error) {
	return s.repo.List(ctx)
}

// Delete removes a memory by its ID.
func (s *MemoryService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
