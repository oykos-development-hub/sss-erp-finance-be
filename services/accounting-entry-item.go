package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type AccountingEntryItemServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.AccountingEntryItem
}

func NewAccountingEntryItemServiceImpl(app *celeritas.Celeritas, repo data.AccountingEntryItem) AccountingEntryItemService {
	return &AccountingEntryItemServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *AccountingEntryItemServiceImpl) CreateAccountingEntryItem(input dto.AccountingEntryItemDTO) (*dto.AccountingEntryItemResponseDTO, error) {
	dataToInsert := input.ToAccountingEntryItem()

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

	res := dto.ToAccountingEntryItemResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *AccountingEntryItemServiceImpl) UpdateAccountingEntryItem(id int, input dto.AccountingEntryItemDTO) (*dto.AccountingEntryItemResponseDTO, error) {
	dataToInsert := input.ToAccountingEntryItem()
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

	response := dto.ToAccountingEntryItemResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *AccountingEntryItemServiceImpl) DeleteAccountingEntryItem(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *AccountingEntryItemServiceImpl) GetAccountingEntryItem(id int) (*dto.AccountingEntryItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToAccountingEntryItemResponseDTO(*data)

	return &response, nil
}

func (h *AccountingEntryItemServiceImpl) GetAccountingEntryItemList(filter dto.AccountingEntryItemFilterDTO) ([]dto.AccountingEntryItemResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.EntryID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"entry_id": *filter.EntryID})
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
	response := dto.ToAccountingEntryItemListResponseDTO(data)

	return response, total, nil
}
