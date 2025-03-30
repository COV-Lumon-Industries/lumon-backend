package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ContentSummary string     `gorm:"type:text;unique;not null" json:"content_summary"`
	Type           string     `gorm:"size:50;not null" json:"type"`
	UploadedAt     *time.Time `json:"uploaded_at"`

	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
}

func (b *Document) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
