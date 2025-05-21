package user_repository

import (
	"github.com/VaheMuradyan/CodeSignal2/todoapp/models"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

func ValidCredentials(creds models.Credentials) error {
	return validate.Struct(creds)
}

func CreateUser(db *gorm.DB, user *models.User) error {
	result := db.Create(user)
	return result.Error
}

func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	return &user, err
}
