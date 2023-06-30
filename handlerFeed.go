package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/valentin-popov/rss-aggregator/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var feedController dbController = dbController{
	collection: getDB().Collection(
		collections["feed"],
	),
}

func (feedController *dbController) createFeed(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		sendError(w, http.StatusUnauthorized, ERR_UNAUTHORIZED)
		return
	}
	user, err := userController.getUserByKey(apiKey)

	if err != nil {
		e := errorObject{
			err.Error(),
		}
		e.handleError(w, r)
		return
	}

	feedData := struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&feedData)

	if err != nil {
		e := errorObject{
			err.Error(),
		}
		e.handleError(w, r)
		return
	}

	feed := Feed{
		Id:        primitive.NewObjectID(),
		Name:      feedData.Name,
		Url:       feedData.Url,
		UserId:    user.Id,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = feedController.collection.InsertOne(
		context.TODO(),
		feed,
	)

	if err != nil {
		e := errorObject{
			ERR_DB_INSERT,
		}
		e.handleError(w, r)
		return
	}

	sendJSON(w, http.StatusOK, feed)

}
