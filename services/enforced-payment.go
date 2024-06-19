package services

import (
	"context"
	"fmt"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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
			return newErrors.Wrap(err, "repo enforced payment insert")
		}

		for _, item := range input.Items {
			itemToInsert := item.ToEnforcedPaymentItem()
			itemToInsert.PaymentOrderID = id
			_, err = h.itemsRepo.Insert(tx, *itemToInsert)
			if err != nil {
				return newErrors.Wrap(err, "repo enforced payment item insert")
			}

			if item.InvoiceID != nil {
				err = updateInvoiceStatusForEnforcedPayment(ctx, *item.InvoiceID, tx, h)

				if err != nil {
					return newErrors.Wrap(err, "update invoice status for enforced payment")
				}
			}
		}

		if err != nil {
			return newErrors.Wrap(err, "upper tx")
		}

		for _, item := range input.Items {
			currentBudget, _, err := h.currentBudget.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
				UnitID:    &input.OrganizationUnitID,
				AccountID: &item.AccountID,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo current budget get all")
			}

			if len(currentBudget) > 0 {
				amount := 0.0

				if len(input.Items) == 1 {
					amount = input.Amount
				} else {
					amount, err = h.getInvoiceAmount(*item.InvoiceID)

					if err != nil {
						return newErrors.Wrap(err, "get invoice amount")
					}
				}

				currentAmount := currentBudget[0].Balance.Sub(decimal.NewFromFloat32(float32(amount)))
				if currentAmount.LessThan(decimal.NewFromInt(0)) {
					return newErrors.Wrap(errors.ErrInsufficientFunds, "repo current budget update balance")
				} else {
					err = h.currentBudget.UpdateBalance(ctx, tx, currentBudget[0].ID, currentAmount)
					if err != nil {
						return newErrors.Wrap(err, "repo current budget update balance")
					}
				}
			} else {
				return newErrors.Wrap(errors.ErrInsufficientFunds, "repo current budget update balance")
			}

			amount := input.AmountForAgent + input.AmountForBank + input.AmountForLawyer

			currentBudget, _, err = h.currentBudget.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
				UnitID:    &input.OrganizationUnitID,
				AccountID: &input.AccountIDForExpenses,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo current budget get")
			}
			if len(currentBudget) > 0 {
				currentAmount := currentBudget[0].Balance.Sub(decimal.NewFromFloat32(float32(amount)))
				if currentAmount.LessThan(decimal.NewFromInt(0)) {
					return newErrors.Wrap(errors.ErrInsufficientFunds, "repo current budget update balance")
				} else {
					err = h.currentBudget.UpdateBalance(ctx, tx, currentBudget[0].ID, currentAmount)
					if err != nil {
						return newErrors.Wrap(errors.ErrInsufficientFunds, "repo current budget update balance")
					}
				}
			} else {
				return newErrors.Wrap(errors.ErrNotFound, "repo current budget get all")
			}

		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(errors.ErrInsufficientFunds, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(errors.ErrInsufficientFunds, "repo enforced payment get")
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
			return newErrors.Wrap(errors.ErrInsufficientFunds, "repo enforced payment update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment get")
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
			return newErrors.Wrap(err, "repo enforced payment return enforced payment")
		}

		paymentOrder, err := h.GetEnforcedPayment(id)

		if err != nil {
			return newErrors.Wrap(err, "get enforced payment")
		}

		for _, item := range paymentOrder.Items {
			currentBudget, _, err := h.currentBudget.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
				UnitID:    &paymentOrder.OrganizationUnitID,
				AccountID: &item.AccountID,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo current budget get all")
			}

			if len(currentBudget) > 0 {
				amount := 0.0

				if len(paymentOrder.Items) == 1 {
					amount = *input.ReturnAmount
				} else {
					amount, err = h.getInvoiceAmount(*item.InvoiceID)

					if err != nil {
						return newErrors.Wrap(err, "get invoice amount")
					}
				}

				currentAmount := currentBudget[0].Balance.Add(decimal.NewFromFloat32(float32(amount)))

				err = h.currentBudget.UpdateBalance(ctx, tx, currentBudget[0].ID, currentAmount)
				if err != nil {
					return newErrors.Wrap(err, "repo current budget update balance")

				}
			} else {
				return newErrors.Wrap(errors.ErrNotFound, "repo current budget get all")
			}

		}

		return nil
	})
	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}

