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

// SpendingReleaseRequestHandler is a concrete type that implements SpendingReleaseRequestHandler
type spendingreleaserequestHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.SpendingReleaseRequestService
}

// NewSpendingReleaseRequestHandler initializes a new SpendingReleaseRequestHandler with its dependencies
func NewSpendingReleaseRequestHandler(app *celeritas.Celeritas, spendingreleaserequestService services.SpendingReleaseRequestService) SpendingReleaseRequestHandler {
	return &spendingreleaserequestHandlerImpl{
		App:     app,
		service: spendingreleaserequestService,
	}
}

func (h *spendingreleaserequestHandlerImpl) CreateSpendingReleaseRequest(w http.ResponseWriter, r *http.Request) {
	var input dto.SpendingReleaseRequestDTO
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

	res, err := h.service.CreateSpendingReleaseRequest(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "SpendingReleaseRequest created successfuly", res)
}

func (h *spendingreleaserequestHandlerImpl) UpdateSpendingReleaseRequest(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.SpendingReleaseRequestDTO
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

	res, err := h.service.UpdateSpendingReleaseRequest(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "SpendingReleaseRequest updated successfuly", res)
}

func (h *spendingreleaserequestHandlerImpl) DeleteSpendingReleaseRequest(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteSpendingReleaseRequest(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "SpendingReleaseRequest deleted successfuly")
}

func (h *spendingreleaserequestHandlerImpl) GetSpendingReleaseRequestById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetSpendingReleaseRequest(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *spendingreleaserequestHandlerImpl) GetSpendingReleaseRequestList(w http.ResponseWriter, r *http.Request) {
	var filter dto.SpendingReleaseRequestFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetSpendingReleaseRequestList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}

func (h *spendingreleaserequestHandlerImpl) AcceptSSSRequest(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.SpendingReleaseRequestDTO
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

	err = h.service.AcceptSSSRequest(r.Context(), id, input.SSSFileID)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "SpendingReleaseRequest updated successfuly", nil)
}
