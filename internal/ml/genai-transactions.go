package ml

import (
	"context"
	"encoding/json"

	"lumon-backend/internal/config"
	"lumon-backend/internal/domain/schemas"
)

func GetMoMoTransactionData(ctx context.Context, fileUrl string) (*schemas.MTNMoMoTransactionScrape, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client, err := GetClientWithContext(ctx, *cfg)
	if err != nil {
		return nil, err
	}

	schema := GetTransactionSchema()

	response, err := GenerateContentWithFilesJSON(ctx, client, fileUrl, "Can you grab me all the transaction information from the file?", schema)
	if err != nil {
		return nil, err
	}

	var transaction schemas.MTNMoMoTransactionScrape
	if err := ResponseToStructure(response, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func GetMoMoTransactions(
	ctx context.Context,
	fileUrl string,
) (data []*schemas.MTNMoMoTransactionScrape, batchSize int, err error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, 0, err
	}

	client, err := GetClientWithContext(ctx, *cfg)
	if err != nil {
		return nil, 0, err
	}

	schema := GetTransactionSchema()

	response, err := GenerateContentStreamWithFilesJSON(ctx, client, fileUrl, "Extract the first 20 transactions from the table and return them as a JSON array.", schema)
	if err != nil {
		return nil, 0, err
	}

	blob := IterResponseToString(response)

	var transactions []*schemas.MTNMoMoTransactionScrape
	if err := json.Unmarshal([]byte(blob), &transactions); err != nil {
		return nil, 0, err
	}

	return transactions, len(transactions), nil
}
