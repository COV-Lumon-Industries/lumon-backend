package service

import (
	"math"
	"sort"
	"time"

	"lumon-backend/internal/domain/models"
)

type CreditScoreCalculator struct {
	Transactions []models.Transaction
	ScoreRange   struct {
		Min, Max float64
	}
	Weights struct {
		PaymentHistory, IncomeStability, CashFlow,
		TransactionHabits, CreditHistory float64
	}
}

func NewCreditScoreCalculator(tx []models.Transaction) *CreditScoreCalculator {
	return &CreditScoreCalculator{
		Transactions: tx,
		ScoreRange:   struct{ Min, Max float64 }{300, 850},
		Weights: struct {
			PaymentHistory, IncomeStability, CashFlow,
			TransactionHabits, CreditHistory float64
		}{
			PaymentHistory:    0.35,
			IncomeStability:   0.25,
			CashFlow:          0.20,
			TransactionHabits: 0.10,
			CreditHistory:     0.10,
		},
	}
}

func (c *CreditScoreCalculator) Calculate() float64 {
	income, expenses := c.categorizeTransactions()
	firstTx, lastTx := c.getTimeBounds()

	features := struct {
		IncomeStability   float64
		CashFlow          float64
		PaymentBehavior   float64
		TransactionHabits float64
		CreditHistory     float64
	}{
		IncomeStability:   c.calculateIncomeStability(income),
		CashFlow:          c.calculateCashFlow(income, expenses),
		PaymentBehavior:   c.calculatePaymentBehavior(expenses),
		TransactionHabits: c.calculateTransactionHabits(),
		CreditHistory:     c.calculateCreditHistory(firstTx, lastTx),
	}

	rawScore := (features.PaymentBehavior * c.Weights.PaymentHistory) +
		(features.IncomeStability * c.Weights.IncomeStability) +
		(features.CashFlow * c.Weights.CashFlow) +
		(features.TransactionHabits * c.Weights.TransactionHabits) +
		(features.CreditHistory * c.Weights.CreditHistory)

	return c.normalizeScore(rawScore)
}

func (c *CreditScoreCalculator) categorizeTransactions() ([]models.Transaction, []models.Transaction) {
	var income, expenses []models.Transaction
	for _, tx := range c.Transactions {
		if tx.TransactionType == "CASH_IN" {
			income = append(income, tx)
		} else {
			expenses = append(expenses, tx)
		}
	}
	return income, expenses
}

func (c *CreditScoreCalculator) getTimeBounds() (time.Time, time.Time) {
	dates := make([]time.Time, len(c.Transactions))
	for i, tx := range c.Transactions {
		dates[i] = tx.TransactionDate
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })
	return dates[0], dates[len(dates)-1]
}

func (c *CreditScoreCalculator) calculateIncomeStability(income []models.Transaction) float64 {
	monthlyIncome := make(map[time.Month]float64)
	for _, tx := range income {
		month := tx.TransactionDate.Month()
		monthlyIncome[month] += tx.Amount
	}

	var amounts []float64
	for _, amt := range monthlyIncome {
		amounts = append(amounts, amt)
	}

	avg := average(amounts)
	stdDev := standardDeviation(amounts, avg)

	if stdDev == 0 {
		return 100
	}
	return math.Max(0, 100-(stdDev/avg)*100)
}

func (c *CreditScoreCalculator) calculateCashFlow(income, expenses []models.Transaction) float64 {
	totalIncome := 0.0
	for _, tx := range income {
		totalIncome += tx.Amount
	}

	totalExpenses := 0.0
	for _, tx := range expenses {
		totalExpenses += tx.Amount + tx.Fees
	}

	netCashFlow := totalIncome - totalExpenses
	savingsRate := (netCashFlow / totalIncome) * 100

	return math.Min(savingsRate, 100)
}

func (c *CreditScoreCalculator) calculatePaymentBehavior(expenses []models.Transaction) float64 {
	payeeCount := make(map[string]int)
	for _, tx := range expenses {
		payeeCount[tx.ToNumber]++
	}

	recurringPayments := 0
	for _, count := range payeeCount {
		if count > 1 {
			recurringPayments++
		}
	}

	return math.Min(float64(recurringPayments*20), 100)
}

func (c *CreditScoreCalculator) calculateTransactionHabits() float64 {
	typeCount := make(map[string]int)
	for _, tx := range c.Transactions {
		typeCount[tx.TransactionType]++
	}
	return math.Min(float64(len(typeCount))*25, 100)
}

func (c *CreditScoreCalculator) calculateCreditHistory(first, last time.Time) float64 {
	days := last.Sub(first).Hours() / 24
	return math.Min(float64(days)/365*100, 100)
}

func (c *CreditScoreCalculator) normalizeScore(raw float64) float64 {
	return (raw/100)*(c.ScoreRange.Max-c.ScoreRange.Min) + c.ScoreRange.Min
}

func average(nums []float64) float64 {
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	return sum / float64(len(nums))
}

func standardDeviation(nums []float64, mean float64) float64 {
	if len(nums) < 2 {
		return 0
	}

	sum := 0.0
	for _, n := range nums {
		sum += math.Pow(n-mean, 2)
	}
	return math.Sqrt(sum / float64(len(nums)-1))
}
