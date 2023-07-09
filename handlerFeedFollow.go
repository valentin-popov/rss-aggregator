package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/valentin-popov/rss-aggregator/auth"
	"github.com/valentin-popov/rss-aggregator/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func followFeed(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		handleClientError(ERR_CODE_EMPTY_KEY, w, r)
		return
	}

	user, err := db.GetUserByKey(apiKey)
	if err != nil {
		handleClientError(ERR_CODE_UNAUTHORIZED, w, r)
		return
	}

	feedData := struct {
		FeedId primitive.ObjectID `json:"feedId"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&feedData)

	if err != nil {
		handleClientError(ERR_CODE_JSON, w, r)
		return
	}

	feedFollow := FeedFollow{
		Id:        primitive.NewObjectID(),
		UserId:    user.Id,
		FeedId:    feedData.FeedId,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	err = db.AddFeedFollow(feedFollowToDTO(feedFollow))

	if err != nil {
		errType, ok := customErrors["insert"][err.Error()]

		if ok {
			if errType == CLIENT {
				handleClientError(err.Error(), w, r)
				return
			}
			if errType == SERVER {
				handleServerError(err.Error(), w, r)
				return
			}
		}
		handleServerError(ERR_CODE_INS_OBJ, w, r)
		return
	}

	sendJSON(w, http.StatusOK, feedFollow)

}

func unfollowFeed(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		handleClientError(ERR_CODE_EMPTY_KEY, w, r)
		return
	}

	user, err := db.GetUserByKey(apiKey)
	if err != nil {
		handleClientError(ERR_CODE_UNAUTHORIZED, w, r)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	feedId, err := primitive.ObjectIDFromHex(pathParts[len(pathParts)-1])

	if err != nil {
		handleClientError(ERR_CODE_INV_ID, w, r)
		return
	}

	fields := db.FeedFollowDeleteFilter{
		FeedId: feedId,
		UserId: user.Id,
	}

	err = db.DeleteFeedFollow(fields)

	if err != nil {
		if err.Error() == ERR_CODE_NO_DOC {
			handleClientError(err.Error(), w, r)
		}

		handleServerError(ERR_CODE_DEL, w, r)
		return
	}

	sendJSON(w, http.StatusNoContent, nil)

}

func findFollowedFeeds(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		handleClientError(ERR_CODE_EMPTY_KEY, w, r)
		return
	}

	user, err := db.GetUserByKey(apiKey)
	if err != nil {
		handleClientError(ERR_CODE_UNAUTHORIZED, w, r)
		return
	}

	if err != nil {
		handleServerError(ERR_CODE_INTERNAL_SRV, w, r)
		return
	}

	feedFollowDTOs, err := db.GetFollowedFeeds(user.Id)
	feedFollows := []FeedFollow{}

	for _, dto := range feedFollowDTOs {
		feedFollows = append(feedFollows, dtoToFeedFollow(dto))
	}

	if err != nil {
		handleServerError(ERR_CODE_FETCH_DOCS, w, r)
		return
	}
	sendJSON(w, http.StatusOK, feedFollows)

}

func feedFollowToDTO(feedFollow FeedFollow) db.FeedFollowDTO {
	return db.FeedFollowDTO(feedFollow)
}

func dtoToFeedFollow(dto db.FeedFollowDTO) FeedFollow {
	return FeedFollow(dto)
}
