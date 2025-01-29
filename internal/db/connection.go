package db

import (
	"log"
	"os"

	"github.com/hitpads/reado_ap/internal/models"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db

	db.AutoMigrate(&models.Role{}, &models.User{}, &models.Post{}, &models.Comment{}, &models.Like{}, &models.VerificationToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedRoles()

	log.Println("Database connected and migrations applied successfully!")
}

func seedRoles() {
	defaultRoles := []string{"Admin", "Writer", "Reader"}
	for _, roleName := range defaultRoles {
		var role models.Role
		if err := DB.Where("name = ?", roleName).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				DB.Create(&models.Role{Name: roleName})
				log.Printf("Role %s created", roleName)
			} else {
				log.Printf("Error checking role %s: %v", roleName, err)
			}
		}
	}
}
