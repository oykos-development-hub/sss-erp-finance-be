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

type FineServiceImpl struct {
	App                    *celeritas.Celeritas
	repo                   data.Fine
	fineSharedLogicService FineSharedLogicService
}

// NewFineServiceImpl creates a new instance of FineService
func NewFineServiceImpl(app *celeritas.Celeritas, repo data.Fine, fineSharedLogicService FineSharedLogicService) FineService {
	return &FineServiceImpl{
		App:                    app,
		repo:                   repo,
		fineSharedLogicService: fineSharedLogicService,
	}
}

// CreateFine creates a new fine
func (h *FineServiceImpl) CreateFine(ctx context.Context, input dto.FineDTO) (*dto.FineResponseDTO, error) {
	fine := input.ToFine()
	fine.Status = data.UnpaidFineStatus

	id, err := h.repo.Insert(ctx, *fine)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine insert")
	}

	fine, err = fine.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine get")
	}

	return h.createFineResponse(ctx, fine)
}

// GetFine returns a fine by id
func (h *FineServiceImpl) GetFine(id int) (*dto.FineResponseDTO, error) {
	fine, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine get")
	}

	return h.createFineResponse(context.Background(), fine)
}

// UpdateFine updates a fine
func (h *FineServiceImpl) UpdateFine(ctx context.Context, id int, input dto.FineDTO) (*dto.FineResponseDTO, error) {
	fine := input.ToFine()
	fine.ID = id

	err := h.repo.Update(ctx, *fine)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine update")
	}

	fine, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine get")
	}

	return h.createFineResponse(ctx, fine)
}

// DeleteFine deletes a fine by its id
func (h *FineServiceImpl) DeleteFine(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo fine delete")
	}

	return nil
}

// GetFineList returns a list of fines
func (h *FineServiceImpl) GetFineList(input dto.FineFilterDTO) ([]dto.FineResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}

	if input.Subject != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Subject)
		subject := up.Or(
			up.Cond{"subject ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, subject)
	}

	if input.FilterByActTypeID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"act_type": *input.FilterByActTypeID})
	}

	if input.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *input.OrganizationUnitID})
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

	fines, total, err := h.repo.GetAll(input.Page, input.Size, conditionAndExp)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo fine get all")
	}

	var finesList []data.Fine
	for _, fine := range fines {
		finesList = append(finesList, *fine)
	}

	response, err := h.convertFinesToResponses(finesList)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "convert fines to response")
	}

	return response, total, nil
}

// convertFinesToResponses is a helper method that converts a list of fines to a list of response DTOs.
func (h *FineServiceImpl) convertFinesToResponses(fines []data.Fine) ([]dto.FineResponseDTO, error) {
	var responses []dto.FineResponseDTO
	for _, fee := range fines {
		response, err := h.createFineResponse(context.Background(), &fee)
		if err != nil {
			return nil, newErrors.Wrap(err, "create fine response")
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

// createFineResponse creates a FineResponseDTO from a Fine
func (h *FineServiceImpl) createFineResponse(ctx context.Context, fine *data.Fine) (*dto.FineResponseDTO, error) {
	response := dto.ToFineResponseDTO(*fine)
	var newStatus data.FineStatus
	var err error
	response.FineFeeDetailsDTO, newStatus, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(ctx, fine.ID)
	if err != nil {
		return nil, newErrors.Wrap(err, "fine shared logic calculate fine details and update status")
	}
	response.Status = newStatus

	return &response, nil
}
