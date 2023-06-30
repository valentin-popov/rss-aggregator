package main

import "net/http"

var errorMap = map[string]int{
	ERR_UNAUTHORIZED: http.StatusUnauthorized,
	ERR_INTERNAL_SRV: http.StatusInternalServerError,
}

type errorObject struct {
	errorMessage string
}

func (errorObj *errorObject) handleError(w http.ResponseWriter, r *http.Request) {

	errorCode, ok := errorMap[errorObj.errorMessage]

	if ok {
		sendError(w, errorCode, errorObj.errorMessage)
		return
	}

	sendError(w, errorMap[ERR_INTERNAL_SRV], ERR_INTERNAL_SRV)

}
