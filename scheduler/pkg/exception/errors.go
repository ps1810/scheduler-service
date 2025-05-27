package exception

import "net/http"

var (
	InvalidRequestBodyError = createFixedExceptionErrors(
		http.StatusBadRequest,
		ERROR_TYPE_BAD_REQUEST,
		"invalid request body",
	)

	DataNotFoundError = createFixedExceptionErrors(
		http.StatusNotFound,
		ERROR_TYPE_NOT_FOUND,
		"data is not found",
	)

	ValidationFailedError = createFixedExceptionErrors(
		http.StatusUnprocessableEntity,
		ERROR_TYPE_VALIDATION_ERROR,
		"validation failed",
	)

	UpdateFailedError = createFixedExceptionErrors(
		http.StatusInternalServerError,
		ERROR_TYPE_UPDATE_FAILED,
		"Database updation failed",
	)

	FailedAddJobError = createFixedExceptionErrors(
		http.StatusInternalServerError,
		ERROR_TYPE_CREATE_FAILED,
		"Unable to create entry",
	)

	InvalidCronExpression = createFixedExceptionErrors(
		http.StatusBadRequest,
		ERROR_TYPE_BAD_REQUEST,
		"Invalid cron expression",
	)

	NotScheduledJob = createFixedExceptionErrors(
		http.StatusUnprocessableEntity,
		ERROR_TYPE_NOT_FOUND,
		"Job is not scheduled",
	)
)
