package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FixedDepositDispatchServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.FixedDepositDispatch
}

func NewFixedDepositDispatchServiceImpl(app *celeritas.Celeritas, repo data.FixedDepositDispatch) FixedDepositDispatchService {
	return &FixedDepositDispatchServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FixedDepositDispatchServiceImpl) CreateFixedDepositDispatch(input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error) {
	dataToInsert := input.ToFixedDepositDispatch()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
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

	res := dto.ToFixedDepositDispatchResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *FixedDepositDispatchServiceImpl) UpdateFixedDepositDispatch(id int, input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error) {
	dataToInsert := input.ToFixedDepositDispatch()
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

	response := dto.ToFixedDepositDispatchResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *FixedDepositDispatchServiceImpl) DeleteFixedDepositDispatch(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *FixedDepositDispatchServiceImpl) GetFixedDepositDispatch(id int) (*dto.FixedDepositDispatchResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFixedDepositDispatchResponseDTO(*data)

	return &response, nil
}

func (h *FixedDepositDispatchServiceImpl) GetFixedDepositDispatchList(filter dto.FixedDepositDispatchFilterDTO) ([]dto.FixedDepositDispatchResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.DepositID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"deposit_id": *filter.DepositID})
	}
	/*
		if filter.SortByTitle != nil {
			if *filter.SortByTitle == "asc" {
				orders = append(orders, "-title")
			} else {
				orders = append(orders, "title")
			}
		}
	*/
	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToFixedDepositDispatchListResponseDTO(data)

	return response, total, nil
}
