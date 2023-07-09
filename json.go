package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendJSON(writer http.ResponseWriter, statusCode int, payload interface{}) {

	writer.Header().Add("Content-Type", "application/json")

	if payload == nil && statusCode == http.StatusNoContent {
		// payload sent empty intentionally
		writer.WriteHeader(statusCode)
		return
	}

	data, err := json.Marshal(payload)

	if err != nil {

		fmt.Printf("Unable to marshal JSON: %v", payload)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(statusCode)
	writer.Write(data)
}

func sendError(w http.ResponseWriter, statusCode int, errorMessage string) {
	if statusCode >= 500 {
		log.Println("5XX Server Error.")
	}
	sendJSON(w, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: errorMessage,
	})
}
