package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-api/middlewares"
	"go-api/models"
	"go-api/services"
)

func TestJWTSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	token, _ := services.GenerateToken(&models.User{ID: primitive.NewObjectID(), Role: "user"})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Context().Value(middlewares.ContextRole) != "user" {
			t.Errorf("role ausente")
		}
	})

	middlewares.JWT(next).ServeHTTP(w, r)

	if !called || w.Code != 200 {
		t.Fatalf("middleware falhou: %d", w.Code)
	}
}

func TestJWTMissing(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	middlewares.JWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("esperado %d, obteve %d", http.StatusUnauthorized, w.Code)
	}
}
