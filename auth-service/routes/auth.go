package routes

import (
	"auth-service/database"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"

	"auth-service/middleware"
	"auth-service/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register User
func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// Assign default role
	user.RoleID = 2 // 'Reader' role

	// Save to DB
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
	}

	// Generate JWT Token
	token, _ := middleware.GenerateToken(user.ID, user.RoleID)

	return c.JSON(fiber.Map{"message": "User registered successfully", "token": token})
}

// Login User
func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	err := database.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Ensure email is verified before allowing login
	if !user.IsEmailVerified {
		return c.Status(403).JSON(fiber.Map{"error": "Email not verified. Please check your inbox."})
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, _ := middleware.GenerateToken(user.ID, user.RoleID)

	return c.JSON(fiber.Map{"token": token})
}

// Forgot Password
func ForgotPassword(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if user exists
	var user models.User
	err := database.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
	user.ResetToken = hex.EncodeToString(tokenBytes)

	// Update reset token in DB
	database.DB.Model(&user).Update("reset_token", user.ResetToken)

	// Send reset email
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), os.Getenv("SMTP_HOST"))
	msg := []byte(fmt.Sprintf("Subject: Password Reset\n\nClick here: http://localhost:8081/reset-password?token=%s", user.ResetToken))
	err = smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, os.Getenv("SMTP_USER"), []string{input.Email}, msg)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Reset link sent!"})
}

// Reset Password
func ResetPassword(c *fiber.Ctx) error {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find user with the reset token
	var user models.User
	err := database.DB.Where("reset_token = ?", input.Token).First(&user).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid reset token"})
	}

	// Hash new password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	// Update user password and clear reset token
	database.DB.Model(&user).Updates(map[string]interface{}{
		"password":    string(hashedPassword),
		"reset_token": "",
	})

	return c.JSON(fiber.Map{"message": "Password reset successful"})
}
