package auth

import (
	"errors"

	"github.com/hitpads/reado_ap/internal/utils"
	"gorm.io/gorm"
)

// UserModel represents the subset of the users table needed for authentication
type UserModel struct {
	ID       uint
	Password string
}

// Authenticate validates user credentials and generates a JWT token
func Authenticate(email, password string, db *gorm.DB) (uint, string, error) {
	var user UserModel
	if err := db.Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		return 0, "", errors.New("invalid credentials")
	}

	// Compare hashed password
	if err := utils.ComparePasswords(user.Password, password); err != nil {
		return 0, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}
