package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
)

type BudgetServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.Budget
}

func NewBudgetServiceImpl(app *celeritas.Celeritas, repo data.Budget) BudgetService {
	return &BudgetServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *BudgetServiceImpl) CreateBudget(input dto.BudgetDTO) (*dto.BudgetResponseDTO, error) {
	data := input.ToBudget()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToBudgetResponseDTO(*data)

	return &res, nil
}

func (h *BudgetServiceImpl) UpdateBudget(id int, input dto.BudgetDTO) (*dto.BudgetResponseDTO, error) {
	data := input.ToBudget()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToBudgetResponseDTO(*data)

	return &response, nil
}

func (h *BudgetServiceImpl) DeleteBudget(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *BudgetServiceImpl) GetBudget(id int) (*dto.BudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToBudgetResponseDTO(*data)

	return &response, nil
}

func (h *BudgetServiceImpl) GetBudgetList(input dto.GetBudgetListInput) ([]dto.BudgetResponseDTO, error) {
	cond := db.Cond{}
	if input.Year != nil {
		cond["year"] = input.Year
	}
	if input.BudgetType != nil {
		cond["budget_type"] = input.BudgetType
	}

	data, err := h.repo.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}
	response := dto.ToBudgetListResponseDTO(data)

	return response, nil
}