package usecase

import (
	"context"
	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

// SaveMemoryUseCase handles the logic for storing new memories in Ursus.
type SaveMemoryUseCase struct {
	service *service.MemoryService
}

// NewSaveMemoryUseCase creates a new SaveMemoryUseCase instance.
func NewSaveMemoryUseCase(service *service.MemoryService) *SaveMemoryUseCase {
	return &SaveMemoryUseCase{service: service}
}

// Execute performs the save operation.
func (u *SaveMemoryUseCase) Execute(ctx context.Context, content, metadata string) (*entity.Ursus, error) {
	return u.service.Store(ctx, content, metadata)
}
