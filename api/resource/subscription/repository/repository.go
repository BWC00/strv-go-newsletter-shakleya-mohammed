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
	if err := r.postgresDB.Where("id = ?", SubscriptionInfo.NewsletterID).First(&newsletter.Newsletter{}).Error; err != nil {
		return "", err
	}

	subscriptionsRef := r.firebaseDB.Child("subscriptions")

	//check if already subscribed (this is a blocking operation)
	// if err := subscriptionsRef.Get(context.Background(), &subscription.Subscription{}); err == nil {
	//         log.Fatalln("Error reading value:", err)
	// }

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	newSubscriptionRef, err := subscriptionsRef.Push(context.Background(), SubscriptionInfo)
	if err != nil {
		return "", err
	}

	// Get the unique key generated by Push()
	subscriptionID := newSubscriptionRef.Key

	return subscriptionID, nil
}

func (r *Repository) Unsubscribe(subscriptionID string) error {

	subscriptionsRef := r.firebaseDB.Child("subscriptions").Child(subscriptionID)

	// Remove subscription
	err := subscriptionsRef.Set(context.Background(), nil)

	return err
}