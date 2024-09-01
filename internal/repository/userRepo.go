package repository

import (
	"github.com/google/uuid"
	"github.com/just-umyt/task/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	Database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Database: db}
}

// Create User in Database
func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.Database.Create(&user).Error
}

// Get user by id
func (repo *UserRepository) GetUserById(id uuid.UUID) (models.User, error) {
	var user models.User

	if err := repo.Database.First(&user, id).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Get user by email
func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}

	if err := repo.Database.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Update users refresh
func (repo *UserRepository) UpdateUserToken(user *models.User) error {
	return repo.Database.Save(&user).Error
}
