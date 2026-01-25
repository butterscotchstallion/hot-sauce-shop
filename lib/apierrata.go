package lib

type InternalServerError struct {
	StatusCode int
	Message    string
}

type StatusForbiddenError struct {
	StatusCode int
	Message    string
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (e *StatusForbiddenError) Error() string {
	return e.Message
}

const ErrorCodeInsufficientKarma = "ERR_INSUFFICIENT_KARMA"
const ErrorCodePermissionDenied = "ERR_PERMISSION_DENIED"
const ErrorCodeUserNotFound = "ERR_USER_NOT_FOUND"
const ErrorCodeUserExists = "ERR_USER_EXISTS"
