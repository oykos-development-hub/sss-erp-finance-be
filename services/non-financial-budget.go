package services

import (
	"context"

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

func (h *NonFinancialBudgetServiceImpl) CreateNonFinancialBudget(ctx context.Context, input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error) {
	data := input.ToNonFinancialBudget()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	res := dto.ToNonFinancialBudgetResponseDTO(*data)

	return &res, nil
}

func (h *NonFinancialBudgetServiceImpl) UpdateNonFinancialBudget(ctx context.Context, id int, input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error) {
	data := input.ToNonFinancialBudget()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	response := dto.ToNonFinancialBudgetResponseDTO(*data)

	return &response, nil
}

func (h *NonFinancialBudgetServiceImpl) DeleteNonFinancialBudget(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
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

	if filter.RequestIDList != nil && len(*filter.RequestIDList) > 0 {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"request_id": up.In(*filter.RequestIDList...)})
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
