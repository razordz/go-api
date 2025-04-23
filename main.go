package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-api/config"
	"go-api/database"
	"go-api/routes"
)

func main() {
	config.LoadEnv()
	database.Init()
	router := routes.Setup()

	port := os.Getenv("PORT")
	fmt.Println("Servidor rodando na porta:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
