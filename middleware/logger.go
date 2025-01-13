package middleware

import (
	"log"
	"net/http"
)

func LoggerMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s, Method: %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}
