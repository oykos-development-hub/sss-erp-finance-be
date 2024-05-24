package services

import (
	goerrors "errors"
	"fmt"
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

func (h *SpendingDynamicServiceImpl) CreateSpendingDynamic(inputData []dto.SpendingDynamicDTO) ([]dto.SpendingDynamicWithEntryResponseDTO, error) {
	res := make([]dto.SpendingDynamicWithEntryResponseDTO, len(inputData))
	for i, input := range inputData {
		inputData := input.ToSpendingDynamic()
		entriesInputData := input.ToSpendingDynamicEntry()

		spendingDynamic, err := h.repo.GetBy(up.And(
			up.Cond{"budget_id": inputData.BudgetID},
			up.Cond{"unit_id": inputData.UnitID},
			up.Cond{"account_id": inputData.AccountID},
		), nil)
		if err != nil {
			if !goerrors.Is(err, errors.ErrNotFound) {
				return nil, err
			}

			actual, err := h.repoBudgetRequests.GetActual(inputData.BudgetID, inputData.UnitID, input.AccountID)
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

		fmt.Println(entriesInputData.SumOfMonths(), spendingDynamic.PlannedTotal)

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

		entriesData, err := h.repoEntries.FindBy(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		res[i] = *dto.ToSpendingDynamicWithEntryResponseDTO(spendingDynamic, entriesData)
	}

	return res, nil
}

func (h *SpendingDynamicServiceImpl) GetSpendingDynamicHistory(budgetID, unitID int) ([]dto.SpendingDynamicHistoryResponseDTO, error) {
	spendingDynamic, err := h.repo.GetBy(up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
	), nil)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	history, err := h.repoEntries.FindAll(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
	if err != nil {
		return nil, err
	}

	res := make([]dto.SpendingDynamicHistoryResponseDTO, len(history))

	for i, item := range history {
		res[i] = dto.SpendingDynamicHistoryResponseDTO{
			BudgetID:  budgetID,
			UnitID:    unitID,
			CreatedAt: item.CreatedAt,
			Username:  item.Username,
		}
	}

	return res, nil
}

func (h *SpendingDynamicServiceImpl) GetSpendingDynamic(budgetID, unitID int) ([]dto.SpendingDynamicWithEntryResponseDTO, error) {
	spendingDynamicList, err := h.repo.List(up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
	), nil)
	if err != nil {
		return nil, err
	}

	res := make([]dto.SpendingDynamicWithEntryResponseDTO, len(spendingDynamicList))

	for i, spendingDynamic := range spendingDynamicList {
		entry, err := h.repoEntries.FindBy(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
		if err != nil {
			return nil, err
		}
		res[i] = *dto.ToSpendingDynamicWithEntryResponseDTO(&spendingDynamic, entry)
	}

	return res, nil
}

func (h *SpendingDynamicServiceImpl) GetActual(budgetID, unitID, accountID int) (decimal.NullDecimal, error) {
	actual, err := h.repoBudgetRequests.GetActual(budgetID, unitID, accountID)
	if err != nil {
		return decimal.NullDecimal{}, errors.ErrInternalServer
	}

	return actual, nil
}
