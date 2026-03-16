package usecase

import (
	"context"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// GetMemoryDetailUseCase handles the logic for retrieving full details of a memory.
type GetMemoryDetailUseCase struct {
	repo repository.UrsusRepository
}

// NewGetMemoryDetailUseCase creates a new GetMemoryDetailUseCase instance.
func NewGetMemoryDetailUseCase(repo repository.UrsusRepository) *GetMemoryDetailUseCase {
	return &GetMemoryDetailUseCase{
		repo: repo,
	}
}

// Execute performs the retrieval operation.
func (u *GetMemoryDetailUseCase) Execute(ctx context.Context, id string) (*entity.Ursus, error) {
	return u.repo.GetByID(ctx, id)
}
