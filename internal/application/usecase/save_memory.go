package usecase

import (
	"context"
	"time"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
	"github.com/JoseGusnay/ursus/internal/domain/service"
	"github.com/google/uuid"
)

// SaveMemoryUseCase handles the logic for storing new memories in Ursus.
type SaveMemoryUseCase struct {
	repo        repository.UrsusRepository
	sessionRepo repository.SessionRepository
	privacySvc  *service.PrivacyService
	promptRepo  repository.PromptRepository
}

// NewSaveMemoryUseCase creates a new SaveMemoryUseCase instance.
func NewSaveMemoryUseCase(repo repository.UrsusRepository, sessionRepo repository.SessionRepository, privacySvc *service.PrivacyService, promptRepo repository.PromptRepository) *SaveMemoryUseCase {
	return &SaveMemoryUseCase{
		repo:        repo,
		sessionRepo: sessionRepo,
		privacySvc:  privacySvc,
		promptRepo:  promptRepo,
	}
}

// Execute performs the save operation with hygiene logic and optional prompt logging.
func (u *SaveMemoryUseCase) Execute(ctx context.Context, content, metadata, topicKey, scope, promptText string) (*entity.Ursus, error) {
	// 0. Redaction Logic
	content = u.privacySvc.Redact(content)

	if scope == "" {
		scope = entity.ScopeProject
	}

	var sessionID string
	if active, err := u.sessionRepo.GetActive(ctx); err == nil && active != nil {
		sessionID = active.ID
	}

	// 0.5 Prompt Logging Logic
	var promptID string
	if promptText != "" {
		p := &entity.Prompt{
			ID:        uuid.New().String(),
			Input:     promptText,
			SessionID: sessionID,
			CreatedAt: time.Now(),
		}
		if err := u.promptRepo.Save(ctx, p); err == nil {
			promptID = p.ID
		}
	}

	// 1. Topic Upsert Logic
	if topicKey != "" {
		existing, err := u.repo.GetByTopicKey(ctx, topicKey)
		if err == nil && existing != nil {
			existing.Content = content
			existing.Metadata = metadata
			existing.RevisionCount++
			existing.PromptID = promptID
			existing.UpdatedAt = time.Now()
			existing.LastSeenAt = time.Now()
			if err := u.repo.Update(ctx, existing); err != nil {
				return nil, err
			}
			return existing, nil
		}
	}

	// 2. Exact Deduplication (Dynamic: increment DuplicateCount)
	memories, err := u.repo.List(ctx)
	if err == nil {
		// Look for duplicate in recent memories
		for _, m := range memories {
			if m.Content == content && m.Scope == scope {
				m.DuplicateCount++
				m.UpdatedAt = time.Now()
				m.LastSeenAt = time.Now()
				m.PromptID = promptID
				if err := u.repo.Update(ctx, m); err != nil {
					return nil, err
				}
				return m, nil
			}
		}
	}

	// 3. Normal Save
	mem := &entity.Ursus{
		ID:             uuid.New().String(),
		Content:        content,
		Metadata:       metadata,
		SessionID:      sessionID,
		TopicKey:       topicKey,
		PromptID:       promptID,
		Scope:          scope,
		DuplicateCount: 1,
		RevisionCount:  0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastSeenAt:     time.Now(),
	}

	if err := u.repo.Save(ctx, mem); err != nil {
		return nil, err
	}

	return mem, nil
}