func (h *EnforcedPaymentServiceImpl) DeleteEnforcedPayment(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo enforced payment delete")
	}

	return nil
}

func (h *EnforcedPaymentServiceImpl) GetEnforcedPayment(id int) (*dto.EnforcedPaymentResponseDTO, error) {
	paymentData, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment get")
	}

	response := dto.ToEnforcedPaymentResponseDTO(*paymentData)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment items get all")
	}

	for _, item := range items {

		if item.InvoiceID != nil && *item.InvoiceID != 0 {
			invoice, err := h.invoicesRepo.Get(*item.InvoiceID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo invoice get")
			}

			if invoice.Type == data.TypeInvoice {
				conditionAndExp := &up.AndExpr{}
				conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": item.InvoiceID})
				articles, _, err := h.invoiceArticlesRepo.GetAll(nil, nil, conditionAndExp, nil)
				if err != nil {
					return nil, newErrors.Wrap(err, "repo article get all")
				}
				for _, article := range articles {
					price := (article.NetPrice + article.NetPrice*float64(article.VatPercentage)/100) * float64(article.Amount)
					item.Amount += float32(price)
				}
				item.Title = invoice.InvoiceNumber
			} else {
				conditionAndExp := &up.AndExpr{}
				conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": item.InvoiceID})
				articles, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)
				if err != nil {
					return nil, newErrors.Wrap(err, "repo additional expenses get all")
				}
				for _, article := range articles {
					if article.Title == "Neto" {
						item.Amount += float32(article.Price)
					}
				}

				item.Title = invoice.InvoiceNumber

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
		return nil, nil, newErrors.Wrap(err, "repo enforced payment get all")
	}
	response := dto.ToEnforcedPaymentListResponseDTO(data)

	return response, total, nil
}

func updateInvoiceStatusForEnforcedPayment(ctx context.Context, id int, tx up.Session, h *EnforcedPaymentServiceImpl) error {
	invoice, err := h.invoicesRepo.Get(id)

	if err != nil {
		return newErrors.Wrap(err, "repo invoice get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})

	additionalExpenses, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo additional expenses get")
	}

	for _, item := range additionalExpenses {
		if item.Title == "Neto" {
			item.Status = data.InvoiceStatusFull
			err = h.additionalExpensesRepo.Update(tx, *item)
			if err != nil {
				return newErrors.Wrap(err, "repo additional expenses update")
			}
		}
	}

	invoice.Status = data.InvoiceStatusFull

	err = h.invoicesRepo.Update(ctx, tx, *invoice)

	if err != nil {
		return newErrors.Wrap(err, "repo invoice update")
	}

	return nil
}

func (h *EnforcedPaymentServiceImpl) getInvoiceAmount(id int) (float64, error) {
	invoice, err := h.invoicesRepo.Get(id)

	amount := 0.0

	if err != nil {
		return 0.0, newErrors.Wrap(err, "repo invoice get")
	}

	if invoice.Type == data.TypeInvoice {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
		articles, _, err := h.invoiceArticlesRepo.GetAll(nil, nil, conditionAndExp, nil)
		if err != nil {
			return 0.0, newErrors.Wrap(err, "repo article get all")
		}

		for _, article := range articles {
			price := (article.NetPrice + article.NetPrice*float64(article.VatPercentage)/100) * float64(article.Amount)
			amount += price
		}

		return amount, err
	} else {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})

		additionalExpenses, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

		if err != nil {
			return 0.0, newErrors.Wrap(err, "repo additional expenses get all")
		}

		for _, item := range additionalExpenses {
			if item.Title == "Neto" {
				return float64(item.Price), err
			}
		}
	}

	return amount, nil
}
