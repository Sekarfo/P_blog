package main

import (
	"fmt"
	"log"
	"net/http"

	"personal_blog/config"
	"personal_blog/routes"
)

func main() {
	db := config.InitDB()

	config.MigrateDB(db)

	router := routes.SetupRoutes(db)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
