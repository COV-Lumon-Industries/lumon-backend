package ml

import (
	"context"

	"lumon-backend/internal/config"
)

func GetDocumentSummary(ctx context.Context, fileUrl string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	client, err := GetClientWithContext(ctx, *cfg)
	if err != nil {
		return "", err
	}

	response, err := GenerateContentWithFilesText(
		ctx,
		client,
		fileUrl,
		"Can you generate a summary of this document in less than 100 words?",
	)
	if err != nil {
		return "", err
	}

	txt := ResponseToPartString(response)

	return txt, nil
}

func GetDocumentType(ctx context.Context, fileUrl string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	client, err := GetClientWithContext(ctx, *cfg)
	if err != nil {
		return "", err
	}

	response, err := GenerateContentWithFilesText(
		ctx,
		client,
		fileUrl,
		"Can you give me a category word to describe this document?",
	)
	if err != nil {
		return "", err
	}

	txt := ResponseToPartString(response)

	return txt, nil
}
