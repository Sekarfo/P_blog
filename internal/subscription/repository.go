package subscription

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (repo *SubscriptionRepository) CreateSubscription(userID uint) (*Subscription, error) {
	subscription := Subscription{
		UserID:      userID,
		Status:      "pending",
		RequestedAt: time.Now(),
	}
	if err := repo.DB.Create(&subscription).Error; err != nil {
		log.Printf("Error creating subscription for user %d: %v\n", userID, err)
		return nil, err
	}
	log.Printf("Subscription created for user %d\n", userID)
	return &subscription, nil
}

func (repo *SubscriptionRepository) GetPendingSubscriptions() ([]Subscription, error) {
	var subscriptions []Subscription
	if err := repo.DB.Where("status = ?", "pending").Find(&subscriptions).Error; err != nil {
		log.Println("Error fetching pending subscriptions:", err)
		return nil, err
	}
	return subscriptions, nil
}

func (repo *SubscriptionRepository) ApproveSubscription(id uint) error {
	var subscription Subscription
	if err := repo.DB.First(&subscription, id).Error; err != nil {
		log.Printf("Subscription ID %d not found\n", id)
		return errors.New("subscription not found")
	}

	now := time.Now()
	subscription.Status = "approved"
	subscription.ApprovedAt = &now
	subscription.SubscriptionStart = &now
	expiry := now.AddDate(0, 1, 0) // 1-month premium duration
	subscription.SubscriptionEnd = &expiry

	if err := repo.DB.Save(&subscription).Error; err != nil {
		log.Printf("Error saving approved subscription ID %d: %v\n", id, err)
		return err
	}
	log.Printf("Subscription ID %d approved, expiry on %s\n", id, expiry.Format("2006-01-02"))
	return nil
}

func (repo *SubscriptionRepository) RejectSubscription(id uint) error {
	if err := repo.DB.Model(&Subscription{}).Where("id = ?", id).Update("status", "rejected").Error; err != nil {
		log.Printf("Error rejecting subscription ID %d: %v\n", id, err)
		return err
	}
	log.Printf("Subscription ID %d rejected\n", id)
	return nil
}

func (repo *SubscriptionRepository) GetSubscriptionByID(id uint, subscription *Subscription) error {
	return repo.DB.First(subscription, id).Error
}

func (repo *SubscriptionRepository) SaveSubscription(subscription *Subscription) error {
	return repo.DB.Save(subscription).Error
}
