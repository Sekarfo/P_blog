package config

import (
	
	log "github.com/sirupsen/logrus"
	"personal_blog/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
<<<<<<< Updated upstream
    dsn := "host=localhost user=postgres password=0000 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Almaty"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.WithFields(log.Fields{
            "dsn": dsn,
        }).Fatal("Failed to connect to database:", err)
    }
    log.Info("Database connected successfully!")
    return db
=======
	dsn := "host=localhost user=postgres password=postgres dbname=ap_blog port=5432 sslmode=disable TimeZone=Asia/Almaty"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connected successfully!")
	return db
>>>>>>> Stashed changes
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

    if (!exists) {
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

    err = db.AutoMigrate(&models.User{})
    if err != nil {
        log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Failed to migrate database")
    }
    log.Info("Database migrated successfully!")
}
