package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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

func (userController *dbController) createUser(w http.ResponseWriter, r *http.Request) {

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

	user := User{
		Id:        primitive.NewObjectID(),
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Secret:    hex.EncodeToString(byte32Arr[:]),
	}

	_, err = userController.collection.InsertOne(
		context.TODO(),
		user,
	)

	if err != nil {
		sendError(w, http.StatusInternalServerError, ERR_DB_INSERT)
		return
	}

	sendJSON(w, 200, user)
}

// uses the API key extracted from one of the headers
// to fetch an user
func (userController *dbController) getAuthUser(w http.ResponseWriter, r *http.Request) {
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

	sendJSON(w, http.StatusOK, user)
}

func (userController *dbController) getUserByKey(apiKey string) (*User, error) {

	var user User
	err := userController.collection.FindOne(context.TODO(), bson.M{
		"secret": apiKey,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(ERR_UNAUTHORIZED)
		}
		return nil, errors.New(ERR_INTERNAL_SRV)
	}
	return &user, nil
}
