package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDTO struct {
	Id        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
	Secret    string
}

type FeedDTO struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string
	Url       string
	UserId    primitive.ObjectID `bson:"user_id"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}

type FeedFollowDTO struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	FeedId    primitive.ObjectID `bson:"feed_id"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}
