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

// ExternalReallocationItemHandler is a concrete type that implements ExternalReallocationItemHandler
type externalreallocationitemHandlerImpl struct {
	App             *celeritas.Celeritas
	service         services.ExternalReallocationItemService
	errorLogService services.ErrorLogService
}

// NewExternalReallocationItemHandler initializes a new ExternalReallocationItemHandler with its dependencies
func NewExternalReallocationItemHandler(app *celeritas.Celeritas, externalreallocationitemService services.ExternalReallocationItemService, errorLogService services.ErrorLogService) ExternalReallocationItemHandler {
	return &externalreallocationitemHandlerImpl{
		App:             app,
		service:         externalreallocationitemService,
		errorLogService: errorLogService,
	}
}

func (h *externalreallocationitemHandlerImpl) CreateExternalReallocationItem(w http.ResponseWriter, r *http.Request) {
	var input dto.ExternalReallocationItemDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
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

	res, err := h.service.CreateExternalReallocationItem(input)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "ExternalReallocationItem created successfuly", res)
}

func (h *externalreallocationitemHandlerImpl) DeleteExternalReallocationItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteExternalReallocationItem(id)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "ExternalReallocationItem deleted successfuly")
}

func (h *externalreallocationitemHandlerImpl) GetExternalReallocationItemById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetExternalReallocationItem(id)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *externalreallocationitemHandlerImpl) GetExternalReallocationItemList(w http.ResponseWriter, r *http.Request) {
	var filter dto.ExternalReallocationItemFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetExternalReallocationItemList(filter)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
