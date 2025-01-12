package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/Sekarfo/P_blog/models"
	"github.com/Sekarfo/P_blog/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func FetchArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

/*
// CreateUser handles the creation of a new user.
func CreateUser(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if err := r.ParseForm(); err != nil {
            log.WithFields(log.Fields{
                "error": err,
            }).Error("Invalid input")
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        user := models.User{
            Name:     r.FormValue("name"),
            Email:    r.FormValue("email"),
            Password: r.FormValue("password"),
        }

<<<<<<< Updated upstream
        if err := utils.ValidateUserInput(&user); err != nil {
            log.WithFields(log.Fields{
                "user": user,
                "error": err,
            }).Error("Validation failed")
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            log.WithFields(log.Fields{
                "error": err,
            }).Error("Error hashing password")
            http.Error(w, "Error hashing password", http.StatusInternalServerError)
            return
        }
        user.Password = string(hashedPassword)

        if result := db.Create(&user); result.Error != nil {
            log.WithFields(log.Fields{
                "user": user,
                "error": result.Error,
            }).Error("Failed to create user")
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
            return
        }
=======
		// Send to service
		createdUser, err :=
>>>>>>> Stashed changes

        log.WithFields(log.Fields{
            "user": user,
        }).Info("User created successfully")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}
*/

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// mock
	}
}

// LoginUser handles user login.
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		log.Println("Attempting login with email:", email)

		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			log.Println("Email not found:", email)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		log.Println("User found:", user)

		// Compare password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			log.Println("Password mismatch for email:", email)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		log.Println("Login successful for email:", email)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Status:  "success",
			Message: "Login successful",
		})
	}
}

// GetAllUsers retrieves all users, with optional filters, sorting, and pagination.
func GetAllUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var users []models.User
		query := db.Model(&models.User{})

		// Optional filters
		if name := r.URL.Query().Get("name"); name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}
		if email := r.URL.Query().Get("email"); email != "" {
			query = query.Where("email LIKE ?", "%"+email+"%")
		}
		if age := r.URL.Query().Get("age"); age != "" {
			query = query.Where("age = ?", age)
		}

		// Sorting
		sortBy, sortOrder := r.URL.Query().Get("sort_by"), r.URL.Query().Get("sort_order")
		if sortBy != "" {
			order := "ASC"
			if sortOrder == "desc" {
				order = "DESC"
			}
			query = query.Order(sortBy + " " + order)
		}

		// Pagination
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
				query = query.Limit(limit)
			}
		}
		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
				query = query.Offset(offset)
			}
		}

		if err := query.Find(&users).Error; err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var user models.User
		if result := db.First(&user, id); result.Error != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// UpdateUser updates an existing user's details.
func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var user models.User
		if result := db.First(&user, id); result.Error != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		var updatedData models.User
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := utils.ValidateUserInput(&updatedData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user.Name = updatedData.Name
		user.Email = updatedData.Email
		if result := db.Save(&user); result.Error != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// DeleteUser deletes a user by their ID.
func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if result := db.Delete(&models.User{}, id); result.Error != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "home.html"))
	}
}

func JSONMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			var data map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(Response{
					Status:  "fail",
					Message: "Invalid message. Could not parse JSON.",
				})
				return
			}

			message, ok := data["message"].(string)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(Response{
					Status:  "fail",
					Message: "JSON does not contain 'message' field or not a string.",
				})
				return
			}

			fmt.Println("Received message:", message)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Data successfully received",
			})

		case http.MethodGet:
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Hello! This is a GET response from the server.",
			})

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(Response{
				Status:  "fail",
				Message: "Only GET and POST requests are allowed.",
			})
		}
	}
}
