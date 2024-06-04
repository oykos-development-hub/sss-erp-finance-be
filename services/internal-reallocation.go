package services

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type InternalReallocationServiceImpl struct {
	App               *celeritas.Celeritas
	repo              data.InternalReallocation
	itemsRepo         data.InternalReallocationItem
	currentBudgetRepo data.CurrentBudget
}

func NewInternalReallocationServiceImpl(app *celeritas.Celeritas, repo data.InternalReallocation, itemsRepo data.InternalReallocationItem, currentBudgetRepo data.CurrentBudget) InternalReallocationService {
	return &InternalReallocationServiceImpl{
		App:               app,
		repo:              repo,
		itemsRepo:         itemsRepo,
		currentBudgetRepo: currentBudgetRepo,
	}
}

func (h *InternalReallocationServiceImpl) CreateInternalReallocation(input dto.InternalReallocationDTO) (*dto.InternalReallocationResponseDTO, error) {
	dataToInsert := input.ToInternalReallocation()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToInternalReallocationItem()
			itemToInsert.ReallocationID = id

			_, err = h.itemsRepo.Insert(tx, *itemToInsert)

			if err != nil {
				return errors.ErrInternalServer
			}

			if item.SourceAccountID != 0 {

				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": dataToInsert.BudgetID},
					up.Cond{"unit_id": dataToInsert.OrganizationUnitID},
					up.Cond{"account_id": itemToInsert.SourceAccountID},
				))

				if err != nil {
					return errors.ErrInternalServer
				}

				currentBudget.Actual.Sub(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateActual(currentBudget.ID, currentBudget.Actual)

				if err != nil {
					return errors.ErrInternalServer
				}

			}
			if item.DestinationAccountID != 0 {

				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": dataToInsert.BudgetID},
					up.Cond{"unit_id": dataToInsert.OrganizationUnitID},
					up.Cond{"account_id": itemToInsert.DestinationAccountID},
				))

				if err != nil {
					return errors.ErrInternalServer
				}
				currentBudget.Actual.Add(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateActual(currentBudget.ID, currentBudget.Actual)

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

	res := dto.ToInternalReallocationResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *InternalReallocationServiceImpl) DeleteInternalReallocation(id int) error {
	reallocation, err := h.GetInternalReallocation(id)

	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	for _, item := range reallocation.Items {
		if item.DestinationAccountID != 0 {

			currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
				up.Cond{"budget_id": reallocation.BudgetID},
				up.Cond{"unit_id": reallocation.OrganizationUnitID},
				up.Cond{"account_id": item.DestinationAccountID},
			))

			if err != nil {
				return errors.ErrInternalServer
			}

			value := currentBudget.Actual.Sub(item.Amount)

			if value.Cmp(decimal.NewFromFloat(0)) < 0 {
				return errors.ErrInvalidInput
			}
		}
	}

	err = h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	for _, item := range reallocation.Items {
		if item.SourceAccountID != 0 {

			currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
				up.Cond{"budget_id": reallocation.BudgetID},
				up.Cond{"unit_id": reallocation.OrganizationUnitID},
				up.Cond{"account_id": item.SourceAccountID},
			))

			if err != nil {
				return errors.ErrInternalServer
			}

			currentBudget.Actual.Add(item.Amount)

			err = h.currentBudgetRepo.UpdateActual(currentBudget.ID, currentBudget.Actual)

			if err != nil {
				return errors.ErrInternalServer
			}

		}
		if item.DestinationAccountID != 0 {

			currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
				up.Cond{"budget_id": reallocation.BudgetID},
				up.Cond{"unit_id": reallocation.OrganizationUnitID},
				up.Cond{"account_id": item.DestinationAccountID},
			))

			if err != nil {
				return errors.ErrInternalServer
			}

			currentBudget.Actual.Sub(item.Amount)

			err = h.currentBudgetRepo.UpdateActual(currentBudget.ID, currentBudget.Actual)

			if err != nil {
				return errors.ErrInternalServer
			}
		}
	}

	return nil
}

func (h *InternalReallocationServiceImpl) GetInternalReallocation(id int) (*dto.InternalReallocationResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}

	condition := up.And(
		up.Cond{"reallocation_id": data.ID},
	)

	items, _, err := h.itemsRepo.GetAll(nil, nil, condition, nil)

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	response := dto.ToInternalReallocationResponseDTO(*data)

	responseItems := dto.ToInternalReallocationItemListResponseDTO(items)

	response.Items = responseItems

	return &response, nil
}

func (h *InternalReallocationServiceImpl) GetInternalReallocationList(filter dto.InternalReallocationFilterDTO) ([]dto.InternalReallocationResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.Year != nil {
		year := *filter.Year
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_request": up.Between(startOfYear, endOfYear)})
	}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.BudgetID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_id": *filter.BudgetID})
	}

	if filter.RequestedBy != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"requested_by": *filter.RequestedBy})
	}

	if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToInternalReallocationListResponseDTO(data)

	return response, total, nil
}
