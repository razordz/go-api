package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestCreateUserSuccess(t *testing.T) {
	payload := models.User{Name: "Usuário Teste POST", Email: fmt.Sprintf("post_test_%d@example.com", time.Now().UnixNano())}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
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
	email := "duplicado@example.com"
	database.DB.Create(&models.User{Name: "Usuário", Email: email})

	payload := models.User{Name: "Novo Usuário", Email: email}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperado status %d para e-mail duplicado, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreateUserMissingFields(t *testing.T) {
	payload := models.User{Name: "", Email: ""}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Esperado status %d para campos inválidos, recebeu %d", http.StatusBadRequest, rr.Code)
	}
}
