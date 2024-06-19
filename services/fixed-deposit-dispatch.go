package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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

func (h *FixedDepositDispatchServiceImpl) CreateFixedDepositDispatch(ctx context.Context, input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error) {
	dataToInsert := input.ToFixedDepositDispatch()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo fixed deposit dispatch insert")
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit dispatch get")
	}

	res := dto.ToFixedDepositDispatchResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *FixedDepositDispatchServiceImpl) UpdateFixedDepositDispatch(ctx context.Context, id int, input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error) {
	dataToInsert := input.ToFixedDepositDispatch()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo fixed deposit dispatch update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit dispatch get")
	}

	response := dto.ToFixedDepositDispatchResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *FixedDepositDispatchServiceImpl) DeleteFixedDepositDispatch(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo fixed deposit dispatch delete")
	}

	return nil
}

func (h *FixedDepositDispatchServiceImpl) GetFixedDepositDispatch(id int) (*dto.FixedDepositDispatchResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit dispatch get")
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
		return nil, nil, newErrors.Wrap(err, "repo fixed deposit dispatch get all")
	}
	response := dto.ToFixedDepositDispatchListResponseDTO(data)

	return response, total, nil
}
