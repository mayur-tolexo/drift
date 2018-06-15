package lib

const PACKAGE_NAME = "github.com/mayur-tolexo/drift/lib"

//Error Messages
const INVALID_REQUEST_MSG = "Invalid Request, Please provide correct input"

// Error Code. PLEASE Map New Error Code To The HTTP Code Map Below.
const (
	VALIDATE_ERROR = 101 //Primary Validation fail
	NO_ERROR       = 0
)

// HTTP Code
const (
	STATUS_OK          = 200
	STATUS_BAD_REQUEST = 400
)

var (
	//ErrorHTTPCode : Error Code to Http Code map
	ErrorHTTPCode = map[int]int{
		NO_ERROR:       STATUS_OK,
		VALIDATE_ERROR: STATUS_BAD_REQUEST,
	}
)
