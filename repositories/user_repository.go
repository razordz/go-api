package repositories

import (
	"go-api/database"
	"go-api/models"
)

func FindAllUsers() ([]models.User, error) {
	var users []models.User
	result := database.DB.Find(&users)
	return users, result.Error
}

func CreateUser(user *models.User) error {
	result := database.DB.Create(user)
	return result.Error
}
