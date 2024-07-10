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

// InternalReallocationItemHandler is a concrete type that implements InternalReallocationItemHandler
type internalreallocationitemHandlerImpl struct {
	App             *celeritas.Celeritas
	service         services.InternalReallocationItemService
	errorLogService services.ErrorLogService
}

// NewInternalReallocationItemHandler initializes a new InternalReallocationItemHandler with its dependencies
func NewInternalReallocationItemHandler(app *celeritas.Celeritas, internalreallocationitemService services.InternalReallocationItemService, errorLogService services.ErrorLogService) InternalReallocationItemHandler {
	return &internalreallocationitemHandlerImpl{
		App:             app,
		service:         internalreallocationitemService,
		errorLogService: errorLogService,
	}
}

func (h *internalreallocationitemHandlerImpl) CreateInternalReallocationItem(w http.ResponseWriter, r *http.Request) {
	var input dto.InternalReallocationItemDTO
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

	res, err := h.service.CreateInternalReallocationItem(input)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "InternalReallocationItem created successfuly", res)
}

func (h *internalreallocationitemHandlerImpl) DeleteInternalReallocationItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteInternalReallocationItem(id)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "InternalReallocationItem deleted successfuly")
}

func (h *internalreallocationitemHandlerImpl) GetInternalReallocationItemById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetInternalReallocationItem(id)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *internalreallocationitemHandlerImpl) GetInternalReallocationItemList(w http.ResponseWriter, r *http.Request) {
	var filter dto.InternalReallocationItemFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		h.App.ErrorLog.Print(validator.Errors)
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetInternalReallocationItemList(filter)
	if err != nil {
		h.errorLogService.CreateErrorLog(err)
		h.App.ErrorLog.Print(err)
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
