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

type procedurecostPaymentHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.ProcedureCostPaymentService
}

// NewProcedureCostPaymentHandler is a factory function that returns a new instance of ProcedureCostPaymentHandler
func NewProcedureCostPaymentHandler(app *celeritas.Celeritas, procedurecostPaymentService services.ProcedureCostPaymentService) ProcedureCostPaymentHandler {
	return &procedurecostPaymentHandlerImpl{
		App:     app,
		service: procedurecostPaymentService,
	}
}

// CreateProcedureCostPayment creates a new procedurecost payment
func (h *procedurecostPaymentHandlerImpl) CreateProcedureCostPayment(w http.ResponseWriter, r *http.Request) {
	var input dto.ProcedureCostPaymentDTO
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

	res, err := h.service.CreateProcedureCostPayment(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "ProcedureCost payment created successfuly", res)
}

// DeleteProcedureCostPayment deletes a procedurecost payment by its id
func (h *procedurecostPaymentHandlerImpl) DeleteProcedureCostPayment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteProcedureCostPayment(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "ProcedureCost payment deleted successfuly")
}

// UpdateProcedureCostPayment updates a procedurecost payment
func (h *procedurecostPaymentHandlerImpl) UpdateProcedureCostPayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var input dto.ProcedureCostPaymentDTO
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

	res, err := h.service.UpdateProcedureCostPayment(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "ProcedureCost payment updated successfuly", res)
}

// GetProcedureCostPaymentById returns a procedurecost payment by its id
func (h *procedurecostPaymentHandlerImpl) GetProcedureCostPaymentById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetProcedureCostPayment(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

// GetProcedureCostPaymentList returns a list of procedurecost payments
func (h *procedurecostPaymentHandlerImpl) GetProcedureCostPaymentList(w http.ResponseWriter, r *http.Request) {
	var filter dto.ProcedureCostPaymentFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetProcedureCostPaymentList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
