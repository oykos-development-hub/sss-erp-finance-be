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

func (h *SpendingDynamicServiceImpl) CreateSpendingDynamic(inputDataDTO []dto.SpendingDynamicDTO) error {
	for _, inputDTO := range inputDataDTO {
		currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
			up.Cond{"budget_id": inputDTO.BudgetID},
			up.Cond{"unit_id": inputDTO.UnitID},
			up.Cond{"account_id": inputDTO.AccountID},
		))
		if err != nil {
			return errors.Wrap(err, "CreateSpendingDynamic")
		}

		entriesInputData := inputDTO.ToSpendingDynamicEntry()

		// Validate that the sum of the months matches the planned total
		if !entriesInputData.SumOfMonths().Equal(currentBudget.Actual) {
			return errors.NewBadRequestError("sum must match actual of current budget")
		}

		entries, err := h.repoEntries.FindAll(&currentBudget.ID, nil, nil, nil)
		if err != nil {
			if !errors.IsErr(err, errors.NotFoundCode) {
				return errors.Wrap(err, "CreateSpendingDynamic")
			}
		}

		// validate months if there are entries already
		if len(entries) > 0 {
			ok := entriesInputData.ValidateNewEntry(&entries[0])
			if !ok {
				return errors.NewBadRequestError("cannot change months in past")
			}
		}

		entriesInputData.CurrentBudgetID = currentBudget.ID
		entriesInputData.Version = len(entries) + 1

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

func (h *SpendingDynamicServiceImpl) GetSpendingDynamic(budgetID, unitID int, version *int) ([]dto.SpendingDynamicWithEntryResponseDTO, error) {
	entries, err := h.repoEntries.FindAll(nil, version, &budgetID, &unitID)
	if err != nil {
		return nil, errors.Wrap(err, "GetSpendingDynamic")
	}

	return dto.ToSpendingDynamicWithEntryListResponseDTO(entries), nil
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
