package handlers

import (
	"net/http"
	"strconv"

	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	"gitlab.sudovi.me/erp/finance-api/services"

	"github.com/oykos-development-hub/celeritas"
	"github.com/go-chi/chi/v5"
)

// SalaryAdditionalExpenseHandler is a concrete type that implements SalaryAdditionalExpenseHandler
type salaryadditionalexpenseHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.SalaryAdditionalExpenseService
}

// NewSalaryAdditionalExpenseHandler initializes a new SalaryAdditionalExpenseHandler with its dependencies
func NewSalaryAdditionalExpenseHandler(app *celeritas.Celeritas, salaryadditionalexpenseService services.SalaryAdditionalExpenseService) SalaryAdditionalExpenseHandler {
	return &salaryadditionalexpenseHandlerImpl{
		App:     app,
		service: salaryadditionalexpenseService,
	}
}

func (h *salaryadditionalexpenseHandlerImpl) CreateSalaryAdditionalExpense(w http.ResponseWriter, r *http.Request) {
	var input dto.SalaryAdditionalExpenseDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, err := h.service.CreateSalaryAdditionalExpense(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "SalaryAdditionalExpense created successfuly", res)
}

func (h *salaryadditionalexpenseHandlerImpl) UpdateSalaryAdditionalExpense(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.SalaryAdditionalExpenseDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, err := h.service.UpdateSalaryAdditionalExpense(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "SalaryAdditionalExpense updated successfuly", res)
}

func (h *salaryadditionalexpenseHandlerImpl) DeleteSalaryAdditionalExpense(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteSalaryAdditionalExpense(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "SalaryAdditionalExpense deleted successfuly")
}

func (h *salaryadditionalexpenseHandlerImpl) GetSalaryAdditionalExpenseById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetSalaryAdditionalExpense(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *salaryadditionalexpenseHandlerImpl) GetSalaryAdditionalExpenseList(w http.ResponseWriter, r *http.Request) {
	var filter dto.SalaryAdditionalExpenseFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetSalaryAdditionalExpenseList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
