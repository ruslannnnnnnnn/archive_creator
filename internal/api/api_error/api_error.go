package api_error

import "net/http"

type IApiError interface {
	HttpStatusCode() int
	Error() string
}

type InternalError struct{}

func (t *InternalError) HttpStatusCode() int {
	return http.StatusInternalServerError
}

func (t *InternalError) Error() string {
	return "Internal Server Error"
}

type BadRequest struct {
	Message string
}

func (t *BadRequest) HttpStatusCode() int {
	return http.StatusBadRequest
}

func (t *BadRequest) Error() string {
	return t.Message
}

type NotFound struct {
	Message string
}

func (t *NotFound) HttpStatusCode() int {
	return http.StatusNotFound
}

func (t *NotFound) Error() string {
	return t.Message
}

type UnsupportedMediaType struct {
	Message string
}

func (t *UnsupportedMediaType) HttpStatusCode() int {
	return http.StatusUnsupportedMediaType
}
func (t *UnsupportedMediaType) Error() string {
	return t.Message
}
