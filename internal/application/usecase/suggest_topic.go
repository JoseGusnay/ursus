package usecase

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

type SuggestTopicUseCase struct {
	repo repository.UrsusRepository
}

func NewSuggestTopicUseCase(repo repository.UrsusRepository) *SuggestTopicUseCase {
	return &SuggestTopicUseCase{repo: repo}
}

func (u *SuggestTopicUseCase) Execute(ctx context.Context) ([]string, error) {
	memories, err := u.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch memories: %w", err)
	}

	if len(memories) == 0 {
		return []string{}, nil
	}

	// Simple frequency analysis
	wordFreq := make(map[string]int)
	stopWords := map[string]bool{
		"the": true, "and": true, "a": true, "to": true, "of": true, "in": true, "is": true, "it": true,
		"you": true, "that": true, "he": true, "was": true, "for": true, "on": true, "are": true, "with": true,
		"as": true, "i": true, "his": true, "they": true, "be": true, "at": true, "one": true, "have": true,
		"this": true, "from": true, "or": true, "had": true, "by": true, "hot": true, "word": true, "but": true,
		"some": true, "what": true, "there": true, "we": true, "can": true, "out": true, "other": true, "were": true,
		"all": true, "when": true, "up": true, "use": true, "your": true, "how": true, "said": true,
		"an": true, "each": true, "she": true,
	}

	re := regexp.MustCompile(`\w+`)
	for _, m := range memories {
		text := strings.ToLower(m.Content)
		words := re.FindAllString(text, -1)
		for _, w := range words {
			if len(w) > 3 && !stopWords[w] {
				wordFreq[w]++
			}
		}
	}

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range wordFreq {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	limit := 10
	if len(ss) < limit {
		limit = len(ss)
	}

	var result []string
	for i := 0; i < limit; i++ {
		result = append(result, ss[i].Key)
	}

	return result, nil
}
