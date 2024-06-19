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

// FixedDepositHandler is a concrete type that implements FixedDepositHandler
type fixeddepositHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.FixedDepositService
}

// NewFixedDepositHandler initializes a new FixedDepositHandler with its dependencies
func NewFixedDepositHandler(app *celeritas.Celeritas, fixeddepositService services.FixedDepositService) FixedDepositHandler {
	return &fixeddepositHandlerImpl{
		App:     app,
		service: fixeddepositService,
	}
}

func (h *fixeddepositHandlerImpl) CreateFixedDeposit(w http.ResponseWriter, r *http.Request) {
	var input dto.FixedDepositDTO
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

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	res, err := h.service.CreateFixedDeposit(ctx, input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FixedDeposit created successfuly", res)
}

func (h *fixeddepositHandlerImpl) UpdateFixedDeposit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.FixedDepositDTO
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

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	res, err := h.service.UpdateFixedDeposit(ctx, id, input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FixedDeposit updated successfuly", res)
}

func (h *fixeddepositHandlerImpl) DeleteFixedDeposit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	userIDString := r.Header.Get("UserID")

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(errors.ErrUnauthorized), errors.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	ctx = contextutil.SetUserIDInContext(ctx, userID)

	err = h.service.DeleteFixedDeposit(ctx, id)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "FixedDeposit deleted successfuly")
}

func (h *fixeddepositHandlerImpl) GetFixedDepositById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFixedDeposit(id)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *fixeddepositHandlerImpl) GetFixedDepositList(w http.ResponseWriter, r *http.Request) {
	var filter dto.FixedDepositFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetFixedDepositList(filter)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
