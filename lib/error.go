package lib

//BadReqError : error occured while validating request
//like while typecasting request, fk in request dosn't exists
func BadReqError(err error, debugMsg ...string) error {
	return newError(INVALID_REQUEST_MSG, err, VALIDATE_ERROR, debugMsg...)
}

//UnmarshalError : error occured while unmarshal
func UnmarshalError(err error, debugMsg ...string) error {
	return newError(INVALID_REQUEST_MSG, err, VALIDATE_ERROR, debugMsg...)
}
