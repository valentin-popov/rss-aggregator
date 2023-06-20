package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userController dbController = dbController{
	collection: getDB().Collection(
		collections["user"],
	),
}

func (userController *dbController) create(w http.ResponseWriter, r *http.Request) {

	// userData := user{}
	userData := struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)

	if err != nil {
		sendError(w, BAD_REQUEST, ERR_JSON)
		return
	}

	user := User{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := userController.collection.InsertOne(
		context.TODO(),
		user,
	)

	if err != nil {
		sendError(w, INTERNAL_SERVER_ERROR, ERR_DB_INSERT)
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
