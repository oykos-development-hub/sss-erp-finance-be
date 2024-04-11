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

// TaxAuthorityCodebookHandler is a concrete type that implements TaxAuthorityCodebookHandler
type taxauthoritycodebookHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.TaxAuthorityCodebookService
}

// NewTaxAuthorityCodebookHandler initializes a new TaxAuthorityCodebookHandler with its dependencies
func NewTaxAuthorityCodebookHandler(app *celeritas.Celeritas, taxauthoritycodebookService services.TaxAuthorityCodebookService) TaxAuthorityCodebookHandler {
	return &taxauthoritycodebookHandlerImpl{
		App:     app,
		service: taxauthoritycodebookService,
	}
}

func (h *taxauthoritycodebookHandlerImpl) CreateTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request) {
	var input dto.TaxAuthorityCodebookDTO
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

	res, err := h.service.CreateTaxAuthorityCodebook(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "TaxAuthorityCodebook created successfuly", res)
}

func (h *taxauthoritycodebookHandlerImpl) UpdateTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.TaxAuthorityCodebookDTO
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

	res, err := h.service.UpdateTaxAuthorityCodebook(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "TaxAuthorityCodebook updated successfuly", res)
}

func (h *taxauthoritycodebookHandlerImpl) DeactivateTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.TaxAuthorityCodebookDTO
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

	err = h.service.DeactivateTaxAuthorityCodebook(id, input.Active)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "TaxAuthorityCodebook updated successfuly", nil)
}

func (h *taxauthoritycodebookHandlerImpl) DeleteTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteTaxAuthorityCodebook(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "TaxAuthorityCodebook deleted successfuly")
}

func (h *taxauthoritycodebookHandlerImpl) GetTaxAuthorityCodebookById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetTaxAuthorityCodebook(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *taxauthoritycodebookHandlerImpl) GetTaxAuthorityCodebookList(w http.ResponseWriter, r *http.Request) {
	var filter dto.TaxAuthorityCodebookFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetTaxAuthorityCodebookList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
