package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "regexp"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"personal_blog/models"
	"personal_blog/utils"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := utils.ValidateUserInput(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if result := db.Create(&user); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func GetAllUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var users []models.User
		db.Find(&users)
		json.NewEncoder(w).Encode(users)
	}
}

func GetUserByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
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

func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
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
		db.Save(&user)
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
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

// ===============================
// Misc Handlers
// ===============================
func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Status:  "success",
			Message: "Welcome to the API!",
		})
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
