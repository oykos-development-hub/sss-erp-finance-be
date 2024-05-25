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

func (h *SpendingReleaseServiceImpl) CreateSpendingRelease(inputDTO dto.SpendingReleaseDTO) (*dto.SpendingReleaseResponseDTO, error) {
	currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
		up.Cond{"budget_id": inputDTO.BudgetID},
		up.Cond{"unit_id": inputDTO.UnitID},
		up.Cond{"account_id": inputDTO.AccountID},
	))
	if err != nil {
		return nil, err
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
	if inputData.ValidateNewRelease() {
		return nil, errors.NewBadRequestError("service.CreateSpendingRelease: release is possible only in the first 5 days of current month")
	}

	if currentBudget.Actual.Sub(inputData.Value).LessThan(decimal.Zero) {
		return nil, errors.NewBadRequestError("service.CreateSpendingRelease: not enough funds")
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

	res := dto.ToSpendingReleaseResponseDTO(*item)

	return &res, nil
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
	response := dto.ToSpendingReleaseResponseDTO(*data)

	return &response, nil
}

func (h *SpendingReleaseServiceImpl) GetSpendingReleaseList(filter dto.SpendingReleaseFilterDTO) ([]dto.SpendingReleaseResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.Year != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"year": *filter.Year})
	}
	if filter.Month != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"month": *filter.Month})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, errors.Wrap(err, "service.spending-release.GetSpendingReleaseList")
	}
	response := dto.ToSpendingReleaseListResponseDTO(data)

	return response, total, nil
}
