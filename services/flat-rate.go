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

type FlatRateServiceImpl struct {
	App                        *celeritas.Celeritas
	repo                       data.FlatRate
	FlatRateSharedLogicService FlatRateSharedLogicService
}

// NewFlatRateServiceImpl creates a new instance of FlatRateService
func NewFlatRateServiceImpl(app *celeritas.Celeritas, repo data.FlatRate, FlatRateSharedLogicService FlatRateSharedLogicService) FlatRateService {
	return &FlatRateServiceImpl{
		App:                        app,
		repo:                       repo,
		FlatRateSharedLogicService: FlatRateSharedLogicService,
	}
}

// CreateFlatRate creates a new FlatRate
func (h *FlatRateServiceImpl) CreateFlatRate(ctx context.Context, input dto.FlatRateDTO) (*dto.FlatRateResponseDTO, error) {
	FlatRate := input.ToFlatRate()
	FlatRate.Status = data.UnpaidFlatRateStatus

	id, err := h.repo.Insert(ctx, *FlatRate)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo flat rate insert")
	}

	FlatRate, err = FlatRate.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo flat rate insert")
	}

	return h.createFlatRateResponse(ctx, FlatRate)
}

// GetFlatRate returns a FlatRate by id
func (h *FlatRateServiceImpl) GetFlatRate(id int) (*dto.FlatRateResponseDTO, error) {
	FlatRate, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo flat rate get")
	}

	return h.createFlatRateResponse(context.Background(), FlatRate)
}

// UpdateFlatRate updates a FlatRate
func (h *FlatRateServiceImpl) UpdateFlatRate(ctx context.Context, id int, input dto.FlatRateDTO) (*dto.FlatRateResponseDTO, error) {
	FlatRate := input.ToFlatRate()
	FlatRate.ID = id

	err := h.repo.Update(ctx, *FlatRate)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo flat rate update")
	}

	FlatRate, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo flat rate get")
	}

	return h.createFlatRateResponse(ctx, FlatRate)
}

// DeleteFlatRate deletes a FlatRate by its id
func (h *FlatRateServiceImpl) DeleteFlatRate(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo flat rate delete")
	}

	return nil
}

// GetFlatRateList returns a list of FlatRates
func (h *FlatRateServiceImpl) GetFlatRateList(input dto.FlatRateFilterDTO) ([]dto.FlatRateResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}

	if input.Subject != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Subject)
		subject := up.Or(
			up.Cond{"subject ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, subject)
	}

	if input.FilterByTypeID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"flat_rate_type": *input.FilterByTypeID})
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

	FlatRates, total, err := h.repo.GetAll(input.Page, input.Size, conditionAndExp)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo flat rate get all")
	}

	var FlatRatesList []data.FlatRate
	for _, FlatRate := range FlatRates {
		FlatRatesList = append(FlatRatesList, *FlatRate)
	}

	response, err := h.convertFlatRatesToResponses(FlatRatesList)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "convert flat rates to responses")
	}

	return response, total, nil
}

// convertFlatRatesToResponses is a helper method that converts a list of FlatRates to a list of response DTOs.
func (h *FlatRateServiceImpl) convertFlatRatesToResponses(FlatRates []data.FlatRate) ([]dto.FlatRateResponseDTO, error) {
	var responses []dto.FlatRateResponseDTO
	for _, fee := range FlatRates {
		response, err := h.createFlatRateResponse(context.Background(), &fee)
		if err != nil {
			return nil, newErrors.Wrap(err, "create flat rate response")
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

// createFlatRateResponse creates a FlatRateResponseDTO from a FlatRate
func (h *FlatRateServiceImpl) createFlatRateResponse(ctx context.Context, FlatRate *data.FlatRate) (*dto.FlatRateResponseDTO, error) {
	response := dto.ToFlatRateResponseDTO(*FlatRate)
	var newStatus data.FlatRateStatus
	var err error
	response.FlatRateDetails, newStatus, err = h.FlatRateSharedLogicService.CalculateFlatRateDetailsAndUpdateStatus(ctx, FlatRate.ID)
	if err != nil {
		return nil, newErrors.Wrap(err, "flat rate shared logic service calcualte flat rate details and update status")
	}
	response.Status = newStatus

	return &response, nil
}
