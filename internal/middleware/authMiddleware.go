package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// TokenClaims struct to store user ID and role from JWT response
type TokenClaims struct {
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}

// Middleware to verify JWT via `auth-service`
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Verify token with `auth-service`
		authServiceURL := "http://localhost:8081/verify-token"
		req, _ := http.NewRequest("GET", authServiceURL, nil)
		req.Header.Set("Authorization", token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		// Parse JSON response
		var claims TokenClaims
		if err := json.Unmarshal(body, &claims); err != nil {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Set user ID in request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
