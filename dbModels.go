package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	FirstName string             `bson:"first_name" json:"firstName"`
	LastName  string             `bson:"last_name" json:"lastName"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updatedAt"`
	Secret    string             `bson:"secret" json:"secret"`
}

type Feed struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name"`
	Url       string             `json:"url"`
	UserId    primitive.ObjectID `bson:"user_id" json:"userId"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updatedAt"`
}
