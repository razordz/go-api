package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"go-api/controllers"
	"go-api/database"
	"go-api/models"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Erro ao carregar .env para testes")
	}
	database.Init()
}

// Helper para limpar a coleção antes de cada teste
func clearUsersCollection(t *testing.T) {
	ctx := context.TODO()

	// Dropa a coleção (isso apaga os índices também)
	err := database.MongoDB.Collection("users").Drop(ctx)
	if err != nil {
		t.Fatalf("Erro ao limpar a collection: %v", err)
	}

	// Recria o índice único no campo email
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = database.MongoDB.Collection("users").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		t.Fatalf("Erro ao recriar índice de email: %v", err)
	}
}

func TestCreateUserSuccess(t *testing.T) {
	clearUsersCollection(t)

	payload := models.User{
		Name:     "Usuário Teste POST",
		Email:    fmt.Sprintf("post_test_%d@example.com", time.Now().UnixNano()),
		Password: "123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Esperado status %d, recebeu %d", http.StatusCreated, rr.Code)
	}

	var response models.User
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.Email != payload.Email {
		t.Errorf("Esperado email %s, recebeu %s", payload.Email, response.Email)
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	clearUsersCollection(t)

	email := "duplicado@example.com"
	ctx := context.TODO()
	hash, _ := bcrypt.GenerateFromPassword([]byte("abc"), bcrypt.DefaultCost)
	_, err := database.MongoDB.Collection("users").InsertOne(ctx, models.User{Name: "Usuário", Email: email, Password: string(hash)})
	if err != nil {
		t.Fatalf("Erro ao inserir user duplicado: %v", err)
	}

	payload := models.User{Name: "Novo Usuário", Email: email, Password: "abc"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperado status %d para e-mail duplicado, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreateUserMissingFields(t *testing.T) {
	clearUsersCollection(t)

	payload := models.User{Name: "", Email: "", Password: ""}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperado status %d para campos inválidos, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}
