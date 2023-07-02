package main

import "github.com/valentin-popov/rss-aggregator/db"

// error codes used internally

// client
const ERR_CODE_BAD_REQ string = "bad_request"
const ERR_CODE_UNAUTHORIZED string = "unauthorized"
const ERR_CODE_INV_ID string = db.ERR_INV_ID
const ERR_CODE_JSON string = "json_err"
const ERR_CODE_EMPTY_KEY string = "empty_api_key"

// server
const ERR_CODE_INTERNAL_SRV string = db.ERR_INTERNAL_SRV
const ERR_CODE_INS_OBJ string = "err_ins_obj"
const ERR_CODE_FETCH_DOCS string = "err_fetch_docs"

// error messages
const ERR_MSG_BAD_REQ string = "Bad Request."
const ERR_MSG_UNAUTHORIZED string = "Unauthorized."
const ERR_MSG_INV_ID string = "Invalid Object ID."
const ERR_MSG_JSON string = "Could not parse JSON."
const ERR_MSG_EMPTY_KEY string = "API Key is empty."

const ERR_MSG_INTERNAL_SRV string = "Internal Server Error."
const ERR_MSG_INS_OBJ string = "Could not insert object into database."
const ERR_MSG_FETCH_DOCS string = "Could not retrieve documents from the database."

// fatal error messages
const ERR_MSG_PORT_UNDEF string = "Port is undefined."
const ERR_MSG_START_SRV string = "Could not start server."
