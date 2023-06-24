package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	FirstName string             `bson:"first_name" json:"firstName"`
	LastName  string             `bson:"last_name" json:"lastName"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updatedAt"`
	Secret    string             `bson:"secret" json:"secret"`
}

type ObjectId struct {
	Id string `json:"id"`
}
