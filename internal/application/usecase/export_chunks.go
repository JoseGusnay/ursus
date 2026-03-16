package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"compress/gzip"

	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// ExportChunksUseCase handles exporting local memories to git-friendly chunks.
type ExportChunksUseCase struct {
	repo repository.UrsusRepository
}

// Manifest represents the index of synchronized chunks.
type Manifest struct {
	Chunks map[string]ChunkInfo `json:"chunks"`
}

// ChunkInfo contains metadata about a specific chunk.
type ChunkInfo struct {
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name"`
}

func NewExportChunksUseCase(repo repository.UrsusRepository) *ExportChunksUseCase {
	return &ExportChunksUseCase{repo: repo}
}

// Execute exports unsynced memories to a new chunk in the specified project directory.
func (u *ExportChunksUseCase) Execute(ctx context.Context, projectDir string, userName string) error {
	ursusDir := filepath.Join(projectDir, ".ursus")
	chunksDir := filepath.Join(ursusDir, "chunks")
	manifestPath := filepath.Join(ursusDir, "manifest.json")

	if err := os.MkdirAll(chunksDir, 0755); err != nil {
		return err
	}

	// 1. Get all local memories
	memories, err := u.repo.List(ctx)
	if err != nil {
		return err
	}

	if len(memories) == 0 {
		return nil
	}

	// 2. Generate content for the chunk
	var content []byte
	for _, m := range memories {
		line, _ := json.Marshal(m)
		content = append(content, line...)
		content = append(content, '\n')
	}

	// 3. Generate content hash for filename
	hash := fmt.Sprintf("%x", sha256.Sum256(content))[:12]
	chunkFilename := fmt.Sprintf("%s.jsonl.gz", hash)
	chunkPath := filepath.Join(chunksDir, chunkFilename)

	// 4. Save gzipped chunk file
	f, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	if _, err := gw.Write(content); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}

	// 5. Update manifest
	manifest := Manifest{Chunks: make(map[string]ChunkInfo)}
	if _, err := os.Stat(manifestPath); err == nil {
		data, _ := os.ReadFile(manifestPath)
		json.Unmarshal(data, &manifest)
	}

	manifest.Chunks[chunkFilename] = ChunkInfo{
		CreatedAt: time.Now(),
		UserName:  userName,
	}

	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	return os.WriteFile(manifestPath, manifestData, 0644)
}
