package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB indisponível:", err)
	}

	MongoClient = client
	MongoDB = client.Database(dbName)

	// ✅ Criar índice único no campo email (somente uma vez)
	userCollection := MongoDB.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = userCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("❌ Erro ao criar índice único no campo email:", err)
	}

	log.Println("Conectado ao MongoDB com sucesso!")
}
