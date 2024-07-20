package services

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ErrorLogServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.ErrorLog
}

func NewErrorLogServiceImpl(app *celeritas.Celeritas, repo data.ErrorLog) ErrorLogService {
	return &ErrorLogServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *ErrorLogServiceImpl) CreateErrorLog(err error) {
	dataToInsert := &data.ErrorLog{
		Error:     err.Error(),
		CreatedAt: time.Now(),
	}

	_ = data.Upper.Tx(func(tx up.Session) error {
		var err error
		_, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		return nil
	})

}

func (h *ErrorLogServiceImpl) UpdateErrorLog(id int, input dto.ErrorLogDTO) (*dto.ErrorLogResponseDTO, error) {
	dataToInsert := input.ToErrorLog()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}
		return nil
	})
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToErrorLogResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *ErrorLogServiceImpl) DeleteErrorLog(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *ErrorLogServiceImpl) GetErrorLog(id int) (*dto.ErrorLogResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToErrorLogResponseDTO(*data)

	return &response, nil
}

func (h *ErrorLogServiceImpl) GetErrorLogList(filter dto.ErrorLogFilterDTO) ([]dto.ErrorLogResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.Entity != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"entity": *filter.Entity})
	}

	if filter.DateOfStart != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"created_at > ": *filter.DateOfStart})
	}

	if filter.DateOfEnd != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"created_at < ": *filter.DateOfEnd})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToErrorLogListResponseDTO(data)

	return response, total, nil
}
