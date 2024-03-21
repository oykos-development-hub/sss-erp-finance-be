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

type propbenconfHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.PropBenConfService
}

// NewPropBenConfHandler is a factory function that returns a new instance of PropBenConfHandler
func NewPropBenConfHandler(app *celeritas.Celeritas, propbenconfService services.PropBenConfService) PropBenConfHandler {
	return &propbenconfHandlerImpl{
		App:     app,
		service: propbenconfService,
	}
}

// CreatePropBenConf creates a new propbenconf
func (h *propbenconfHandlerImpl) CreatePropBenConf(w http.ResponseWriter, r *http.Request) {
	var input dto.PropBenConfDTO
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

	res, err := h.service.CreatePropBenConf(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "Property Benefits Confiscation created successfuly", res)
}

// GetPropBenConfById returns a propbenconf by its id
func (h *propbenconfHandlerImpl) GetPropBenConfById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetPropBenConf(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

// UpdatePropBenConf updates a propbenconf
func (h *propbenconfHandlerImpl) UpdatePropBenConf(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.PropBenConfDTO
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

	res, err := h.service.UpdatePropBenConf(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "Property Benefits Confiscation updated successfuly", res)
}

// DeletePropBenConf deletes a propbenconf
func (h *propbenconfHandlerImpl) DeletePropBenConf(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeletePropBenConf(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "Property Benefits Confiscation deleted successfuly")
}

// GetPropBenConfList returns a list of propbenconfs
func (h *propbenconfHandlerImpl) GetPropBenConfList(w http.ResponseWriter, r *http.Request) {
	var filter dto.PropBenConfFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetPropBenConfList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
