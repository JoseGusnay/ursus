package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// StatsReport represents the aggregated statistics of the Ursus system.
type StatsReport struct {
	TotalMemories     int
	TotalSessions     int
	TotalPrompts      int
	TopTopics         []string
	Last7DaysActivity map[string]int
}

type GetStatsUseCase struct {
	repo        repository.UrsusRepository
	sessionRepo repository.SessionRepository
	promptRepo  repository.PromptRepository
}

func NewGetStatsUseCase(repo repository.UrsusRepository, sessionRepo repository.SessionRepository, promptRepo repository.PromptRepository) *GetStatsUseCase {
	return &GetStatsUseCase{
		repo:        repo,
		sessionRepo: sessionRepo,
		promptRepo:  promptRepo,
	}
}

func (uc *GetStatsUseCase) Execute(ctx context.Context) (*StatsReport, error) {
	memories, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	sessions, err := uc.sessionRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Note: PromptRepository doesn't have a List method yet in the interface, 
	// but we can estimate or skip for now if not critical. 
	// Let's assume we want it, so we might need to add it or just count memories with PromptID.
	promptCount := 0
	topicCounts := make(map[string]int)
	activity := make(map[string]int)
	
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	for _, m := range memories {
		if m.PromptID != "" {
			promptCount++
		}
		if m.TopicKey != "" {
			topicCounts[m.TopicKey]++
		}
		
		if m.CreatedAt.After(sevenDaysAgo) {
			dateStr := m.CreatedAt.Format("2006-01-02")
			activity[dateStr]++
		}
	}

	// Get Top Topics
	type topicStat struct {
		name  string
		count int
	}
	var topics []topicStat
	for k, v := range topicCounts {
		topics = append(topics, topicStat{k, v})
	}
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].count > topics[j].count
	})

	topTopics := []string{}
	for i := 0; i < len(topics) && i < 5; i++ {
		topTopics = append(topTopics, topics[i].name)
	}

	return &StatsReport{
		TotalMemories:     len(memories),
		TotalSessions:     len(sessions),
		TotalPrompts:      promptCount, // Approximate from memories for now
		TopTopics:         topTopics,
		Last7DaysActivity: activity,
	}, nil
}
