package admin

import (
	"encoding/json"
	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
	"net/http"
	"strconv"
)

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type PaginatedResponse struct {
	Data        []UserResponse `json:"data"`
	CurrentPage int            `json:"current_page"`
	NextPage    bool           `json:"next_page"`
	PrevPage    bool           `json:"prev_page"`
}

// GetAllUsers fetches all users and formats the response
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	search := r.URL.Query().Get("search")

	limit := 5
	offset := (page - 1) * limit

	query := db.DB.Preload("Role").Offset(offset).Limit(limit)
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Find(&users)

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role.Name,
		})
	}

	totalUsers := int64(0)
	db.DB.Model(&models.User{}).Count(&totalUsers)

	response := PaginatedResponse{
		Data:        userResponses,
		CurrentPage: page,
		NextPage:    int64(page*limit) < totalUsers,
		PrevPage:    page > 1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUserRole updates a user's role
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID uint `json:"user_id"`
		RoleID uint `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, payload.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.RoleID = payload.RoleID
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User role updated successfully"))
}
