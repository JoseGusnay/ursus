package service_test

import (
	"context"
	"testing"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/application/usecase"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
	dservice "github.com/JoseGusnay/ursus/internal/domain/service"
)

type mockUrsusRepository struct {
	repository.UrsusRepository
}

func (m *mockUrsusRepository) Save(ctx context.Context, u *entity.Ursus) error { return nil }
func (m *mockUrsusRepository) Search(ctx context.Context, q string) ([]*entity.Ursus, error) {
	return nil, nil
}
func (m *mockUrsusRepository) List(ctx context.Context) ([]*entity.Ursus, error) { return nil, nil }
func (m *mockUrsusRepository) ListBySession(ctx context.Context, s string) ([]*entity.Ursus, error) {
	return nil, nil
}
func (m *mockUrsusRepository) GetByID(ctx context.Context, id string) (*entity.Ursus, error) { return nil, nil }
func (m *mockUrsusRepository) GetByTopicKey(ctx context.Context, t string) (*entity.Ursus, error) {
	return nil, nil
}
func (m *mockUrsusRepository) Update(ctx context.Context, u *entity.Ursus) error { return nil }
func (m *mockUrsusRepository) Delete(ctx context.Context, id string) error           { return nil }

type mockSessionRepository struct {
	repository.SessionRepository
}

func (m *mockSessionRepository) GetActive(ctx context.Context) (*entity.Session, error) {
	return nil, nil
}

func (m *mockSessionRepository) DeactivateAll(ctx context.Context) error {
	return nil
}

func (m *mockSessionRepository) Save(ctx context.Context, s *entity.Session) error {
	return nil
}

type mockPromptRepository struct {
	repository.PromptRepository
}

func (m *mockPromptRepository) Save(ctx context.Context, p *entity.Prompt) error { return nil }

func TestMemoryService_Store(t *testing.T) {
	mockRepo := &mockUrsusRepository{}
	mockSessionRepo := &mockSessionRepository{}
	mockPromptRepo := &mockPromptRepository{}
	privacySvc := dservice.NewPrivacyService()
	saveUseCase := usecase.NewSaveMemoryUseCase(mockRepo, mockSessionRepo, privacySvc, mockPromptRepo)
	searchUseCase := usecase.NewSearchMemoryUseCase(mockRepo)
	svc := service.NewMemoryService(mockRepo, saveUseCase, searchUseCase)

	content := "test content"
	metadata := "test metadata"
	u, err := svc.Store(context.Background(), content, metadata, "", "", "")
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	if u.Content != content {
		t.Errorf("Expected content %s, got %s", content, u.Content)
	}
	
	if u.ID == "" {
		t.Error("Expected generated ID, got empty string")
	}
}
