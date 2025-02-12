package users

import (
	"io"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	authServiceURL := "http://localhost:8081/register"

	resp, err := http.Post(authServiceURL, "application/json", r.Body)
	if err != nil {
		http.Error(w, "Auth service unavailable", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Forward response from `auth-service` back to the client
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	authServiceURL := "http://localhost:8081/login"

	resp, err := http.Post(authServiceURL, "application/json", r.Body)
	if err != nil {
		http.Error(w, "Auth service unavailable", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Forward response from `auth-service` back to the client
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
