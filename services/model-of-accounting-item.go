package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ModelOfAccountingItemServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.ModelOfAccountingItem
}

func NewModelOfAccountingItemServiceImpl(app *celeritas.Celeritas, repo data.ModelOfAccountingItem) ModelOfAccountingItemService {
	return &ModelOfAccountingItemServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *ModelOfAccountingItemServiceImpl) CreateModelOfAccountingItem(input dto.ModelOfAccountingItemDTO) (*dto.ModelOfAccountingItemResponseDTO, error) {
	dataToInsert := input.ToModelOfAccountingItem()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo model of accounting item insert")
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting item get")
	}

	res := dto.ToModelOfAccountingItemResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ModelOfAccountingItemServiceImpl) UpdateModelOfAccountingItem(id int, input dto.ModelOfAccountingItemDTO) (*dto.ModelOfAccountingItemResponseDTO, error) {
	dataToInsert := input.ToModelOfAccountingItem()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo model of accounting item update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting item get")
	}

	response := dto.ToModelOfAccountingItemResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *ModelOfAccountingItemServiceImpl) GetModelOfAccountingItem(id int) (*dto.ModelOfAccountingItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting item get")
	}

	response := dto.ToModelOfAccountingItemResponseDTO(*data)

	return &response, nil
}

func (h *ModelOfAccountingItemServiceImpl) GetModelOfAccountingItemList(filter dto.ModelOfAccountingItemFilterDTO) ([]dto.ModelOfAccountingItemResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	// if filter.Year != nil {
	// 	conditionAndExp = up.And(conditionAndExp, &up.Cond{"year": *filter.Year})
	// }

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
		return nil, nil, newErrors.Wrap(err, "repo model of accounting item get all")
	}
	response := dto.ToModelOfAccountingItemListResponseDTO(data)

	return response, total, nil
}
