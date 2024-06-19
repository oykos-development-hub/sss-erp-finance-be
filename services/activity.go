package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ActivityServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.Activity
}

func NewActivityServiceImpl(app *celeritas.Celeritas, repo data.Activity) ActivityService {
	return &ActivityServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *ActivityServiceImpl) CreateActivity(ctx context.Context, input dto.ActivityDTO) (*dto.ActivityResponseDTO, error) {
	data := input.ToActivity()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo activity insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo activity get")
	}

	res := dto.ToActivityResponseDTO(*data)

	return &res, nil
}

func (h *ActivityServiceImpl) UpdateActivity(ctx context.Context, id int, input dto.ActivityDTO) (*dto.ActivityResponseDTO, error) {
	data := input.ToActivity()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo activity update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo activity get")
	}

	response := dto.ToActivityResponseDTO(*data)

	return &response, nil
}

func (h *ActivityServiceImpl) DeleteActivity(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo activity delete")
	}

	return nil
}

func (h *ActivityServiceImpl) GetActivity(id int) (*dto.ActivityResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo activity get")
	}

	response := dto.ToActivityResponseDTO(*data)

	return &response, nil
}

func (h *ActivityServiceImpl) GetActivityList(filter dto.ActivityFilterDTO) ([]dto.ActivityResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.FilterBySubProgramID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"sub_program_id": *filter.FilterBySubProgramID})
	}
	if filter.FilterByOrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.FilterByOrganizationUnitID})
	}

	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		searchCond := up.Or(
			up.Cond{"title ILIKE": likeCondition},
			up.Cond{"code ILIKE": likeCondition},
			up.Cond{"description ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, searchCond)
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
		return nil, nil, newErrors.Wrap(err, "repo activity get all")
	}
	response := dto.ToActivityListResponseDTO(data)

	return response, total, nil
}
