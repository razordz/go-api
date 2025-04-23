package services

import (
	"errors"
	"strings"

	"go-api/models"
	"go-api/repositories"
)

func GetUsers() ([]models.User, error) {
	return repositories.FindAllUsers()
}

func RegisterUser(user *models.User) error {
	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" {
		return errors.New("Nome e email são obrigatórios")
	}

	return repositories.CreateUser(user)
}
