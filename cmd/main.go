package main

import (
	"net/http"

	"github.com/hitpads/reado_ap/internal/admin"
	"github.com/hitpads/reado_ap/internal/comments"
	"github.com/hitpads/reado_ap/internal/posts"

	"github.com/go-chi/chi/v5"
	"github.com/hitpads/reado_ap/internal/auth"
	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/users"
)

func main() {
	db.ConnectDB()

	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	r.Handle("/*", http.FileServer(http.Dir("./web/templates")))

	// Public Routes
	r.Post("/register", users.RegisterHandler)
	r.Get("/verify", users.VerifyEmailHandler)
	r.Post("/login", users.LoginHandler)

	// Admin routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)                 // Ensure user is authenticated
		r.Use(auth.RoleMiddleware(auth.AdminRole)) // Ensure user is admin

		r.Get("/admin/users", admin.GetAllUsers)
		r.Patch("/admin/users/role", admin.UpdateUserRole)
	})
	// Posts routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Post("/posts", posts.CreatePost)
		r.Get("/posts", posts.GetPosts)
		r.Delete("/posts/{id}", posts.DeletePost)

		r.Get("/posts/{id}", posts.GetPostDetails)
		r.Post("/posts/{id}/like", posts.LikePost)
		r.Post("/posts/{id}/unlike", posts.UnlikePost)
		r.Get("/posts/{id}/likes", posts.GetLikeStatus)
	})

	// Comments routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Post("/comments", comments.AddComment)
		r.Get("/posts/{id}/comments", comments.GetComments)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Profile page - User authenticated!"))
		})
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
