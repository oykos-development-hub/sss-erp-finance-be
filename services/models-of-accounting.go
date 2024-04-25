package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ModelsOfAccountingServiceImpl struct {
	App       *celeritas.Celeritas
	repo      data.ModelsOfAccounting
	itemsRepo data.ModelOfAccountingItem
}

func NewModelsOfAccountingServiceImpl(app *celeritas.Celeritas, repo data.ModelsOfAccounting, itemsRepo data.ModelOfAccountingItem) ModelsOfAccountingService {
	return &ModelsOfAccountingServiceImpl{
		App:       app,
		repo:      repo,
		itemsRepo: itemsRepo,
	}
}

func (h *ModelsOfAccountingServiceImpl) CreateModelsOfAccounting(input dto.ModelsOfAccountingDTO) (*dto.ModelsOfAccountingResponseDTO, error) {
	dataToInsert := input.ToModelsOfAccounting()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToModelOfAccountingItem()
			itemToInsert.ModelID = id
			_, err = h.itemsRepo.Insert(tx, *itemToInsert)

			if err != nil {
				return err
			}
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

	res := dto.ToModelsOfAccountingResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ModelsOfAccountingServiceImpl) UpdateModelsOfAccounting(id int, input dto.ModelsOfAccountingDTO) (*dto.ModelsOfAccountingResponseDTO, error) {
	dataToInsert := input.ToModelsOfAccounting()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToModelOfAccountingItem()
			itemToInsert.ModelID = id
			err = h.itemsRepo.Update(tx, *itemToInsert)

			if err != nil {
				return err
			}
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

	response := dto.ToModelsOfAccountingResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *ModelsOfAccountingServiceImpl) GetModelsOfAccounting(id int) (*dto.ModelsOfAccountingResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToModelsOfAccountingResponseDTO(*data)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"model_id": &id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	for _, item := range items {
		builtItem := dto.ModelOfAccountingItemResponseDTO{
			ID:              item.ID,
			Title:           item.Title,
			DebitAccountID:  item.DebitAccountID,
			CreditAccountID: item.CreditAccountID,
		}
		response.Items = append(response.Items, builtItem)
	}

	return &response, nil
}

func (h *ModelsOfAccountingServiceImpl) GetModelsOfAccountingList(filter dto.ModelsOfAccountingFilterDTO) ([]dto.ModelsOfAccountingResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.Type != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"type": *filter.Type})
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
	response := dto.ToModelsOfAccountingListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"model_id": &response[i].ID})

		items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return nil, nil, errors.ErrInternalServer
		}

		for _, item := range items {
			builtItem := dto.ModelOfAccountingItemResponseDTO{
				ID:              item.ID,
				Title:           item.Title,
				DebitAccountID:  item.DebitAccountID,
				CreditAccountID: item.CreditAccountID,
			}
			response[i].Items = append(response[i].Items, builtItem)
		}
	}

	return response, total, nil
}
