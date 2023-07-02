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

func createFeed(w http.ResponseWriter, r *http.Request) {

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
		Name string `json:"name"`
		Url  string `json:"url"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&feedData)

	if err != nil {
		handleClientError(ERR_CODE_JSON, w, r)
		return
	}

	feed := Feed{
		Id:        primitive.NewObjectID(),
		Name:      feedData.Name,
		Url:       feedData.Url,
		UserId:    user.Id,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	err = db.AddFeed(feedToDTO(feed))

	if err != nil {
		handleServerError(ERR_CODE_INS_OBJ, w, r)
		return
	}

	sendJSON(w, http.StatusOK, feed)

}

func findFeeds(w http.ResponseWriter, r *http.Request) {

	feedDTOs, err := db.GetFeeds("", "")
	feeds := []Feed{}
	for _, dto := range feedDTOs {
		feeds = append(feeds, dtoToFeed(dto))
	}
	if err != nil {
		handleServerError(ERR_CODE_FETCH_DOCS, w, r)
		return
	}

	sendJSON(w, http.StatusOK, feeds)
}

func findFeed(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	feedDTO, err := db.GetFeedById(pathParts[len(pathParts)-1])

	if err != nil {
		if err.Error() == ERR_CODE_INV_ID {
			handleClientError(ERR_CODE_INV_ID, w, r)
			return
		}
		handleServerError(ERR_CODE_FETCH_DOCS, w, r)

		return
	}

	sendJSON(w, http.StatusOK, dtoToFeed(feedDTO))
}

func feedToDTO(feed Feed) db.FeedDTO {
	return db.FeedDTO{
		Id:        feed.Id,
		Name:      feed.Name,
		Url:       feed.Url,
		UserId:    feed.UserId,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.CreatedAt,
	}
}

func dtoToFeed(dto db.FeedDTO) Feed {
	return Feed{
		Id:        dto.Id,
		Name:      dto.Name,
		Url:       dto.Url,
		UserId:    dto.UserId,
		CreatedAt: dto.CreatedAt,
	}
}
