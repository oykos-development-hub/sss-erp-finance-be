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

// NonFinancialBudgetGoalHandler is a concrete type that implements NonFinancialBudgetGoalHandler
type nonfinancialbudgetgoalHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.NonFinancialBudgetGoalService
}

// NewNonFinancialBudgetGoalHandler initializes a new NonFinancialBudgetGoalHandler with its dependencies
func NewNonFinancialBudgetGoalHandler(app *celeritas.Celeritas, nonfinancialbudgetgoalService services.NonFinancialBudgetGoalService) NonFinancialBudgetGoalHandler {
	return &nonfinancialbudgetgoalHandlerImpl{
		App:     app,
		service: nonfinancialbudgetgoalService,
	}
}

func (h *nonfinancialbudgetgoalHandlerImpl) CreateNonFinancialBudgetGoal(w http.ResponseWriter, r *http.Request) {
	var input dto.NonFinancialBudgetGoalDTO
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

	res, err := h.service.CreateNonFinancialBudgetGoal(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "NonFinancialBudgetGoal created successfuly", res)
}

func (h *nonfinancialbudgetgoalHandlerImpl) UpdateNonFinancialBudgetGoal(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.NonFinancialBudgetGoalDTO
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

	res, err := h.service.UpdateNonFinancialBudgetGoal(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "NonFinancialBudgetGoal updated successfuly", res)
}

func (h *nonfinancialbudgetgoalHandlerImpl) DeleteNonFinancialBudgetGoal(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteNonFinancialBudgetGoal(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "NonFinancialBudgetGoal deleted successfuly")
}

func (h *nonfinancialbudgetgoalHandlerImpl) GetNonFinancialBudgetGoalById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetNonFinancialBudgetGoal(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *nonfinancialbudgetgoalHandlerImpl) GetNonFinancialBudgetGoalList(w http.ResponseWriter, r *http.Request) {
	var filter dto.NonFinancialBudgetGoalFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetNonFinancialBudgetGoalList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
