package main

import (
	"flag"
	"log"

	"go-api/config"
	"go-api/database"
	"go-api/models"
	"go-api/services"
)

func main() {
	name := flag.String("name", "", "Nome do usuário")
	email := flag.String("email", "", "Email do usuário")
	password := flag.String("password", "", "Senha")
	admin := flag.Bool("admin", false, "Criar como administrador")
	flag.Parse()

	if *name == "" || *email == "" || *password == "" {
		log.Fatal("Informe name, email e password")
	}

	config.LoadEnv()
	database.Init()

	role := "user"
	if *admin {
		role = "admin"
	}

	user := models.User{Name: *name, Email: *email, Password: *password, Role: role}
	if err := services.RegisterUser(&user); err != nil {
		log.Fatal(err)
	}
	log.Println("Usuário criado com sucesso")
}
