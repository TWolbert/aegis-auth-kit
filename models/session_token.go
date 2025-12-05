package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (s *SessionToken) OwnedByIp(c *fiber.Ctx) bool {
	return s.IPAddress == c.IP()
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func CreateToken(ctx context.Context, db *gorm.DB, user User, expiry time.Time, ipAddr string, userAgent string) (*SessionToken, error) {
	token := generateSecureToken(32)
	sessionToken := SessionToken{UserID: user.ID, Token: token, ExpiresAt: expiry, IPAddress: ipAddr, UserAgent: userAgent}

	err := gorm.G[SessionToken](db).Create(ctx, &sessionToken)

	if err != nil {
		return nil, err
	}

	return &sessionToken, nil
}

func (s *SessionToken) Delete(ctx context.Context, db *gorm.DB) {
	gorm.G[SessionToken](db).Where("token = ?", s.Token).Delete(ctx)
}

func GetUserByToken(ctx context.Context, db *gorm.DB, token string) (*User, *SessionToken, error) {
	data, err := gorm.G[SessionToken](db).Preload("User", nil).Where("token = ?", token).First(ctx)

	if err != nil {
		return nil, nil, err
	}

	return &data.User, &data, nil
}
