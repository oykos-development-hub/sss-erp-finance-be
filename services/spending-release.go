package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type SpendingReleaseServiceImpl struct {
	App               *celeritas.Celeritas
	repo              data.SpendingRelease
	repoCurrentBudget data.CurrentBudget
}

func NewSpendingReleaseServiceImpl(app *celeritas.Celeritas, repo data.SpendingRelease) SpendingReleaseService {
	return &SpendingReleaseServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *SpendingReleaseServiceImpl) CreateSpendingRelease(inputDTO dto.SpendingReleaseDTO) (*dto.SpendingReleaseResponseDTO, error) {
	currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
		up.Cond{"budget_id": inputDTO.BudgetID},
		up.Cond{"unit_id": inputDTO.UnitID},
		up.Cond{"account_id": inputDTO.AccountID},
	))
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	inputData := data.SpendingRelease{
		CurrentBudgetID: currentBudget.ID,
		Month:           inputDTO.Month,
		Value:           inputDTO.Value,
	}
	if inputData.ValidateNewRelease() {
		h.App.ErrorLog.Println("release is possible only in the first 5 days of current month")
		return nil, errors.ErrBadRequest
	}

	id, err := h.repo.Insert(inputData)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	item, err := h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	err = h.repoCurrentBudget.UpdateBalance(currentBudget.ID, currentBudget.Balance.Add(item.Value))
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToSpendingReleaseResponseDTO(*item)

	return &res, nil
}

func (h *SpendingReleaseServiceImpl) DeleteSpendingRelease(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *SpendingReleaseServiceImpl) GetSpendingRelease(id int) (*dto.SpendingReleaseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToSpendingReleaseResponseDTO(*data)

	return &response, nil
}

func (h *SpendingReleaseServiceImpl) GetSpendingReleaseList(filter dto.SpendingReleaseFilterDTO) ([]dto.SpendingReleaseResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	// if filter.Year != nil {
	// 	conditionAndExp = up.And(conditionAndExp, &up.Cond{"year": *filter.Year})
	// }

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
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToSpendingReleaseListResponseDTO(data)

	return response, total, nil
}
