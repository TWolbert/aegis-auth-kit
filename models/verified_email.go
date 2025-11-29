package models

import (
	"time"

	"gorm.io/gorm"
)

type VerifiedEmail struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null;index" json:"user_id"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	VerifiedAt     time.Time      `gorm:"not null" json:"verified_at"`
	VerificationToken string      `gorm:"index" json:"-"` // Token used for verification
	TokenExpiresAt *time.Time     `json:"-"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
