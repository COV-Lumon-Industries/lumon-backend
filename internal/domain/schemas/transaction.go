package schemas

import "time"

type TransactionResponse struct {
	TransactionDate time.Time `json:"transaction_date"`
	FromName        string    `json:"from_name"`
	FromNumber      string    `json:"from_number"`
	TransactionType string    `json:"transaction_type"`
	BalanceBefore   float64   `json:"balance_before"`
	BalanceAfter    float64   `json:"balance_after"`
	Amount          float64   `json:"amount"`
	ToNumber        string    `json:"to_number"`
	ToName          string    `json:"to_name"`
	Reference       string    `json:"reference"`
	UserID          string    `json:"user_id"`
}

// type MoMoTimeFormat struct {
// 	time.Time
// }

// func (ct *MoMoTimeFormat) UnmarshalJSON(data []byte) error {
// 	str := string(data[1 : len(data)-1])

// 	layout := "02-Jan-2006 03:04:05 PM"

// 	parsedTime, err := time.Parse(layout, str)
// 	if err != nil {
// 		return err
// 	}

// 	ct.Time = parsedTime
// 	return nil
// }

type MTNMoMoTransactionScrape struct {
	TransactionDate string  `json:"transaction_date"`
	FromAccount     string  `json:"from_account"`
	FromName        string  `json:"from_name"`
	FromNumber      string  `json:"from_number"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	Fees            float64 `json:"fees"`
	BalanceBefore   float64 `json:"balance_before"`
	BalanceAfter    float64 `json:"balance_after"`
	ToNumber        string  `json:"to_number"`
	ToName          string  `json:"to_name"`
	ToAccount       string  `json:"to_account"`
	Reference       string  `json:"reference"`
}
