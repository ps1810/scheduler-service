package exception

type errorType string

const (
	ERROR_TYPE_BAD_REQUEST      errorType = "BadRequest"
	ERROR_TYPE_NOT_FOUND        errorType = "NotFound"
	ERROR_TYPE_VALIDATION_ERROR errorType = "ValidationError"
	ERROR_TYPE_UPDATE_FAILED    errorType = "UpdateFailed"
	ERROR_TYPE_CREATE_FAILED    errorType = "CreateFailed"
)
