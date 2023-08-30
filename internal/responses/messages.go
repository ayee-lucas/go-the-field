package responses

/** Message Responses */

const (
	OK                        = "success"
	DATA_NOT_FOUND_MESSAGE    = "Oh no! The data you're seeking isn't available at the moment. Feel free to double-check your details, and keep an eye out for updates."
	BODY_PARSE_MESSAGE        = "Oops! Something got mixed up while processing your request. Please double-check the data format and try again."
	UNAUTHORIZED_MESSAGE      = "It seems you don't have the necessary authorization to access this resource."
	REQUIRED_FIELD            = "Field required."
	INVALID_ID                = "ID does not match any records."
	SAVE_DATA_MESSAGE_ERROR   = "We couldn't save your data this time."
	DELETE_DATA_MESSAGE_ERROR = "Sorry, we couldn't delete the data as expected"

	SIGN_UP_MESSAGE_ERROR = "Something went wrong during the sign-up process."

	/** User error messages */
	P_FINISHED_MESSAGE_ERROR = "The user's profile has already been completed."
	P_UPDATE_ERROR           = "Oh no! We encountered an error while trying to update your profile. This might be due to a technical issue on our end"
	U_ORG_MESSAGE_ERROR      = "This user already has an org or athlete attached"
)
