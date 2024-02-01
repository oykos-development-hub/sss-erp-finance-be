package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ProgramServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.Program
}

func NewProgramServiceImpl(app *celeritas.Celeritas, repo data.Program) ProgramService {
	return &ProgramServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *ProgramServiceImpl) CreateProgram(input dto.ProgramDTO) (*dto.ProgramResponseDTO, error) {
	data := input.ToProgram()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToProgramResponseDTO(*data)

	return &res, nil
}

func (h *ProgramServiceImpl) UpdateProgram(id int, input dto.ProgramDTO) (*dto.ProgramResponseDTO, error) {
	data := input.ToProgram()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToProgramResponseDTO(*data)

	return &response, nil
}

func (h *ProgramServiceImpl) DeleteProgram(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *ProgramServiceImpl) GetProgram(id int) (*dto.ProgramResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToProgramResponseDTO(*data)

	return &response, nil
}

func (h *ProgramServiceImpl) GetProgramList(filter dto.ProgramFilterDTO) ([]dto.ProgramResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.IsProgram != nil {
		if *filter.IsProgram {
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"parent_id IS NOT": nil})
		} else {
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"parent_id IS": nil})
		}
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
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToProgramListResponseDTO(data)

	return response, total, nil
}
