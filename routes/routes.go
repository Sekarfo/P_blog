// routes/routes.go
package routes

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
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

	// Home route serving index.html
	router.HandleFunc("/", controllers.HomeHandler()).Methods("GET")

	// JSON parse route
	router.HandleFunc("/api/message", controllers.JSONMessageHandler()).Methods("POST", "GET")

	// Serve static files from /static/
	staticDir := http.Dir("./static/")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticDir)))

	return router
}
