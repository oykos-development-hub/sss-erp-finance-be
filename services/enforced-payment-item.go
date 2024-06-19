package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type EnforcedPaymentItemServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.EnforcedPaymentItem
}

func NewEnforcedPaymentItemServiceImpl(app *celeritas.Celeritas, repo data.EnforcedPaymentItem) EnforcedPaymentItemService {
	return &EnforcedPaymentItemServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *EnforcedPaymentItemServiceImpl) CreateEnforcedPaymentItem(input dto.EnforcedPaymentItemDTO) (*dto.EnforcedPaymentItemResponseDTO, error) {
	dataToInsert := input.ToEnforcedPaymentItem()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo enforced payment item insert")
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment item get")
	}

	res := dto.ToEnforcedPaymentItemResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *EnforcedPaymentItemServiceImpl) UpdateEnforcedPaymentItem(id int, input dto.EnforcedPaymentItemDTO) (*dto.EnforcedPaymentItemResponseDTO, error) {
	dataToInsert := input.ToEnforcedPaymentItem()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo enforced payment item update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment item get")
	}

	response := dto.ToEnforcedPaymentItemResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *EnforcedPaymentItemServiceImpl) DeleteEnforcedPaymentItem(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		return newErrors.Wrap(err, "repo enforced payment item delete")
	}

	return nil
}

func (h *EnforcedPaymentItemServiceImpl) GetEnforcedPaymentItem(id int) (*dto.EnforcedPaymentItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment item get")
	}

	response := dto.ToEnforcedPaymentItemResponseDTO(*data)

	return &response, nil
}

func (h *EnforcedPaymentItemServiceImpl) GetEnforcedPaymentItemList(filter dto.EnforcedPaymentItemFilterDTO) ([]dto.EnforcedPaymentItemResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.PaymentOrderID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": *filter.PaymentOrderID})
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
		return nil, nil, newErrors.Wrap(err, "repo enforced payment item get all")
	}
	response := dto.ToEnforcedPaymentItemListResponseDTO(data)

	return response, total, nil
}
