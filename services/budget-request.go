package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type BudgetRequestServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.BudgetRequest
}

func NewBudgetRequestServiceImpl(app *celeritas.Celeritas, repo data.BudgetRequest) BudgetRequestService {
	return &BudgetRequestServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *BudgetRequestServiceImpl) CreateBudgetRequest(ctx context.Context, input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error) {
	data := input.ToBudgetRequest()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo budget request insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo budget request get")
	}

	res := dto.ToBudgetRequestResponseDTO(*data)

	return &res, nil
}

func (h *BudgetRequestServiceImpl) UpdateBudgetRequest(ctx context.Context, id int, input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error) {
	data := input.ToBudgetRequest()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo budget request update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo budget request get")
	}

	response := dto.ToBudgetRequestResponseDTO(*data)

	return &response, nil
}

func (h *BudgetRequestServiceImpl) DeleteBudgetRequest(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo budget request delete")
	}

	return nil
}

func (h *BudgetRequestServiceImpl) GetBudgetRequest(id int) (*dto.BudgetRequestResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo budget request get")
	}

	response := dto.ToBudgetRequestResponseDTO(*data)

	return &response, nil
}

func (h *BudgetRequestServiceImpl) GetBudgetRequestList(filter dto.BudgetRequestFilterDTO) ([]dto.BudgetRequestResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.BudgetID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_id": *filter.BudgetID})
	}
	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}
	if filter.RequestType != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"request_type": *filter.RequestType})
	}
	if filter.RequestTypes != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"request_type": up.In(*filter.RequestTypes...)})
	}
	if filter.ParentID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"parent_id": *filter.ParentID})
	}
	if len(filter.Statuses) > 0 {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": up.In(filter.Statuses...)})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo budget request get all")
	}
	response := dto.ToBudgetRequestListResponseDTO(data)

	return response, total, nil
}
