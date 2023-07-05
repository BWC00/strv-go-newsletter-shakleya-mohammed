package repository

import (
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/newsletter"
)

// Repository represents a repository for interacting with the database for newsletters.
type Repository struct {
	postgresDB *gorm.DB
}

// NewRepository creates a new instance of the Repository with the provided PostgresDB connection.
func NewRepository(postgresDB *gorm.DB) *Repository {
	return &Repository{
		postgresDB: postgresDB,
	}
}

// ListNewsletters retrieves all newsletters associated with a specific user.
// It returns a slice of Newsletter pointers and an error if any occurred.
func (r *Repository) ListNewsletters(userId uint32) (newsletter.Newsletters, error) {
	newsletters := make([]*newsletter.Newsletter, 0)

	// Retrieve newsletters from the database for the given user ID
	if err := r.postgresDB.Find(&newsletters, "editor_id = ?", userId).Error; err != nil {
		return nil, err
	}

	// Check if no newsletters were found for the user
	if len(newsletters) == 0 {
		return nil, nil
	}

	return newsletters, nil
}

// CreateNewsletter creates a new newsletter in the database.
// It takes a Newsletter pointer as input and returns the created newsletter and an error if any occurred.
func (r *Repository) CreateNewsletter(newsletter *newsletter.Newsletter) (*newsletter.Newsletter, error) {
	// Create the newsletter in the database
	if err := r.postgresDB.Create(newsletter).Error; err != nil {
		return nil, err
	}

	return newsletter, nil
}

// ReadNewsletter retrieves a specific newsletter by its ID and associated with a specific user.
// It returns a Newsletter pointer and an error if any occurred.
func (r *Repository) ReadNewsletter(newsletterId uint32, userId uint32) (*newsletter.Newsletter, error) {
	newsletter := newsletter.Newsletter{}

	// Retrieve the newsletter from the database based on its ID and associated user ID
	if err := r.postgresDB.Where("id = ? AND editor_id = ?", newsletterId, userId).First(&newsletter).Error; err != nil {
		return nil, err
	}

	return &newsletter, nil
}

// UpdateNewsletter updates an existing newsletter in the database.
// It takes a Newsletter pointer with the updated information as input and returns an error if any occurred.
func (r *Repository) UpdateNewsletter(newsletterInfo *newsletter.Newsletter) error {
	// Check if the newsletter exists for the provided ID and user ID
	if _, err := r.ReadNewsletter(newsletterInfo.ID, newsletterInfo.EditorId); err != nil {
		return err
	}

	// Update the newsletter in the database with the provided information
	if err := r.postgresDB.Model(&newsletter.Newsletter{}).Where("id = ? AND editor_id = ?", newsletterInfo.ID, newsletterInfo.EditorId).Updates(&newsletter.Newsletter{Name: newsletterInfo.Name, Description: newsletterInfo.Description}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteNewsletter deletes a specific newsletter by its ID and associated with a specific user.
// It returns an error if any occurred.
func (r *Repository) DeleteNewsletter(newsletterId uint32, userId uint32) error {
	// Check if the newsletter exists for the provided ID and user ID
	if _, err := r.ReadNewsletter(newsletterId, userId); err != nil {
		return err
	}

	// Delete the newsletter from the database
	if err := r.postgresDB.Where("id = ? AND editor_id = ?", newsletterId, userId).Delete(&newsletter.Newsletter{}).Error; err != nil {
		return err
	}

	return nil
}
