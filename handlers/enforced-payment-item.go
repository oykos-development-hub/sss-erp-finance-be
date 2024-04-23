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

// EnforcedPaymentItemHandler is a concrete type that implements EnforcedPaymentItemHandler
type enforcedpaymentitemHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.EnforcedPaymentItemService
}

// NewEnforcedPaymentItemHandler initializes a new EnforcedPaymentItemHandler with its dependencies
func NewEnforcedPaymentItemHandler(app *celeritas.Celeritas, enforcedpaymentitemService services.EnforcedPaymentItemService) EnforcedPaymentItemHandler {
	return &enforcedpaymentitemHandlerImpl{
		App:     app,
		service: enforcedpaymentitemService,
	}
}

func (h *enforcedpaymentitemHandlerImpl) CreateEnforcedPaymentItem(w http.ResponseWriter, r *http.Request) {
	var input dto.EnforcedPaymentItemDTO
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

	res, err := h.service.CreateEnforcedPaymentItem(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "EnforcedPaymentItem created successfuly", res)
}

func (h *enforcedpaymentitemHandlerImpl) UpdateEnforcedPaymentItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.EnforcedPaymentItemDTO
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

	res, err := h.service.UpdateEnforcedPaymentItem(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "EnforcedPaymentItem updated successfuly", res)
}

func (h *enforcedpaymentitemHandlerImpl) DeleteEnforcedPaymentItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteEnforcedPaymentItem(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "EnforcedPaymentItem deleted successfuly")
}

func (h *enforcedpaymentitemHandlerImpl) GetEnforcedPaymentItemById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetEnforcedPaymentItem(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *enforcedpaymentitemHandlerImpl) GetEnforcedPaymentItemList(w http.ResponseWriter, r *http.Request) {
	var filter dto.EnforcedPaymentItemFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetEnforcedPaymentItemList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
