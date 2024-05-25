package errors

import "net/http"

const (
	BadRequestCode          = 100
	NotFoundCode            = 101
	InternalCode            = 104
	MicroserviceRequestCode = 105
)

func httpStatusCode(code int) int {
	switch code {
	case BadRequestCode:
		return http.StatusBadRequest
	case NotFoundCode:
		return http.StatusNotFound
	case InternalCode:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func prettyMsg(code int) string {
	switch code {
	case BadRequestCode:
		return "Bad request"
	case NotFoundCode:
		return "Not found"
	default:
		return "Internal server error"
	}
}
