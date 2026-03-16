package usecase

import (
	"context"
	"regexp"
	"strings"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

// PassiveCaptureUseCase extracts learnings from text and saves them.
type PassiveCaptureUseCase struct {
	saveUseCase *SaveMemoryUseCase
}

// NewPassiveCaptureUseCase creates a new PassiveCaptureUseCase instance.
func NewPassiveCaptureUseCase(saveUseCase *SaveMemoryUseCase) *PassiveCaptureUseCase {
	return &PassiveCaptureUseCase{
		saveUseCase: saveUseCase,
	}
}

// Execute parses the text for learnings and saves them as memories.
func (uc *PassiveCaptureUseCase) Execute(ctx context.Context, text string) ([]*entity.Ursus, error) {
	var savedMemories []*entity.Ursus

	// 1. Extract learnings using markers
	learnings := uc.extractLearnings(text)

	// 2. Save each learning as a memory
	for _, learning := range learnings {
		mem, err := uc.saveUseCase.Execute(ctx, learning, "passive-capture", "", "", "")
		if err == nil && mem != nil {
			savedMemories = append(savedMemories, mem)
		}
	}

	return savedMemories, nil
}

func (uc *PassiveCaptureUseCase) extractLearnings(text string) []string {
	var results []string

	// Pattern A: <learning>Content</learning>
	tagRegex := regexp.MustCompile(`(?s)<learning>(.*?)</learning>`)
	matches := tagRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 {
			content := strings.TrimSpace(match[1])
			if content != "" {
				results = append(results, content)
			}
		}
	}

	// Pattern B: ### Aprendizajes (Markdown Header)
	// Looks for the section until the next header or end of file
	headerMarker := "### Aprendizajes"
	if idx := strings.Index(text, headerMarker); idx != -1 {
		section := text[idx+len(headerMarker):]
		// Stop at next header or end
		nextHeaderIdx := strings.Index(section, "\n#")
		if nextHeaderIdx != -1 {
			section = section[:nextHeaderIdx]
		}

		// Split by lines and find bullet points or meaningful lines
		lines := strings.Split(section, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// Remove common bullet point markers
			line = strings.TrimPrefix(line, "-")
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
			
			if line != "" {
				results = append(results, line)
			}
		}
	}

	return results
}
