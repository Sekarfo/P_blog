package utils

import (
	"fmt"
	"personal_blog/models"
	"regexp"
)

func ValidateUserInput(user *models.User) error {
	if user.Name == "" || user.Email == "" {
		return fmt.Errorf("Both 'name' and 'email' are required")
	}

	// Check name contains only letters (English/Cyrillic)
	nameRegex := `^[a-zA-Zа-яА-ЯёЁ]+$`
	matchedName, err := regexp.MatchString(nameRegex, user.Name)
	if err != nil || !matchedName {
		return fmt.Errorf("Name must contain only letters")
	}

	// Email format
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matchedEmail, err := regexp.MatchString(emailRegex, user.Email)
	if err != nil || !matchedEmail {
		return fmt.Errorf("Invalid email format")
	}

	return nil
}
