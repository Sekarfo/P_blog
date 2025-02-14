package subscription

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hitpads/reado_ap/internal/db"
)

func SetupRoutes(app *fiber.App) {
	db.ConnectDB() // Ensuring DB is properly initialized

	repo := NewSubscriptionRepository(db.DB)
	service := NewSubscriptionService(repo)
	handler := NewSubscriptionHandler(service)

	api := app.Group("/api/subscriptions")
	api.Post("/subscribe", handler.RequestSubscription)
	api.Get("/admin/pending", handler.GetPendingSubscriptions)
	api.Patch("/admin/approve/:id", handler.ApproveSubscription)
	api.Patch("/admin/reject/:id", handler.RejectSubscription)

	log.Println("Subscription routes initialized successfully")
}
