package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"

	"github.com/oykos-development-hub/celeritas"
)

type SpendingDynamicServiceImpl struct {
	App               *celeritas.Celeritas
	repo              data.SpendingDynamic
	repoEntries       data.SpendingDynamicEntry
	repoCurrentBudget data.CurrentBudget
}

func NewSpendingDynamicServiceImpl(
	app *celeritas.Celeritas,
	repo data.SpendingDynamic,
	repoEntries data.SpendingDynamicEntry,
	repoCurrentBudget data.CurrentBudget,
) SpendingDynamicService {
	return &SpendingDynamicServiceImpl{
		App:               app,
		repo:              repo,
		repoEntries:       repoEntries,
		repoCurrentBudget: repoCurrentBudget,
	}
}

func (h *SpendingDynamicServiceImpl) CreateSpendingDynamic(inputDataDTO []dto.SpendingDynamicDTO) ([]dto.SpendingDynamicWithEntryResponseDTO, error) {
	res := make([]dto.SpendingDynamicWithEntryResponseDTO, len(inputDataDTO))
	for i, inputDTO := range inputDataDTO {
		currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
			up.Cond{"budget_id": inputDTO.BudgetID},
			up.Cond{"unit_id": inputDTO.UnitID},
			up.Cond{"account_id": inputDTO.AccountID},
		))
		if err != nil {
			return nil, errors.Wrap(err, "CreateSpendingDynamic")
		}

		spendingDynamic, err := h.repo.GetBy(up.And(
			up.Cond{"current_budget_id": currentBudget.ID},
		), nil)
		if err != nil {
			if !errors.IsErr(err, errors.NotFoundCode) {
				return nil, err
			}

			inputData := data.SpendingDynamic{
				CurrentBudgetID: currentBudget.ID,
				PlannedTotal:    currentBudget.Actual,
			}

			id, err := h.repo.Insert(inputData)
			if err != nil {
				return nil, errors.Wrap(err, "CreateSpendingDynamic")
			}

			spendingDynamic, err = h.repo.Get(id)
			if err != nil {
				return nil, errors.Wrap(err, "CreateSpendingDynamic")
			}
		}

		entriesInputData := inputDTO.ToSpendingDynamicEntry()

		// Validate that the sum of the months matches the planned total
		if !entriesInputData.SumOfMonths().Equal(spendingDynamic.PlannedTotal) {
			return nil, errors.NewBadRequestError("sum must match actual of current budget")
		}

		entries, err := h.repoEntries.FindAll(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
		if err != nil {
			if !errors.IsErr(err, errors.NotFoundCode) {
				return nil, err
			}
		}

		version := len(entries) + 1

		// validate months if there are entries already
		if len(entries) > 0 {
			ok := entriesInputData.ValidateNewEntry(&entries[0])
			if !ok {
				return nil, errors.NewBadRequestError("cannot change months in past")
			}
		}

		entriesInputData.SpendingDynamicID = spendingDynamic.ID
		entriesInputData.Version = version

		_, err = h.repoEntries.Insert(*entriesInputData)
		if err != nil {
			return nil, errors.Wrap(err, "CreateSpendingDynamic")
		}

		entriesData, err := h.repoEntries.FindBy(up.And(up.Cond{"spending_dynamic_id": spendingDynamic.ID}))
		if err != nil {
			return nil, errors.Wrap(err, "CreateSpendingDynamic")
		}

		res[i] = *dto.ToSpendingDynamicWithEntryResponseDTO(spendingDynamic, entriesData, currentBudget)
	}

	return res, nil
}

func (h *SpendingDynamicServiceImpl) CreateInititalSpendingDynamicFromCurrentBudget(currentBudget *data.CurrentBudget) error {
	inputData := data.SpendingDynamic{
		CurrentBudgetID: currentBudget.ID,
		PlannedTotal:    currentBudget.Actual,
	}

	spendingDynamicID, err := h.repo.Insert(inputData)
	if err != nil {
		return errors.Wrap(err, "CreateInititalSpendingDynamicFromCurrentBudget")
	}

	spendingDynamicEntry := h.generateInitialSpendingDynamicEntry(currentBudget)

	spendingDynamicEntry.SpendingDynamicID = spendingDynamicID

	_, err = h.repoEntries.Insert(*spendingDynamicEntry)
	if err != nil {
		return errors.Wrap(err, "CreateInititalSpendingDynamicFromCurrentBudget")
	}

	return err
}

