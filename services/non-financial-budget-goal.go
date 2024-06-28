package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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

func (h *NonFinancialBudgetGoalServiceImpl) CreateNonFinancialBudgetGoal(ctx context.Context, input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error) {
	data := input.ToNonFinancialBudgetGoal()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo non financial budget goal insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo non financial budget goal get")
	}

	res := dto.ToNonFinancialBudgetGoalResponseDTO(*data)

	return &res, nil
}

func (h *NonFinancialBudgetGoalServiceImpl) UpdateNonFinancialBudgetGoal(ctx context.Context, id int, input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error) {
	data := input.ToNonFinancialBudgetGoal()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo non financial budget goal update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo non financial budget goal get")
	}

	response := dto.ToNonFinancialBudgetGoalResponseDTO(*data)

	return &response, nil
}

func (h *NonFinancialBudgetGoalServiceImpl) DeleteNonFinancialBudgetGoal(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo non financial budget goal delete")
	}

	return nil
}

func (h *NonFinancialBudgetGoalServiceImpl) GetNonFinancialBudgetGoal(id int) (*dto.NonFinancialBudgetGoalResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo non financial budget goal get")
	}

	response := dto.ToNonFinancialBudgetGoalResponseDTO(*data)

	return &response, nil
}

func (h *NonFinancialBudgetGoalServiceImpl) GetNonFinancialBudgetGoalList(filter dto.NonFinancialBudgetGoalFilterDTO) ([]dto.NonFinancialBudgetGoalResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.NonFinancialBudgetID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"non_financial_budget_id": *filter.NonFinancialBudgetID})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(nil, nil, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo non financial budget goal get all")
	}
	response := dto.ToNonFinancialBudgetGoalListResponseDTO(data)

	return response, total, nil
}
