package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type PropBenConfServiceImpl struct {
	App                           *celeritas.Celeritas
	repo                          data.PropBenConf
	propbenconfSharedLogicService PropBenConfSharedLogicService
}

// NewPropBenConfServiceImpl creates a new instance of PropBenConfService
func NewPropBenConfServiceImpl(app *celeritas.Celeritas, repo data.PropBenConf, propbenconfSharedLogicService PropBenConfSharedLogicService) PropBenConfService {
	return &PropBenConfServiceImpl{
		App:                           app,
		repo:                          repo,
		propbenconfSharedLogicService: propbenconfSharedLogicService,
	}
}

// CreatePropBenConf creates a new propbenconf
func (h *PropBenConfServiceImpl) CreatePropBenConf(input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error) {
	propbenconf := input.ToPropBenConf()
	propbenconf.Status = data.UnpaidPropBenConfStatus

	id, err := h.repo.Insert(*propbenconf)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	propbenconf, err = propbenconf.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return h.createPropBenConfResponse(propbenconf)
}

// GetPropBenConf returns a propbenconf by id
func (h *PropBenConfServiceImpl) GetPropBenConf(id int) (*dto.PropBenConfResponseDTO, error) {
	propbenconf, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}

	return h.createPropBenConfResponse(propbenconf)
}

// UpdatePropBenConf updates a propbenconf
func (h *PropBenConfServiceImpl) UpdatePropBenConf(id int, input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error) {
	propbenconf := input.ToPropBenConf()
	propbenconf.ID = id

	err := h.repo.Update(*propbenconf)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	propbenconf, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return h.createPropBenConfResponse(propbenconf)
}

// DeletePropBenConf deletes a propbenconf by its id
func (h *PropBenConfServiceImpl) DeletePropBenConf(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

// GetPropBenConfList returns a list of propbenconfs
func (h *PropBenConfServiceImpl) GetPropBenConfList(input dto.PropBenConfFilterDTO) ([]dto.PropBenConfResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}

	if input.Subject != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Subject)
		subject := up.Or(
			up.Cond{"subject ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, subject)
	}

	if input.FilterByPropBenConfTypeID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"property_benefits_confiscation_type": *input.FilterByPropBenConfTypeID})
	}

	// combine search by subject, jmbg and description with filter by decision number
	if input.Search != nil && *input.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Search)
		stringConditions := up.Or(
			up.Cond{"subject ILIKE": likeCondition},
			up.Cond{"description ILIKE": likeCondition},
			up.Cond{"jmbg ILIKE": likeCondition},
			up.Cond{"decision_number ILIKE": likeCondition},
		)

		conditionAndExp = up.And(conditionAndExp, stringConditions)
	}

	propbenconfs, total, err := h.repo.GetAll(input.Page, input.Size, conditionAndExp)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}

	var propbenconfsList []data.PropBenConf
	for _, propbenconf := range propbenconfs {
		propbenconfsList = append(propbenconfsList, *propbenconf)
	}

	response, err := h.convertPropBenConfsToResponses(propbenconfsList)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}

	return response, total, nil
}

// convertPropBenConfsToResponses is a helper method that converts a list of propbenconfs to a list of response DTOs.
func (h *PropBenConfServiceImpl) convertPropBenConfsToResponses(propbenconfs []data.PropBenConf) ([]dto.PropBenConfResponseDTO, error) {
	var responses []dto.PropBenConfResponseDTO
	for _, fee := range propbenconfs {
		response, err := h.createPropBenConfResponse(&fee)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

// createPropBenConfResponse creates a PropBenConfResponseDTO from a PropBenConf
func (h *PropBenConfServiceImpl) createPropBenConfResponse(propbenconf *data.PropBenConf) (*dto.PropBenConfResponseDTO, error) {
	response := dto.ToPropBenConfResponseDTO(*propbenconf)
	var newStatus data.PropBenConfStatus
	var err error
	response.PropBenConfDetailsDTO, newStatus, err = h.propbenconfSharedLogicService.CalculatePropBenConfDetailsAndUpdateStatus(propbenconf.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}
	response.Status = newStatus

	return &response, nil
}
