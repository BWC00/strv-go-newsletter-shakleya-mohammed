package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/user"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/auth"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
)

// Repository provides methods to interact with the user data in the database.
type Repository struct {
	postgresDB *gorm.DB
}

// NewRepository creates a new instance of the Repository with the provided PostgreSQL database connection.
func NewRepository(postgresDB *gorm.DB) *Repository {
	return &Repository{
		postgresDB: postgresDB,
	}
}

// RegisterUser registers a new user in the database.
// It performs email uniqueness validation, hashes the password, and creates a user record.
// Returns the generated authentication token on success, or an error if the registration fails.
func (r *Repository) RegisterUser(userInfo *user.User) (string, error) {
	// Check if the email is already taken
	if err := r.postgresDB.Where("email = ?", userInfo.Email).First(&user.User{}).Error; err == nil {
		return "", errors.New(e.FieldNotUnique)
	}

	// Hash the user's password
	hashedPassword, err := auth.Hash(userInfo.Password)
	if err != nil {
		return "", err
	}
	userInfo.Password = string(hashedPassword)

	// Create the user record in the database
	if err := r.postgresDB.Create(userInfo).Error; err != nil {
		return "", err
	}

	// Generate an authentication token for the user
	return auth.CreateToken(userInfo.ID)
}

// LoginUser performs user login authentication.
// It retrieves the user record from the database based on the provided email,
// verifies the password, and returns the authentication token on success.
// Returns an error if the login fails.
func (r *Repository) LoginUser(userInfo *user.User) (string, error) {
	// Retrieve the user record from the database based on the email
	authUser := user.User{}
	if err := r.postgresDB.Where("email = ?", userInfo.Email).First(&authUser).Error; err != nil {
		return "", err
	}

	// Verify the password
	if err := auth.VerifyPassword(authUser.Password, userInfo.Password); err != nil {
		return "", err
	}

	// Generate an authentication token for the user
	return auth.CreateToken(authUser.ID)
}