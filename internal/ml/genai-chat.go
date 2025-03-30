package ml

import (
	"context"
	"iter"

	"lumon-backend/internal/config"

	"google.golang.org/genai"
)

func GetChatResponseStream(ctx context.Context, prompt string) (iter.Seq2[*genai.GenerateContentResponse, error], error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client, err := GetClientWithContext(ctx, *cfg)
	if err != nil {
		return nil, err
	}

	stream := GenerateContentStreamText(ctx, client, prompt)

	return stream, nil
}

func GetChatResponse(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client, err := GetClientWithContext(ctx, *cfg)
	if err != nil {
		return nil, err
	}

	resp, err := GenerateContentText(ctx, client, prompt)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetSearchResponse(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client, err := GetClientWithContext(ctx, *cfg)
	if err != nil {
		return nil, err
	}

	resp, err := Search(ctx, client, prompt)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
