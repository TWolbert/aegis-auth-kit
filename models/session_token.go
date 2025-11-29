package models

import (
	"time"

	"gorm.io/gorm"
)

type SessionToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expires_at"`
	IPAddress string         `gorm:"size:45" json:"ip_address"` // IPv4 or IPv6
	UserAgent string         `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// IsExpired checks if the session token has expired
func (s *SessionToken) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
