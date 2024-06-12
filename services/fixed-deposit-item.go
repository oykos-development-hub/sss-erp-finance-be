package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FixedDepositItemServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.FixedDepositItem
}

func NewFixedDepositItemServiceImpl(app *celeritas.Celeritas, repo data.FixedDepositItem) FixedDepositItemService {
	return &FixedDepositItemServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FixedDepositItemServiceImpl) CreateFixedDepositItem(ctx context.Context, input dto.FixedDepositItemDTO) (*dto.FixedDepositItemResponseDTO, error) {
	dataToInsert := input.ToFixedDepositItem()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
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

	res := dto.ToFixedDepositItemResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *FixedDepositItemServiceImpl) UpdateFixedDepositItem(ctx context.Context, id int, input dto.FixedDepositItemDTO) (*dto.FixedDepositItemResponseDTO, error) {
	dataToInsert := input.ToFixedDepositItem()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
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

	response := dto.ToFixedDepositItemResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *FixedDepositItemServiceImpl) DeleteFixedDepositItem(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *FixedDepositItemServiceImpl) GetFixedDepositItem(id int) (*dto.FixedDepositItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFixedDepositItemResponseDTO(*data)

	return &response, nil
}

func (h *FixedDepositItemServiceImpl) GetFixedDepositItemList(filter dto.FixedDepositItemFilterDTO) ([]dto.FixedDepositItemResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.DepositID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"deposit_id": *filter.DepositID})
	}

	/*if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}*/

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToFixedDepositItemListResponseDTO(data)

	return response, total, nil
}
