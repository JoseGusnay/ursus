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
	ListBySession(ctx context.Context, sessionID string) ([]*entity.Ursus, error)
	GetByID(ctx context.Context, id string) (*entity.Ursus, error)
	GetByTopicKey(ctx context.Context, topicKey string) (*entity.Ursus, error)
	Update(ctx context.Context, ursus *entity.Ursus) error
	Delete(ctx context.Context, id string) error
}
