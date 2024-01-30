package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
)

type FinancialBudgetServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.FinancialBudget
}

func NewFinancialBudgetServiceImpl(app *celeritas.Celeritas, repo data.FinancialBudget) FinancialBudgetService {
	return &FinancialBudgetServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FinancialBudgetServiceImpl) CreateFinancialBudget(input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error) {
	data := input.ToFinancialBudget()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFinancialBudgetResponseDTO(*data)

	return &res, nil
}

func (h *FinancialBudgetServiceImpl) UpdateFinancialBudget(id int, input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error) {
	data := input.ToFinancialBudget()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *FinancialBudgetServiceImpl) DeleteFinancialBudget(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *FinancialBudgetServiceImpl) GetFinancialBudget(id int) (*dto.FinancialBudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *FinancialBudgetServiceImpl) GetFinancialBudgetByBudgetID(id int) (*dto.FinancialBudgetResponseDTO, error) {
	cond := db.Cond{"budget_id": id}
	data, err := h.repo.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFinancialBudgetResponseDTO(*data[0])

	return &response, nil
}

func (h *FinancialBudgetServiceImpl) GetFinancialBudgetList() ([]dto.FinancialBudgetResponseDTO, error) {
	data, err := h.repo.GetAll(nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}
	response := dto.ToFinancialBudgetListResponseDTO(data)

	return response, nil
}
