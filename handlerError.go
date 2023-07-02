package main

import "net/http"

const ERR_SRV string = "server"
const ERR_CLIENT string = "client"

var errorMap = map[string]struct {
	status  int
	message string
}{
	ERR_CODE_BAD_REQ: {
		status:  http.StatusBadRequest,
		message: ERR_MSG_BAD_REQ,
	},
	ERR_CODE_UNAUTHORIZED: {
		status:  http.StatusUnauthorized,
		message: ERR_MSG_UNAUTHORIZED,
	},
	ERR_CODE_INV_ID: {
		status:  http.StatusBadRequest,
		message: ERR_MSG_INV_ID,
	},
	ERR_CODE_JSON: {
		status:  http.StatusBadRequest,
		message: ERR_MSG_JSON,
	},
	ERR_CODE_EMPTY_KEY: {
		status:  http.StatusBadRequest,
		message: ERR_MSG_EMPTY_KEY,
	},

	ERR_CODE_INTERNAL_SRV: {
		status:  http.StatusInternalServerError,
		message: ERR_MSG_INTERNAL_SRV,
	},
	ERR_CODE_INS_OBJ: {
		status:  http.StatusInternalServerError,
		message: ERR_MSG_INS_OBJ,
	},
	ERR_CODE_FETCH_DOCS: {
		status:  http.StatusInternalServerError,
		message: ERR_MSG_FETCH_DOCS,
	},
}

// Maps an error string to a http error code.
func handleError(errMsg string, errType string, w http.ResponseWriter, r *http.Request) {

	err, ok := errorMap[errMsg]

	if ok {
		sendError(w, err.status, err.message)
		return
	}

	// if no custom error definition was found
	switch errType {
	case ERR_CLIENT:
		sendError(w, http.StatusBadRequest, ERR_MSG_BAD_REQ)
		return

	case ERR_SRV:
	default:
		sendError(w, http.StatusInternalServerError, ERR_MSG_INTERNAL_SRV)
		return
	}

}

func handleClientError(errCode string, w http.ResponseWriter, r *http.Request) {
	handleError(errCode, ERR_CLIENT, w, r)
}

func handleServerError(errCode string, w http.ResponseWriter, r *http.Request) {
	handleError(errCode, ERR_SRV, w, r)
}
