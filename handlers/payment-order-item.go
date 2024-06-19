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

// PaymentOrderItemHandler is a concrete type that implements PaymentOrderItemHandler
type paymentorderitemHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.PaymentOrderItemService
}

// NewPaymentOrderItemHandler initializes a new PaymentOrderItemHandler with its dependencies
func NewPaymentOrderItemHandler(app *celeritas.Celeritas, paymentorderitemService services.PaymentOrderItemService) PaymentOrderItemHandler {
	return &paymentorderitemHandlerImpl{
		App:     app,
		service: paymentorderitemService,
	}
}

func (h *paymentorderitemHandlerImpl) CreatePaymentOrderItem(w http.ResponseWriter, r *http.Request) {
	var input dto.PaymentOrderItemDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, err := h.service.CreatePaymentOrderItem(input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "PaymentOrderItem created successfuly", res)
}

func (h *paymentorderitemHandlerImpl) UpdatePaymentOrderItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.PaymentOrderItemDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, err := h.service.UpdatePaymentOrderItem(id, input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "PaymentOrderItem updated successfuly", res)
}

func (h *paymentorderitemHandlerImpl) DeletePaymentOrderItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeletePaymentOrderItem(id)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "PaymentOrderItem deleted successfuly")
}

func (h *paymentorderitemHandlerImpl) GetPaymentOrderItemById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetPaymentOrderItem(id)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *paymentorderitemHandlerImpl) GetPaymentOrderItemList(w http.ResponseWriter, r *http.Request) {
	var filter dto.PaymentOrderItemFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetPaymentOrderItemList(filter)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
