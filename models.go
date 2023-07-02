package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	CreatedAt primitive.DateTime `json:"createdAt"`
	Secret    string             `json:"secret"`
}

type Feed struct {
	Id        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Url       string             `json:"url"`
	UserId    primitive.ObjectID `json:"userId"`
	CreatedAt primitive.DateTime `json:"createdAt"`
}

type FeedFollow struct {
	Id        primitive.ObjectID `json:"id"`
	UserId    primitive.ObjectID `json:"userId"`
	FeedId    primitive.ObjectID `json:"feedId"`
	CreatedAt primitive.DateTime `json:"createdAt"`
}
