package fielderr

import "net/http"

const (
	CodeBadRequest = iota
	CodeNotFound
	CodeInternal
	CodeUnauthorized
	CodeConflict
	CodeForbidden
	CodeNoContent
)

var httpCodes = map[int]int{
	CodeBadRequest:   http.StatusBadRequest,
	CodeNotFound:     http.StatusNotFound,
	CodeInternal:     http.StatusInternalServerError,
	CodeUnauthorized: http.StatusUnauthorized,
	CodeForbidden:    http.StatusForbidden,
	CodeConflict:     http.StatusConflict,
	CodeNoContent:    http.StatusNoContent,
}
