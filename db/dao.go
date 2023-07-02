package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Retrieves a document by its id.
func findById(itemType, itemId string, result interface{}) error {

	dbId, _ := primitive.ObjectIDFromHex(itemId)
	err := getDB().Collection(collections[itemType]).FindOne(context.TODO(), bson.M{
		"_id": dbId,
	}).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(ERR_INV_ID)
		}
		return errors.New(ERR_INTERNAL_SRV)
	}
	return nil
}

// Retrieves a slice of documents
func findByField(itemType, field, value string, dto interface{}) error {

	filter := bson.M{}
	if field != "" && value != "" {
		filter = bson.M{
			field: value,
		}
	}
	cursor, err := getDB().Collection(collections[itemType]).Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.M{"created_at": -1}),
	)

	if err != nil {
		return err
	}

	defer cursor.Close(context.TODO())
	err = cursor.All(context.TODO(), dto)

	return err
}

// Retrieves a document using a custom field.
func findOneByField(itemType, field, value string, dto interface{}) error {

	err := getDB().Collection(collections[itemType]).FindOne(
		context.TODO(),
		bson.M{
			field: value,
		},
	).Decode(dto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(ERR_UNAUTHORIZED)
		}
		return errors.New(ERR_INTERNAL_SRV)
	}

	return nil
}

// Inserts a document.
func addItem(item interface{}, itemType string) error {
	_, err := getDB().Collection(collections[itemType]).InsertOne(
		context.TODO(),
		item,
	)
	return err
}

// Retrieves an user using an API key.
func GetUserByKey(apiKey string) (UserDTO, error) {
	user := UserDTO{}
	err := findOneByField("user", "secret", apiKey, &user)

	return user, err
}

func GetFeedById(id string) (FeedDTO, error) {
	feed := FeedDTO{}
	err := findById("feed", id, &feed)
	return feed, err
}

func GetFeeds(field, value string) ([]FeedDTO, error) {
	feeds := []FeedDTO{}
	err := findByField("feed", field, value, &feeds)

	return feeds, err
}

// Inserts an user.
func AddUser(user UserDTO) error {
	return addItem(user, "user")
}

// Inserts a feed.
func AddFeed(feed interface{}) error {
	return addItem(feed, "feed")
}
