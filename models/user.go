package models

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // Password hash, excluded from JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	VerifiedEmails []VerifiedEmail `gorm:"foreignKey:UserID" json:"verified_emails,omitempty"`
	SessionTokens  []SessionToken  `gorm:"foreignKey:UserID" json:"-"`
}

func (user *User) Update(ctx context.Context, db *gorm.DB, username string, email string, password string) (bool, error) {
	newUser := User{}

	if username != "" {
		newUser.Username = username
	}

	if email != "" {
		newUser.Email = email
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err == nil {
			newUser.Password = string(hash)
		}
	}

	_, err := gorm.G[User](db).Where("id = ?", user.ID).Updates(ctx, newUser)

	if err != nil {
		return false, err
	}
	return true, nil
}
