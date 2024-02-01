package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

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

func (h *ActivityServiceImpl) CreateActivity(input dto.ActivityDTO) (*dto.ActivityResponseDTO, error) {
	data := input.ToActivity()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToActivityResponseDTO(*data)

	return &res, nil
}

func (h *ActivityServiceImpl) UpdateActivity(id int, input dto.ActivityDTO) (*dto.ActivityResponseDTO, error) {
	data := input.ToActivity()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToActivityResponseDTO(*data)

	return &response, nil
}

func (h *ActivityServiceImpl) DeleteActivity(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *ActivityServiceImpl) GetActivity(id int) (*dto.ActivityResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
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
	response := dto.ToActivityListResponseDTO(data)

	return response, total, nil
}
