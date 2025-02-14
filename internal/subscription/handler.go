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

	return c.Status(201).JSON(subscription)
}

func (h *SubscriptionHandler) GetPendingSubscriptions(c *fiber.Ctx) error {
	subscriptions, err := h.Service.GetPendingSubscriptions()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pending subscriptions"})
	}

	response := []map[string]interface{}{}
	for _, sub := range subscriptions {
		var user models.User
		if err := db.DB.First(&user, sub.UserID).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user details"})
		}

		response = append(response, map[string]interface{}{
			"id":           sub.ID,
			"user_id":      sub.UserID,
			"user_name":    user.Name, // ✅ Include user name
			"status":       sub.Status,
			"requested_at": sub.RequestedAt,
		})
	}

	return c.JSON(response)
}

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

	if err := c.BodyParser(&payment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment data"})
	}

	// Mock Payment Validation (if card starts with '4' → success, otherwise fail)
	success := strings.HasPrefix(payment.CardNumber, "4")

	status := "declined"
	if success {
		status = "paid"
	}

	// Update subscription status in the database
	if err := h.Service.UpdateSubscriptionStatus(payment.SubscriptionID, status); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update subscription status"})
	}

	return c.JSON(fiber.Map{"payment_success": success})
}
