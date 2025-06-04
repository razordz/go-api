package services

import (
	"errors"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"

	"go-api/models"
	"go-api/repositories"
)

func RegisterUser(user *models.User) error {
	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		return errors.New("Nome, email e senha são obrigatórios")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	if user.Role == "" {
		user.Role = "user"
	}

	err = repositories.CreateUser(user)
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

func Authenticate(email, password string) (*models.User, error) {
	user, err := repositories.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("Credenciais inválidas")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("Credenciais inválidas")
	}
	return user, nil
}

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func UpdateUser(id primitive.ObjectID, data map[string]string) error {
	update := make(map[string]interface{})
	if name, ok := data["name"]; ok {
		update["name"] = name
	}
	if email, ok := data["email"]; ok {
		update["email"] = email
	}
	if password, ok := data["password"]; ok && password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		update["password"] = string(hashed)
	}
	if role, ok := data["role"]; ok {
		update["role"] = role
	}
	if len(update) == 0 {
		return errors.New("nenhum dado para atualizar")
	}
	return repositories.UpdateUser(id, update)
}
