package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
	"github.com/hitpads/reado_ap/internal/utils"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	var role models.Role
	if err := db.DB.Where("name = ?", "Reader").First(&role).Error; err != nil {
		http.Error(w, "Role not found", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		RoleID:   role.ID,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Generate verification token
	token := utils.GenerateRandomToken(32)
	verificationToken := models.VerificationToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Token valid for 24 hours
	}

	if err := db.DB.Create(&verificationToken).Error; err != nil {
		http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
		return
	}

	// Send verification email
	go utils.SendVerificationEmail(user.Email, token)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Registration successful! Please check your email to verify your account."))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Preload("Role").Where("email = ?", creds.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := utils.ComparePasswords(user.Password, creds.Password); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token": token,
		"role":  user.Role.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
