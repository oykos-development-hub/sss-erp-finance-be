package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ExternalReallocationServiceImpl struct {
	App               *celeritas.Celeritas
	repo              data.ExternalReallocation
	itemsRepo         data.ExternalReallocationItem
	currentBudgetRepo data.CurrentBudget
}

func NewExternalReallocationServiceImpl(app *celeritas.Celeritas, repo data.ExternalReallocation, itemsRepo data.ExternalReallocationItem, currentBudgetRepo data.CurrentBudget) ExternalReallocationService {
	return &ExternalReallocationServiceImpl{
		App:               app,
		repo:              repo,
		itemsRepo:         itemsRepo,
		currentBudgetRepo: currentBudgetRepo,
	}
}

func (h *ExternalReallocationServiceImpl) CreateExternalReallocation(input dto.ExternalReallocationDTO) (*dto.ExternalReallocationResponseDTO, error) {
	dataToInsert := input.ToExternalReallocation()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToExternalReallocationItem()
			itemToInsert.ReallocationID = id

			_, err = h.itemsRepo.Insert(tx, *itemToInsert)

			if err != nil {
				return errors.ErrInternalServer
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToExternalReallocationResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ExternalReallocationServiceImpl) DeleteExternalReallocation(id int) error {
	reallocation, err := h.GetExternalReallocation(id)

	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	if reallocation.Status == data.ReallocationStatusCreated {
		err = h.repo.Delete(id)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return errors.ErrInternalServer
		}
	} else {
		return errors.ErrBadRequest
	}

	return nil
}

func (h *ExternalReallocationServiceImpl) GetExternalReallocation(id int) (*dto.ExternalReallocationResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToExternalReallocationResponseDTO(*data)

	condition := up.And(
		up.Cond{"reallocation_id": data.ID},
	)

	items, _, err := h.itemsRepo.GetAll(nil, nil, condition, nil)

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	responseItems := dto.ToExternalReallocationItemListResponseDTO(items)

	response.Items = responseItems

	return &response, nil
}

func (h *ExternalReallocationServiceImpl) GetExternalReallocationList(filter dto.ExternalReallocationFilterDTO) ([]dto.ExternalReallocationResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	// if filter.Year != nil {
	// 	conditionAndExp = up.And(conditionAndExp, &up.Cond{"year": *filter.Year})
	// }

	if filter.SourceOrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"source_organization_unit_id": *filter.SourceOrganizationUnitID})
	}

	if filter.DestinationOrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"destination_organization_unit_id": *filter.DestinationOrganizationUnitID})
	}

	if filter.BudgetID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_id": *filter.BudgetID})
	}

	if filter.RequestedBy != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"requested_by": *filter.RequestedBy})
	}

	if filter.RequestedBy != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"requested_by": *filter.RequestedBy})
	}

	if filter.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.RequestedBy})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToExternalReallocationListResponseDTO(data)

	return response, total, nil
}

func (h *ExternalReallocationServiceImpl) AcceptOUExternalReallocation(input dto.ExternalReallocationDTO) (*dto.ExternalReallocationResponseDTO, error) {
	dataToInsert := input.ToExternalReallocation()

	id := input.ID
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		err = h.repo.AcceptOUExternalReallocation(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToExternalReallocationItem()
			itemToInsert.ReallocationID = id

			_, err = h.itemsRepo.Insert(tx, *itemToInsert)

			if err != nil {
				return errors.ErrInternalServer
			}

			if item.DestinationAccountID != 0 {

				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": dataToInsert.BudgetID},
					up.Cond{"unit_id": dataToInsert.DestinationOrganizationUnitID},
					up.Cond{"account_id": itemToInsert.DestinationAccountID},
				))

				if err != nil {
					return errors.ErrInternalServer
				}

				value := currentBudget.Actual.Sub(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateActual(currentBudget.ID, value)

				if err != nil {
					return errors.ErrInternalServer
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToExternalReallocationResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ExternalReallocationServiceImpl) RejectOUExternalReallocation(id int) error {

	err := h.repo.RejectOUExternalReallocation(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}
