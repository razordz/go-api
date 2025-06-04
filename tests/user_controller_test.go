package tests

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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

// Helper para limpar e inserir user
func prepareTestUser(t *testing.T) {
	ctx := context.TODO()

	// 🔥 Dropa a collection (isso apaga índices!)
	err := database.MongoDB.Collection("users").Drop(ctx)
	if err != nil {
		t.Fatalf("Erro ao limpar coleção: %v", err)
	}

	// ✅ Recria índice único no campo email
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = database.MongoDB.Collection("users").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		t.Fatalf("Erro ao recriar índice de email: %v", err)
	}

	// ✅ Insere usuário de teste
	hash, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	user := models.User{Name: "Josuel Test", Email: "teste@exemplo.com", Password: string(hash), Role: "user"}
	_, err = database.MongoDB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		t.Fatalf("Erro ao inserir user: %v", err)
	}
}

// Testa o endpoint GET /users
func TestGetUsers(t *testing.T) {
	prepareTestUser(t)

	// Cria uma request GET falsa
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Cria um gravador para capturar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetUsers)

	// Executa a requisição
	handler.ServeHTTP(rr, req)

	// Valida status HTTP
	if rr.Code != http.StatusOK {
		t.Errorf("Status esperado %d, recebido %d", http.StatusOK, rr.Code)
	}

	// Valida conteúdo da resposta (pelo menos 1 usuário)
	var users []models.User
	json.Unmarshal(rr.Body.Bytes(), &users)

	if len(users) == 0 {
		t.Errorf("Esperado ao menos 1 usuário, recebido %d", len(users))
	}
}
