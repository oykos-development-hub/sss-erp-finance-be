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

// GoalIndicatorHandler is a concrete type that implements GoalIndicatorHandler
type goalindicatorHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.GoalIndicatorService
}

// NewGoalIndicatorHandler initializes a new GoalIndicatorHandler with its dependencies
func NewGoalIndicatorHandler(app *celeritas.Celeritas, goalindicatorService services.GoalIndicatorService) GoalIndicatorHandler {
	return &goalindicatorHandlerImpl{
		App:     app,
		service: goalindicatorService,
	}
}

func (h *goalindicatorHandlerImpl) CreateGoalIndicator(w http.ResponseWriter, r *http.Request) {
	var input dto.GoalIndicatorDTO
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

	res, err := h.service.CreateGoalIndicator(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "GoalIndicator created successfuly", res)
}

func (h *goalindicatorHandlerImpl) UpdateGoalIndicator(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.GoalIndicatorDTO
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

	res, err := h.service.UpdateGoalIndicator(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "GoalIndicator updated successfuly", res)
}

func (h *goalindicatorHandlerImpl) DeleteGoalIndicator(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.service.DeleteGoalIndicator(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "GoalIndicator deleted successfuly")
}

func (h *goalindicatorHandlerImpl) GetGoalIndicatorById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetGoalIndicator(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *goalindicatorHandlerImpl) GetGoalIndicatorList(w http.ResponseWriter, r *http.Request) {
	var filter dto.GoalIndicatorFilterDTO

	_ = h.App.ReadJSON(w, r, &filter)

	validator := h.App.Validator().ValidateStruct(&filter)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, total, err := h.service.GetGoalIndicatorList(filter)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponseWithTotal(w, http.StatusOK, "", res, int(*total))
}
