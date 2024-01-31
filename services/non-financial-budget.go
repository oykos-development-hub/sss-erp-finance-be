package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type NonFinancialBudgetServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.NonFinancialBudget
}

func NewNonFinancialBudgetServiceImpl(app *celeritas.Celeritas, repo data.NonFinancialBudget) NonFinancialBudgetService {
	return &NonFinancialBudgetServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *NonFinancialBudgetServiceImpl) CreateNonFinancialBudget(input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error) {
	data := input.ToNonFinancialBudget()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToNonFinancialBudgetResponseDTO(*data)

	return &res, nil
}

func (h *NonFinancialBudgetServiceImpl) UpdateNonFinancialBudget(id int, input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error) {
	data := input.ToNonFinancialBudget()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToNonFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *NonFinancialBudgetServiceImpl) DeleteNonFinancialBudget(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *NonFinancialBudgetServiceImpl) GetNonFinancialBudget(id int) (*dto.NonFinancialBudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToNonFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *NonFinancialBudgetServiceImpl) GetNonFinancialBudgetList(filter dto.NonFinancialBudgetFilterDTO) ([]dto.NonFinancialBudgetResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.BudgetID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_id": *filter.BudgetID})
	}

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
	response := dto.ToNonFinancialBudgetListResponseDTO(data)

	return response, total, nil
}
