package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanRequest struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Amount       int64     `gorm:"not null" json:"amount"`
	BorrowerID   uuid.UUID `gorm:"type:uuid;not null" json:"borrower_id"`
	Borrower     User      `gorm:"foreignKey:BorrowerID" json:"borrower"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	InterestRate float64   `gorm:"type:decimal(5,2);not null" json:"interest_rate"`
	LoanDuration int       `gorm:"not null" json:"loan_duration"`
	Purpose      string    `gorm:"type:varchar(255);not null" json:"purpose"`
	Status       string    `gorm:"type:varchar(50);default:'pending'" json:"status"`
}

func (b *LoanRequest) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
