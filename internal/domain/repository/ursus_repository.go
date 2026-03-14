package repository

import (
	"context"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

// UrsusRepository defines the interface for persisting and retrieving Ursus memories.
type UrsusRepository interface {
	Save(ctx context.Context, ursus *entity.Ursus) error
	Search(ctx context.Context, query string) ([]*entity.Ursus, error)
	List(ctx context.Context) ([]*entity.Ursus, error)
	Delete(ctx context.Context, id string) error
}
