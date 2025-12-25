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
