package services

import (
	"context"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type InternalReallocationServiceImpl struct {
	App                 *celeritas.Celeritas
	repo                data.InternalReallocation
	itemsRepo           data.InternalReallocationItem
	currentBudgetRepo   data.CurrentBudget
	spendingDynamicRepo data.SpendingDynamicEntry
}

func NewInternalReallocationServiceImpl(app *celeritas.Celeritas, repo data.InternalReallocation, itemsRepo data.InternalReallocationItem, currentBudgetRepo data.CurrentBudget, spendingDynamic data.SpendingDynamicEntry) InternalReallocationService {
	return &InternalReallocationServiceImpl{
		App:                 app,
		repo:                repo,
		itemsRepo:           itemsRepo,
		currentBudgetRepo:   currentBudgetRepo,
		spendingDynamicRepo: spendingDynamic,
	}
}

func (h *InternalReallocationServiceImpl) CreateInternalReallocation(ctx context.Context, input dto.InternalReallocationDTO) (*dto.InternalReallocationResponseDTO, error) {
	dataToInsert := input.ToInternalReallocation()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo internal reallocation insert")
		}

		spendingDynamic, err := h.spendingDynamicRepo.FindAll(nil, nil, &input.BudgetID, &input.OrganizationUnitID)

		if err != nil {
			return newErrors.Wrap(err, "repo spending dynamic get spending dynamic")
		}

		for _, item := range input.Items {
			itemToInsert := item.ToInternalReallocationItem()
			itemToInsert.ReallocationID = id

			_, err = h.itemsRepo.Insert(tx, *itemToInsert)

			if err != nil {
				return newErrors.Wrap(err, "repo internal reallocation item insert")
			}

			if item.SourceAccountID != 0 {

				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": dataToInsert.BudgetID},
					up.Cond{"unit_id": dataToInsert.OrganizationUnitID},
					up.Cond{"account_id": itemToInsert.SourceAccountID},
				))

				if err != nil {
					return newErrors.Wrap(err, "repo current budget get by")
				}

				value := currentBudget.Actual.Sub(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateActual(ctx, currentBudget.ID, value)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update actual")
				}

				currentAmountValue := currentBudget.CurrentAmount.Sub(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateCurrentAmount(ctx, currentBudget.ID, currentAmountValue)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update current amount")
				}

				spendingDynamic = updateDynamicSub(itemToInsert.SourceAccountID, itemToInsert.Amount, spendingDynamic)

			}
			if item.DestinationAccountID != 0 {

				currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
					up.Cond{"budget_id": dataToInsert.BudgetID},
					up.Cond{"unit_id": dataToInsert.OrganizationUnitID},
					up.Cond{"account_id": itemToInsert.DestinationAccountID},
				))

				if err != nil {
					return newErrors.Wrap(err, "repo current budget get by")
				}

				value := currentBudget.Actual.Add(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateActual(ctx, currentBudget.ID, value)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update actual")
				}

				currentAmountValue := currentBudget.CurrentAmount.Add(itemToInsert.Amount)

				err = h.currentBudgetRepo.UpdateCurrentAmount(ctx, currentBudget.ID, currentAmountValue)

				if err != nil {
					return newErrors.Wrap(err, "repo current budget update current amount")
				}

				spendingDynamic = updateDynamicAdd(itemToInsert.DestinationAccountID, itemToInsert.Amount, spendingDynamic)

			}
		}

		for _, item := range spendingDynamic {
			_, err = h.spendingDynamicRepo.Insert(ctx, data.SpendingDynamicEntry{
				ID:              item.ID,
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
		return nil, newErrors.Wrap(err, "repo internal reallocation get")
	}

	res := dto.ToInternalReallocationResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *InternalReallocationServiceImpl) DeleteInternalReallocation(ctx context.Context, id int) error {
	reallocation, err := h.GetInternalReallocation(id)

	if err != nil {
		return newErrors.Wrap(err, "repo internal reallocation get")
	}

	spendingDynamic, err := h.spendingDynamicRepo.FindAll(nil, nil, &reallocation.BudgetID, &reallocation.OrganizationUnitID)

	if err != nil {
		return newErrors.Wrap(err, "repo spending dynamic get spending dynamic")
	}

	for _, item := range reallocation.Items {
		if item.DestinationAccountID != 0 {

			currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
				up.Cond{"budget_id": reallocation.BudgetID},
				up.Cond{"unit_id": reallocation.OrganizationUnitID},
				up.Cond{"account_id": item.DestinationAccountID},
			))

			if err != nil {
				return newErrors.Wrap(err, "repo current budget get by")
			}

			value := currentBudget.Actual.Sub(item.Amount)

			if value.Cmp(decimal.NewFromFloat(0)) < 0 {
				return newErrors.Wrap(errors.ErrInsufficientFunds, "update actual")
			}

		}
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo internal reallocation delete")
	}

	for _, item := range reallocation.Items {
		if item.SourceAccountID != 0 {

			currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
				up.Cond{"budget_id": reallocation.BudgetID},
				up.Cond{"unit_id": reallocation.OrganizationUnitID},
				up.Cond{"account_id": item.SourceAccountID},
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

			spendingDynamic = updateDynamicAdd(item.SourceAccountID, item.Amount, spendingDynamic)

		}
		if item.DestinationAccountID != 0 {

			currentBudget, err := h.currentBudgetRepo.GetBy(*up.And(
				up.Cond{"budget_id": reallocation.BudgetID},
				up.Cond{"unit_id": reallocation.OrganizationUnitID},
				up.Cond{"account_id": item.DestinationAccountID},
			))

			if err != nil {
				return newErrors.Wrap(err, "repo current budget get by")
			}

			value := currentBudget.Actual.Sub(item.Amount)

			err = h.currentBudgetRepo.UpdateActual(ctx, currentBudget.ID, value)

			if err != nil {
				return newErrors.Wrap(err, "repo current budget update actual")
			}

			currentAmountValue := currentBudget.CurrentAmount.Sub(item.Amount)

			err = h.currentBudgetRepo.UpdateCurrentAmount(ctx, currentBudget.ID, currentAmountValue)

			if err != nil {
				return newErrors.Wrap(err, "repo current budget update current amount")
			}

			spendingDynamic = updateDynamicSub(item.DestinationAccountID, item.Amount, spendingDynamic)
		}
	}

	for _, item := range spendingDynamic {
		_, err = h.spendingDynamicRepo.Insert(ctx, data.SpendingDynamicEntry{
			ID:              item.ID,
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
}

func (h *InternalReallocationServiceImpl) GetInternalReallocation(id int) (*dto.InternalReallocationResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo internal reallocation get")
	}

	condition := up.And(
		up.Cond{"reallocation_id": data.ID},
	)

	items, _, err := h.itemsRepo.GetAll(nil, nil, condition, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo internal reallocation item get")
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
		return nil, nil, newErrors.Wrap(err, "repo internal reallocation get all")
	}
	response := dto.ToInternalReallocationListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		condition := up.And(
			up.Cond{"reallocation_id": response[i].ID},
		)

		items, _, err := h.itemsRepo.GetAll(nil, nil, condition, nil)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo internal reallocation item get all")
		}

		responseItems := dto.ToInternalReallocationItemListResponseDTO(items)

		response[0].Items = responseItems
		var amount decimal.Decimal
		for _, item := range items {
			if item.DestinationAccountID != 0 {
				amount = amount.Sub(item.Amount)
			}
		}
		response[0].Sum = amount
	}

	return response, total, nil
}

