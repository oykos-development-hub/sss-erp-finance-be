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

// CurrentBudgetHandler is a concrete type that implements CurrentBudgetHandler
type currentbudgetHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.CurrentBudgetService
}

// NewCurrentBudgetHandler initializes a new CurrentBudgetHandler with its dependencies
func NewCurrentBudgetHandler(app *celeritas.Celeritas, currentbudgetService services.CurrentBudgetService) CurrentBudgetHandler {
	return &currentbudgetHandlerImpl{
		App:     app,
		service: currentbudgetService,
	}
}

func (h *currentbudgetHandlerImpl) CreateCurrentBudget(w http.ResponseWriter, r *http.Request) {
	var input dto.CurrentBudgetDTO
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

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	res, err := h.service.CreateCurrentBudget(ctx, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "CurrentBudget created successfuly", res)
}

func (h *currentbudgetHandlerImpl) GetCurrentBudgetById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetCurrentBudget(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *currentbudgetHandlerImpl) GetCurrentBudgetList(w http.ResponseWriter, r *http.Request) {
	var filter dto.CurrentBudgetFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetCurrentBudgetList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}

func (h *currentbudgetHandlerImpl) GetAcctualCurrentBudget(w http.ResponseWriter, r *http.Request) {
	organizationUnitID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetAcctualCurrentBudget(organizationUnitID)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}
