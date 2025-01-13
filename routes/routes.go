package routes

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/Sekarfo/P_blog/controllers"
	"github.com/Sekarfo/P_blog/controllers/articles"
	"github.com/Sekarfo/P_blog/controllers/users"
	"github.com/Sekarfo/P_blog/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
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

	// Articles route
	router.HandleFunc("/api/articles", controllers.FetchArticles()).Methods("GET")

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

// By svdness.LXR371
func SetupRouter2(
	usersC *users.Controller,
	articlesC *articles.Controller,
) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", usersC.CreateUser).Methods("POST")
	router.HandleFunc("/api/users", usersC.GetByParams).Methods("GET")
	router.HandleFunc("/api/users/{id}", usersC.GetUserByID).Methods("GET")
	router.HandleFunc("/api/users/{id}", usersC.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", usersC.DeleteUser).Methods("DELETE")

	router.HandleFunc("/api/articles", articlesC.FetchArticles).Methods("GET")

	staticDir := http.Dir("./static/")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Home route serving home.html
	router.HandleFunc("/", controllers.HomeHandler()).Methods("GET")

	// Profile route serving profile.html
	router.HandleFunc("/profile", ProfileHandler()).Methods("GET")

	// Users route serving users.html
	router.HandleFunc("/users", UsersHandler()).Methods("GET")

	// Login route
	router.HandleFunc("/login", usersC.LoginUser).Methods("POST")

	// Register route
	router.HandleFunc("/register", usersC.CreateUser).Methods("POST")

	return router
}

func AcceptMiddlewares(h http.Handler) http.Handler {
	h = middleware.LoggerMiddlware(h)
	// Middlewares
	limiter := middleware.NewRateLimiter(rate.Every(1*time.Second), 5)
	h = limiter.Limit(h)
	return h
}
