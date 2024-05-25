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
			return nil, errors.ErrInternalServer
		}

		spendingDynamic, err := h.repo.GetBy(up.And(
			up.Cond{"current_budget_id": currentBudget.ID},
		), nil)
		if err != nil {
			if !goerrors.Is(err, errors.ErrNotFound) {
				return nil, err
			}

			inputData := data.SpendingDynamic{
				CurrentBudgetID: currentBudget.ID,
				PlannedTotal:    currentBudget.Actual,
			}

			id, err := h.repo.Insert(inputData)
			if err != nil {
				return nil, errors.ErrInternalServer
			}

			spendingDynamic, err = h.repo.Get(id)
			if err != nil {
				return nil, errors.ErrInternalServer
			}
		}

		entriesInputData := inputDTO.ToSpendingDynamicEntry()

		// Validate that the sum of the months matches the planned total
		if !entriesInputData.SumOfMonths().Equal(spendingDynamic.PlannedTotal) {
			return nil, errors.ErrInvalidInput
		}

		entries, err := h.repoEntries.FindAll(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})
		if err != nil {
			if !goerrors.Is(err, up.ErrNoMoreRows) {
				return nil, err
			}
		}

		version := len(entries) + 1

		// validate months if there are entries already
		if len(entries) > 0 {
			ok := entriesInputData.ValidateNewEntry(&entries[0])
			if !ok {
				log.Println("cannot change months in past")
				return nil, errors.ErrBadRequest
			}
		}

		entriesInputData.SpendingDynamicID = spendingDynamic.ID
		entriesInputData.Version = version

		_, err = h.repoEntries.Insert(*entriesInputData)
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		entriesData, err := h.repoEntries.FindBy(up.And(up.Cond{"spending_dynamic_id": spendingDynamic.ID}))
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		res[i] = *dto.ToSpendingDynamicWithEntryResponseDTO(spendingDynamic, entriesData)
	}

	return res, nil
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

	currentBudgetIDList := make([]int, len(currentBudgets))
	for i, currentBudget := range currentBudgets {
		currentBudgetIDList[i] = currentBudget.ID
	}

	spendingDynamicList, err := h.repo.List(up.And(
		up.Cond{"current_budget_id IN": currentBudgetIDList},
	), nil)
	if err != nil {
		return nil, err
	}

	res := make([]dto.SpendingDynamicWithEntryResponseDTO, len(spendingDynamicList))

	for i, spendingDynamic := range spendingDynamicList {
		conditionAndExp := up.And(&up.Cond{"spending_dynamic_id": spendingDynamic.ID})

		if version != nil {
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"version": version})
		}

		entry, err := h.repoEntries.FindBy(conditionAndExp)
		if err != nil {
			return nil, err
		}
		res[i] = *dto.ToSpendingDynamicWithEntryResponseDTO(&spendingDynamic, entry)
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
		return decimal.Zero, errors.ErrInternalServer
	}

	return currentBudget.Actual, nil
}
