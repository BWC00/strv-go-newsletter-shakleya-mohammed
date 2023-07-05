package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash generates a bcrypt hash for the given password.
// It takes the password as input and returns the hashed password as a byte slice and any error encountered.
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword compares a hashed password with a plain-text password to check for a match.
// It takes the hashed password and the plain-text password as input and returns an error if the passwords don't match.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}