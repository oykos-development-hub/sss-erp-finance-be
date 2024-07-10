package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type ExternalReallocationServiceImpl struct {
	App                 *celeritas.Celeritas
	repo                data.ExternalReallocation
	itemsRepo           data.ExternalReallocationItem
	currentBudgetRepo   data.CurrentBudget
	spendingDynamicRepo data.SpendingDynamicEntry
}

func NewExternalReallocationServiceImpl(app *celeritas.Celeritas, repo data.ExternalReallocation, itemsRepo data.ExternalReallocationItem, currentBudgetRepo data.CurrentBudget, spendingDynamicRepo data.SpendingDynamicEntry) ExternalReallocationService {
	return &ExternalReallocationServiceImpl{
		App:                 app,
		repo:                repo,
		itemsRepo:           itemsRepo,
		currentBudgetRepo:   currentBudgetRepo,
		spendingDynamicRepo: spendingDynamicRepo,
	}
}

func (h *ExternalReallocationServiceImpl) CreateExternalReallocation(ctx context.Context, input dto.ExternalReallocationDTO) (*dto.ExternalReallocationResponseDTO, error) {
	dataToInsert := input.ToExternalReallocation()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo external reallocation insert")
		}

		for _, item := range input.Items {
			itemToInsert := item.ToExternalReallocationItem()
			itemToInsert.ReallocationID = id

			_, err = h.itemsRepo.Insert(tx, *itemToInsert)

			if err != nil {
				return newErrors.Wrap(err, "repo external reallocation item insert")
			}
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo external reallocation get")
	}

	res := dto.ToExternalReallocationResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ExternalReallocationServiceImpl) DeleteExternalReallocation(ctx context.Context, id int) error {
	reallocation, err := h.GetExternalReallocation(id)

	if err != nil {
		return newErrors.Wrap(err, "repo external reallocation get")
	}

	if reallocation.Status == data.ReallocationStatusCreated {
		err = h.repo.Delete(ctx, id)
		if err != nil {
			return newErrors.Wrap(err, "repo external reallocation delete")
		}
	} else {
		return newErrors.Wrap(errors.ErrBadRequest, "repo external reallocation delete")
	}

	return nil
}

func (h *ExternalReallocationServiceImpl) GetExternalReallocation(id int) (*dto.ExternalReallocationResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo external reallocation get")
	}

	response := dto.ToExternalReallocationResponseDTO(*data)

	condition := up.And(
		up.Cond{"reallocation_id": data.ID},
	)

	items, _, err := h.itemsRepo.GetAll(nil, nil, condition, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo external reallocation items get all")
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

	if filter.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	if filter.OrganizationUnitID != nil && *filter.OrganizationUnitID != 0 {
		search := up.Or(
			up.Cond{"source_organization_unit_id": *filter.OrganizationUnitID},
			up.Cond{"destination_organization_unit_id": *filter.OrganizationUnitID},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo external reallocation get all")
	}
	response := dto.ToExternalReallocationListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		condition := up.And(
			up.Cond{"reallocation_id": response[i].ID},
		)

		items, _, err := h.itemsRepo.GetAll(nil, nil, condition, nil)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo external reallocation items get all")
		}

		responseItems := dto.ToExternalReallocationItemListResponseDTO(items)

		response[i].Items = responseItems
	}
	return response, total, nil
}

func (h *ExternalReallocationServiceImpl) AcceptOUExternalReallocation(ctx context.Context, input dto.ExternalReallocationDTO) (*dto.ExternalReallocationResponseDTO, error) {
	dataToInsert := input.ToExternalReallocation()

	reallocation, err := h.GetExternalReallocation(input.ID)

	if err != nil {
		return nil, newErrors.Wrap(err, "get external reallocation")
	}

	if reallocation.Status != data.ReallocationStatusCreated {
		return nil, newErrors.Wrap(errors.ErrAlreadyDone, "get external reallocation")
	}

	id := input.ID
	err = data.Upper.Tx(func(tx up.Session) error {
		var err error
		dataToInsert.ID = id
		err = h.repo.AcceptOUExternalReallocation(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo external reallocation accept ou")
		}

		spendingDynamic, err := h.spendingDynamicRepo.FindAll(nil, nil, &reallocation.BudgetID, &reallocation.SourceOrganizationUnitID)

		if err != nil {
			return newErrors.Wrap(err, "repo spending dynamic get spending dynamic")
		}

		for _, item := range input.Items {

			if item.SourceAccountID != 0 {
				itemToInsert := item.ToExternalReallocationItem()
				itemToInsert.ReallocationID = id

				_, err = h.itemsRepo.Insert(tx, *itemToInsert)

				if err != nil {
					return newErrors.Wrap(err, "repo external reallocation item insert")
				}

				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": reallocation.BudgetID},
					up.Cond{"unit_id": reallocation.SourceOrganizationUnitID},
					up.Cond{"account_id": itemToInsert.SourceAccountID},
					up.Cond{"type": 1}, // preusmejrenja se rade samo kod tekuceg budzeta
				))

				if err != nil {
					return newErrors.Wrap(err, "repo current budget get by")
				}

				value := currentBudget.Actual.Sub(itemToInsert.Amount)

				if value.Compare(decimal.NewFromInt(0)) < 0 {
					return newErrors.Wrap(errors.ErrInsufficientFunds, "repo current budget update actual")
				}

				err = h.currentBudgetRepo.UpdateActual(ctx, currentBudget.ID, value)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update actual")
				}

				currentAmountValue := currentBudget.CurrentAmount.Sub(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateCurrentAmount(ctx, currentBudget.ID, currentAmountValue)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update current amount")
				}

				spendingDynamic = updateDynamicSub(currentBudget.ID, itemToInsert.Amount, spendingDynamic)

			}
		}

		for _, item := range spendingDynamic {
			_, err = h.spendingDynamicRepo.Insert(ctx, data.SpendingDynamicEntry{
				CurrentBudgetID: item.CurrentBudgetID,
				//Username:        item.Username,
				January:   item.January,
				February:  item.February,
				March:     item.March,
				April:     item.April,
				May:       item.May,
				June:      item.June,
				July:      item.July,
				August:    item.August,
				September: item.September,
				October:   item.October,
				November:  item.November,
				December:  item.December,
				Version:   item.Version + 1,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo spending dynamic repo insert")
			}
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo external realocation get")
	}

	res := dto.ToExternalReallocationResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ExternalReallocationServiceImpl) RejectOUExternalReallocation(ctx context.Context, id int) error {

	err := h.repo.RejectOUExternalReallocation(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo external reallocation reject ou")
	}

	return nil
}

