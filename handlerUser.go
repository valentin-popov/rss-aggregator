package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/valentin-popov/rss-aggregator/auth"
	"github.com/valentin-popov/rss-aggregator/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createUser(w http.ResponseWriter, r *http.Request) {

	userData := struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)

	if err != nil {
		handleClientError(ERR_CODE_JSON, w, r)
		return
	}

	byte32Arr := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))

	user := User{
		Id:        primitive.NewObjectID(),
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Secret:    hex.EncodeToString(byte32Arr[:]),
	}

	err = db.AddUser(userToDTO(user))
	if err != nil {
		handleServerError(ERR_CODE_INS_OBJ, w, r)
	}
	sendJSON(w, http.StatusOK, user)

}

// Uses the API key extracted from one of the headers
// to fetch an user.
func getAuthUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		handleClientError(ERR_CODE_EMPTY_KEY, w, r)
		return
	}
	userDTO, err := db.GetUserByKey(apiKey)

	if err != nil {
		handleClientError(ERR_CODE_UNAUTHORIZED, w, r)
		return
	}
	sendJSON(w, http.StatusOK, dtoToUser(userDTO))
}

// Maps a DB user to an user, excludes some unwanted properties.
func dtoToUser(dto db.UserDTO) User {
	return User{
		Id:        dto.Id,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Secret:    dto.Secret,
		CreatedAt: dto.CreatedAt,
	}
}

// Probably needs to be moved to db.
func userToDTO(user User) db.UserDTO {
	return db.UserDTO{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Secret:    user.Secret,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
	}
}
