package handlers

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
	"gitlab.sudovi.me/erp/finance-api/services"

	"github.com/go-chi/chi/v5"
	"github.com/oykos-development-hub/celeritas"
)

// SpendingReleaseHandler is a concrete type that implements SpendingReleaseHandler
type spendingreleaseHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.SpendingReleaseService
}

// NewSpendingReleaseHandler initializes a new SpendingReleaseHandler with its dependencies
func NewSpendingReleaseHandler(app *celeritas.Celeritas, spendingreleaseService services.SpendingReleaseService) SpendingReleaseHandler {
	return &spendingreleaseHandlerImpl{
		App:     app,
		service: spendingreleaseService,
	}
}

func (h *spendingreleaseHandlerImpl) CreateSpendingRelease(w http.ResponseWriter, r *http.Request) {
	budgetID, _ := strconv.Atoi(chi.URLParam(r, "budget_id"))
	unitID, _ := strconv.Atoi(chi.URLParam(r, "unit_id"))

	var input []dto.SpendingReleaseDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.BadRequestCode, errors.NewBadRequestError("input validation"), validator.Errors)
		return
	}

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.ErrUnauthorized, err)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	res, err := h.service.CreateSpendingRelease(ctx, budgetID, unitID, input)
	if err != nil {
		if errors.IsErr(err, errors.BadRequestCode) {
			_ = h.App.WriteErrorResponse(w, errors.BadRequestCode, err)
			return
		}
		if errors.IsErr(err, errors.NotFoundCode) {
			_ = h.App.WriteErrorResponse(w, errors.NotFoundCode, err)
			return
		}
		_ = h.App.WriteErrorResponse(w, errors.InternalCode, err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "SpendingRelease created successfuly", res)
}

func (h *spendingreleaseHandlerImpl) DeleteSpendingRelease(w http.ResponseWriter, r *http.Request) {
	var input dto.DeleteSpendingReleaseInput
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.BadRequestCode, errors.NewBadRequestError("input validation"), validator.Errors)
		return
	}

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.ErrUnauthorized, err)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	err = h.service.DeleteSpendingRelease(ctx, &input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.ErrInternalServerError, err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "SpendingRelease deleted successfuly")
}

func (h *spendingreleaseHandlerImpl) GetSpendingReleaseById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetSpendingRelease(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.InternalCode, err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *spendingreleaseHandlerImpl) GetSpendingReleaseList(w http.ResponseWriter, r *http.Request) {
	var filter data.SpendingReleaseFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.BadRequestCode, errors.NewBadRequestError("input validation"), validator.Errors)
		return
	}

	res, err := h.service.GetSpendingReleaseList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.InternalCode, err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

// GetSpendingReleaseOverview implements SpendingReleaseHandler.
func (h *spendingreleaseHandlerImpl) GetSpendingReleaseOverview(w http.ResponseWriter, r *http.Request) {
	var filter dto.SpendingReleaseOverviewFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.BadRequestCode, errors.NewBadRequestError("input validation"), validator.Errors)
		return
	}

	res, err := h.service.GetSpendingReleaseOverview(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}
