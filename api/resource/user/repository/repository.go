package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/user"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/auth"
)

type Repository struct {
	postgresDB *gorm.DB
}

func NewRepository(postgresDB *gorm.DB) *Repository {
	return &Repository{
		postgresDB: postgresDB,
	}
}


func (r *Repository) RegisterUser(userInfo *user.User) (string, error) {

	// email must be unique
	if err := r.postgresDB.Where("email = ?", userInfo.Email).First(&user.User{}).Error; err == nil {
		return "", errors.New("email already exists")
	}

	hashedPassword, err := auth.Hash(userInfo.Password)
	if err != nil {
		return "", err
	}
	userInfo.Password = string(hashedPassword)

	if err := r.postgresDB.Create(userInfo).Error; err != nil {
		return "", err
	}

	return auth.CreateToken(userInfo.ID)
}

func (r *Repository) LoginUser(userInfo *user.User) (string, error) {

	authUser := user.User{}

	if err := r.postgresDB.Where("email = ?", userInfo.Email).First(&authUser).Error; err != nil {
		return "", err
	}

	if err := auth.VerifyPassword(authUser.Password, userInfo.Password); err != nil {
		return "", err
	}

	return auth.CreateToken(authUser.ID)
}