package routes

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sekarfo/P_blog/controllers/articles"
	"github.com/Sekarfo/P_blog/controllers/users"
	"github.com/Sekarfo/P_blog/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

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

// NEW ROUTES
func SetupRouter2(
	usersC *users.Controller,
	articlesC *articles.Controller,
) *mux.Router {
	router := mux.NewRouter()

	// Static files handler with proper MIME types
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Add content type middleware for JavaScript files
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			}
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/api/users", usersC.CreateUser).Methods("POST")

	router.HandleFunc("/api/articles", articlesC.FetchArticles).Methods("GET")

	staticDir := http.Dir("./static/")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Home route serving home.html
	router.HandleFunc("/", HomeHandler()).Methods("GET")

	// Profile route serving profile.html
	router.HandleFunc("/profile", ProfileHandler()).Methods("GET")
	// Login route
	router.HandleFunc("/login", usersC.LoginUser).Methods("POST")

	// Register route
	router.HandleFunc("/register", usersC.CreateUser).Methods("POST")

	// Profile API route
	router.HandleFunc("/api/profile", usersC.GetProfile).Methods("GET")

	router.HandleFunc("/api/profile", usersC.UpdateProfile).Methods("PUT")

	router.HandleFunc("/api/support", usersC.SendSupportRequest).Methods("POST")

	return router
}

func AcceptMiddlewares(h http.Handler) http.Handler {
	h = middleware.LoggerMiddlware(h)
	// Middlewares
	limiter := middleware.NewRateLimiter(rate.Every(1*time.Second), 5)
	h = limiter.Limit(h)
	return h
}
