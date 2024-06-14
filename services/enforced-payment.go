package services

import (
	"context"
	"fmt"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type EnforcedPaymentServiceImpl struct {
	App                    *celeritas.Celeritas
	repo                   data.EnforcedPayment
	currentBudget          CurrentBudgetService
	itemsRepo              data.EnforcedPaymentItem
	invoicesRepo           data.Invoice
	invoiceArticlesRepo    data.Article
	additionalExpensesRepo data.AdditionalExpense
}

func NewEnforcedPaymentServiceImpl(app *celeritas.Celeritas, repo data.EnforcedPayment, currentBudget CurrentBudgetService, itemsRepo data.EnforcedPaymentItem, invoicesRepo data.Invoice, invoiceArticlesRepo data.Article, additionalExpensesRepo data.AdditionalExpense) EnforcedPaymentService {
	return &EnforcedPaymentServiceImpl{
		App:                    app,
		repo:                   repo,
		currentBudget:          currentBudget,
		itemsRepo:              itemsRepo,
		invoicesRepo:           invoicesRepo,
		invoiceArticlesRepo:    invoiceArticlesRepo,
		additionalExpensesRepo: additionalExpensesRepo,
	}
}

func (h *EnforcedPaymentServiceImpl) CreateEnforcedPayment(ctx context.Context, input dto.EnforcedPaymentDTO) (*dto.EnforcedPaymentResponseDTO, error) {
	dataToInsert := input.ToEnforcedPayment()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
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
				err = updateInvoiceStatusForEnforcedPayment(ctx, *item.InvoiceID, tx, h)

				if err != nil {
					return err
				}
			}
		}

		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			currentBudget, _, err := h.currentBudget.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
				UnitID:    &input.OrganizationUnitID,
				AccountID: &item.AccountID,
			})

			if err != nil {
				return err
			}

			if len(currentBudget) > 0 {
				currentAmount := currentBudget[0].Balance.Sub(decimal.NewFromFloat32(float32(item.Amount)))
				if currentAmount.LessThan(decimal.NewFromInt(0)) {
					return errors.ErrInsufficientFunds
				} else {
					err = h.currentBudget.UpdateBalance(ctx, tx, currentBudget[0].ID, currentAmount)
					if err != nil {
						return err
					}
				}
			} else {
				return errors.ErrInsufficientFunds
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		fmt.Println(err)
		return nil, errors.ErrInternalServer
	}

	res := dto.ToEnforcedPaymentResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *EnforcedPaymentServiceImpl) UpdateEnforcedPayment(ctx context.Context, id int, input dto.EnforcedPaymentDTO) (*dto.EnforcedPaymentResponseDTO, error) {
	dataToInsert := input.ToEnforcedPayment()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
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

func (h *EnforcedPaymentServiceImpl) ReturnEnforcedPayment(ctx context.Context, id int, input dto.EnforcedPaymentDTO) error {
	dataToInsert := input.ToEnforcedPayment()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.ReturnEnforcedPayment(ctx, tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		paymentOrder, err := h.GetEnforcedPayment(id)

		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range paymentOrder.Items {
			currentBudget, _, err := h.currentBudget.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
				UnitID:    &paymentOrder.OrganizationUnitID,
				AccountID: &item.AccountID,
			})

			if err != nil {
				return err
			}

			if len(currentBudget) > 0 {
				currentAmount := currentBudget[0].Balance.Add(decimal.NewFromFloat32(float32(item.Amount)))

				err = h.currentBudget.UpdateBalance(ctx, tx, currentBudget[0].ID, currentAmount)
				if err != nil {
					return err
				}

			} else {
				return errors.ErrNotFound
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *EnforcedPaymentServiceImpl) DeleteEnforcedPayment(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *EnforcedPaymentServiceImpl) GetEnforcedPayment(id int) (*dto.EnforcedPaymentResponseDTO, error) {
	paymentData, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToEnforcedPaymentResponseDTO(*paymentData)

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

			if invoice.Type == data.TypeInvoice {
				conditionAndExp := &up.AndExpr{}
				conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": item.InvoiceID})
				articles, _, err := h.invoiceArticlesRepo.GetAll(nil, nil, conditionAndExp, nil)
				if err != nil {
					return nil, err
				}
				for _, article := range articles {
					price := (article.NetPrice + article.NetPrice*float64(article.VatPercentage)/100) * float64(article.Amount)
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
					if article.Title == "Neto" {
						item.Amount += float32(article.Price)
					}
				}

				if invoice.Type == data.TypeDecision {
					item.Title = "Rješenje broj " + invoice.InvoiceNumber
				} else {
					item.Title = "Ugovor broj " + invoice.InvoiceNumber
				}
			}
		}
		var amount float32
		if len(items) == 1 {
			amount = float32(paymentData.Amount)
		} else {
			amount = item.Amount
		}

		response.Items = append(response.Items, dto.EnforcedPaymentItemResponseDTO{
			ID:             item.ID,
			PaymentOrderID: item.PaymentOrderID,
			Title:          item.Title,
			Amount:         amount,
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

	if filter.Registred != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"registred": *filter.Registred})
	}

	if filter.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"sap_id ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
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

func updateInvoiceStatusForEnforcedPayment(ctx context.Context, id int, tx up.Session, h *EnforcedPaymentServiceImpl) error {
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
		if item.Title == "Neto" {
			item.Status = data.InvoiceStatusFull
			err = h.additionalExpensesRepo.Update(tx, *item)
			if err != nil {
				return err
			}
		}
	}

	invoice.Status = data.InvoiceStatusFull

	err = h.invoicesRepo.Update(ctx, tx, *invoice)

	if err != nil {
		return err
	}

	return nil
}