func updateDynamicAdd(accountID int, amount decimal.Decimal, spendingDynamic []data.SpendingDynamicEntryWithCurrentBudget) []data.SpendingDynamicEntryWithCurrentBudget {

	monthMap := map[time.Month]int{
		time.January:   0,
		time.February:  1,
		time.March:     2,
		time.April:     3,
		time.May:       4,
		time.June:      5,
		time.July:      6,
		time.August:    7,
		time.September: 8,
		time.October:   9,
		time.November:  10,
		time.December:  11,
	}

	for i := 0; i < len(spendingDynamic); i++ {
		if spendingDynamic[i].CurrentBudgetID == accountID {
			currentMonth := time.Now().Month()
			monthIndex := monthMap[currentMonth]

			parts := decimal.NewFromInt(int64(12 - monthIndex))
			amountPart := amount.Div(parts)

			for j := monthIndex; j < 12; j++ {
				switch j {
				case 0:
					spendingDynamic[i].January = spendingDynamic[i].January.Add(amountPart)
				case 1:
					spendingDynamic[i].February = spendingDynamic[i].February.Add(amountPart)
				case 2:
					spendingDynamic[i].March = spendingDynamic[i].March.Add(amountPart)
				case 3:
					spendingDynamic[i].April = spendingDynamic[i].April.Add(amountPart)
				case 4:
					spendingDynamic[i].May = spendingDynamic[i].May.Add(amountPart)
				case 5:
					spendingDynamic[i].June = spendingDynamic[i].June.Add(amountPart)
				case 6:
					spendingDynamic[i].July = spendingDynamic[i].July.Add(amountPart)
				case 7:
					spendingDynamic[i].August = spendingDynamic[i].August.Add(amountPart)
				case 8:
					spendingDynamic[i].September = spendingDynamic[i].September.Add(amountPart)
				case 9:
					spendingDynamic[i].October = spendingDynamic[i].October.Add(amountPart)
				case 10:
					spendingDynamic[i].November = spendingDynamic[i].November.Add(amountPart)
				case 11:
					spendingDynamic[i].December = spendingDynamic[i].December.Add(amountPart)
				}
			}
		}
	}

	return spendingDynamic
}

