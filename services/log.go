package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type LogServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.Log
}

func NewLogServiceImpl(app *celeritas.Celeritas, repo data.Log) LogService {
	return &LogServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *LogServiceImpl) CreateLog(input dto.LogDTO) (*dto.LogResponseDTO, error) {
	data := input.ToLog()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToLogResponseDTO(*data)

	return &res, nil
}

func (h *LogServiceImpl) UpdateLog(id int, input dto.LogDTO) (*dto.LogResponseDTO, error) {
	data := input.ToLog()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToLogResponseDTO(*data)

	return &response, nil
}

func (h *LogServiceImpl) DeleteLog(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *LogServiceImpl) GetLog(id int) (*dto.LogResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToLogResponseDTO(*data)

	return &response, nil
}

func (h *LogServiceImpl) GetLogList(filter dto.LogFilterDTO) ([]dto.LogResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.Entity != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"entity": *filter.Entity})
	}

	if filter.Operation != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"operation": *filter.Operation})
	}

	if filter.UserID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"user_id": *filter.UserID})
	}

	if filter.ItemID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"item_id": *filter.ItemID})
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
	response := dto.ToLogListResponseDTO(data)

	return response, total, nil
}
