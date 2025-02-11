package database

import (
	"log"

	"github.com/ctf/api/internal/config"
	"github.com/ctf/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes and returns a GORM database instance
func InitDB(config config.Config ) *gorm.DB {
	dsn := config.DATABASEURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// AutoMigrate will create the teams table if it doesn’t exist
	err = db.AutoMigrate(
		&models.Team{},
	    &models.Challenge{},
		&models.Submission{},
		&models.Score{},
       )
	if err != nil {
		log.Fatalf("❌ AutoMigration failed: %v", err)
	}

	log.Println("✅ Database connected and migrated successfully")
	return db
}
