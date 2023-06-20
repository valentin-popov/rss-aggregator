package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendJSON(writer http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {

		fmt.Printf("Unable to marshal JSON: %v", payload)
		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
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
