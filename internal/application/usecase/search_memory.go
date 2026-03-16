package usecase

import (
	"context"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// SearchMemoryUseCase handles the logic for searching memories.
type SearchMemoryUseCase struct {
	repo repository.UrsusRepository
}

// NewSearchMemoryUseCase creates a new SearchMemoryUseCase instance.
func NewSearchMemoryUseCase(repo repository.UrsusRepository) *SearchMemoryUseCase {
	return &SearchMemoryUseCase{
		repo: repo,
	}
}

// Execute performs the search and returns compact results for token efficiency.
func (u *SearchMemoryUseCase) Execute(ctx context.Context, query string) ([]*entity.Ursus, error) {
	results, err := u.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	// Truncate content for token efficiency (Discover Layer)
	for _, m := range results {
		if len(m.Content) > 200 {
			m.Content = m.Content[:200] + "..."
		}
	}

	return results, nil
}
