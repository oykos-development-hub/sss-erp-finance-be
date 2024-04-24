package services

import (
	"fmt"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type EnforcedPaymentServiceImpl struct {
	App                    *celeritas.Celeritas
	repo                   data.EnforcedPayment
	itemsRepo              data.EnforcedPaymentItem
	invoicesRepo           data.Invoice
	invoiceArticlesRepo    data.Article
	additionalExpensesRepo data.AdditionalExpense
}

func NewEnforcedPaymentServiceImpl(app *celeritas.Celeritas, repo data.EnforcedPayment, itemsRepo data.EnforcedPaymentItem, invoicesRepo data.Invoice, invoiceArticlesRepo data.Article, additionalExpensesRepo data.AdditionalExpense) EnforcedPaymentService {
	return &EnforcedPaymentServiceImpl{
		App:                    app,
		repo:                   repo,
		itemsRepo:              itemsRepo,
		invoicesRepo:           invoicesRepo,
		invoiceArticlesRepo:    invoiceArticlesRepo,
		additionalExpensesRepo: additionalExpensesRepo,
	}
}

func (h *EnforcedPaymentServiceImpl) CreateEnforcedPayment(input dto.EnforcedPaymentDTO) (*dto.EnforcedPaymentResponseDTO, error) {
	dataToInsert := input.ToEnforcedPayment()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			fmt.Println(err)
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToEnforcedPaymentItem()
			itemToInsert.PaymentOrderID = id
			_, err = h.itemsRepo.Insert(tx, *itemToInsert)
			if err != nil {
				return err
			}

			if item.InvoiceID != nil {
				err = updateInvoiceStatusForEnforcedPayment(*item.InvoiceID, dataToInsert.Amount, len(input.Items), tx, h)

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

	res := dto.ToEnforcedPaymentResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *EnforcedPaymentServiceImpl) UpdateEnforcedPayment(id int, input dto.EnforcedPaymentDTO) (*dto.EnforcedPaymentResponseDTO, error) {
	dataToInsert := input.ToEnforcedPayment()
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

	response := dto.ToEnforcedPaymentResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *EnforcedPaymentServiceImpl) ReturnEnforcedPayment(id int, input dto.EnforcedPaymentDTO) error {
	dataToInsert := input.ToEnforcedPayment()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.ReturnEnforcedPayment(tx, *dataToInsert)
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

func (h *EnforcedPaymentServiceImpl) DeleteEnforcedPayment(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *EnforcedPaymentServiceImpl) GetEnforcedPayment(id int) (*dto.EnforcedPaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToEnforcedPaymentResponseDTO(*data)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, err
	}

	for _, item := range items {

		if item.InvoiceID != nil && *item.InvoiceID != 0 {
			invoice, err := h.invoicesRepo.Get(*item.InvoiceID)

			if err != nil {
				return nil, err
			}

			if invoice.Type == "invoice" {
				conditionAndExp := &up.AndExpr{}
				conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": item.InvoiceID})
				articles, _, err := h.invoiceArticlesRepo.GetAll(nil, nil, conditionAndExp, nil)
				if err != nil {
					return nil, err
				}
				for _, article := range articles {
					price := article.NetPrice + article.NetPrice*float64(article.VatPercentage)/100
					item.Amount += float32(price)
				}
				item.Title = "Faktura broj " + invoice.InvoiceNumber
			} else {
				conditionAndExp := &up.AndExpr{}
				conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": item.InvoiceID})
				articles, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)
				if err != nil {
					return nil, err
				}
				for _, article := range articles {
					item.Amount += float32(article.Price)
				}
				item.Amount += float32(invoice.GrossPrice)
				if invoice.Type == "decisions" {
					item.Title = "Rješenje broj " + invoice.InvoiceNumber
				} else {
					item.Title = "Ugovor broj " + invoice.InvoiceNumber
				}
			}
		}

		response.Items = append(response.Items, dto.EnforcedPaymentItemResponseDTO{
			ID:             item.ID,
			PaymentOrderID: item.PaymentOrderID,
			Title:          item.Title,
			Amount:         item.Amount,
			InvoiceID:      item.InvoiceID,
			AccountID:      item.AccountID,
		})
	}

	return &response, nil
}

func (h *EnforcedPaymentServiceImpl) GetEnforcedPaymentList(filter dto.EnforcedPaymentFilterDTO) ([]dto.EnforcedPaymentResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.Year != nil {
		year := *filter.Year
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_order": up.Between(startOfYear, endOfYear)})
	}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.SupplierID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"supplier_id": *filter.SupplierID})
	}

	if filter.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
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
	response := dto.ToEnforcedPaymentListResponseDTO(data)

	return response, total, nil
}

func updateInvoiceStatusForEnforcedPayment(id int, amount float64, lenOfArray int, tx up.Session, h *EnforcedPaymentServiceImpl) error {
	invoice, err := h.invoicesRepo.Get(id)

	if err != nil {
		return err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})

	additionalExpenses, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return err
	}

	for _, item := range additionalExpenses {
		item.Status = data.InvoiceStatusFull
		err = h.additionalExpensesRepo.Update(tx, *item)
		if err != nil {
			return err
		}
	}

	invoice.Status = data.InvoiceStatusFull

	err = h.invoicesRepo.Update(tx, *invoice)

	if err != nil {
		return err
	}

	return nil
}
