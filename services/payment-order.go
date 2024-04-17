package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type PaymentOrderServiceImpl struct {
	App                          *celeritas.Celeritas
	repo                         data.PaymentOrder
	itemsRepo                    data.PaymentOrderItem
	invoiceRepo                  data.Invoice
	invoiceArticlesRepo          data.Article
	additionalExpensesRepo       data.AdditionalExpense
	salaryAdditionalExpensesRepo data.SalaryAdditionalExpense
}

func NewPaymentOrderServiceImpl(app *celeritas.Celeritas, repo data.PaymentOrder, itemsRepo data.PaymentOrderItem, invoiceRepo data.Invoice, invoiceArticleRepo data.Article, additionalExpensesRepo data.AdditionalExpense, salaryAdditionalExpensesRepo data.SalaryAdditionalExpense) PaymentOrderService {
	return &PaymentOrderServiceImpl{
		App:                          app,
		repo:                         repo,
		itemsRepo:                    itemsRepo,
		invoiceRepo:                  invoiceRepo,
		invoiceArticlesRepo:          invoiceArticleRepo,
		additionalExpensesRepo:       additionalExpensesRepo,
		salaryAdditionalExpensesRepo: salaryAdditionalExpensesRepo,
	}
}

func (h *PaymentOrderServiceImpl) CreatePaymentOrder(input dto.PaymentOrderDTO) (*dto.PaymentOrderResponseDTO, error) {
	dataToInsert := input.ToPaymentOrder()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			fmt.Println(err)
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToPaymentOrderItem()
			itemToInsert.PaymentOrderID = id
			_, err = h.itemsRepo.Insert(tx, *itemToInsert)
			if err != nil {
				return err
			}

			if item.InvoiceID != nil {
				err = updateInvoiceStatus(*item.InvoiceID, dataToInsert.Amount, len(input.Items), tx, h)

				if err != nil {
					return err
				}
			} else if item.AdditionalExpenseID != nil {
				err = updateAdditionalExpenseStatus(*item.AdditionalExpenseID, dataToInsert.Amount, len(input.Items), tx, h)

				if err != nil {
					return err
				}
			} else if item.SalaryAdditionalExpenseID != nil {
				err = updateSalaryAdditionalExpenseStatus(*item.SalaryAdditionalExpenseID, dataToInsert.Amount, len(input.Items), tx, h)

				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, errors.ErrInternalServer
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		fmt.Println(err)
		return nil, errors.ErrInternalServer
	}

	res := dto.ToPaymentOrderResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *PaymentOrderServiceImpl) UpdatePaymentOrder(id int, input dto.PaymentOrderDTO) (*dto.PaymentOrderResponseDTO, error) {
	dataToInsert := input.ToPaymentOrder()
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

	response := dto.ToPaymentOrderResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *PaymentOrderServiceImpl) DeletePaymentOrder(id int) error {

	input, err := h.GetPaymentOrder(id)

	if err != nil {
		return err
	}

	for _, item := range input.Items {
		if item.InvoiceID != nil {
			err = updateInvoiceStatusOnDelete(*item.InvoiceID, input.Amount, len(input.Items), data.Upper, h)

			if err != nil {
				h.App.ErrorLog.Println(err)
				return err
			}
		} else if item.AdditionalExpenseID != nil {
			err = updateAdditionalExpenseStatusOnDelete(*item.AdditionalExpenseID, input.Amount, len(input.Items), data.Upper, h)

			if err != nil {
				h.App.ErrorLog.Println(err)
				return err
			}
		} else if item.SalaryAdditionalExpenseID != nil {
			err = updateSalaryAdditionalExpenseStatusOnDelete(*item.SalaryAdditionalExpenseID, input.Amount, len(input.Items), data.Upper, h)

			if err != nil {
				h.App.ErrorLog.Println(err)
				return err
			}
		}
	}

	err = h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *PaymentOrderServiceImpl) GetPaymentOrder(id int) (*dto.PaymentOrderResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToPaymentOrderResponseDTO(*data)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, err
	}

	for _, item := range items {
		builtItem := dto.PaymentOrderItemResponseDTO{
			ID:                        item.ID,
			PaymentOrderID:            item.PaymentOrderID,
			InvoiceID:                 item.InvoiceID,
			AdditionalExpenseID:       item.AdditionalExpenseID,
			SalaryAdditionalExpenseID: item.SalaryAdditionalExpenseID,
			AccountID:                 item.AccountID,
			CreatedAt:                 item.CreatedAt,
			UpdatedAt:                 item.UpdatedAt,
		}

		if item.InvoiceID != nil {
			builtItem.Type = "invoice"
		} else if item.AdditionalExpenseID != nil {
			builtItem.Type = "additional_expense"
		} else {
			builtItem.Type = "salary_additional_expense"
		}

		response.Items = append(response.Items, builtItem)
	}

	return &response, nil
}

