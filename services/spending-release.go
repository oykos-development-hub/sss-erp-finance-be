package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type SpendingReleaseServiceImpl struct {
	App               *celeritas.Celeritas
	repo              data.SpendingRelease
	repoBudget        data.Budget
	repoCurrentBudget data.CurrentBudget
}

func NewSpendingReleaseServiceImpl(app *celeritas.Celeritas, repo data.SpendingRelease, repoCurrentBudget data.CurrentBudget, repoBudget data.Budget) SpendingReleaseService {
	return &SpendingReleaseServiceImpl{
		App:               app,
		repo:              repo,
		repoBudget:        repoBudget,
		repoCurrentBudget: repoCurrentBudget,
	}
}

func (h *SpendingReleaseServiceImpl) CreateSpendingRelease(budgetID, unitID int, inputDTOList []dto.SpendingReleaseDTO) ([]dto.SpendingReleaseResponseDTO, error) {
	res := make([]dto.SpendingReleaseResponseDTO, 0, len(inputDTOList))

	for _, inputDTO := range inputDTOList {
		currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
			up.Cond{"budget_id": budgetID},
			up.Cond{"unit_id": unitID},
			up.Cond{"account_id": inputDTO.AccountID},
		))
		if err != nil {
			return nil, err
		}

		_, err = h.repo.GetBy(*up.And(up.Cond{"current_budget_id": currentBudget.ID}, up.Cond{"month": inputDTO.Month}))
		if !errors.IsErr(err, errors.NotFoundCode) {
			return nil, errors.NewWithCode(errors.SingleMonthSpendingReleaseCode, "service.spending-release.CreateSpendingRelease: only single release is allowed per month")
		}

		budget, err := h.repoBudget.Get(currentBudget.BudgetID)
		if err != nil {
			return nil, errors.Wrap(err, "service.spending-release.CreateSpendingRelease")
		}

		inputData := data.SpendingRelease{
			CurrentBudgetID: currentBudget.ID,
			Year:            budget.Year,
			Month:           inputDTO.Month,
			Value:           inputDTO.Value,
		}
		if !inputData.ValidateNewRelease() {
			return nil, errors.NewWithCode(errors.ReleaseInCurrentMonthCode, "service.CreateSpendingRelease: release is possible only in the current month")
		}

		if currentBudget.Vault().Sub(inputData.Value).LessThan(decimal.Zero) {
			return nil, errors.NewWithCode(errors.NotEnoughFundsCode, "service.CreateSpendingRelease: not enough funds")
		}

		id, err := h.repo.Insert(inputData)
		if err != nil {
			return nil, errors.Wrap(err, "service.spending-release.CreateSpendingRelease")
		}

		item, err := h.repo.Get(id)
		if err != nil {
			return nil, errors.Wrap(err, "service.spending-release.CreateSpendingRelease")
		}

		err = h.repoCurrentBudget.UpdateBalance(currentBudget.ID, currentBudget.Balance.Add(item.Value))
		if err != nil {
			return nil, errors.Wrap(err, "service.spending-release.CreateSpendingRelease")
		}

		resItem := dto.ToSpendingReleaseResponseDTO(item)

		res = append(res, resItem)
	}

	return res, nil
}

func (h *SpendingReleaseServiceImpl) DeleteSpendingRelease(id int) error {
	spendingRelease, err := h.repo.Get(id)
	if err != nil {
		return errors.Wrap(err, "service.spending-release.DeleteSpendingRelease")
	}

	err = h.repo.Delete(id)
	if err != nil {
		return errors.Wrap(err, "service.spending-release.DeleteSpendingRelease")
	}

	currentBudget, err := h.repoCurrentBudget.Get(spendingRelease.CurrentBudgetID)
	if err != nil {
		return errors.Wrap(err, "service.spending-release.DeleteSpendingRelease")
	}

	err = h.repoCurrentBudget.UpdateBalance(currentBudget.ID, currentBudget.Balance.Sub(spendingRelease.Value))
	if err != nil {
		return errors.Wrap(err, "service.spending-release.DeleteSpendingRelease")
	}

	return nil
}

func (h *SpendingReleaseServiceImpl) GetSpendingRelease(id int) (*dto.SpendingReleaseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "service.spending-release.GetSpendingRelease")
	}
	response := dto.ToSpendingReleaseResponseDTO(data)

	return &response, nil
}

func (h *SpendingReleaseServiceImpl) GetSpendingReleaseList(filter data.SpendingReleaseFilterDTO) ([]dto.SpendingReleaseResponseDTO, error) {
	data, err := h.repo.GetAll(filter)
	if err != nil {
		return nil, errors.Wrap(err, "service.spending-release.GetSpendingReleaseList")
	}
	response := dto.ToSpendingReleaseListResponseDTO(data)

	return response, nil
}

// GetSpendingReleaseOverview implements SpendingReleaseService.
func (h *SpendingReleaseServiceImpl) GetSpendingReleaseOverview(filter dto.SpendingReleaseOverviewFilterDTO) ([]dto.SpendingReleaseOverview, error) {
	data, err := h.repo.GetAllSum(filter.Month, filter.Year, filter.BudgetID, filter.UnitID)
	if err != nil {
		return nil, errors.Wrap(err, "service.spending-release.GetSpendingReleaseList")
	}

	response := dto.ToSpendingReleaseOverviewDTO(data)

	return response, nil
}
