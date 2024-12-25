// routes/routes.go
package routes

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"personal_blog/controllers"
)

func SetupRoutes(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// CRUD routes
	router.HandleFunc("/api/users", controllers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/api/users", controllers.GetAllUsers(db)).Methods("GET")
	router.HandleFunc("/api/users/{id}", controllers.GetUserByID(db)).Methods("GET")
	router.HandleFunc("/api/users/{id}", controllers.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/api/users/{id}", controllers.DeleteUser(db)).Methods("DELETE")

	// Home route
	router.HandleFunc("/", controllers.HomeHandler()).Methods("GET")

	// JSON parse route
	router.HandleFunc("/api/message", controllers.JSONMessageHandler()).Methods("POST", "GET")

	// Serve static files
	staticDir := http.Dir("./static/")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Serve index.html for root path
	router.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	}).Methods("GET")

	return router
}
