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

// FixedDepositWillDispatchHandler is a concrete type that implements FixedDepositWillDispatchHandler
type fixeddepositwilldispatchHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.FixedDepositWillDispatchService
}

// NewFixedDepositWillDispatchHandler initializes a new FixedDepositWillDispatchHandler with its dependencies
func NewFixedDepositWillDispatchHandler(app *celeritas.Celeritas, fixeddepositwilldispatchService services.FixedDepositWillDispatchService) FixedDepositWillDispatchHandler {
	return &fixeddepositwilldispatchHandlerImpl{
		App:     app,
		service: fixeddepositwilldispatchService,
	}
}

func (h *fixeddepositwilldispatchHandlerImpl) CreateFixedDepositWillDispatch(w http.ResponseWriter, r *http.Request) {
	var input dto.FixedDepositWillDispatchDTO
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

	res, err := h.service.CreateFixedDepositWillDispatch(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FixedDepositWillDispatch created successfuly", res)
}

func (h *fixeddepositwilldispatchHandlerImpl) UpdateFixedDepositWillDispatch(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.FixedDepositWillDispatchDTO
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

	res, err := h.service.UpdateFixedDepositWillDispatch(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FixedDepositWillDispatch updated successfuly", res)
}

func (h *fixeddepositwilldispatchHandlerImpl) DeleteFixedDepositWillDispatch(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteFixedDepositWillDispatch(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "FixedDepositWillDispatch deleted successfuly")
}

func (h *fixeddepositwilldispatchHandlerImpl) GetFixedDepositWillDispatchById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFixedDepositWillDispatch(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *fixeddepositwilldispatchHandlerImpl) GetFixedDepositWillDispatchList(w http.ResponseWriter, r *http.Request) {
	var filter dto.FixedDepositWillDispatchFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetFixedDepositWillDispatchList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
