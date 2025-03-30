package ml

import (
	"context"
	"encoding/json"
	"io"
	"iter"
	"net/http"

	"lumon-backend/internal/config"

	"github.com/pkg/errors"
	"google.golang.org/genai"
)

const BASEMODEL = "gemini-2.0-flash"

func GetClientWithContext(ctx context.Context, cfg config.Config) (*genai.Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.GeminiAPIKey,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getGenerateConfigJSON(responseSchema *genai.Schema) *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		Temperature:      genai.Ptr[float32](0.2),
		ResponseMIMEType: "application/json",
		StopSequences:    []string{"\n\n"},
		Seed:             genai.Ptr[int32](420),
		ResponseSchema:   responseSchema,
	}
}

func getGenerateConfigText() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		Temperature:       genai.Ptr[float32](0.5),
		TopP:              genai.Ptr[float32](0.5),
		TopK:              genai.Ptr[float32](2.0),
		ResponseMIMEType:  "text/plain",
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "You are a helpful assistant."}}},
	}
}

func getGenerateConfigSearch() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		Temperature:      genai.Ptr[float32](0.2),
		TopP:             genai.Ptr[float32](0.9),
		TopK:             genai.Ptr[float32](2.0),
		ResponseMIMEType: "text/plain",
		Tools:            []*genai.Tool{{GoogleSearch: &genai.GoogleSearch{}}},
	}
}

func GenerateContentStreamText(
	ctx context.Context,
	client *genai.Client,
	prompt string,
) iter.Seq2[*genai.GenerateContentResponse, error] {
	parts := []*genai.Part{
		{Text: prompt},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigText()
	result := client.Models.GenerateContentStream(ctx, BASEMODEL, contents, config)

	return result
}

func GenerateContentText(
	ctx context.Context,
	client *genai.Client,
	prompt string,
) (*genai.GenerateContentResponse, error) {
	parts := []*genai.Part{
		{Text: prompt},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigText()
	result, err := client.Models.GenerateContent(ctx, BASEMODEL, contents, config)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GenerateContentWithFilesText(
	ctx context.Context,
	client *genai.Client,
	fileUrl, prompt string,
) (*genai.GenerateContentResponse, error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return nil, errors.Errorf("Error fetching Document: %s", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Error failed loading Document content: %s", err)
	}

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: data, MIMEType: "application/pdf"}},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigText()
	result, err := client.Models.GenerateContent(ctx, BASEMODEL, contents, config)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GenerateContentWithFilesJSON(
	ctx context.Context,
	client *genai.Client,
	fileUrl, prompt string,
	responseSchema *genai.Schema,
) (*genai.GenerateContentResponse, error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return nil, errors.Errorf("Error fetching Document: %s", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Error failed loading Document content: %s", err)
	}

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: data, MIMEType: "application/json"}},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigText()
	result, err := client.Models.GenerateContent(ctx, BASEMODEL, contents, config)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GenerateContentStreamWithFilesJSON(
	ctx context.Context,
	client *genai.Client,
	fileUrl, prompt string,
	responseSchema *genai.Schema,
) (iter.Seq2[*genai.GenerateContentResponse, error], error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return nil, errors.Errorf("Error fetching Document: %s", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Error failed loading Document content: %s", err)
	}

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: data, MIMEType: "application/pdf"}},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigJSON(responseSchema)
	result := client.Models.GenerateContentStream(ctx, BASEMODEL, contents, config)

	return result, nil
}

func GenerateContentWithImagesText(
	ctx context.Context,
	client *genai.Client,
	imageUrl, prompt string,
) (*genai.GenerateContentResponse, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, errors.Errorf("Error fetching Image: %s", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Error failed loading Image content: %s", err)
	}

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: data, MIMEType: "image/jpeg"}},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigText()
	result, err := client.Models.GenerateContent(ctx, BASEMODEL, contents, config)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GenerateContentWithImagesJSON(
	ctx context.Context,
	client *genai.Client,
	imageUrl, prompt string,
	responseSchema *genai.Schema,
) (*genai.GenerateContentResponse, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, errors.Errorf("Error fetching Image: %s", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Error failed loading Image content: %s", err)
	}

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: data, MIMEType: "image/jpeg"}},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigJSON(responseSchema)
	result, err := client.Models.GenerateContent(ctx, BASEMODEL, contents, config)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ResponseToPartString(response *genai.GenerateContentResponse) string {
	var raw string = ""
	for _, cand := range response.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				raw += part.Text
			}
		}
	}
	return raw
}

func IterResponseToString(response iter.Seq2[*genai.GenerateContentResponse, error]) string {
	next, stop := iter.Pull2(response)
	defer stop()

	var data string = ""
	for {
		resp, err, ok := next()
		if !ok {
			break
		}

		if err != nil {
			break
		}

		data += ResponseToPartString(resp)
	}

	return data
}

func ResponseToStructure(response *genai.GenerateContentResponse, object any) error {
	content := ResponseToPartString(response)

	err := json.Unmarshal([]byte(content), &object)
	if err != nil {
		return err
	}

	return nil
}

func Search(
	ctx context.Context,
	client *genai.Client,
	prompt string,
) (*genai.GenerateContentResponse, error) {
	parts := []*genai.Part{
		{Text: prompt},
	}
	contents := []*genai.Content{{Parts: parts}}

	config := getGenerateConfigSearch()

	result, err := client.Models.GenerateContent(ctx, BASEMODEL, contents, config)
	if err != nil {
		return nil, err
	}

	return result, nil
}
