package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"507f1f77bcf86cd799439011"`
	Name  string             `json:"name" bson:"name" example:"Josuel"`
	Email string             `json:"email" bson:"email" example:"josuel@example.com"`
}
