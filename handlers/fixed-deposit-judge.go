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

// FixedDepositJudgeHandler is a concrete type that implements FixedDepositJudgeHandler
type fixeddepositjudgeHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.FixedDepositJudgeService
}

// NewFixedDepositJudgeHandler initializes a new FixedDepositJudgeHandler with its dependencies
func NewFixedDepositJudgeHandler(app *celeritas.Celeritas, fixeddepositjudgeService services.FixedDepositJudgeService) FixedDepositJudgeHandler {
	return &fixeddepositjudgeHandlerImpl{
		App:     app,
		service: fixeddepositjudgeService,
	}
}

func (h *fixeddepositjudgeHandlerImpl) CreateFixedDepositJudge(w http.ResponseWriter, r *http.Request) {
	var input dto.FixedDepositJudgeDTO
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

	res, err := h.service.CreateFixedDepositJudge(input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FixedDepositJudge created successfuly", res)
}

func (h *fixeddepositjudgeHandlerImpl) UpdateFixedDepositJudge(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.FixedDepositJudgeDTO
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

	res, err := h.service.UpdateFixedDepositJudge(id, input)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "FixedDepositJudge updated successfuly", res)
}

func (h *fixeddepositjudgeHandlerImpl) DeleteFixedDepositJudge(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteFixedDepositJudge(id)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "FixedDepositJudge deleted successfuly")
}

func (h *fixeddepositjudgeHandlerImpl) GetFixedDepositJudgeById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFixedDepositJudge(id)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *fixeddepositjudgeHandlerImpl) GetFixedDepositJudgeList(w http.ResponseWriter, r *http.Request) {
	var filter dto.FixedDepositJudgeFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetFixedDepositJudgeList(filter)
	if err != nil {
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
