package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hitpads/reado_ap/internal/admin"
	"github.com/hitpads/reado_ap/internal/auth"
	"github.com/hitpads/reado_ap/internal/comments"
	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/posts"
	"github.com/hitpads/reado_ap/internal/users"
)

func main() {
	db.ConnectDB()
	r := chi.NewRouter()

	// Public Routes (forward to `auth-service`)
	r.Post("/register", users.RegisterHandler)
	r.Post("/login", users.LoginHandler)

	// Protected Routes (Use `AuthMiddleware`)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Profile page - User authenticated!"))
		})

		// Admin routes
		r.Get("/admin/users", admin.GetAllUsers)
		r.Patch("/admin/users/role", admin.UpdateUserRole)

		// Posts routes
		r.Post("/posts", posts.CreatePost)
		r.Get("/posts", posts.GetPosts)
		r.Delete("/posts/{id}", posts.DeletePost)
		r.Get("/posts/{id}", posts.GetPostDetails)
		r.Post("/posts/{id}/like", posts.LikePost)
		r.Post("/posts/{id}/unlike", posts.UnlikePost)
		r.Get("/posts/{id}/likes", posts.GetLikeStatus)

		// Comments routes
		r.Post("/comments", comments.AddComment)
		r.Get("/posts/{id}/comments", comments.GetComments)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
