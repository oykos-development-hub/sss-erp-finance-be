package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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

func (h *FinancialBudgetServiceImpl) CreateFinancialBudget(ctx context.Context, input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error) {
	data := input.ToFinancialBudget()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo financial budget insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo financial budget get")
	}

	res := dto.ToFinancialBudgetResponseDTO(*data)

	return &res, nil
}

func (h *FinancialBudgetServiceImpl) UpdateFinancialBudget(ctx context.Context, id int, input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error) {
	data := input.ToFinancialBudget()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo financial budget update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo financial budget get")
	}

	response := dto.ToFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *FinancialBudgetServiceImpl) DeleteFinancialBudget(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo financial budget delete")
	}

	return nil
}

func (h *FinancialBudgetServiceImpl) GetFinancialBudget(id int) (*dto.FinancialBudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo financial budget get")
	}

	response := dto.ToFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *FinancialBudgetServiceImpl) GetFinancialBudgetByBudgetID(id int) (*dto.FinancialBudgetResponseDTO, error) {
	cond := db.Cond{"budget_id": id}
	data, err := h.repo.GetAll(&cond)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo financial budget get all")
	}
	if len(data) == 0 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo financial budget get all")
	}
	response := dto.ToFinancialBudgetResponseDTO(*data[0])

	return &response, nil
}

func (h *FinancialBudgetServiceImpl) GetFinancialBudgetList() ([]dto.FinancialBudgetResponseDTO, error) {
	data, err := h.repo.GetAll(nil)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo financial budget get all")
	}
	response := dto.ToFinancialBudgetListResponseDTO(data)

	return response, nil
}
