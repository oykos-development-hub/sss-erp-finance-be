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
	repoEntries       data.SpendingDynamicEntry
	repoCurrentBudget data.CurrentBudget
	repoRelease       data.SpendingRelease
}

func NewSpendingDynamicServiceImpl(
	app *celeritas.Celeritas,
	repoEntries data.SpendingDynamicEntry,
	repoCurrentBudget data.CurrentBudget,
) SpendingDynamicService {
	return &SpendingDynamicServiceImpl{
		App:               app,
		repoEntries:       repoEntries,
		repoCurrentBudget: repoCurrentBudget,
	}
}

func (h *SpendingDynamicServiceImpl) CreateSpendingDynamic(budgetID, unitID int, inputDataDTO []dto.SpendingDynamicDTO) error {
	latestVersion, err := h.repoEntries.FindLatestVersion()
	if err != nil {
		return errors.Wrap(err, "find latest version")
	}

	for _, inputDTO := range inputDataDTO {
		currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
			up.Cond{"budget_id": budgetID},
			up.Cond{"unit_id": unitID},
			up.Cond{"account_id": inputDTO.AccountID},
		))
		if err != nil {
			return errors.Wrap(err, "CreateSpendingDynamic")
		}

		entriesInputData := inputDTO.ToSpendingDynamicEntry()

		oldDynamic, err := h.GetSpendingDynamic(&currentBudget.ID, nil, nil, nil)
		if err != nil {
			return errors.Wrap(err, "CreateSpendingDynamic")
		}
		if len(oldDynamic) > 0 {
			if entriesInputData.SumOfMonths().GreaterThan(currentBudget.Actual.Add(oldDynamic[0].TotalSavings)) {
				return errors.NewBadRequestError("sum cannot be greater than actual plus all the savings")
			}
		} else {
			if !entriesInputData.SumOfMonths().Equal(currentBudget.Actual) {
				return errors.NewBadRequestError("sum must match actual of current budget")
			}
		}

		entriesInputData.CurrentBudgetID = currentBudget.ID
		entriesInputData.Version = latestVersion

		_, err = h.repoEntries.Insert(*entriesInputData)
		if err != nil {
			return errors.Wrap(err, "CreateSpendingDynamic")
		}
	}

	return nil
}

func (h *SpendingDynamicServiceImpl) CreateInititalSpendingDynamicFromCurrentBudget(currentBudget *data.CurrentBudget) error {
	spendingDynamicEntry := h.generateInitialSpendingDynamicEntry(currentBudget)

	_, err := h.repoEntries.Insert(*spendingDynamicEntry)
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
		January:         monthlyAmount,
		February:        monthlyAmount,
		March:           monthlyAmount,
		April:           monthlyAmount,
		May:             monthlyAmount,
		June:            monthlyAmount,
		July:            monthlyAmount,
		August:          monthlyAmount,
		September:       monthlyAmount,
		October:         monthlyAmount,
		November:        monthlyAmount,
		December:        decemberAmount,
		Version:         1,
		CurrentBudgetID: currentBudget.ID,
	}
}

func (h *SpendingDynamicServiceImpl) GetSpendingDynamicHistory(budgetID, unitID int) ([]dto.SpendingDynamicHistoryResponseDTO, error) {
	history, err := h.repoEntries.FindHistoryChanges(budgetID, unitID)
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

func (h *SpendingDynamicServiceImpl) GetSpendingDynamic(currentBudgetID, budgetID, unitID *int, version *int) ([]dto.SpendingDynamicWithEntryResponseDTO, error) {
	entries, err := h.repoEntries.FindAll(currentBudgetID, version, budgetID, unitID)
	if err != nil {
		return nil, errors.Wrap(err, "GetSpendingDynamic")
	}

	entriesListRes := make([]dto.SpendingDynamicWithEntryResponseDTO, 0)

	for _, entry := range entries {
		filter := data.SpendingReleaseFilterDTO{
			CurrentBudgetID: &entry.CurrentBudgetID,
		}
		releases, err := h.repoRelease.GetAll(filter)
		if err != nil {
			return nil, errors.Wrap(err, "GetSpendingDynamic")
		}
		entryRes := dto.SpendingDynamicWithEntryResponseDTO{
			ID:              entry.ID,
			CurrentBudgetID: entry.CurrentBudgetID,
			AccountID:       entry.AccountID,
			BudgetID:        entry.BudgetID,
			UnitID:          entry.UnitID,
			Actual:          entry.Actual,
			Username:        entry.Username,
			CreatedAt:       entry.CreatedAt,
			January:         dto.MonthEntry{Value: entry.January},
			February:        dto.MonthEntry{Value: entry.February},
			March:           dto.MonthEntry{Value: entry.March},
			April:           dto.MonthEntry{Value: entry.April},
			May:             dto.MonthEntry{Value: entry.May},
			June:            dto.MonthEntry{Value: entry.June},
			July:            dto.MonthEntry{Value: entry.July},
			August:          dto.MonthEntry{Value: entry.August},
			September:       dto.MonthEntry{Value: entry.September},
			October:         dto.MonthEntry{Value: entry.October},
			November:        dto.MonthEntry{Value: entry.November},
			December:        dto.MonthEntry{Value: entry.December},
		}
		for _, release := range releases {
			switch release.Month {
			case 1:
				entryRes.January.Savings = entryRes.January.Value.Sub(release.Value)
			case 2:
				entryRes.February.Savings = entryRes.February.Value.Sub(release.Value)
			case 3:
				entryRes.March.Savings = entryRes.March.Value.Sub(release.Value)
			case 4:
				entryRes.April.Savings = entryRes.April.Value.Sub(release.Value)
			case 5:
				entryRes.May.Savings = entryRes.May.Value.Sub(release.Value)
			case 6:
				entryRes.June.Savings = entryRes.June.Value.Sub(release.Value)
			case 7:
				entryRes.July.Savings = entryRes.July.Value.Sub(release.Value)
			case 8:
				entryRes.August.Savings = entryRes.August.Value.Sub(release.Value)
			case 9:
				entryRes.September.Savings = entryRes.September.Value.Sub(release.Value)
			case 10:
				entryRes.October.Savings = entryRes.October.Value.Sub(release.Value)
			case 11:
				entryRes.November.Savings = entryRes.November.Value.Sub(release.Value)
			case 12:
				entryRes.December.Savings = entryRes.December.Value.Sub(release.Value)
			}
		}

		entryRes.TotalSavings = entryRes.GetTotalSavings()

		entriesListRes = append(entriesListRes, entryRes)
	}

	return entriesListRes, nil
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
