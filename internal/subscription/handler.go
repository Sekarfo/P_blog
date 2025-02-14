package subscription

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
)

type SubscriptionHandler struct {
	Service *SubscriptionService
}

func NewSubscriptionHandler(service *SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{Service: service}
}

func (h *SubscriptionHandler) RequestSubscription(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Query("user_id")) // Get user_id from query param
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	subscription, err := h.Service.RequestSubscription(uint(userID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to request subscription"})
	}

	// ✅ Ensure subscription_id is included in response
	return c.Status(201).JSON(fiber.Map{
		"subscription_id": subscription.ID,
		"message":         "Subscription request successful!",
	})
}

// func (h *SubscriptionHandler) GetPendingSubscriptions(c *fiber.Ctx) error {
// 	subscriptions, err := h.Service.GetPendingSubscriptions()
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pending subscriptions"})
// 	}

// 	response := []map[string]interface{}{}
// 	for _, sub := range subscriptions {
// 		var user models.User
// 		if err := db.DB.First(&user, sub.UserID).Error; err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user details"})
// 		}

// 		response = append(response, map[string]interface{}{
// 			"id":           sub.ID,
// 			"user_id":      sub.UserID,
// 			"user_name":    user.Name, // ✅ Include user name
// 			"status":       sub.Status,
// 			"requested_at": sub.RequestedAt,
// 		})
// 	}

// 	return c.JSON(response)
// }

func (h *SubscriptionHandler) ApproveSubscription(c *fiber.Ctx) error {
	idParam := c.Params("id")
	log.Printf("Received approval request for ID: %s", idParam) // ✅ Debugging log

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid subscription ID received: %s", idParam)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid subscription ID"})
	}

	if err := h.Service.ApproveSubscription(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to approve subscription"})
	}

	return c.JSON(fiber.Map{"message": "Subscription approved successfully"})
}

func (h *SubscriptionHandler) RejectSubscription(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid subscription ID"})
	}

	if err := h.Service.RejectSubscription(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject subscription"})
	}

	return c.JSON(fiber.Map{"message": "Subscription rejected successfully"})
}

func (h *SubscriptionHandler) GetAllSubscriptions(c *fiber.Ctx) error {
	subscriptions, err := h.Service.GetAllSubscriptions() // ✅ Now correctly linked
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch subscriptions"})
	}

	response := []map[string]interface{}{}
	for _, sub := range subscriptions {
		var user models.User
		if err := db.DB.First(&user, sub.UserID).Error; err != nil {
			log.Printf("Failed to fetch user details for subscription %d: %v", sub.ID, err)
			continue // Skip this entry instead of failing the entire request
		}

		response = append(response, map[string]interface{}{
			"id":                 sub.ID,
			"user_id":            sub.UserID,
			"user_name":          user.Name,
			"status":             sub.Status,
			"requested_at":       sub.RequestedAt,
			"approved_at":        sub.ApprovedAt,
			"subscription_start": sub.SubscriptionStart,
			"subscription_end":   sub.SubscriptionEnd,
		})
	}

	return c.JSON(response)
}

// PAYMENT

func (h *SubscriptionHandler) ProcessPayment(c *fiber.Ctx) error {
	var payment struct {
		CardNumber     string `json:"cardNumber"`
		ExpirationDate string `json:"expirationDate"`
		CVV            string `json:"cvv"`
		Name           string `json:"name"`
		Address        string `json:"address"`
		SubscriptionID uint   `json:"subscription_id"`
	}

	log.Printf("Received payment request: %s", string(c.Body())) // Debugging log

	if err := c.BodyParser(&payment); err != nil {
		log.Printf("Failed to parse payment data: %v", err) // Debugging log
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment data"})
	}

	// Simulate payment validation (if card starts with '4', payment is valid)
	paymentValid := strings.HasPrefix(payment.CardNumber, "4")

	status := "rejected"
	if paymentValid {
		status = "pending_approval" // ✅ Mark as "pending_approval", not "paid"
	}

	// Update subscription status in the database
	if err := h.Service.UpdateSubscriptionStatus(payment.SubscriptionID, status); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update subscription status"})
	}

	return c.JSON(fiber.Map{"payment_success": paymentValid})
}
