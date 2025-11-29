package database

import (
	"log"

	"aegis.wlbt.nl/aegis-auth/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(silent bool) {
	var err error
	DB, err = gorm.Open(sqlite.Open("aegis.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if !silent {
		log.Println("Database connection established")
	}
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.VerifiedEmail{},
		&models.SessionToken{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}
