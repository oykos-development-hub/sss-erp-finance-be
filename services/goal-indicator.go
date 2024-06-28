package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type GoalIndicatorServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.GoalIndicator
}

func NewGoalIndicatorServiceImpl(app *celeritas.Celeritas, repo data.GoalIndicator) GoalIndicatorService {
	return &GoalIndicatorServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *GoalIndicatorServiceImpl) CreateGoalIndicator(ctx context.Context, input dto.GoalIndicatorDTO) (*dto.GoalIndicatorResponseDTO, error) {
	data := input.ToGoalIndicator()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo goal indicator insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo goal indicator get")
	}

	res := dto.ToGoalIndicatorResponseDTO(*data)

	return &res, nil
}

func (h *GoalIndicatorServiceImpl) UpdateGoalIndicator(ctx context.Context, id int, input dto.GoalIndicatorDTO) (*dto.GoalIndicatorResponseDTO, error) {
	data := input.ToGoalIndicator()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo goal indicator update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo goal indicator get")
	}

	response := dto.ToGoalIndicatorResponseDTO(*data)

	return &response, nil
}

func (h *GoalIndicatorServiceImpl) DeleteGoalIndicator(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo goal indicator delete")
	}

	return nil
}

func (h *GoalIndicatorServiceImpl) GetGoalIndicator(id int) (*dto.GoalIndicatorResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo goal indicator get")
	}

	response := dto.ToGoalIndicatorResponseDTO(*data)

	return &response, nil
}

func (h *GoalIndicatorServiceImpl) GetGoalIndicatorList(filter dto.GoalIndicatorFilterDTO) ([]dto.GoalIndicatorResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.GoalID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"goal_id": *filter.GoalID})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo goal indicator get all")
	}
	response := dto.ToGoalIndicatorListResponseDTO(data)

	return response, total, nil
}
