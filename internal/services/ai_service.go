package services

import (
	"context"

	"github.com/yourname/moodle/internal/ai"
)

type AIService struct {
	client *ai.GeminiClient
}

func NewAIService(client *ai.GeminiClient) *AIService {
	return &AIService{client: client}
}

func (s *AIService) Ask(ctx context.Context, query string) (string, error) {
	return s.client.Ask(ctx, query)
}
