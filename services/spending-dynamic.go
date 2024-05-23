package services

import (
	goerrors "errors"
	"log"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"

	"github.com/oykos-development-hub/celeritas"
)

type SpendingDynamicServiceImpl struct {
	App                *celeritas.Celeritas
	repo               data.SpendingDynamic
	repoEntries        data.SpendingDynamicEntry
	repoBudgetRequests data.BudgetRequest
}

func NewSpendingDynamicServiceImpl(
	app *celeritas.Celeritas,
	repo data.SpendingDynamic,
	repoEntries data.SpendingDynamicEntry,
	repoBudgetRequests data.BudgetRequest,
) SpendingDynamicService {
	return &SpendingDynamicServiceImpl{
		App:                app,
		repo:               repo,
		repoEntries:        repoEntries,
		repoBudgetRequests: repoBudgetRequests,
	}
}

func (h *SpendingDynamicServiceImpl) CreateSpendingDynamic(input dto.SpendingDynamicDTO) (*dto.SpendingDynamicWithEntriesResponseDTO, error) {
	inputData := input.ToSpendingDynamic()
	entriesInputData := input.ToSpendingDynamicEntry()

	spendingDynamic, err := h.repo.GetBy(up.And(
		up.Cond{"budget_id": inputData.BudgetID},
		up.Cond{"unit_id": inputData.UnitID},
	), nil)
	if err != nil {
		if !goerrors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		actual, err := h.repoBudgetRequests.GetActual(inputData.BudgetID, inputData.UnitID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		if !actual.Valid {
			log.Printf("No actual for budget with id: %d and unit id: %d", input.BudgetID, input.UnitID)
			return nil, errors.ErrInternalServer
		}

		inputData.PlannedTotal = actual.Decimal

		id, err := h.repo.Insert(*inputData)
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		spendingDynamic, err = h.repo.Get(id)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
	}

	// Validate that the sum of the months matches the planned total
	if !entriesInputData.SumOfMonths().Equal(spendingDynamic.PlannedTotal) {
		return nil, errors.ErrInvalidInput
	}

	entry, err := h.repoEntries.FindBy(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
	if err != nil {
		if !goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, err
		}
	}

	// validate months if there are entries already
	if entry != nil {
		ok := entriesInputData.ValidateNewEntry(entry)
		if !ok {
			log.Println("cannot change months in past")
			return nil, errors.ErrBadRequest
		}
	}

	entriesInputData.SpendingDynamicID = spendingDynamic.ID

	_, err = h.repoEntries.Insert(*entriesInputData)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	entriesData, err := h.repoEntries.FindAll(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return dto.ToSpendingDynamicWithEntryResponseDTO(spendingDynamic, entriesData), nil
}

func (h *SpendingDynamicServiceImpl) GetSpendingDynamicHistory(budgetID, unitID int) (*dto.SpendingDynamicWithEntriesResponseDTO, error) {
	data, err := h.repo.GetBy(up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
	), nil)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	entries, err := h.repoEntries.FindAll(&up.Cond{"spending_dynamic_id": data.ID})
	if err != nil {
		return nil, err
	}

	return dto.ToSpendingDynamicWithEntryResponseDTO(data, entries), nil
}

func (h *SpendingDynamicServiceImpl) GetSpendingDynamic(budgetID, unitID int) (*dto.SpendingDynamicWithEntriesResponseDTO, error) {
	spendingDynamic, err := h.repo.GetBy(up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
	), nil)
	if err != nil {
		return nil, err
	}

	entry, err := h.repoEntries.FindBy(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
	if err != nil {
		return nil, err
	}

	return dto.ToSpendingDynamicWithEntryResponseDTO(spendingDynamic, []data.SpendingDynamicEntry{*entry}), nil
}

func (h *SpendingDynamicServiceImpl) GetActual(budgetID, unitID int) (decimal.NullDecimal, error) {

	actual, err := h.repoBudgetRequests.GetActual(budgetID, unitID)
	if err != nil {
		return decimal.NullDecimal{}, errors.ErrInternalServer
	}

	return actual, nil
}
