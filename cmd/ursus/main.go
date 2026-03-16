package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/application/usecase"
	dservice "github.com/JoseGusnay/ursus/internal/domain/service"
	"github.com/JoseGusnay/ursus/internal/infrastructure/mcp"
	"github.com/JoseGusnay/ursus/internal/infrastructure/storage"
	"github.com/JoseGusnay/ursus/internal/interfaces/cli"
	_ "modernc.org/sqlite"
)

func main() {
	// Ensure data directory exists
	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".ursus")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	dbPath := filepath.Join(dbDir, "ursus.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	// Initialize storage
	repo := storage.NewSQLiteUrsusRepository(db)
	sessionRepo := storage.NewSQLiteSessionRepository(db)
	promptRepo := storage.NewSQLitePromptRepository(db)
	
	jsonlPath := filepath.Join(home, ".ursus", "memories.jsonl")
	jsonlRepo := storage.NewJSONLUrsusRepository(jsonlPath)

	if err := repo.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	privacySvc := dservice.NewPrivacyService()

	// Initialize use cases and services
	saveUC := usecase.NewSaveMemoryUseCase(repo, sessionRepo, privacySvc, promptRepo)
	searchUC := usecase.NewSearchMemoryUseCase(repo)
	syncUC := usecase.NewSyncMemoriesUseCase(repo, jsonlRepo)
	suggestUC := usecase.NewSuggestTopicUseCase(repo)
	timelineUC := usecase.NewGetTimelineUseCase(repo)
	summarizeUC := usecase.NewSummarizeSessionUseCase(repo, sessionRepo)
	getDetailUC := usecase.NewGetMemoryDetailUseCase(repo)
	passiveUC := usecase.NewPassiveCaptureUseCase(saveUC)
	statsUC := usecase.NewGetStatsUseCase(repo, sessionRepo, promptRepo)

	// Initialize Services
	memorySvc := service.NewMemoryService(repo, saveUC, searchUC)
	sessionSvc := service.NewSessionService(sessionRepo)

	// MCP Server Mode
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		mcpServer := mcp.NewUrsusMCPServer(memorySvc, sessionSvc, suggestUC, timelineUC, summarizeUC, getDetailUC, passiveUC, statsUC)
		log.Println("Ursus MCP Server starting on stdio...")
		if err := mcpServer.Serve(); err != nil {
			log.Fatalf("MCP server error: %v", err)
		}
		return
	}

	// CLI Mode
	cli.Execute(memorySvc, sessionSvc, syncUC, suggestUC, timelineUC, summarizeUC, getDetailUC, statsUC)
}
