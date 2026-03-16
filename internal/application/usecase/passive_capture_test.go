package usecase_test

import (
	"context"
	"testing"

	"github.com/JoseGusnay/ursus/internal/application/usecase"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
	"github.com/JoseGusnay/ursus/internal/domain/service"
)

type mockUrsusRepo struct {
	repository.UrsusRepository
}
func (m *mockUrsusRepo) Save(ctx context.Context, u *entity.Ursus) error { return nil }
func (m *mockUrsusRepo) List(ctx context.Context) ([]*entity.Ursus, error) { return nil, nil }
func (m *mockUrsusRepo) GetByTopicKey(ctx context.Context, k string) (*entity.Ursus, error) { return nil, nil }

type mockSessionRepo struct {
	repository.SessionRepository
}
func (m *mockSessionRepo) GetActive(ctx context.Context) (*entity.Session, error) { return nil, nil }

type mockPromptRepo struct {
	repository.PromptRepository
}
func (m *mockPromptRepo) Save(ctx context.Context, p *entity.Prompt) error { return nil }

func TestPassiveCapture_Execute(t *testing.T) {
	repo := &mockUrsusRepo{}
	sessRepo := &mockSessionRepo{}
	promptRepo := &mockPromptRepo{}
	privacy := service.NewPrivacyService()
	
	saveUC := usecase.NewSaveMemoryUseCase(repo, sessRepo, privacy, promptRepo)
	pc := usecase.NewPassiveCaptureUseCase(saveUC)

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "Tags extraction",
			text:     "Some text <learning>Important fact</learning> more text <learning>Another fact</learning>",
			expected: 2,
		},
		{
			name:     "Markdown extraction",
			text:     "Summary of work.\n### Aprendizajes\n- Need to use Gzip\n- SQLite is fast\n\n### Next Steps",
			expected: 2,
		},
		{
			name:     "No learnings",
			text:     "Just plain text here.",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mems, err := pc.Execute(context.Background(), tt.text)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}
			if len(mems) != tt.expected {
				t.Errorf("Expected %d memories, got %d", tt.expected, len(mems))
			}
		})
	}
}
