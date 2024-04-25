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

// ModelsOfAccountingHandler is a concrete type that implements ModelsOfAccountingHandler
type modelsofaccountingHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.ModelsOfAccountingService
}

// NewModelsOfAccountingHandler initializes a new ModelsOfAccountingHandler with its dependencies
func NewModelsOfAccountingHandler(app *celeritas.Celeritas, modelsofaccountingService services.ModelsOfAccountingService) ModelsOfAccountingHandler {
	return &modelsofaccountingHandlerImpl{
		App:     app,
		service: modelsofaccountingService,
	}
}

func (h *modelsofaccountingHandlerImpl) CreateModelsOfAccounting(w http.ResponseWriter, r *http.Request) {
	var input dto.ModelsOfAccountingDTO
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

	res, err := h.service.CreateModelsOfAccounting(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "ModelsOfAccounting created successfuly", res)
}

func (h *modelsofaccountingHandlerImpl) GetModelsOfAccountingById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetModelsOfAccounting(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *modelsofaccountingHandlerImpl) GetModelsOfAccountingList(w http.ResponseWriter, r *http.Request) {
	var filter dto.ModelsOfAccountingFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetModelsOfAccountingList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
