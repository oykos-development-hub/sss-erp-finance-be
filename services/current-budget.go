package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type CurrentBudgetServiceImpl struct {
	App             *celeritas.Celeritas
	repo            data.CurrentBudget
	spendingService SpendingDynamicService
}

func NewCurrentBudgetServiceImpl(
	app *celeritas.Celeritas,
	repo data.CurrentBudget,
	spendingService SpendingDynamicService,
) CurrentBudgetService {
	return &CurrentBudgetServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *CurrentBudgetServiceImpl) CreateCurrentBudget(input dto.CurrentBudgetDTO) (*dto.CurrentBudgetResponseDTO, error) {
	data := input.ToCurrentBudget()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.Wrap(err, "CreateCurrentBudget")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "CreateCurrentBudget")
	}

	res := dto.ToCurrentBudgetResponseDTO(data)

	err = h.spendingService.CreateInititalSpendingDynamicFromCurrentBudget(data)
	if err != nil {
		return nil, errors.Wrap(err, "CreateCurrentBudget")
	}

	return &res, nil
}

func (h *CurrentBudgetServiceImpl) UpdateActual(unitID, budgetID, accountID int, actual decimal.Decimal) (*dto.CurrentBudgetResponseDTO, error) {
	currentBudget, err := h.repo.GetBy(*up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
		up.Cond{"account_id": accountID},
	))
	if err != nil {
		return nil, errors.Wrap(err, "UpdateActual")
	}

	err = h.repo.UpdateActual(currentBudget.ID, actual)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateActual")
	}

	response := dto.ToCurrentBudgetResponseDTO(currentBudget)

	return &response, nil
}

func (h *CurrentBudgetServiceImpl) GetCurrentBudget(id int) (*dto.CurrentBudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.WrapNotFoundError(err, "GetCurrentBudget")
	}
	response := dto.ToCurrentBudgetResponseDTO(data)

	return &response, nil
}

func (h *CurrentBudgetServiceImpl) GetCurrentBudgetList(filter dto.CurrentBudgetFilterDTO) ([]dto.CurrentBudgetResponseDTO, *uint64, error) {
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
		return nil, nil, errors.Wrap(err, "GetCurrentBudgetList")
	}
	response := dto.ToCurrentBudgetListResponseDTO(data)

	return response, total, nil
}
