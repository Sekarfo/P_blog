package main

import (
	"auth-service/database"
	"auth-service/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.ConnectDB()

	app := fiber.New()

	// Authentication routes
	app.Post("/register", routes.Register)
	app.Post("/login", routes.Login)
	app.Post("/forgot-password", routes.ForgotPassword)
	app.Post("/reset-password", routes.ResetPassword)

	app.Listen(":8081")
}
