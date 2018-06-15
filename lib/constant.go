package lib

const (
	PACKAGE_NAME = "github.com/mayur-tolexo/drift/lib"
)

//Error Messages
const (
	REMOVE_MAPPING            = "Please remove it's all mapping before delete"
	INVALID_REQUEST_MSG       = "Invalid Request, Please contact system adminstrator for further clarification."
	NO_ACCESS_MSG             = "You are not allowed to perform this operation. Please contact system adminstrator."
	BLANK_REQ_MSG             = "Blank request, Please provide input to process"
	INVALID_USER_MSG          = "Unknown User, Please login again"
	NOT_FOUND_MSG             = "Record not found"
	SERVER_ERROR_MSG          = "Sorry unable to process this request. Please Try Again"
	DUPLICATE_ENTRY_ERROR_MSG = "Record Already Exists"
	SCOUT_RUNNING_ERROR_MSG   = "Scout Already Running"
)

// Error Code. PLEASE Map New Error Code To The HTTP Code Map Below.
const (
	NO_ERROR              = 0
	VALIDATE_ERROR        = 101 //Primary Validation fail
	FORMAT_ERROR          = 102 //Formatter Error
	DB_ERROR              = 103 //Database Connection error
	SELECT_QUERY_ERROR    = 104 //Select Query error
	CREATE_QUERY_ERROR    = 105 //Create Query error
	INSERT_QUERY_ERROR    = 106 //Insert Query error
	UPDATE_QUERY_ERROR    = 107 //Update Query error
	DELETE_QUERY_ERROR    = 108 //Delete Query error
	CUSTOM_ERROR          = 109 //Custom error
	TRANSACTION_ERROR     = 110 //Transaction error
	NOT_FOUND_ERROR       = 111 //Request ID not found error
	NO_ACCESS             = 112 //Unauthorized access
	DUPLICATE_ENTRY_ERROR = 113 //Duplicate entry error
	MISC_ERROR            = 114 //Mislanious error
	SCOUT_RUNNING_ERROR   = 115 //Elasticsearch scout already running error
	SEARCH_QUERY_ERROR    = 116 //Search Query error

)

// HTTP Code
const (
	STATUS_OK                    = 200
	STATUS_BAD_REQUEST           = 400
	STATUS_NOT_FOUND             = 404
	STATUS_INTERNAL_SERVER_ERROR = 500
	STATUS_SERVICE_UNAVAILABLE   = 503
	STATUS_FORBIDDEN             = 403
)

var (
	//ErrorHTTPCode : Error Code to Http Code map
	ErrorHTTPCode = map[int]int{
		NO_ERROR:              STATUS_OK,
		VALIDATE_ERROR:        STATUS_BAD_REQUEST,
		FORMAT_ERROR:          STATUS_OK,
		DB_ERROR:              STATUS_INTERNAL_SERVER_ERROR,
		SELECT_QUERY_ERROR:    STATUS_BAD_REQUEST,
		CREATE_QUERY_ERROR:    STATUS_BAD_REQUEST,
		INSERT_QUERY_ERROR:    STATUS_BAD_REQUEST,
		UPDATE_QUERY_ERROR:    STATUS_BAD_REQUEST,
		DELETE_QUERY_ERROR:    STATUS_BAD_REQUEST,
		CUSTOM_ERROR:          STATUS_BAD_REQUEST,
		TRANSACTION_ERROR:     STATUS_INTERNAL_SERVER_ERROR,
		NOT_FOUND_ERROR:       STATUS_OK,
		NO_ACCESS:             STATUS_FORBIDDEN,
		DUPLICATE_ENTRY_ERROR: STATUS_BAD_REQUEST,
		MISC_ERROR:            STATUS_INTERNAL_SERVER_ERROR,
		SCOUT_RUNNING_ERROR:   STATUS_SERVICE_UNAVAILABLE,
		SEARCH_QUERY_ERROR:    STATUS_BAD_REQUEST,
	}
)
