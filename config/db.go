package config

import (
	"fmt"

	"github.com/Sekarfo/P_blog/models"

	log "github.com/sirupsen/logrus"

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
	var exists bool
	err := db.Raw(`
        SELECT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'uni_users_email'
        );
    `).Scan(&exists).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to check for unique constraint")
	}

	if !exists {
		err = db.Exec("ALTER TABLE users ADD CONSTRAINT uni_users_email UNIQUE (email)").Error
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Failed to add unique constraint")
		} else {
			log.Info("Unique constraint 'uni_users_email' added successfully")
		}
	} else {
		log.Info("Unique constraint 'uni_users_email' already exists")
	}

	err = db.AutoMigrate(&models.User{}, &models.Blog{}, &models.Commentary{})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to migrate database")
	}
	log.Info("Database migrated successfully!")
}
