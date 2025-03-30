package schemas

type CreateLoanRequestDetails struct {
	Amount       int64   `gorm:"not null" json:"amount"`
	InterestRate float64 `gorm:"type:decimal(5,2);not null" json:"interest_rate"`
	LoanDuration int     `gorm:"not null" json:"loan_duration"`
	Purpose      string  `gorm:"type:varchar(255);not null" json:"purpose"`
}