func updateDynamicSub(accountID int, amount decimal.Decimal, spendingDynamic []data.SpendingDynamicEntryWithCurrentBudget) []data.SpendingDynamicEntryWithCurrentBudget {

	monthMap := map[time.Month]int{
		time.January:   0,
		time.February:  1,
		time.March:     2,
		time.April:     3,
		time.May:       4,
		time.June:      5,
		time.July:      6,
		time.August:    7,
		time.September: 8,
		time.October:   9,
		time.November:  10,
		time.December:  11,
	}

	for i := 0; i < len(spendingDynamic); i++ {
		if spendingDynamic[i].CurrentBudgetID == accountID {
			currentMonth := time.Now().Month()
			monthIndex := monthMap[currentMonth]

			parts := decimal.NewFromInt(int64(12 - monthIndex))
			amountPart := amount.Div(parts)

			for j := monthIndex; j < 12; j++ {
				switch j {
				case 0:
					spendingDynamic[i].January = spendingDynamic[i].January.Sub(amountPart)
				case 1:
					spendingDynamic[i].February = spendingDynamic[i].February.Sub(amountPart)
				case 2:
					spendingDynamic[i].March = spendingDynamic[i].March.Sub(amountPart)
				case 3:
					spendingDynamic[i].April = spendingDynamic[i].April.Sub(amountPart)
				case 4:
					spendingDynamic[i].May = spendingDynamic[i].May.Sub(amountPart)
				case 5:
					spendingDynamic[i].June = spendingDynamic[i].June.Sub(amountPart)
				case 6:
					spendingDynamic[i].July = spendingDynamic[i].July.Sub(amountPart)
				case 7:
					spendingDynamic[i].August = spendingDynamic[i].August.Sub(amountPart)
				case 8:
					spendingDynamic[i].September = spendingDynamic[i].September.Sub(amountPart)
				case 9:
					spendingDynamic[i].October = spendingDynamic[i].October.Sub(amountPart)
				case 10:
					spendingDynamic[i].November = spendingDynamic[i].November.Sub(amountPart)
				case 11:
					spendingDynamic[i].December = spendingDynamic[i].December.Sub(amountPart)
				}
			}
		}
	}

	return spendingDynamic
}
