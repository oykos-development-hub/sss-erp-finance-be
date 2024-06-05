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

// ExternalReallocationHandler is a concrete type that implements ExternalReallocationHandler
type externalreallocationHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.ExternalReallocationService
}

// NewExternalReallocationHandler initializes a new ExternalReallocationHandler with its dependencies
func NewExternalReallocationHandler(app *celeritas.Celeritas, externalreallocationService services.ExternalReallocationService) ExternalReallocationHandler {
	return &externalreallocationHandlerImpl{
		App:     app,
		service: externalreallocationService,
	}
}

func (h *externalreallocationHandlerImpl) CreateExternalReallocation(w http.ResponseWriter, r *http.Request) {
	var input dto.ExternalReallocationDTO
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

	res, err := h.service.CreateExternalReallocation(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "ExternalReallocation created successfuly", res)
}

func (h *externalreallocationHandlerImpl) DeleteExternalReallocation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteExternalReallocation(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "ExternalReallocation deleted successfuly")
}

func (h *externalreallocationHandlerImpl) GetExternalReallocationById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetExternalReallocation(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *externalreallocationHandlerImpl) GetExternalReallocationList(w http.ResponseWriter, r *http.Request) {
	var filter dto.ExternalReallocationFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetExternalReallocationList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}

func (h *externalreallocationHandlerImpl) AcceptOUExternalReallocation(w http.ResponseWriter, r *http.Request) {
	var input dto.ExternalReallocationDTO
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

	res, err := h.service.AcceptOUExternalReallocation(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "ExternalReallocation accepted successfuly", res)
}

func (h *externalreallocationHandlerImpl) RejectOUExternalReallocation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.RejectOUExternalReallocation(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "ExternalReallocation rejected successfuly")
}

func (h *externalreallocationHandlerImpl) AcceptSSSExternalReallocation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.AcceptSSSExternalReallocation(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "ExternalReallocation accepted successfuly")
}

func (h *externalreallocationHandlerImpl) RejectSSSExternalReallocation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.RejectSSSExternalReallocation(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "ExternalReallocation rejected successfuly")
}
