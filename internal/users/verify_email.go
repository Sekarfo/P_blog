package users

import (
	"net/http"
	"time"

	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
)

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	var verificationToken models.VerificationToken
	if err := db.DB.Where("token = ?", token).First(&verificationToken).Error; err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	if verificationToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Token has expired", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, verificationToken.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	user.IsEmailVerified = true
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}

	// Delete the verification token
	db.DB.Delete(&verificationToken)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully! You can now log in."))
}
