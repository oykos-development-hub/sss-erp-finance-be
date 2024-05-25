package errors

import "net/http"

const (
	BadRequestCode = 100
	NotFoundCode   = 101
	InternalCode   = 104

	SingleMonthSpendingReleaseCode = 200
	ReleaseInCurrentMonthCode      = 201
	NotEnoughFundsCode             = 202
)

func httpStatusCode(code int) int {
	switch code {
	case BadRequestCode, SingleMonthSpendingReleaseCode, ReleaseInCurrentMonthCode, NotEnoughFundsCode:
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
	case SingleMonthSpendingReleaseCode:
		return "Only single release is allowed per month"
	case ReleaseInCurrentMonthCode:
		return "Release is possible only in the current month"
	case NotEnoughFundsCode:
		return "Not enough funds"
	default:
		return "Internal server error"
	}
}
