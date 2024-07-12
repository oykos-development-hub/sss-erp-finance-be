package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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

func (h *SpendingDynamicServiceImpl) CreateSpendingDynamic(ctx context.Context, budgetID, unitID int, inputDataDTO []dto.SpendingDynamicDTO) error {
	latestVersion, err := h.repoEntries.FindLatestVersion(nil, &budgetID, &unitID)
	if err != nil {
		return newErrors.Wrap(err, "repo entries find latest version")
	}

	for _, inputDTO := range inputDataDTO {
		currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
			up.Cond{"budget_id": budgetID},
			up.Cond{"unit_id": unitID},
			up.Cond{"account_id": inputDTO.AccountID},
			up.Cond{"type": 1}))
		if err != nil {
			return newErrors.Wrap(err, "repo current budget get by")
		}

		entriesInputData := inputDTO.ToSpendingDynamicEntry()

		oldDynamic, err := h.GetSpendingDynamic(&currentBudget.ID, &budgetID, &unitID, nil)
		if err != nil {
			return newErrors.Wrap(err, "get spending dynamic")
		}
		if len(oldDynamic) > 0 {
			if entriesInputData.SumOfMonths().GreaterThan(currentBudget.CurrentAmount.Add(oldDynamic[0].TotalSavings)) {
				return newErrors.NewBadRequestError("sum cannot be greater than actual plus all the savings")
			}
		} else {
			if !entriesInputData.SumOfMonths().Equal(currentBudget.CurrentAmount) {
				return newErrors.NewBadRequestError("sum must match actual of current budget")
			}
		}

		entriesInputData.CurrentBudgetID = currentBudget.ID
		entriesInputData.Version = latestVersion + 1

		_, err = h.repoEntries.Insert(ctx, *entriesInputData)
		if err != nil {
			return newErrors.Wrap(err, "repo entries insert")
		}
	}

	return nil
}

func (h *SpendingDynamicServiceImpl) CreateInititalSpendingDynamicFromCurrentBudget(ctx context.Context, currentBudget *data.CurrentBudget) error {
	spendingDynamicEntry := h.generateInitialSpendingDynamicEntry(currentBudget)

	_, err := h.repoEntries.Insert(ctx, *spendingDynamicEntry)
	if err != nil {
		return newErrors.Wrap(err, "repo entries insert")
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
		return nil, newErrors.Wrap(err, "repo entries find history changes")
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
		return nil, newErrors.Wrap(err, "repo entries find all")
	}

	entriesListRes := make([]dto.SpendingDynamicWithEntryResponseDTO, 0)

	for _, entry := range entries {
		filter := data.SpendingReleaseFilterDTO{
			CurrentBudgetID: &entry.CurrentBudgetID,
		}
		releases, err := h.repoRelease.GetAll(filter)
		if err != nil {
			return nil, newErrors.Wrap(err, "repo release get all")
		}

		currentAmount, err := h.GetCurrentAmount(*budgetID, *unitID, entry.AccountID, 1)
		if err != nil {
			return nil, newErrors.Wrap(err, "get current amount")
		}

		entryRes := dto.SpendingDynamicWithEntryResponseDTO{
			ID:              entry.ID,
			CurrentBudgetID: entry.CurrentBudgetID,
			AccountID:       entry.AccountID,
			BudgetID:        entry.BudgetID,
			UnitID:          entry.UnitID,
			Actual:          entry.Actual,
			CurrentAmount:   currentAmount,
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

func (h *SpendingDynamicServiceImpl) GetActual(budgetID, unitID, accountID, Type int) (decimal.Decimal, error) {
	currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
		up.Cond{"account_id": accountID},
		up.Cond{"type": Type},
	))
	if err != nil {
		return decimal.Zero, newErrors.Wrap(err, "repo current budget get by")
	}

	return currentBudget.Actual, nil
}

func (h *SpendingDynamicServiceImpl) GetCurrentAmount(budgetID, unitID, accountID, Type int) (decimal.Decimal, error) {
	currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
		up.Cond{"account_id": accountID},
		up.Cond{"type": Type},
	))
	if err != nil {
		return decimal.Zero, newErrors.Wrap(err, "repo current budget get by")
	}

	return currentBudget.CurrentAmount, nil
}
