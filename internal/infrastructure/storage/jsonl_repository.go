package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

type JSONLUrsusRepository struct {
	filePath string
}

func NewJSONLUrsusRepository(filePath string) *JSONLUrsusRepository {
	return &JSONLUrsusRepository{filePath: filePath}
}

func (r *JSONLUrsusRepository) Save(ctx context.Context, u *entity.Ursus) error {
	f, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	_, err = f.Write(append(data, '\n'))
	return err
}

func (r *JSONLUrsusRepository) SaveAll(ctx context.Context, memories []*entity.Ursus) error {
	f, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, u := range memories {
		data, err := json.Marshal(u)
		if err != nil {
			return err
		}
		if _, err := f.Write(append(data, '\n')); err != nil {
			return err
		}
	}
	return nil
}

func (r *JSONLUrsusRepository) List(ctx context.Context) ([]*entity.Ursus, error) {
	f, err := os.Open(r.filePath)
	if os.IsNotExist(err) {
		return []*entity.Ursus{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var memories []*entity.Ursus
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var u entity.Ursus
		if err := json.Unmarshal(line, &u); err != nil {
			continue // Skip malformed lines
		}
		memories = append(memories, &u)
	}

	return memories, nil
}

func (r *JSONLUrsusRepository) ListBySession(ctx context.Context, sessionID string) ([]*entity.Ursus, error) {
	memories, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var results []*entity.Ursus
	for _, m := range memories {
		if m.SessionID == sessionID {
			results = append(results, m)
		}
	}
	return results, nil
}

func (r *JSONLUrsusRepository) GetByID(ctx context.Context, id string) (*entity.Ursus, error) {
	memories, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range memories {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, nil
}

func (r *JSONLUrsusRepository) GetByTopicKey(ctx context.Context, topicKey string) (*entity.Ursus, error) {
	memories, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range memories {
		if m.TopicKey == topicKey {
			return m, nil
		}
	}
	return nil, nil
}

func (r *JSONLUrsusRepository) Update(ctx context.Context, u *entity.Ursus) error {
	memories, err := r.List(ctx)
	if err != nil {
		return err
	}
	found := false
	for i, m := range memories {
		if m.ID == u.ID {
			memories[i] = u
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("memory not found for update")
	}
	return r.SaveAll(ctx, memories)
}

func (r *JSONLUrsusRepository) Search(ctx context.Context, query string) ([]*entity.Ursus, error) {
	// For simplicity, JSONL search will be a basic filter for now.
	// In a real scenario, we'd probably rely on the SQLite index.
	return nil, fmt.Errorf("search not implemented for JSONL repository")
}

func (r *JSONLUrsusRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("delete not implemented for JSONL repository directly")
}

func (r *JSONLUrsusRepository) Migrate(ctx context.Context) error {
	return nil
}
