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

// SpendingDynamicHandler is a concrete type that implements SpendingDynamicHandler
type spendingdynamicHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.SpendingDynamicService
}

// NewSpendingDynamicHandler initializes a new SpendingDynamicHandler with its dependencies
func NewSpendingDynamicHandler(app *celeritas.Celeritas, spendingdynamicService services.SpendingDynamicService) SpendingDynamicHandler {
	return &spendingdynamicHandlerImpl{
		App:     app,
		service: spendingdynamicService,
	}
}

func (h *spendingdynamicHandlerImpl) CreateSpendingDynamic(w http.ResponseWriter, r *http.Request) {
	var input dto.SpendingDynamicDTO
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

	res, err := h.service.CreateSpendingDynamic(input)
	if err != nil {
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
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *spendingdynamicHandlerImpl) GetBudgetSpendingDynamic(w http.ResponseWriter, r *http.Request) {
	budgetID, _ := strconv.Atoi(chi.URLParam(r, "budget_id"))
	unitID, _ := strconv.Atoi(chi.URLParam(r, "unit_id"))

	res, err := h.service.GetSpendingDynamic(budgetID, unitID)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *spendingdynamicHandlerImpl) GetActual(w http.ResponseWriter, r *http.Request) {
	budgetID, _ := strconv.Atoi(chi.URLParam(r, "budget_id"))
	unitID, _ := strconv.Atoi(chi.URLParam(r, "unit_id"))
	accountID, _ := strconv.Atoi(chi.URLParam(r, "account_id"))

	res, err := h.service.GetActual(budgetID, unitID, accountID)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}
