package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Status:  "fail",
				Message: "Invalid message empty",
			})
			return
		}

		message, ok := data["message"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Status:  "fail",
				Message: "Invalid JSON message",
			})
			return
		}

		fmt.Println("Received message:", message)
		json.NewEncoder(w).Encode(Response{
			Status:  "success",
			Message: "Data successfully received",
		})
		return
	}

	if r.Method == http.MethodGet {
		{
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Hello! This is a GET response from the server.",
			})
		}
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(Response{
		Status:  "fail",
		Message: "Only GET and POST requests are allowed",
	})
}

func main() {
	http.HandleFunc("/", handleRequests)
	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
