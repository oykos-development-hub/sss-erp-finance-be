package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type NonFinancialBudgetGoalServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.NonFinancialBudgetGoal
}

func NewNonFinancialBudgetGoalServiceImpl(app *celeritas.Celeritas, repo data.NonFinancialBudgetGoal) NonFinancialBudgetGoalService {
	return &NonFinancialBudgetGoalServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *NonFinancialBudgetGoalServiceImpl) CreateNonFinancialBudgetGoal(input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error) {
	data := input.ToNonFinancialBudgetGoal()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToNonFinancialBudgetGoalResponseDTO(*data)

	return &res, nil
}

func (h *NonFinancialBudgetGoalServiceImpl) UpdateNonFinancialBudgetGoal(id int, input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error) {
	data := input.ToNonFinancialBudgetGoal()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToNonFinancialBudgetGoalResponseDTO(*data)

	return &response, nil
}

func (h *NonFinancialBudgetGoalServiceImpl) DeleteNonFinancialBudgetGoal(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *NonFinancialBudgetGoalServiceImpl) GetNonFinancialBudgetGoal(id int) (*dto.NonFinancialBudgetGoalResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToNonFinancialBudgetGoalResponseDTO(*data)

	return &response, nil
}

func (h *NonFinancialBudgetGoalServiceImpl) GetNonFinancialBudgetGoalList(filter dto.NonFinancialBudgetGoalFilterDTO) ([]dto.NonFinancialBudgetGoalResponseDTO, *uint64, error) {
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
	response := dto.ToNonFinancialBudgetGoalListResponseDTO(data)

	return response, total, nil
}
