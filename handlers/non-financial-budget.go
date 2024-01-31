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

// NonFinancialBudgetHandler is a concrete type that implements NonFinancialBudgetHandler
type nonfinancialbudgetHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.NonFinancialBudgetService
}

// NewNonFinancialBudgetHandler initializes a new NonFinancialBudgetHandler with its dependencies
func NewNonFinancialBudgetHandler(app *celeritas.Celeritas, nonfinancialbudgetService services.NonFinancialBudgetService) NonFinancialBudgetHandler {
	return &nonfinancialbudgetHandlerImpl{
		App:     app,
		service: nonfinancialbudgetService,
	}
}

func (h *nonfinancialbudgetHandlerImpl) CreateNonFinancialBudget(w http.ResponseWriter, r *http.Request) {
	var input dto.NonFinancialBudgetDTO
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

	res, err := h.service.CreateNonFinancialBudget(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "NonFinancialBudget created successfuly", res)
}

func (h *nonfinancialbudgetHandlerImpl) UpdateNonFinancialBudget(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.NonFinancialBudgetDTO
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

	res, err := h.service.UpdateNonFinancialBudget(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "NonFinancialBudget updated successfuly", res)
}

func (h *nonfinancialbudgetHandlerImpl) DeleteNonFinancialBudget(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteNonFinancialBudget(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "NonFinancialBudget deleted successfuly")
}

func (h *nonfinancialbudgetHandlerImpl) GetNonFinancialBudgetById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetNonFinancialBudget(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *nonfinancialbudgetHandlerImpl) GetNonFinancialBudgetList(w http.ResponseWriter, r *http.Request) {
	var filter dto.NonFinancialBudgetFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetNonFinancialBudgetList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
