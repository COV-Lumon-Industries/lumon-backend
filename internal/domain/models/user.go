package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`

	Username    string `gorm:"size:50;unique;not null" json:"username"`
	Password    string `gorm:"size:255;not null" json:"password,omitempty"`
	Email       string `gorm:"size:100;unique;not null" json:"email"`
	UserRole    string `gorm:"type:varchar(20);default:'common'" json:"user_role"`
	PhoneNumber string `gorm:"column:phone_number;unique;not null" json:"phone_number,omitempty"`

	// EmailVerified     bool   `gorm:"default:false" json:"is_verified"`
	// VerificationToken string `gorm:"size:100" json:"verification_token"`

	// PasswordResetToken string     `gorm:"size:100" json:"password_reset_token"`
	// TokenExpiry        *time.Time `json:"token_expiry"`
	// ResetTokenExpiry   *time.Time `json:"reset_token_expiry"`\
	Wallets      []Wallet      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"wallets"`
	Accounts     []Account     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"accounts"`
	LoanRequests []LoanRequest `gorm:"foreignKey:BorrowerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"loan_requests"`
	CreditScore  int64         `gorm:"default:0;not null" json:"credit_score"`
}

func (b *User) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
