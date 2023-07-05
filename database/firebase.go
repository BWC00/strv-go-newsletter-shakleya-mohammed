package database

import (
	"errors"
	"context"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
)

// NewFirebaseDB creates a new Firebase Realtime Database client.
// It initializes the Firebase admin SDK with the provided configuration and
// service account key file path, and returns a reference to the Firebase Realtime Database.
// Returns an error if any initialization or configuration step fails.
func NewFirebaseDB(cfg *config.Firebase) (*db.Ref, error) {
	// Create a new background context
	ctx := context.Background()

	// Prepare the Firebase configuration
	conf := &firebase.Config{
        DatabaseURL: cfg.Location,
    }

    // Get the absolute file path for the service account key file
	serviceAccountKeyFilePath, err := filepath.Abs(cfg.FirebaseCredPath)
	if err != nil {
		return nil, errors.New("unable to load serviceAccountKeys.json file")
	}

	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	// Initialize the Firebase admin SDK with the provided configuration and service account key
	app, err2 := firebase.NewApp(ctx, conf, opt)
	if err2 != nil {
		return nil, err2
	}

	// Create a Firebase Realtime Database client from the initialized app
	client, err3 := app.Database(ctx)
	if err3 != nil {
		return nil, err3
	}

	// Get a reference to the Firebase Realtime Database entry point specified in the configuration
	DBref := client.NewRef(cfg.RefEntryPoint)

	return DBref, nil
}