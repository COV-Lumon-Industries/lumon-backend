package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Type string    `gorm:"type:varchar(255);not null" json:"type"`
	Name string    `gorm:"type:varchar(255);not null" json:"name"`

	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`

	DocumentID uuid.UUID `gorm:"type:uuid;not null" json:"document_id"`
	Document   Document  `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"document"`
}

func (b *Account) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
