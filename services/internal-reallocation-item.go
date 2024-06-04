package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type InternalReallocationItemServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.InternalReallocationItem
}

func NewInternalReallocationItemServiceImpl(app *celeritas.Celeritas, repo data.InternalReallocationItem) InternalReallocationItemService {
	return &InternalReallocationItemServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *InternalReallocationItemServiceImpl) CreateInternalReallocationItem(input dto.InternalReallocationItemDTO) (*dto.InternalReallocationItemResponseDTO, error) {
	dataToInsert := input.ToInternalReallocationItem()

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

	res := dto.ToInternalReallocationItemResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *InternalReallocationItemServiceImpl) DeleteInternalReallocationItem(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *InternalReallocationItemServiceImpl) GetInternalReallocationItem(id int) (*dto.InternalReallocationItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToInternalReallocationItemResponseDTO(*data)

	return &response, nil
}

func (h *InternalReallocationItemServiceImpl) GetInternalReallocationItemList(filter dto.InternalReallocationItemFilterDTO) ([]dto.InternalReallocationItemResponseDTO, *uint64, error) {
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
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToInternalReallocationItemListResponseDTO(data)

	return response, total, nil
}
