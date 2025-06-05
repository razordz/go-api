package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-api/controllers"
	"go-api/middlewares"
	"go-api/models"
	"go-api/services"
)

func TestLoginSuccess(t *testing.T) {
	authOrig := services.Authenticate
	tokenOrig := services.GenerateToken
	defer func() { services.Authenticate = authOrig; services.GenerateToken = tokenOrig }()

	services.Authenticate = func(email, password string) (*models.User, error) {
		return &models.User{ID: primitive.NewObjectID(), Role: "user"}, nil
	}
	services.GenerateToken = func(u *models.User) (string, error) { return "tok", nil }

	creds := map[string]string{"email": "a@a.com", "password": "123"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controllers.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status esperado %d, obteve %d", http.StatusOK, rr.Code)
	}
	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["token"] != "tok" {
		t.Fatalf("token incorreto: %v", resp)
	}
}

func TestLoginInvalid(t *testing.T) {
	authOrig := services.Authenticate
	defer func() { services.Authenticate = authOrig }()
	services.Authenticate = func(email, password string) (*models.User, error) {
		return nil, errors.New("invalid")
	}
	body := []byte(`{"email":"a","password":"b"}`)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	controllers.Login(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("esperado %d, obteve %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestUpdateUserControllerSuccess(t *testing.T) {
	upOrig := services.UpdateUser
	defer func() { services.UpdateUser = upOrig }()

	var called bool
	services.UpdateUser = func(id primitive.ObjectID, data map[string]string) error {
		called = true
		if data["name"] != "Novo" {
			t.Errorf("payload incorreto: %v", data)
		}
		return nil
	}

	id := primitive.NewObjectID()
	payload := map[string]string{"name": "Novo"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/users/"+id.Hex(), bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": id.Hex()})
	ctx := context.WithValue(req.Context(), middlewares.ContextRole, "admin")
	ctx = context.WithValue(ctx, middlewares.ContextUserID, id.Hex())
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	controllers.UpdateUser(rr, req)

	if rr.Code != http.StatusNoContent || !called {
		t.Fatalf("update falhou: status %d called %v", rr.Code, called)
	}
}

func TestUpdateUserForbidden(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	req := httptest.NewRequest("PUT", "/users/"+id, bytes.NewBuffer([]byte(`{}`)))
	req = mux.SetURLVars(req, map[string]string{"id": id})
	ctx := context.WithValue(req.Context(), middlewares.ContextRole, "user")
	ctx = context.WithValue(ctx, middlewares.ContextUserID, "other")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	controllers.UpdateUser(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("esperado %d, obteve %d", http.StatusForbidden, rr.Code)
	}
}

func TestUpdateUserBadID(t *testing.T) {
	req := httptest.NewRequest("PUT", "/users/abc", bytes.NewBuffer([]byte(`{}`)))
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	ctx := context.WithValue(req.Context(), middlewares.ContextRole, "admin")
	ctx = context.WithValue(ctx, middlewares.ContextUserID, "abc")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	controllers.UpdateUser(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("esperado %d, obteve %d", http.StatusBadRequest, rr.Code)
	}
}
