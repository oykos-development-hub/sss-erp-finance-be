package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type PaymentOrderItemServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.PaymentOrderItem
}

func NewPaymentOrderItemServiceImpl(app *celeritas.Celeritas, repo data.PaymentOrderItem) PaymentOrderItemService {
	return &PaymentOrderItemServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *PaymentOrderItemServiceImpl) CreatePaymentOrderItem(input dto.PaymentOrderItemDTO) (*dto.PaymentOrderItemResponseDTO, error) {
	dataToInsert := input.ToPaymentOrderItem()

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

	res := dto.ToPaymentOrderItemResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *PaymentOrderItemServiceImpl) UpdatePaymentOrderItem(id int, input dto.PaymentOrderItemDTO) (*dto.PaymentOrderItemResponseDTO, error) {
	dataToInsert := input.ToPaymentOrderItem()
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

	response := dto.ToPaymentOrderItemResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *PaymentOrderItemServiceImpl) DeletePaymentOrderItem(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *PaymentOrderItemServiceImpl) GetPaymentOrderItem(id int) (*dto.PaymentOrderItemResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToPaymentOrderItemResponseDTO(*data)

	return &response, nil
}

func (h *PaymentOrderItemServiceImpl) GetPaymentOrderItemList(filter dto.PaymentOrderItemFilterDTO) ([]dto.PaymentOrderItemResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.InvoiceID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": *filter.InvoiceID})
	}

	if filter.AdditionalExpenseID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"additional_expense": *filter.AdditionalExpenseID})
	}

	if filter.SalaryAdditionalExpenseID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_additional_expense_id": *filter.SalaryAdditionalExpenseID})
	}

	if filter.PaymentOrderID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": *filter.PaymentOrderID})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToPaymentOrderItemListResponseDTO(data)

	return response, total, nil
}
