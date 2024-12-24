package routes

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"personal_blog/controllers"
)

func SetupRoutes(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// Optional: Migrate route
	// router.HandleFunc("/migrate", func(w http.ResponseWriter, r *http.Request) {
	//     db.AutoMigrate(&models.User{})
	//     // Respond with JSON if you like
	// }).Methods("POST")

	router.HandleFunc("/users", controllers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/users", controllers.GetAllUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.GetUserByID(db)).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.DeleteUser(db)).Methods("DELETE")

	router.HandleFunc("/", controllers.HomeHandler()).Methods("GET")

	router.HandleFunc("/", controllers.JSONMessageHandler()).Methods("POST", "GET")

	return router
}
