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

type feePaymentHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.FeePaymentService
}

// NewFeePaymentHandler is a factory function that returns a new instance of FeePaymentHandler
func NewFeePaymentHandler(app *celeritas.Celeritas, feePaymentService services.FeePaymentService) FeePaymentHandler {
	return &feePaymentHandlerImpl{
		App:     app,
		service: feePaymentService,
	}
}

// CreateFeePayment creates a new fee payment
func (h *feePaymentHandlerImpl) CreateFeePayment(w http.ResponseWriter, r *http.Request) {
	var input dto.FeePaymentDTO
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

	res, err := h.service.CreateFeePayment(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "Fee payment created successfuly", res)
}

// DeleteFeePayment deletes a fee payment by its id
func (h *feePaymentHandlerImpl) DeleteFeePayment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteFeePayment(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "Fee payment deleted successfuly")
}

// UpdateFeePayment updates a fee payment
func (h *feePaymentHandlerImpl) UpdateFeePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var input dto.FeePaymentDTO
	err = h.App.ReadJSON(w, r, &input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, err := h.service.UpdateFeePayment(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "Fee payment updated successfuly", res)
}

// GetFeePaymentById returns a fee payment by its id
func (h *feePaymentHandlerImpl) GetFeePaymentById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFeePayment(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

// GetFeePaymentList returns a list of fee payments
func (h *feePaymentHandlerImpl) GetFeePaymentList(w http.ResponseWriter, r *http.Request) {
	var filter dto.FeePaymentFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetFeePaymentList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}