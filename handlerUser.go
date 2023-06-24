package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/valentin-popov/rss-aggregator/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userController dbController = dbController{
	collection: getDB().Collection(
		collections["user"],
	),
}

func (userController *dbController) create(w http.ResponseWriter, r *http.Request) {

	userData := struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)

	if err != nil {
		sendError(w, http.StatusBadRequest, ERR_JSON)
		return
	}

	byte32Arr := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	byteArr := byte32Arr[:]

	user := User{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Secret:    hex.EncodeToString(byteArr),
	}

	result, err := userController.collection.InsertOne(
		context.TODO(),
		user,
	)

	if err != nil {
		sendError(w, http.StatusInternalServerError, ERR_DB_INSERT)
		return
	}

	userId := ObjectId{
		result.InsertedID.(primitive.ObjectID).Hex(),
	}

	response := struct {
		ObjectId
		User
	}{
		userId,
		user,
	}

	sendJSON(w, 200, response)
}

// uses the API key extracted from one of the headers
// to fetch an user
func (userController *dbController) getUserByKey(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		sendError(w, http.StatusUnauthorized, ERR_UNAUTHORIZED)
		return
	}

	var user User
	err = userController.collection.FindOne(context.TODO(), bson.M{
		"secret": apiKey,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// query matched no results
			sendError(w, http.StatusUnauthorized, ERR_UNAUTHORIZED)
			return
		}
		sendError(w, http.StatusInternalServerError, ERR_INTERNAL_SRV)
		return
	}

	sendJSON(w, http.StatusOK, user)
}
