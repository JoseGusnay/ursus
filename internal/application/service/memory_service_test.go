package service_test

import (
	"context"
	"testing"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

type mockRepo struct{}

func (m *mockRepo) Save(ctx context.Context, u *entity.Ursus) error { return nil }
func (m *mockRepo) Search(ctx context.Context, q string) ([]*entity.Ursus, error) { return nil, nil }
func (m *mockRepo) List(ctx context.Context) ([]*entity.Ursus, error) { return nil, nil }
func (m *mockRepo) Delete(ctx context.Context, id string) error { return nil }

func TestMemoryService_Store(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewMemoryService(repo)

	u, err := svc.Store(context.Background(), "test content", "test metadata")
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	if u.Content != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", u.Content)
	}
	
	if u.ID == "" {
		t.Error("Expected generated ID, got empty string")
	}
}
