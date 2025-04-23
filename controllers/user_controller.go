package controllers

import (
	"encoding/json"
	"net/http"

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

	if err := services.RegisterUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
