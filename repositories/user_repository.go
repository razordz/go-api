package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go-api/database"
	"go-api/models"
)

func FindAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := database.MongoDB.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.MongoDB.Collection("users").InsertOne(ctx, user)
	return err
}
