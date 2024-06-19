package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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
func (h *PropBenConfServiceImpl) CreatePropBenConf(ctx context.Context, input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error) {
	propbenconf := input.ToPropBenConf()
	propbenconf.Status = data.UnpaidPropBenConfStatus

	id, err := h.repo.Insert(ctx, *propbenconf)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf insert")
	}

	propbenconf, err = propbenconf.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf get")
	}

	return h.createPropBenConfResponse(ctx, propbenconf)
}

// GetPropBenConf returns a propbenconf by id
func (h *PropBenConfServiceImpl) GetPropBenConf(id int) (*dto.PropBenConfResponseDTO, error) {
	propbenconf, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf get")
	}

	return h.createPropBenConfResponse(context.Background(), propbenconf)
}

// UpdatePropBenConf updates a propbenconf
func (h *PropBenConfServiceImpl) UpdatePropBenConf(ctx context.Context, id int, input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error) {
	propbenconf := input.ToPropBenConf()
	propbenconf.ID = id

	err := h.repo.Update(ctx, *propbenconf)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf update")
	}

	propbenconf, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf get")
	}

	return h.createPropBenConfResponse(ctx, propbenconf)
}

// DeletePropBenConf deletes a propbenconf by its id
func (h *PropBenConfServiceImpl) DeletePropBenConf(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo prop ben conf delete")
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
		return nil, nil, newErrors.Wrap(err, "repo prop ben conf get all")
	}

	var propbenconfsList []data.PropBenConf
	for _, propbenconf := range propbenconfs {
		propbenconfsList = append(propbenconfsList, *propbenconf)
	}

	response, err := h.convertPropBenConfsToResponses(context.Background(), propbenconfsList)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "convert prop ben confs to responses")
	}

	return response, total, nil
}

// convertPropBenConfsToResponses is a helper method that converts a list of propbenconfs to a list of response DTOs.
func (h *PropBenConfServiceImpl) convertPropBenConfsToResponses(ctx context.Context, propbenconfs []data.PropBenConf) ([]dto.PropBenConfResponseDTO, error) {
	var responses []dto.PropBenConfResponseDTO
	for _, fee := range propbenconfs {
		response, err := h.createPropBenConfResponse(ctx, &fee)
		if err != nil {
			return nil, newErrors.Wrap(err, "create prop ben conf response")
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

// createPropBenConfResponse creates a PropBenConfResponseDTO from a PropBenConf
func (h *PropBenConfServiceImpl) createPropBenConfResponse(ctx context.Context, propbenconf *data.PropBenConf) (*dto.PropBenConfResponseDTO, error) {
	response := dto.ToPropBenConfResponseDTO(*propbenconf)
	var newStatus data.PropBenConfStatus
	var err error
	response.PropBenConfDetailsDTO, newStatus, err = h.propbenconfSharedLogicService.CalculatePropBenConfDetailsAndUpdateStatus(ctx, propbenconf.ID)
	if err != nil {
		return nil, newErrors.Wrap(err, "prop ben conf shared logic service calculate prop ben conf details and update status")
	}
	response.Status = newStatus

	return &response, nil
}
