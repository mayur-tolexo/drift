package lib

import "github.com/go-pg/pg"

//VError : validation error
func VError(debugMsg ...string) error {
	return newError(INVALID_REQUEST_MSG, nil, VALIDATE_ERROR, debugMsg...)
}

//VErrorWitMsg : validation error with user message
func VErrorWitMsg(msg string, debugMsg ...string) error {
	return newError(msg, nil, VALIDATE_ERROR, debugMsg...)
}

//ScoutRunningError : scout running error
func ScoutRunningError(debugMsg ...string) error {
	return newError(SCOUT_RUNNING_ERROR_MSG, nil, SCOUT_RUNNING_ERROR, debugMsg...)
}

//NoUpdateError : error occured when blank input is given to update
func NoUpdateError(debugMsg ...string) error {
	return newError(BLANK_REQ_MSG, nil, VALIDATE_ERROR, debugMsg...)
}

//CustomError : custom error
func CustomError(debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, nil, CUSTOM_ERROR, debugMsg...)
}

//UnmarshalError : error occured while unmarshal
func UnmarshalError(err error, debugMsg ...string) error {
	return newError(INVALID_REQUEST_MSG, err, VALIDATE_ERROR, debugMsg...)
}

//MiscError : error occured while processing
func MiscError(err error, debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, err, MISC_ERROR, debugMsg...)
}

//ConnError : error occured while creating connection
func ConnError(err error, debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, err, DB_ERROR, debugMsg...)
}

//SelectError : error occured while select query
func SelectError(err error, debugMsg ...string) error {
	if err == pg.ErrNoRows {
		return NotFoundError()
	}
	return newError(SERVER_ERROR_MSG, err, SELECT_QUERY_ERROR, debugMsg...)
}

//SearchError : error occured while elastic search query
func SearchError(err error, debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, err, SEARCH_QUERY_ERROR, debugMsg...)
}

//InsertError : error occured while insert query
func InsertError(err error, debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, err, INSERT_QUERY_ERROR, debugMsg...)
}

//UpdateError : error occured while update query
func UpdateError(err error, debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, err, UPDATE_QUERY_ERROR, debugMsg...)
}

//DeleteError : error occured while delete query
func DeleteError(err error, debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, err, DELETE_QUERY_ERROR, debugMsg...)
}

//TxError : error occured while starting transaction
func TxError(err error, debugMsg ...string) error {
	return newError(SERVER_ERROR_MSG, err, TRANSACTION_ERROR, debugMsg...)
}

//NotFoundError : error occured when id not found
func NotFoundError(debugMsg ...string) error {
	if len(debugMsg) == 0 {
		debugMsg = append(debugMsg, "Request ID not found")
	}
	return newError(NOT_FOUND_MSG, nil, NOT_FOUND_ERROR, debugMsg...)
}

//BadReqError : error occured while validating request
//like while typecasting request, fk in request dosn't exists
func BadReqError(err error, debugMsg ...string) error {
	return newError(INVALID_REQUEST_MSG, err, VALIDATE_ERROR, debugMsg...)
}

//ForbiddenErr : unauthorized access
func ForbiddenErr(debugMsg ...string) error {
	return newError(NO_ACCESS_MSG, nil, NO_ACCESS, debugMsg...)
}

//InvalidParamError : Error due to request validation fail
func InvalidParamError(err error) error {
	//TODO change validation error display format
	return VError(err.Error())
}

//MappingError : error occured deleting entity without removing it's mapping
func MappingError(debugMsg ...string) error {
	return newError(REMOVE_MAPPING, nil, VALIDATE_ERROR, debugMsg...)
}
