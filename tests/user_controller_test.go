package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"

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

// Testa o endpoint GET /users
func TestGetUsers(t *testing.T) {
	// Mock: inserindo usuário de teste
	database.DB.Create(&models.User{Name: "Josuel Test", Email: "teste@exemplo.com"})

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
