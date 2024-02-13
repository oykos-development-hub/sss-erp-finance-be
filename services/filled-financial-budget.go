package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FilledFinancialBudgetServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.FilledFinancialBudget
}

func NewFilledFinancialBudgetServiceImpl(app *celeritas.Celeritas, repo data.FilledFinancialBudget) FilledFinancialBudgetService {
	return &FilledFinancialBudgetServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FilledFinancialBudgetServiceImpl) CreateFilledFinancialBudget(input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error) {
	data := input.ToFilledFinancialBudget()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFilledFinancialBudgetResponseDTO(*data)

	return &res, nil
}

func (h *FilledFinancialBudgetServiceImpl) UpdateFilledFinancialBudget(id int, input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error) {
	data := input.ToFilledFinancialBudget()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFilledFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *FilledFinancialBudgetServiceImpl) DeleteFilledFinancialBudget(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *FilledFinancialBudgetServiceImpl) GetFilledFinancialBudget(id int) (*dto.FilledFinancialBudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFilledFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *FilledFinancialBudgetServiceImpl) GetFilledFinancialBudgetList(filter dto.FilledFinancialBudgetFilterDTO) ([]dto.FilledFinancialBudgetResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": filter.OrganizationUnitID})
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_request_id": filter.BudgetRequestID})

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
	response := dto.ToFilledFinancialBudgetListResponseDTO(data)

	return response, total, nil
}
