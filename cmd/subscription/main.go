package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/subscription"
)

func main() {
	db.ConnectDB()

	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080", // ✅ Allow frontend to call the microservice
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Authorization, Content-Type, Cache-Control", // ✅ Explicitly allow Cache-Control
		ExposeHeaders:    "Authorization, Content-Type",
		AllowCredentials: true,
	}))

	// Register subscription routes
	subscription.SetupRoutes(app)

	log.Println("Subscription Service running on :8081")
	err := app.Listen(":8081")
	if err != nil {
		return
	}
}
