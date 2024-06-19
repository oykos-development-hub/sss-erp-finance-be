package handlers

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	"gitlab.sudovi.me/erp/finance-api/services"

	"github.com/go-chi/chi/v5"
	"github.com/oykos-development-hub/celeritas"
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
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	res, err := h.service.CreateNonFinancialBudgetGoal(ctx, input)
	if err != nil {
		h.App.ErrorLog.Print(err)
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
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	res, err := h.service.UpdateNonFinancialBudgetGoal(ctx, id, input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "NonFinancialBudgetGoal updated successfuly", res)
}

func (h *nonfinancialbudgetgoalHandlerImpl) DeleteNonFinancialBudgetGoal(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	err = h.service.DeleteNonFinancialBudgetGoal(ctx, id)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "NonFinancialBudgetGoal deleted successfuly")
}

func (h *nonfinancialbudgetgoalHandlerImpl) GetNonFinancialBudgetGoalById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetNonFinancialBudgetGoal(id)
	if err != nil {
		h.App.ErrorLog.Print(err)
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
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetNonFinancialBudgetGoalList(filter)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
