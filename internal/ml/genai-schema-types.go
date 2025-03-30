package ml

import "google.golang.org/genai"

func GetTransactionSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"transaction_date": {
					Type: genai.TypeString,
				},
				"from_account": {
					Type: genai.TypeString,
				},
				"from_name": {
					Type: genai.TypeString,
				},
				"from_number": {
					Type: genai.TypeString,
				},
				"transaction_type": {
					Type: genai.TypeString,
				},
				"amount": {
					Type: genai.TypeNumber,
				},
				"fees": {
					Type: genai.TypeNumber,
				},
				"balance_before": {
					Type: genai.TypeNumber,
				},
				"balance_after": {
					Type: genai.TypeNumber,
				},
				"to_number": {
					Type: genai.TypeString,
				},
				"to_name": {
					Type: genai.TypeString,
				},
				"to_account": {
					Type: genai.TypeString,
				},
				"reference": {
					Type: genai.TypeString,
				},
			},
			Required: []string{
				"transaction_date",
				"from_account",
				"from_name",
				"from_number",
				"transaction_type",
				"amount",
				"fees",
				"balance_before",
				"balance_after",
				"to_number",
				"to_name",
				"to_account",
				"reference",
			},
		},
	}
}
