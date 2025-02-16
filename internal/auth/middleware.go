package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
	"github.com/hitpads/reado_ap/internal/utils"
)

type key int

const (
	UserIDKey key = iota
	RoleContextKey
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return utils.GetJWTSecret(), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		var user models.User
		if err := db.DB.Preload("Role").First(&user, claims.UserID).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleContextKey, user.Role.Name)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
