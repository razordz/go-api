package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"go-api/middlewares"
	"go-api/models"
	"go-api/services"
)

// GetUsers godoc
// @Summary Lista usuários
// @Description Retorna todos os usuários cadastrados
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Erro interno"
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// CreateUser godoc
// @Summary Cria usuário
// @Description Cria um novo usuário com nome e e-mail
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Usuário a ser criado"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Erro ao criar usuário"
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	err := services.RegisterUser(&user)
	if err != nil {
		if err.Error() == "Nome, email e senha são obrigatórios" || err.Error() == "E-mail já cadastrado" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	user, err := services.Authenticate(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := services.GenerateToken(user)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	role := r.Context().Value(middlewares.ContextRole).(string)
	userID := r.Context().Value(middlewares.ContextUserID).(string)
	if role != "admin" && userID != idParam {
		http.Error(w, "Proibido", http.StatusForbidden)
		return
	}
	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	if err := services.UpdateUser(objID, payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
