package usecase

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"io"
	"compress/gzip"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// ImportChunksUseCase handles importing memories from git chunks.
type ImportChunksUseCase struct {
	repo repository.UrsusRepository
}

func NewImportChunksUseCase(repo repository.UrsusRepository) *ImportChunksUseCase {
	return &ImportChunksUseCase{repo: repo}
}

// Execute imports all memories from chunks in the project directory.
func (u *ImportChunksUseCase) Execute(ctx context.Context, projectDir string) error {
	ursusDir := filepath.Join(projectDir, ".ursus")
	manifestPath := filepath.Join(ursusDir, "manifest.json")
	chunksDir := filepath.Join(ursusDir, "chunks")

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil // Nothing to import
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	for chunkFile := range manifest.Chunks {
		if err := u.importChunk(ctx, filepath.Join(chunksDir, chunkFile)); err != nil {
			return err
		}
	}

	return nil
}

func (u *ImportChunksUseCase) importChunk(ctx context.Context, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file
	if filepath.Ext(path) == ".gz" {
		gzr, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzr.Close()
		reader = gzr
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var m entity.Ursus
		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			continue
		}
		
		// Use repository's hygiene logic (Save handles deduplication at ID level or content level if implemented)
		// Since we want to leverage our new topic/dedup logic, we should use the SaveMemoryUseCase 
		// or at least repository methods that handle existence.
		existing, _ := u.repo.GetByID(ctx, m.ID)
		if existing == nil {
			_ = u.repo.Save(ctx, &m)
		}
	}
	return scanner.Err()
}
