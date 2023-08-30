package xerror

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	status int
	text   string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error: [%d] %s", e.status, e.text)
}

func (e *HTTPError) Code() int {
	return e.status
}

func (e *HTTPError) Message() string {
	return e.text
}

var (
	ErrBadRequest = &HTTPError{
		status: http.StatusBadRequest,
		text:   "400 Bad Request",
	}
)