func (h *ExternalReallocationServiceImpl) AcceptSSSExternalReallocation(ctx context.Context, id int) error {
	reallocation, err := h.GetExternalReallocation(id)

	if err != nil {
		return newErrors.Wrap(err, "repo external reallocation get")
	}

	if reallocation.Status != data.ReallocationStatusOUAccept {
		return newErrors.Wrap(errors.ErrAlreadyDone, "repo external reallocation get")
	}

	err = data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.AcceptSSSExternalReallocation(ctx, tx, id)
		if err != nil {
			return newErrors.Wrap(err, "repo external reallocation accept sss")
		}

		spendingDynamic, err := h.spendingDynamicRepo.FindAll(nil, nil, &reallocation.BudgetID, &reallocation.DestinationOrganizationUnitID)

		if err != nil {
			return newErrors.Wrap(err, "repo spending dynamic get spending dynamic")
		}

		for _, item := range reallocation.Items {

			if item.DestinationAccountID != 0 {
				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": reallocation.BudgetID},
					up.Cond{"unit_id": reallocation.DestinationOrganizationUnitID},
					up.Cond{"account_id": item.DestinationAccountID},
					up.Cond{"type": 1}, // preusmejrenja se rade samo kod tekuceg budzeta
				))

				if err != nil {
					return newErrors.Wrap(err, "repo current budget get by")
				}

				value := currentBudget.Actual.Add(item.Amount)

				err = h.currentBudgetRepo.UpdateActual(ctx, currentBudget.ID, value)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update actual")
				}

				currentAmountValue := currentBudget.CurrentAmount.Add(item.Amount)

				err = h.currentBudgetRepo.UpdateCurrentAmount(ctx, currentBudget.ID, currentAmountValue)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update current amount")
				}

				spendingDynamic = updateDynamicAdd(currentBudget.ID, item.Amount, spendingDynamic)
			}
		}

		for _, item := range spendingDynamic {
			_, err = h.spendingDynamicRepo.Insert(ctx, data.SpendingDynamicEntry{
				CurrentBudgetID: item.CurrentBudgetID,
				//Username:        item.Username,
				January:   item.January,
				February:  item.February,
				March:     item.March,
				April:     item.April,
				May:       item.May,
				June:      item.June,
				July:      item.July,
				August:    item.August,
				September: item.September,
				October:   item.October,
				November:  item.November,
				December:  item.December,
				Version:   item.Version + 1,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo spending dynamic repo insert")
			}
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}

func (h *ExternalReallocationServiceImpl) RejectSSSExternalReallocation(ctx context.Context, id int) error {
	reallocation, err := h.GetExternalReallocation(id)

	if err != nil {
		return newErrors.Wrap(err, "repo external reallocation get")
	}

	if reallocation.Status != data.ReallocationStatusOUAccept {
		return newErrors.Wrap(errors.ErrAlreadyDone, "repo external reallocation get")
	}

	err = data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.RejectSSSExternalReallocation(ctx, tx, id)
		if err != nil {
			return newErrors.Wrap(err, "repo external reallocation reject sss")
		}

		spendingDynamic, err := h.spendingDynamicRepo.FindAll(nil, nil, &reallocation.BudgetID, &reallocation.SourceOrganizationUnitID)

		if err != nil {
			return newErrors.Wrap(err, "repo spending dynamic get spending dynamic")
		}

		for _, item := range reallocation.Items {

			if item.SourceAccountID != 0 {
				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": reallocation.BudgetID},
					up.Cond{"unit_id": reallocation.SourceOrganizationUnitID},
					up.Cond{"account_id": item.SourceAccountID},
					up.Cond{"type": 1}, // preusmejrenja se rade samo kod tekuceg budzeta
				))

				if err != nil {
					return newErrors.Wrap(err, "repo current budget get by")
				}

				value := currentBudget.Actual.Add(item.Amount)

				err = h.currentBudgetRepo.UpdateActual(ctx, currentBudget.ID, value)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update actual")
				}

				currentAmountValue := currentBudget.CurrentAmount.Sub(item.Amount)

				err = h.currentBudgetRepo.UpdateCurrentAmount(ctx, currentBudget.ID, currentAmountValue)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update current amount")
				}

				spendingDynamic = updateDynamicAdd(currentBudget.ID, item.Amount, spendingDynamic)
			}
		}

		for _, item := range spendingDynamic {
			_, err = h.spendingDynamicRepo.Insert(ctx, data.SpendingDynamicEntry{
				CurrentBudgetID: item.CurrentBudgetID,
				//Username:        item.Username,
				January:   item.January,
				February:  item.February,
				March:     item.March,
				April:     item.April,
				May:       item.May,
				June:      item.June,
				July:      item.July,
				August:    item.August,
				September: item.September,
				October:   item.October,
				November:  item.November,
				December:  item.December,
				Version:   item.Version + 1,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo spending dynamic repo insert")
			}
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}
