package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
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
			return newErrors.Wrap(err, "repo accounting entry item insert")
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo accounting entry item delete")
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
			return newErrors.Wrap(err, "repo accounting entry item update")
		}
		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo accounting entry item get")
	}

	response := dto.ToAccountingEntryItemResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *AccountingEntryItemServiceImpl) DeleteAccountingEntryItem(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		return newErrors.Wrap(err, "repo accounting entry item delete")
	}

	return nil
}

func (h *AccountingEntryItemServiceImpl) GetAccountingEntryItem(id int) (*dto.AccountingEntryItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo accounting entry item get")
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
		return nil, nil, newErrors.Wrap(err, "repo accounting entry item get all")
	}
	response := dto.ToAccountingEntryItemListResponseDTO(data)

	return response, total, nil
}
