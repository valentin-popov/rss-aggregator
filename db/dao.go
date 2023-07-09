package db

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Retrieves a document by its id.
func findById(itemType, itemId string, result interface{}) error {

	dbId, _ := primitive.ObjectIDFromHex(itemId)
	err := getDB().Collection(collections[strings.ToLower(itemType)]).FindOne(context.TODO(), bson.M{
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
func findByField(itemType, field string, value interface{}, dto interface{}) error {

	filter := bson.M{}
	if field != "" && value != nil {
		filter = bson.M{
			field: value,
		}
	}
	cursor, err := getDB().Collection(collections[strings.ToLower(itemType)]).Find(
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
func findOneByField(itemType, field string, value interface{}, dto interface{}) error {

	err := getDB().Collection(collections[strings.ToLower(itemType)]).FindOne(
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
	_, err := getDB().Collection(collections[strings.ToLower(itemType)]).InsertOne(
		context.TODO(),
		item,
	)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New(ERR_DUPL_KEY_REC)
	}
	return err
}

func deleteItemByField(fields interface{}, itemType string) error {

	filterBytes, err := bson.Marshal(fields)
	if err != nil {
		return err
	}

	dbFilter := bson.M{}

	err = bson.Unmarshal(filterBytes, &dbFilter)
	if err != nil {
		return err
	}

	delRes, err := getDB().Collection(collections[strings.ToLower(itemType)]).DeleteOne(
		context.TODO(),
		dbFilter,
	)

	// No document matched the filter
	if delRes.DeletedCount == 0 {
		return errors.New(ERR_NO_DOC)
	}

	if err != nil {
		return err
	}

	return nil
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

func GetFollowedFeeds(userId primitive.ObjectID) ([]FeedFollowDTO, error) {
	followedFeeds := []FeedFollowDTO{}
	err := findByField("feedFollow", "user_id", userId, &followedFeeds)
	return followedFeeds, err
}

// Inserts an user.
func AddUser(user UserDTO) error {
	return addItem(user, "user")
}

// Inserts a feed.
func AddFeed(feed interface{}) error {
	return addItem(feed, "feed")
}

// Inserts a feed follow.
func AddFeedFollow(feedFollow interface{}) error {
	return addItem(feedFollow, "feedFollow")
}

// Deletes a feed follow using the feed ID and
// the ID of the following user.
func DeleteFeedFollow(fields FeedFollowDeleteFilter) error {
	return deleteItemByField(struct {
		feed_id primitive.ObjectID
		user_id primitive.ObjectID
	}{
		fields.FeedId,
		fields.UserId,
	}, "feedFollow")
}

type FeedFollowDeleteFilter struct {
	FeedId primitive.ObjectID
	UserId primitive.ObjectID
}
