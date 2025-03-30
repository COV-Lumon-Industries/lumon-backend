package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TransactionDate time.Time `gorm:"type:timestamp;not null" json:"transaction_date"`
	FromAccount     string    `gorm:"type:varchar(255);not null" json:"from_account"`
	FromName        string    `gorm:"type:varchar(255);not null" json:"from_name"`
	FromNumber      string    `gorm:"type:varchar(255);not null" json:"from_number"`
	TransactionType string    `gorm:"type:varchar(50);not null" json:"transaction_type"`
	Amount          float64   `gorm:"type:decimal(20,2);not null" json:"amount"`
	Fees            float64   `gorm:"type:decimal(20,2);not null" json:"fees"`
	ELevy           float64   `gorm:"type:decimal(20,2);not null" json:"e_levy"`
	BalanceBefore   float64   `gorm:"type:decimal(20,2);not null" json:"balance_before"`
	BalanceAfter    float64   `gorm:"type:decimal(20,2);not null" json:"balance_after"`
	ToNumber        string    `gorm:"type:varchar(255);not null" json:"to_number"`
	ToName          string    `gorm:"type:varchar(255);not null" json:"to_name"`
	ToAccount       string    `gorm:"type:varchar(255);not null" json:"to_account"`
	Reference       string    `gorm:"type:varchar(255)" json:"reference"`

	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
}

func (b *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
