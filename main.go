package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func migrateDB() {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully!")
}

var db *gorm.DB

func initDB() {

	dsn := "host=localhost user=postgres password=0000 dbname=postgres port=5433 sslmode=disable TimeZone=Asia/Almaty"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	fmt.Println("Database connected successfully!")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User

	result := db.First(&user, params["id"])
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User

	result := db.First(&user, params["id"])
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result := db.Delete(&User{}, params["id"])

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Find(&users)

	json.NewEncoder(w).Encode(users)
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Status:  "fail",
				Message: "Invalid JSON message",
			})
			return
		}

		message, ok := data["message"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Status:  "fail",
				Message: "Invalid JSON message",
			})
			return
		}

		fmt.Println("Received message:", message)
		json.NewEncoder(w).Encode(Response{
			Status:  "success",
			Message: "Data successfully received",
		})
		return
	}

	if r.Method == http.MethodGet {
		{
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Hello! This is a GET response from the server.",
			})
		}
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(Response{
		Status:  "fail",
		Message: "Only GET and POST requests are allowed",
	})
}

func main() {
	initDB()
	migrateDB()
	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))

	http.HandleFunc("/", handleRequests)
	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