func (h *PaymentOrderServiceImpl) GetPaymentOrderList(filter dto.PaymentOrderFilterDTO) ([]dto.PaymentOrderResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.SupplierID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"supplier_id": *filter.SupplierID})
	}

	if filter.Status != nil {
		switch *filter.Status {
		case "Plaćen":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"id_of_statement is not ": nil})
		case "Na čekanju":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"id_of_statement is ": nil})
		}
	}

	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"payer ILIKE": likeCondition},
			up.Cond{"case_number ILIKE": likeCondition},
			up.Cond{"party_name ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	//if filter.SortByTitle != nil {
	//	if *filter.SortByTitle == "asc" {
	//		orders = append(orders, "-title")
	//	} else {
	//		orders = append(orders, "title")
	//	}
	//}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToPaymentOrderListResponseDTO(data)

	return response, total, nil
}

func (h *PaymentOrderServiceImpl) GetAllObligations(filter dto.GetObligationsFilterDTO) ([]dto.ObligationResponse, *uint64, error) {

	dataFilter := data.ObligationsFilter{
		Page:               filter.Page,
		Size:               filter.Size,
		OrganizationUnitID: filter.OrganizationUnitID,
		SupplierID:         filter.SupplierID,
		Type:               filter.Type,
	}

	items, total, err := h.repo.GetAllObligations(dataFilter)

	if err != nil {
		return nil, nil, err
	}

	var response []dto.ObligationResponse
	for _, item := range items {
		response = append(response, dto.ObligationResponse{
			InvoiceID:                 item.InvoiceID,
			AdditionalExpenseID:       item.AdditionalExpenseID,
			SalaryAdditionalExpenseID: item.SalaryAdditionalExpenseID,
			Type:                      item.Type,
			Title:                     item.Title,
			Status:                    item.Status,
			Price:                     item.Price,
			CreatedAt:                 item.CreatedAt,
		})
	}

	return response, total, nil
}

func (h *PaymentOrderServiceImpl) PayPaymentOrder(id int, input dto.PaymentOrderDTO) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.PayPaymentOrder(tx, id, *input.SAPID, *input.DateOfSAP)
		if err != nil {
			return errors.ErrInternalServer
		}
		return nil
	})

	if err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func updateInvoiceStatus(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	invoice, err := h.invoiceRepo.Get(id)

	if err != nil {
		return err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})

	articles, _, err := h.invoiceArticlesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	var price float64
	for _, article := range articles {
		price += float64(article.NetPrice)
	}

	conditionAndExp = &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	for _, item := range items {
		paymentOrder, err := h.repo.Get(item.PaymentOrderID)

		if err != nil {
			return err
		}

		amount += paymentOrder.Amount
	}

	if amount >= price || lenOfArray > 1 {
		invoice.Status = data.InvoiceStatusFull
	} else {
		invoice.Status = data.InvoiceStatusPart
	}

	err = h.invoiceRepo.Update(tx, *invoice)

	if err != nil {
		return err
	}

	return nil
}

func updateAdditionalExpenseStatus(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.additionalExpensesRepo.Get(id)

	if err != nil {
		return err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"additional_expense_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	for _, item := range items {
		paymentOrder, err := h.repo.Get(item.PaymentOrderID)

		if err != nil {
			return err
		}

		amount += paymentOrder.Amount
	}

	if amount >= float64(item.Price) || lenOfArray > 1 {
		item.Status = data.InvoiceStatusFull
	} else {
		item.Status = data.InvoiceStatusPart
	}

	err = h.additionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return err
	}

	return nil
}

func updateSalaryAdditionalExpenseStatus(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.salaryAdditionalExpensesRepo.Get(id)

	if err != nil {
		return err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_additional_expense_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	for _, item := range items {
		paymentOrder, err := h.repo.Get(item.PaymentOrderID)

		if err != nil {
			return err
		}

		amount += paymentOrder.Amount
	}

	if amount >= item.Amount || lenOfArray > 1 {
		item.Status = data.InvoiceStatusFull
	} else {
		item.Status = data.InvoiceStatusPart
	}

	err = h.salaryAdditionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return err
	}

	return nil
}

func updateInvoiceStatusOnDelete(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	invoice, err := h.invoiceRepo.Get(id)

	if err != nil {
		return err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	//ako smo dobili nazad samo item koji brisemo
	if len(items) == 1 || lenOfArray > 1 {
		invoice.Status = data.InvoiceStatusCreated
	} else {
		invoice.Status = data.InvoiceStatusPart
	}

	err = h.invoiceRepo.Update(tx, *invoice)

	if err != nil {
		return err
	}

	return nil
}

func updateAdditionalExpenseStatusOnDelete(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.additionalExpensesRepo.Get(id)

	if err != nil {
		return err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"additional_expense_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	if len(items) == 1 || lenOfArray > 1 {
		item.Status = data.InvoiceStatusCreated
	} else {
		item.Status = data.InvoiceStatusPart
	}

	err = h.additionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return err
	}

	return nil
}

func updateSalaryAdditionalExpenseStatusOnDelete(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.salaryAdditionalExpensesRepo.Get(id)

	if err != nil {
		return err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_additional_expense_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	if len(items) == 1 || lenOfArray > 1 {
		item.Status = data.InvoiceStatusCreated
	} else {
		item.Status = data.InvoiceStatusPart
	}

	err = h.salaryAdditionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return err
	}

	return nil
}
