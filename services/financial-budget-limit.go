package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FinancialBudgetLimitServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.FinancialBudgetLimit
}

func NewFinancialBudgetLimitServiceImpl(app *celeritas.Celeritas, repo data.FinancialBudgetLimit) FinancialBudgetLimitService {
	return &FinancialBudgetLimitServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FinancialBudgetLimitServiceImpl) CreateFinancialBudgetLimit(input dto.FinancialBudgetLimitDTO) (*dto.FinancialBudgetLimitResponseDTO, error) {
	data := input.ToFinancialBudgetLimit()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFinancialBudgetLimitResponseDTO(*data)

	return &res, nil
}

func (h *FinancialBudgetLimitServiceImpl) UpdateFinancialBudgetLimit(id int, input dto.FinancialBudgetLimitDTO) (*dto.FinancialBudgetLimitResponseDTO, error) {
	data := input.ToFinancialBudgetLimit()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFinancialBudgetLimitResponseDTO(*data)

	return &response, nil
}

func (h *FinancialBudgetLimitServiceImpl) DeleteFinancialBudgetLimit(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *FinancialBudgetLimitServiceImpl) GetFinancialBudgetLimit(id int) (*dto.FinancialBudgetLimitResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFinancialBudgetLimitResponseDTO(*data)

	return &response, nil
}

func (h *FinancialBudgetLimitServiceImpl) GetFinancialBudgetLimitList(filter dto.FinancialBudgetLimitFilterDTO) ([]dto.FinancialBudgetLimitResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_id": filter.BudgetID})

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
	response := dto.ToFinancialBudgetLimitListResponseDTO(data)

	return response, total, nil
}
