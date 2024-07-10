package handlers

import (
	"net/http"
	"strconv"

	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	"gitlab.sudovi.me/erp/finance-api/services"

	"github.com/go-chi/chi/v5"
	"github.com/oykos-development-hub/celeritas"
)

// AdditionalExpenseHandler is a concrete type that implements AdditionalExpenseHandler
type additionalexpenseHandlerImpl struct {
	App             *celeritas.Celeritas
	service         services.AdditionalExpenseService
	errorLogService services.ErrorLogService
}

// NewAdditionalExpenseHandler initializes a new AdditionalExpenseHandler with its dependencies
func NewAdditionalExpenseHandler(app *celeritas.Celeritas, additionalexpenseService services.AdditionalExpenseService, errorLogService services.ErrorLogService) AdditionalExpenseHandler {
	return &additionalexpenseHandlerImpl{
		App:             app,
		service:         additionalexpenseService,
		errorLogService: errorLogService,
	}
}

func (h *additionalexpenseHandlerImpl) DeleteAdditionalExpense(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteAdditionalExpense(id)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "AdditionalExpense deleted successfuly")
}

func (h *additionalexpenseHandlerImpl) GetAdditionalExpenseById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetAdditionalExpense(id)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *additionalexpenseHandlerImpl) GetAdditionalExpenseList(w http.ResponseWriter, r *http.Request) {
	var filter dto.AdditionalExpenseFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetAdditionalExpenseList(filter)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
