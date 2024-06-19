package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FixedDepositWillDispatchServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.FixedDepositWillDispatch
}

func NewFixedDepositWillDispatchServiceImpl(app *celeritas.Celeritas, repo data.FixedDepositWillDispatch) FixedDepositWillDispatchService {
	return &FixedDepositWillDispatchServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FixedDepositWillDispatchServiceImpl) CreateFixedDepositWillDispatch(ctx context.Context, input dto.FixedDepositWillDispatchDTO) (*dto.FixedDepositWillDispatchResponseDTO, error) {
	dataToInsert := input.ToFixedDepositWillDispatch()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo fixed deposit will dispatch insert")
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit will dispatch get")
	}

	res := dto.ToFixedDepositWillDispatchResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *FixedDepositWillDispatchServiceImpl) UpdateFixedDepositWillDispatch(ctx context.Context, id int, input dto.FixedDepositWillDispatchDTO) (*dto.FixedDepositWillDispatchResponseDTO, error) {
	dataToInsert := input.ToFixedDepositWillDispatch()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo fixed deposit will dispatch update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit will dispatch get")
	}

	response := dto.ToFixedDepositWillDispatchResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *FixedDepositWillDispatchServiceImpl) DeleteFixedDepositWillDispatch(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo fixed deposit will dispatch delete")
	}

	return nil
}

func (h *FixedDepositWillDispatchServiceImpl) GetFixedDepositWillDispatch(id int) (*dto.FixedDepositWillDispatchResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit will dispatch get")
	}

	response := dto.ToFixedDepositWillDispatchResponseDTO(*data)

	return &response, nil
}

func (h *FixedDepositWillDispatchServiceImpl) GetFixedDepositWillDispatchList(filter dto.FixedDepositWillDispatchFilterDTO) ([]dto.FixedDepositWillDispatchResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.WillID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"will_id": *filter.WillID})
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
		return nil, nil, newErrors.Wrap(err, "repo fixed deposit will dispatch get all")
	}
	response := dto.ToFixedDepositWillDispatchListResponseDTO(data)

	return response, total, nil
}
