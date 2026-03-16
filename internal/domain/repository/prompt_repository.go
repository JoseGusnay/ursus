package repository

import (
	"context"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

// PromptRepository defines the interface for persisting prompts.
type PromptRepository interface {
	Save(ctx context.Context, p *entity.Prompt) error
	GetByID(ctx context.Context, id string) (*entity.Prompt, error)
}
