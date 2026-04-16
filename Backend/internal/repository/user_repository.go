package repository

import (
	"errors"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// UserRepository abstracts all DB access for the User entity.
type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository constructs a UserRepository backed by *gorm.DB.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByEmail returns the user with the given email or nil if not found.
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID returns the user with the given primary key or nil if not found.
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
