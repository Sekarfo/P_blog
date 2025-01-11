package config

import (
	"fmt"
	"log"
	"personal_blog/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=0000 dbname=postgres port=5433 sslmode=disable TimeZone=Asia/Almaty"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connected successfully!")
	return db
}

func AutoMigrateDB(db *gorm.DB) {
	// Check if the unique constraint exists before adding it
	var exists bool
	err := db.Raw(`
        SELECT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'uni_users_email'
        );
    `).Scan(&exists).Error
	if err != nil {
		log.Fatal("Failed to check for unique constraint:", err)
	}

	if !exists {
		err = db.Exec("ALTER TABLE users ADD CONSTRAINT uni_users_email UNIQUE (email)").Error
		if err != nil {
			log.Fatal("Failed to add unique constraint:", err)
		}
	}

	// Perform the migration
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migrated successfully!")
}
