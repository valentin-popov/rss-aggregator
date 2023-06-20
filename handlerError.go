package main

import "net/http"

func handlerError(writer http.ResponseWriter, request *http.Request) {
	sendError(writer, 400, "Bad Request.")
}
