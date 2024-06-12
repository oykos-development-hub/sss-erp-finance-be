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

type FeeServiceImpl struct {
	App                   *celeritas.Celeritas
	repo                  data.Fee
	feeSharedLogicService FeeSharedLogicService
}

// NewFeeServiceImpl creates a new instance of FeeService
func NewFeeServiceImpl(app *celeritas.Celeritas, repo data.Fee, feeSharedLogicService FeeSharedLogicService) FeeService {
	return &FeeServiceImpl{
		App:                   app,
		repo:                  repo,
		feeSharedLogicService: feeSharedLogicService,
	}
}

// CreateFee creates a new fee
func (h *FeeServiceImpl) CreateFee(ctx context.Context, input dto.FeeDTO) (*dto.FeeResponseDTO, error) {
	fee := input.ToFee()
	fee.Status = data.UnpaidFeeStatus

	id, err := h.repo.Insert(ctx, *fee)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	fee, err = fee.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return h.convertFeeToResponse(fee)
}

// GetFee returns a fee by id
func (h *FeeServiceImpl) GetFee(id int) (*dto.FeeResponseDTO, error) {
	fee, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}

	return h.convertFeeToResponse(fee)
}

// UpdateFee updates a fee
func (h *FeeServiceImpl) UpdateFee(ctx context.Context, id int, input dto.FeeDTO) (*dto.FeeResponseDTO, error) {
	fee := input.ToFee()
	fee.ID = id

	err := h.repo.Update(ctx, *fee)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	fee, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return h.convertFeeToResponse(fee)
}

// DeleteFee deletes a fee by its id
func (h *FeeServiceImpl) DeleteFee(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

// GetFeeList returns a list of fees
func (h *FeeServiceImpl) GetFeeList(input dto.FeeFilterDTO) ([]dto.FeeResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}

	if input.FilterByFeeSubcategoryID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"fee_subcategory_id": *input.FilterByFeeSubcategoryID})
	}

	if input.FilterByFeeTypeID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"fee_type_id": *input.FilterByFeeTypeID})
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

	fees, total, err := h.repo.GetAll(input.Page, input.Size, conditionAndExp)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}

	var feeList []data.Fee
	for _, fee := range fees {
		feeList = append(feeList, *fee)
	}

	response, err := h.convertFeesToResponses(feeList)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}

	return response, total, nil
}

// convertFeesToResponses is a helper method that converts a list of fees to a list of response DTOs.
func (h *FeeServiceImpl) convertFeesToResponses(fees []data.Fee) ([]dto.FeeResponseDTO, error) {
	var responses []dto.FeeResponseDTO
	for _, fee := range fees {
		response, err := h.convertFeeToResponse(&fee)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

// convertFeeToResponse is a helper method that converts a fee to a response DTO
func (h *FeeServiceImpl) convertFeeToResponse(fee *data.Fee) (*dto.FeeResponseDTO, error) {
	response := dto.ToFeeResponseDTO(*fee)
	var newStatus data.FeeStatus
	var err error
	response.FeeDetails, newStatus, err = h.feeSharedLogicService.CalculateFeeDetailsAndUpdateStatus(context.Background(), fee.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}
	response.Status = newStatus
	return &response, nil
}
