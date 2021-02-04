package throw

import (
	"fmt"
)

type ResponseError struct {
	statusCode int
	err        error
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("status %d: err %v", e.statusCode, e.err)
}

func NewJWTError(statusCode int, err error) *ResponseError {
	return &ResponseError{
		statusCode: statusCode,
		err:        err,
	}
}

func (e *ResponseError) GetStatusCode() int {
	return e.statusCode
}

func (e *ResponseError) GetError() error {
	return e.err
}
