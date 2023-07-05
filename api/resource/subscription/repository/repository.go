package repository

import (
	"context"

	"firebase.google.com/go/v4/db"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/subscription"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/newsletter"
)

type Repository struct {
	postgresDB *gorm.DB
	firebaseDB *db.Ref
}

func NewRepository(postgresDB *gorm.DB, firebaseDB *db.Ref) *Repository {
	return &Repository{
		postgresDB: postgresDB,
		firebaseDB: firebaseDB,
	}
}


func (r *Repository) Subscribe(SubscriptionInfo *subscription.Subscription) (string, error) {
	//check if newsletter exists
	if err := r.postgresDB.Where("id = ?", SubscriptionInfo.ID).First(&newsletter.Newsletter{}).Error; err != nil {
		return "", err
	}

	subscriptionsRef := r.firebaseDB.Child("subscriptions")

	// Add subscription
	newSubscriptionRef, err := subscriptionsRef.Push(context.Background(), SubscriptionInfo)
	if err != nil {
		return "", err
	}

	// Get the unique key generated by Push()
	subscriptionID := newSubscriptionRef.Key

	return subscriptionID, nil
}

func (r *Repository) Unsubscribe(subscriptionID string) error {

	subscriptionRef := r.firebaseDB.Child("subscriptions").Child(subscriptionID)

	// Check if there is a subscription
	var subscription subscription.Subscription
	if err := subscriptionRef.Get(context.Background(), &subscription); err != nil {
		return err
	}

	// Remove subscription
	err := subscriptionsRef.Set(context.Background(), nil)

	return err
}