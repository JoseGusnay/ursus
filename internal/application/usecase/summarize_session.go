package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/JoseGusnay/ursus/internal/domain/repository"
)

// SessionReview represents the summarized context of a session.
type SessionReview struct {
	Session  *entity.Session
	Memories []*entity.Ursus
	Summary  string
}

type SummarizeSessionUseCase struct {
	repo        repository.UrsusRepository
	sessionRepo repository.SessionRepository
}

func NewSummarizeSessionUseCase(repo repository.UrsusRepository, sessionRepo repository.SessionRepository) *SummarizeSessionUseCase {
	return &SummarizeSessionUseCase{
		repo:        repo,
		sessionRepo: sessionRepo,
	}
}

func (uc *SummarizeSessionUseCase) Execute(ctx context.Context, sessionID string) (*SessionReview, error) {
	var session *entity.Session
	var err error

	if sessionID == "" {
		// Get active session
		session, err = uc.sessionRepo.GetActive(ctx)
		if err != nil {
			return nil, err
		}
		if session == nil {
			// If no active, get the last one
			sessions, err := uc.sessionRepo.List(ctx)
			if err != nil {
				return nil, err
			}
			if len(sessions) > 0 {
				session = sessions[0] // Assuming List returns ordered by start_time DESC
			}
		}
	} else {
		session, err = uc.sessionRepo.GetByID(ctx, sessionID)
		if err != nil {
			return nil, err
		}
	}

	if session == nil {
		return nil, fmt.Errorf("no session found to summarize")
	}

	memories, err := uc.repo.ListBySession(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	// For the local version, "Summary" is a formatted string of memories
	// that an LLM can easily consume.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Session: %s (Started: %s)\n", session.Title, session.StartTime.Format("2006-01-02 15:04")))
	if len(memories) == 0 {
		sb.WriteString("No memories recorded in this session.")
	} else {
		for _, m := range memories {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", m.CreatedAt.Format("15:04"), m.Content))
		}
	}

	return &SessionReview{
		Session:  session,
		Memories: memories,
		Summary:  sb.String(),
	}, nil
}
