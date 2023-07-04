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

func NewFirebaseDB(cfg *config.Firebase) (*db.Ref, error) {

	ctx := context.Background()
	conf := &firebase.Config{
        DatabaseURL: cfg.Location,
    }

	serviceAccountKeyFilePath, err := filepath.Abs(cfg.FirebaseCredPath)
	if err != nil {
		return nil, errors.New("unable to load serviceAccountKeys.json file")
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err2 := firebase.NewApp(ctx, conf, opt)
	if err2 != nil {
		return nil, err2
	}

	// Create a database client from App
	client, err3 := app.Database(ctx)
	if err3 != nil {
		return nil, err3
	}

	// Get a database reference
	DBref := client.NewRef(cfg.RefEntryPoint)

	return DBref, nil
}