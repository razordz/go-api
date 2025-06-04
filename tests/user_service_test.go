package tests

import (
	"errors"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"go-api/models"
	"go-api/repositories"
	"go-api/services"
)

func TestRegisterUserValidation(t *testing.T) {
	err := services.RegisterUser(&models.User{Name: "", Email: "", Password: ""})
	if err == nil || err.Error() != "Nome, email e senha são obrigatórios" {
		t.Errorf("esperava erro de validação, obteve %v", err)
	}
}

func TestRegisterUserDuplicate(t *testing.T) {
	original := repositories.CreateUser
	defer func() { repositories.CreateUser = original }()

	repositories.CreateUser = func(u *models.User) error {
		return mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000}}}
	}

	err := services.RegisterUser(&models.User{Name: "A", Email: "a@a.com", Password: "123"})
	if err == nil || err.Error() != "E-mail já cadastrado" {
		t.Errorf("esperava erro de duplicidade, obteve %v", err)
	}
}

func TestRegisterUserSuccess(t *testing.T) {
	original := repositories.CreateUser
	defer func() { repositories.CreateUser = original }()

	var created *models.User
	repositories.CreateUser = func(u *models.User) error { created = u; return nil }

	user := models.User{Name: "A", Email: "a@a.com", Password: "123"}
	if err := services.RegisterUser(&user); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if created == nil || created.Password == "123" {
		t.Errorf("senha não foi criptografada")
	}
	if user.Role != "user" {
		t.Errorf("role padrão não aplicado")
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	pwd, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	original := repositories.FindUserByEmail
	defer func() { repositories.FindUserByEmail = original }()
	repositories.FindUserByEmail = func(email string) (*models.User, error) {
		return &models.User{Email: email, Password: string(pwd)}, nil
	}
	user, err := services.Authenticate("a@a.com", "123")
	if err != nil || user.Email != "a@a.com" {
		t.Fatalf("falha na autenticação: %v", err)
	}
}

func TestAuthenticateInvalid(t *testing.T) {
	original := repositories.FindUserByEmail
	defer func() { repositories.FindUserByEmail = original }()
	repositories.FindUserByEmail = func(email string) (*models.User, error) {
		return nil, errors.New("not found")
	}
	if _, err := services.Authenticate("a@a.com", "123"); err == nil {
		t.Fatal("esperava erro para usuário inexistente")
	}
}

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	id := primitive.NewObjectID()
	token, err := services.GenerateToken(&models.User{ID: id, Role: "admin"})
	if err != nil || token == "" {
		t.Fatalf("token inválido: %v", err)
	}
}

func TestUpdateUserNoData(t *testing.T) {
	if err := services.UpdateUser(primitive.NewObjectID(), map[string]string{}); err == nil {
		t.Fatal("esperava erro para payload vazio")
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	original := repositories.UpdateUser
	defer func() { repositories.UpdateUser = original }()

	var gotID primitive.ObjectID
	var gotData interface{}
	repositories.UpdateUser = func(id primitive.ObjectID, data bson.M) error {
		gotID = id
		gotData = data
		return nil
	}

	err := services.UpdateUser(primitive.NewObjectID(), map[string]string{"name": "X", "password": "123"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if gotID.IsZero() || gotData == nil {
		t.Fatal("update não chamado corretamente")
	}
}
