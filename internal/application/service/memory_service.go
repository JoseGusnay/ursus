package service

import (
	"context"
	"time"

	"github.com/JoseGusnay/ursus/internal/application/usecase"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// MemoryService provides high-level operations for managed memories.
type MemoryService struct {
	repo          repository.UrsusRepository
	saveUseCase   *usecase.SaveMemoryUseCase
	searchUseCase *usecase.SearchMemoryUseCase
}

func (s *MemoryService) Repository() repository.UrsusRepository {
	return s.repo
}

// NewMemoryService creates a new MemoryService instance.
func NewMemoryService(repo repository.UrsusRepository, saveUseCase *usecase.SaveMemoryUseCase, searchUseCase *usecase.SearchMemoryUseCase) *MemoryService {
	return &MemoryService{
		repo:          repo,
		saveUseCase:   saveUseCase,
		searchUseCase: searchUseCase,
	}
}

// Store saves a new piece of memory or updates if topicKey matches.
func (s *MemoryService) Store(ctx context.Context, content, metadata, topicKey, scope, prompt string) (*entity.Ursus, error) {
	return s.saveUseCase.Execute(ctx, content, metadata, topicKey, scope, prompt)
}

// GetByID returns a memory by its ID.
func (s *MemoryService) GetByID(ctx context.Context, id string) (*entity.Ursus, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing memory.
func (s *MemoryService) Update(ctx context.Context, u *entity.Ursus) error {
	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, u)
}

// Search searches for memories matching a query.
func (s *MemoryService) Search(ctx context.Context, query string) ([]*entity.Ursus, error) {
	return s.searchUseCase.Execute(ctx, query)
}

// List returns all memories ordered by creation date.
func (s *MemoryService) List(ctx context.Context) ([]*entity.Ursus, error) {
	return s.repo.List(ctx)
}

// Delete removes a memory by its ID.
func (s *MemoryService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
