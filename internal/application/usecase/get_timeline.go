package usecase

import (
	"context"
	"time"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// TimelineDay represents a group of memories for a specific day.
type TimelineDay struct {
	Date     time.Time
	Memories []*entity.Ursus
}

// GetTimelineUseCase fetches all memories and groups them by day.
type GetTimelineUseCase struct {
	repo repository.UrsusRepository
}

// NewGetTimelineUseCase creates a new instance of GetTimelineUseCase.
func NewGetTimelineUseCase(repo repository.UrsusRepository) *GetTimelineUseCase {
	return &GetTimelineUseCase{repo: repo}
}

// Execute fetches and groups memories chronologically.
func (uc *GetTimelineUseCase) Execute(ctx context.Context) ([]TimelineDay, error) {
	memories, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Group by day (YYYY-MM-DD)
	groups := make(map[string][]*entity.Ursus)
	var days []string // To keep track of order if needed, or we can sort later

	for _, m := range memories {
		dateStr := m.CreatedAt.Format("2006-01-02")
		if _, ok := groups[dateStr]; !ok {
			days = append(days, dateStr)
		}
		groups[dateStr] = append(groups[dateStr], m)
	}

	var timeline []TimelineDay
	for _, d := range days {
		date, _ := time.Parse("2006-01-02", d)
		timeline = append(timeline, TimelineDay{
			Date:     date,
			Memories: groups[d],
		})
	}

	return timeline, nil
}
