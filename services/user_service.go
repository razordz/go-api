package services

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"go-api/models"
	"go-api/repositories"
)

func RegisterUser(user *models.User) error {
	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" {
		return errors.New("Nome e email são obrigatórios")
	}

	err := repositories.CreateUser(user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("E-mail já cadastrado")
		}
		return err
	}

	return nil
}

func GetUsers() ([]models.User, error) {
	return repositories.FindAllUsers()
}
