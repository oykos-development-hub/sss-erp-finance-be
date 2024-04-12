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

// DepositAdditionalExpenseHandler is a concrete type that implements DepositAdditionalExpenseHandler
type depositadditionalexpenseHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.DepositAdditionalExpenseService
}

// NewDepositAdditionalExpenseHandler initializes a new DepositAdditionalExpenseHandler with its dependencies
func NewDepositAdditionalExpenseHandler(app *celeritas.Celeritas, depositadditionalexpenseService services.DepositAdditionalExpenseService) DepositAdditionalExpenseHandler {
	return &depositadditionalexpenseHandlerImpl{
		App:     app,
		service: depositadditionalexpenseService,
	}
}

func (h *depositadditionalexpenseHandlerImpl) CreateDepositAdditionalExpense(w http.ResponseWriter, r *http.Request) {
	var input dto.DepositAdditionalExpenseDTO
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

	res, err := h.service.CreateDepositAdditionalExpense(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "DepositAdditionalExpense created successfuly", res)
}

func (h *depositadditionalexpenseHandlerImpl) UpdateDepositAdditionalExpense(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.DepositAdditionalExpenseDTO
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

	res, err := h.service.UpdateDepositAdditionalExpense(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "DepositAdditionalExpense updated successfuly", res)
}

func (h *depositadditionalexpenseHandlerImpl) DeleteDepositAdditionalExpense(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteDepositAdditionalExpense(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "DepositAdditionalExpense deleted successfuly")
}

func (h *depositadditionalexpenseHandlerImpl) GetDepositAdditionalExpenseById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetDepositAdditionalExpense(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *depositadditionalexpenseHandlerImpl) GetDepositAdditionalExpenseList(w http.ResponseWriter, r *http.Request) {
	var filter dto.DepositAdditionalExpenseFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetDepositAdditionalExpenseList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
