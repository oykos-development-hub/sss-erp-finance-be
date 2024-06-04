package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ExternalReallocationItemServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.ExternalReallocationItem
}

func NewExternalReallocationItemServiceImpl(app *celeritas.Celeritas, repo data.ExternalReallocationItem) ExternalReallocationItemService {
	return &ExternalReallocationItemServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *ExternalReallocationItemServiceImpl) CreateExternalReallocationItem(input dto.ExternalReallocationItemDTO) (*dto.ExternalReallocationItemResponseDTO, error) {
	dataToInsert := input.ToExternalReallocationItem()

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

	res := dto.ToExternalReallocationItemResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ExternalReallocationItemServiceImpl) DeleteExternalReallocationItem(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *ExternalReallocationItemServiceImpl) GetExternalReallocationItem(id int) (*dto.ExternalReallocationItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToExternalReallocationItemResponseDTO(*data)

	return &response, nil
}

func (h *ExternalReallocationItemServiceImpl) GetExternalReallocationItemList(filter dto.ExternalReallocationItemFilterDTO) ([]dto.ExternalReallocationItemResponseDTO, *uint64, error) {
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
	response := dto.ToExternalReallocationItemListResponseDTO(data)

	return response, total, nil
}
