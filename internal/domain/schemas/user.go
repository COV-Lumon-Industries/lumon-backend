package schemas

import "github.com/google/uuid"

type UserUpdate struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAccountDetails struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=8"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `gorm:"column:phone_number;unique;not null" json:"phone_number"`
}

// type ForgotPasswordRequest struct {
// 	Email string `json:"email" binding:"required,email"`
// }

// type ResetPasswordRequest struct {
// 	Token       string `json:"token" binding:"required"`
// 	NewPassword string `json:"new_password" binding:"required,min=8"`
// }

// type VerificationDetails struct {
// 	Email string `json:"email" binding:"required,email"`
// }

type UserResponse struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Username    string    `gorm:"size:50;unique;not null" json:"username"`
	Password    string    `gorm:"size:255;not null" json:"password,omitempty"`
	Email       string    `gorm:"size:100;unique;not null" json:"email"`
	UserRole    string    `gorm:"type:varchar(20);default:'common'" json:"user_role"`
	PhoneNumber string    `gorm:"column:phone_number;unique;not null" json:"phone_number"`
	CreditScore int64     `json:"credit_score"`
}