func (h *SpendingDynamicServiceImpl) generateInitialSpendingDynamicEntry(currentBudget *data.CurrentBudget) *data.SpendingDynamicEntry {
	monthlyAmount := currentBudget.Actual.Div(decimal.NewFromInt(12)).Round(2)

	// Sum of the first 11 rounded months
	totalForFirst11Months := monthlyAmount.Mul(decimal.NewFromInt(11))

	// Adjust the December amount to account for rounding differences
	decemberAmount := currentBudget.Actual.Sub(totalForFirst11Months).Round(2)

	return &data.SpendingDynamicEntry{
		January:   monthlyAmount,
		February:  monthlyAmount,
		March:     monthlyAmount,
		April:     monthlyAmount,
		May:       monthlyAmount,
		June:      monthlyAmount,
		July:      monthlyAmount,
		August:    monthlyAmount,
		September: monthlyAmount,
		October:   monthlyAmount,
		November:  monthlyAmount,
		December:  decemberAmount,
		Version:   1,
	}
}

func (h *SpendingDynamicServiceImpl) GetSpendingDynamicHistory(budgetID, unitID int) ([]dto.SpendingDynamicHistoryResponseDTO, error) {
	currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
	))
	if err != nil {
		return nil, err
	}

	spendingDynamic, err := h.repo.GetBy(up.And(
		up.Cond{"current_budget_id": currentBudget.ID},
	), nil)
	if err != nil {
		return nil, errors.Wrap(err, "GetSpendingDynamicHistory")
	}

	history, err := h.repoEntries.FindAll(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
	if err != nil {
		return nil, errors.Wrap(err, "GetSpendingDynamicHistory")
	}

	res := make([]dto.SpendingDynamicHistoryResponseDTO, len(history))

	for i, item := range history {
		res[i] = dto.SpendingDynamicHistoryResponseDTO{
			BudgetID:  budgetID,
			UnitID:    unitID,
			Version:   item.Version,
			CreatedAt: item.CreatedAt,
			Username:  item.Username,
		}
	}

	return res, nil
}

func (h *SpendingDynamicServiceImpl) GetSpendingDynamic(budgetID, unitID int, version *int) ([]dto.SpendingDynamicWithEntryResponseDTO, error) {
	currentBudgets, _, err := h.repoCurrentBudget.GetAll(nil, nil, up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
	), nil)
	if err != nil {
		return nil, err
	}

	currentBudgetMap := make(map[int]*data.CurrentBudget)

	// Populate the map with the current budgets.
	for _, currentBudget := range currentBudgets {
		currentBudgetMap[currentBudget.ID] = currentBudget
	}

	currentBudgetIDList := make([]int, 0, len(currentBudgetMap))
	for id := range currentBudgetMap {
		currentBudgetIDList = append(currentBudgetIDList, id)
	}

	spendingDynamicList, err := h.repo.List(up.And(
		up.Cond{"current_budget_id IN": currentBudgetIDList},
	), nil)
	if err != nil {
		return nil, errors.Wrap(err, "GetSpendingDynamic")
	}

	res := make([]dto.SpendingDynamicWithEntryResponseDTO, len(spendingDynamicList))

	for i, spendingDynamic := range spendingDynamicList {
		conditionAndExp := up.And(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})

		if version != nil {
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"version": version})
		}

		entry, err := h.repoEntries.FindBy(conditionAndExp)
		if err != nil {
			return nil, errors.Wrap(err, "GetSpendingDynamic")
		}
		res[i] = *dto.ToSpendingDynamicWithEntryResponseDTO(&spendingDynamic, entry, currentBudgetMap[spendingDynamic.CurrentBudgetID])
	}

	return res, nil
}

func (h *SpendingDynamicServiceImpl) GetActual(budgetID, unitID, accountID int) (decimal.Decimal, error) {
	currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
		up.Cond{"account_id": accountID},
	))
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "GetActual")
	}

	return currentBudget.Actual, nil
}
