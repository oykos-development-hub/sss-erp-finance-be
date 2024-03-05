package services

import (
	"fmt"
	"strconv"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

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
func (h *FineServiceImpl) CreateFine(input dto.FineDTO) (*dto.FineResponseDTO, error) {
	fine := input.ToFine()
	fine.Status = data.UnpaidFineStatus

	id, err := h.repo.Insert(*fine)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	fine, err = fine.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFineResponseDTO(*fine)

	return &res, nil
}

// GetFine returns a fine by id
func (h *FineServiceImpl) GetFine(id int) (*dto.FineResponseDTO, error) {
	fine, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}

	response := dto.ToFineResponseDTO(*fine)
	var newStatus data.FineStatus
	response.FineFeeDetailsDTO, newStatus, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(fine.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}
	response.Status = newStatus

	return &response, nil
}

// UpdateFine updates a fine
func (h *FineServiceImpl) UpdateFine(id int, input dto.FineDTO) (*dto.FineResponseDTO, error) {
	fine := input.ToFine()
	fine.ID = id

	err := h.repo.Update(*fine)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	fine, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFineResponseDTO(*fine)

	var newStatus data.FineStatus
	response.FineFeeDetailsDTO, newStatus, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(fine.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}
	response.Status = newStatus

	return &response, nil
}

// DeleteFine deletes a fine by its id
func (h *FineServiceImpl) DeleteFine(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
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

	// combine search by subject, jmbg and description with filter by decision number
	if input.Search != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Search)
		stringConditions := up.Or(
			up.Cond{"subject ILIKE": likeCondition},
			up.Cond{"description ILIKE": likeCondition},
			up.Cond{"jmbg ILIKE": likeCondition},
		)

		if num, err := strconv.Atoi(*input.Search); err == nil {
			numericConditions := up.Or(
				up.Cond{"decision_number": num},
			)
			conditionAndExp = up.And(conditionAndExp, up.Or(stringConditions, numericConditions))
		} else {
			conditionAndExp = up.And(conditionAndExp, stringConditions)
		}
	}

	data, total, err := h.repo.GetAll(input.Page, input.Size, conditionAndExp)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}
	response := dto.ToFineListResponseDTO(data)

	return response, total, nil
}
