package repository

import (
	"context"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

type SessionRepository interface {
	Save(ctx context.Context, session *entity.Session) error
	GetActive(ctx context.Context) (*entity.Session, error)
	GetByID(ctx context.Context, id string) (*entity.Session, error)
	List(ctx context.Context) ([]*entity.Session, error)
	DeactivateAll(ctx context.Context) error
}
