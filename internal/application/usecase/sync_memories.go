package usecase

import (
	"context"
	"fmt"

	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

type SyncMemoriesUseCase struct {
	sqliteRepo repository.UrsusRepository
	jsonlRepo  repository.UrsusRepository
}

func NewSyncMemoriesUseCase(sqlite repository.UrsusRepository, jsonl repository.UrsusRepository) *SyncMemoriesUseCase {
	return &SyncMemoriesUseCase{
		sqliteRepo: sqlite,
		jsonlRepo:  jsonl,
	}
}

func (u *SyncMemoriesUseCase) Export(ctx context.Context) error {
	memories, err := u.sqliteRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list memories from sqlite: %w", err)
	}

	// Simple loop for now to stick to interface and ensure all records are saved safely
	for _, m := range memories {
		if err := u.jsonlRepo.Save(ctx, m); err != nil {
			return err
		}
	}

	return nil
}

func (u *SyncMemoriesUseCase) Import(ctx context.Context) error {
	memories, err := u.jsonlRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list memories from jsonl: %w", err)
	}

	for _, m := range memories {
		if err := u.sqliteRepo.Save(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
