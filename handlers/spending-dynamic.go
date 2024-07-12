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

// SpendingDynamicHandler is a concrete type that implements SpendingDynamicHandler
type spendingdynamicHandlerImpl struct {
	App             *celeritas.Celeritas
	service         services.SpendingDynamicService
	errorLogService services.ErrorLogService
}

// NewSpendingDynamicHandler initializes a new SpendingDynamicHandler with its dependencies
func NewSpendingDynamicHandler(app *celeritas.Celeritas, spendingdynamicService services.SpendingDynamicService, errorLogService services.ErrorLogService) SpendingDynamicHandler {
	return &spendingdynamicHandlerImpl{
		App:             app,
		service:         spendingdynamicService,
		errorLogService: errorLogService,
	}
}

func (h *spendingdynamicHandlerImpl) CreateSpendingDynamic(w http.ResponseWriter, r *http.Request) {
	budgetID, budgetErr := strconv.Atoi(chi.URLParam(r, "budget_id"))
	unitID, unitErr := strconv.Atoi(chi.URLParam(r, "unit_id"))

	if budgetErr != nil || unitErr != nil {
		h.App.ErrorLog.Print(errors.ErrInvalidInput)
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var input []dto.SpendingDynamicDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
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
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	err = h.service.CreateSpendingDynamic(ctx, budgetID, unitID, input)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	res, err := h.service.GetSpendingDynamic(nil, &budgetID, &unitID, nil)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "SpendingDynamic created successfuly", res)
}

func (h *spendingdynamicHandlerImpl) GetBudgetSpendingDynamicHistory(w http.ResponseWriter, r *http.Request) {
	budgetID, _ := strconv.Atoi(chi.URLParam(r, "budget_id"))
	unitID, _ := strconv.Atoi(chi.URLParam(r, "unit_id"))

	res, err := h.service.GetSpendingDynamicHistory(budgetID, unitID)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *spendingdynamicHandlerImpl) GetBudgetSpendingDynamic(w http.ResponseWriter, r *http.Request) {
	budgetID, _ := strconv.Atoi(chi.URLParam(r, "budget_id"))
	unitID, _ := strconv.Atoi(chi.URLParam(r, "unit_id"))

	var filter dto.SpendingDynamicFilter

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, err := h.service.GetSpendingDynamic(nil, &budgetID, &unitID, filter.Version)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *spendingdynamicHandlerImpl) GetActual(w http.ResponseWriter, r *http.Request) {
	budgetID, _ := strconv.Atoi(chi.URLParam(r, "budget_id"))
	unitID, _ := strconv.Atoi(chi.URLParam(r, "unit_id"))
	accountID, _ := strconv.Atoi(chi.URLParam(r, "account_id"))

	res, err := h.service.GetActual(budgetID, unitID, accountID, 1)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}
