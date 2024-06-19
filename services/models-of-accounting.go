package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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

func (h *ModelsOfAccountingServiceImpl) CreateModelsOfAccounting(ctx context.Context, input dto.ModelsOfAccountingDTO) (*dto.ModelsOfAccountingResponseDTO, error) {
	dataToInsert := input.ToModelsOfAccounting()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo models of accounting insert")
		}

		for _, item := range input.Items {
			itemToInsert := item.ToModelOfAccountingItem()
			itemToInsert.ModelID = id
			_, err = h.itemsRepo.Insert(tx, *itemToInsert)

			if err != nil {
				return newErrors.Wrap(err, "repo model of accounting item insert")
			}
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo models of accounting get")
	}

	res := dto.ToModelsOfAccountingResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *ModelsOfAccountingServiceImpl) UpdateModelsOfAccounting(ctx context.Context, id int, input dto.ModelsOfAccountingDTO) (*dto.ModelsOfAccountingResponseDTO, error) {
	err := data.Upper.Tx(func(tx up.Session) error {
		for _, item := range input.Items {
			itemToInsert := item.ToModelOfAccountingItem()
			itemToInsert.ModelID = id
			itemToInsert.ID = item.ID
			err := h.itemsRepo.Update(tx, *itemToInsert)

			if err != nil {
				return newErrors.Wrap(err, "repo model of accounting update")
			}
		}

		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo models of accounting get")
	}

	response := dto.ToModelsOfAccountingResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *ModelsOfAccountingServiceImpl) GetModelsOfAccounting(id int) (*dto.ModelsOfAccountingResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo models of accounting get")
	}

	response := dto.ToModelsOfAccountingResponseDTO(*data)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"model_id": &id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting item get all")
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

	if filter.Search != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"title like ": likeCondition})
	}

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
		return nil, nil, newErrors.Wrap(err, "repo models of accounting get all")
	}
	response := dto.ToModelsOfAccountingListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"model_id": &response[i].ID})

		items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)
		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo models of accounting get all")
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
