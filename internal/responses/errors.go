package responses

/** Error Responses */

const (
	DATA_RETRIEVAL       = "Data Retrieval Failure "
	DATA_NOT_FOUND       = "The requested data could not be located "
	PARSE_BODY_ERROR     = "Failed to parse request body "
	UPDATE_DATA_ERROR    = "Failed to update data "
	INVALID_INPUT        = "Invalid input data "
	INVALID_HEADER_ERROR = "Invalid session header "
	INVALID_EMAIL_FORMAT = "Invalid email format "

	/** USER ERRORS */
	PROFILE_FINISHED_ERROR = "Profile Completion Conflict "
	GET_SESSION_ERROR      = "Failed to get Session "
)
