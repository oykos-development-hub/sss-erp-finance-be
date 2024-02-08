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

// FilledFinancialBudgetHandler is a concrete type that implements FilledFinancialBudgetHandler
type filledfinancialbudgetHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.FilledFinancialBudgetService
}

// NewFilledFinancialBudgetHandler initializes a new FilledFinancialBudgetHandler with its dependencies
func NewFilledFinancialBudgetHandler(app *celeritas.Celeritas, filledfinancialbudgetService services.FilledFinancialBudgetService) FilledFinancialBudgetHandler {
	return &filledfinancialbudgetHandlerImpl{
		App:     app,
		service: filledfinancialbudgetService,
	}
}

func (h *filledfinancialbudgetHandlerImpl) CreateFilledFinancialBudget(w http.ResponseWriter, r *http.Request) {
	var input dto.FilledFinancialBudgetDTO
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

	res, err := h.service.CreateFilledFinancialBudget(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FilledFinancialBudget created successfuly", res)
}

func (h *filledfinancialbudgetHandlerImpl) UpdateFilledFinancialBudget(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.FilledFinancialBudgetDTO
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

	res, err := h.service.UpdateFilledFinancialBudget(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FilledFinancialBudget updated successfuly", res)
}

func (h *filledfinancialbudgetHandlerImpl) DeleteFilledFinancialBudget(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteFilledFinancialBudget(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "FilledFinancialBudget deleted successfuly")
}

func (h *filledfinancialbudgetHandlerImpl) GetFilledFinancialBudgetById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFilledFinancialBudget(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *filledfinancialbudgetHandlerImpl) GetFilledFinancialBudgetList(w http.ResponseWriter, r *http.Request) {
	var filter dto.FilledFinancialBudgetFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetFilledFinancialBudgetList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
