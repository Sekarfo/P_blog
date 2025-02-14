package subscription

import (
	"fmt"
	"log"

	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
)

type SubscriptionService struct {
	Repo *SubscriptionRepository
}

func NewSubscriptionService(repo *SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{Repo: repo}
}

func (s *SubscriptionService) RequestSubscription(userID uint) (*Subscription, error) {
	subscription, err := s.Repo.CreateSubscription(userID)
	if err != nil {
		log.Println("Error creating subscription:", err)
		return nil, err
	}
	log.Printf("Subscription request created for User ID %d\n", userID)
	return subscription, nil
}

func (s *SubscriptionService) GetPendingSubscriptions() ([]Subscription, error) {
	subscriptions, err := s.Repo.GetPendingSubscriptions()
	if err != nil {
		log.Println("Error fetching pending subscriptions:", err)
		return nil, err
	}
	return subscriptions, nil
}

func (s *SubscriptionService) ApproveSubscription(id uint) error {
	err := s.Repo.ApproveSubscription(id)
	if err != nil {
		log.Printf("Error approving subscription with ID %d: %v\n", id, err)
		return err
	}

	// Fetch approved subscription details
	var subscription Subscription
	if err := db.DB.First(&subscription, id).Error; err != nil {
		log.Printf("Error fetching approved subscription ID %d: %v\n", id, err)
		return fmt.Errorf("failed to fetch approved subscription: %w", err)
	}

	// Fetch user details
	var user models.User
	if err := db.DB.First(&user, "id = ?", subscription.UserID).Error; err != nil {
		log.Printf("Error fetching user details for User ID %d: %v\n", subscription.UserID, err)
		return fmt.Errorf("failed to find user: %w", err)
	}

	expiryDate := subscription.SubscriptionEnd.Format("2006-01-02")
	log.Printf("User ID %d subscription approved. Expiry: %s\n", subscription.UserID, expiryDate)

	// Send email notification
	err = SendApprovalEmail(user.Email, expiryDate)
	if err != nil {
		log.Printf("Failed to send approval email to %s: %v\n", user.Email, err)
		return err
	}

	log.Printf("Approval email sent to %s successfully.\n", user.Email)
	return nil
}

func (s *SubscriptionService) RejectSubscription(id uint) error {
	err := s.Repo.RejectSubscription(id)
	if err != nil {
		log.Printf("Error rejecting subscription with ID %d: %v\n", id, err)
		return err
	}

	// Fetch user details for email notification
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		log.Printf("Error fetching user details for ID %d: %v\n", id, err)
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Send rejection email
	err = SendRejectionEmail(user.Email)
	if err != nil {
		log.Printf("Failed to send rejection email to %s: %v\n", user.Email, err)
		return err
	}

	log.Printf("Rejection email sent to %s successfully.\n", user.Email)
	return nil
}

func (s *SubscriptionService) UpdateSubscriptionStatus(id uint, status string) error {
	var subscription Subscription
	if err := s.Repo.GetSubscriptionByID(id, &subscription); err != nil {
		return err
	}

	subscription.Status = status
	return s.Repo.SaveSubscription(&subscription)
}
