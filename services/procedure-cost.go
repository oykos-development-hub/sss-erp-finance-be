package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ProcedureCostServiceImpl struct {
	App                             *celeritas.Celeritas
	repo                            data.ProcedureCost
	procedurecostSharedLogicService ProcedureCostSharedLogicService
}

// NewProcedureCostServiceImpl creates a new instance of ProcedureCostService
func NewProcedureCostServiceImpl(app *celeritas.Celeritas, repo data.ProcedureCost, procedurecostSharedLogicService ProcedureCostSharedLogicService) ProcedureCostService {
	return &ProcedureCostServiceImpl{
		App:                             app,
		repo:                            repo,
		procedurecostSharedLogicService: procedurecostSharedLogicService,
	}
}

// CreateProcedureCost creates a new procedurecost
func (h *ProcedureCostServiceImpl) CreateProcedureCost(ctx context.Context, input dto.ProcedureCostDTO) (*dto.ProcedureCostResponseDTO, error) {
	procedurecost := input.ToProcedureCost()
	procedurecost.Status = data.UnpaidProcedureCostStatus

	id, err := h.repo.Insert(ctx, *procedurecost)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	procedurecost, err = procedurecost.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return h.createProcedureCostResponse(context.Background(), procedurecost)
}

// GetProcedureCost returns a procedurecost by id
func (h *ProcedureCostServiceImpl) GetProcedureCost(id int) (*dto.ProcedureCostResponseDTO, error) {
	procedurecost, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}

	return h.createProcedureCostResponse(context.Background(), procedurecost)
}

// UpdateProcedureCost updates a procedurecost
func (h *ProcedureCostServiceImpl) UpdateProcedureCost(ctx context.Context, id int, input dto.ProcedureCostDTO) (*dto.ProcedureCostResponseDTO, error) {
	procedurecost := input.ToProcedureCost()
	procedurecost.ID = id

	err := h.repo.Update(ctx, *procedurecost)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	procedurecost, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return h.createProcedureCostResponse(ctx, procedurecost)
}

// DeleteProcedureCost deletes a procedurecost by its id
func (h *ProcedureCostServiceImpl) DeleteProcedureCost(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

// GetProcedureCostList returns a list of procedurecosts
func (h *ProcedureCostServiceImpl) GetProcedureCostList(input dto.ProcedureCostFilterDTO) ([]dto.ProcedureCostResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}

	if input.Subject != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Subject)
		subject := up.Or(
			up.Cond{"subject ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, subject)
	}

	if input.FilterByProcedureCostTypeID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"procedure_cost_type": *input.FilterByProcedureCostTypeID})
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

	procedureCosts, total, err := h.repo.GetAll(input.Page, input.Size, conditionAndExp)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}

	var procedureCostsList []data.ProcedureCost
	for _, procedureCost := range procedureCosts {
		procedureCostsList = append(procedureCostsList, *procedureCost)
	}

	response, err := h.convertProcedureCostsToResponses(context.Background(), procedureCostsList)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}

	return response, total, nil
}

// convertProcedureCostsToResponses is a helper method that converts a list of procedurecosts to a list of response DTOs.
func (h *ProcedureCostServiceImpl) convertProcedureCostsToResponses(ctx context.Context, procedurecosts []data.ProcedureCost) ([]dto.ProcedureCostResponseDTO, error) {
	var responses []dto.ProcedureCostResponseDTO
	for _, fee := range procedurecosts {
		response, err := h.createProcedureCostResponse(ctx, &fee)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

// createProcedureCostResponse creates a ProcedureCostResponseDTO from a ProcedureCost
func (h *ProcedureCostServiceImpl) createProcedureCostResponse(ctx context.Context, procedurecost *data.ProcedureCost) (*dto.ProcedureCostResponseDTO, error) {
	response := dto.ToProcedureCostResponseDTO(*procedurecost)
	var newStatus data.ProcedureCostStatus
	var err error
	response.ProcedureCostDetails, newStatus, err = h.procedurecostSharedLogicService.CalculateProcedureCostDetailsAndUpdateStatus(ctx, procedurecost.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}
	response.Status = newStatus

	return &response, nil
}
