package routes

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"personal_blog/controllers"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// CRUD routes
	router.HandleFunc("/api/users", controllers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/api/users", controllers.GetAllUsers(db)).Methods("GET")
	router.HandleFunc("/api/users/{id}", controllers.GetUserByID(db)).Methods("GET")
	router.HandleFunc("/api/users/{id}", controllers.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/api/users/{id}", controllers.DeleteUser(db)).Methods("DELETE")

	// Home route serving home.html
	router.HandleFunc("/", controllers.HomeHandler()).Methods("GET")

	// Profile route serving profile.html
	router.HandleFunc("/profile", ProfileHandler()).Methods("GET")

	// Users route serving users.html
	router.HandleFunc("/users", UsersHandler()).Methods("GET")

	// Login route
	router.HandleFunc("/login", controllers.LoginUser(db)).Methods("POST")

	// Register route
	router.HandleFunc("/register", controllers.CreateUser(db)).Methods("POST")

	// JSON parse route
	router.HandleFunc("/api/message", controllers.JSONMessageHandler()).Methods("POST", "GET")

	// Serve static files from /static/
	staticDir := http.Dir("./static/")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticDir)))

	return router
}

func HomeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL.Path,
		}).Info("Serving home.html")
		http.ServeFile(w, r, filepath.Join("static", "home.html"))
    }
}
func ProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving profile.html")
		http.ServeFile(w, r, filepath.Join("static", "profile.html"))
	}
}

func UsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving users.html")
		http.ServeFile(w, r, filepath.Join("static", "users.html"))
	}
}
