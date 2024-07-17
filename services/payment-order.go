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

type PaymentOrderServiceImpl struct {
	App                          *celeritas.Celeritas
	repo                         data.PaymentOrder
	currentBudgetService         CurrentBudgetService
	itemsRepo                    data.PaymentOrderItem
	invoiceRepo                  data.Invoice
	invoiceArticlesRepo          data.Article
	additionalExpensesRepo       data.AdditionalExpense
	salaryAdditionalExpensesRepo data.SalaryAdditionalExpense
	salariesRepo                 data.Salary
}

func NewPaymentOrderServiceImpl(app *celeritas.Celeritas, currentBudgetService CurrentBudgetService, repo data.PaymentOrder, itemsRepo data.PaymentOrderItem,
	invoiceRepo data.Invoice, invoiceArticleRepo data.Article, additionalExpensesRepo data.AdditionalExpense,
	salaryAdditionalExpensesRepo data.SalaryAdditionalExpense, salariesRepo data.Salary) PaymentOrderService {
	return &PaymentOrderServiceImpl{
		App:                          app,
		repo:                         repo,
		currentBudgetService:         currentBudgetService,
		itemsRepo:                    itemsRepo,
		invoiceRepo:                  invoiceRepo,
		invoiceArticlesRepo:          invoiceArticleRepo,
		additionalExpensesRepo:       additionalExpensesRepo,
		salaryAdditionalExpensesRepo: salaryAdditionalExpensesRepo,
		salariesRepo:                 salariesRepo,
	}
}

func (h *PaymentOrderServiceImpl) CreatePaymentOrder(ctx context.Context, input dto.PaymentOrderDTO) (*dto.PaymentOrderResponseDTO, error) {
	dataToInsert := input.ToPaymentOrder()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error

		if dataToInsert.Amount == 0 {
			amount := 0.0
			for _, item := range input.Items {
				amount += item.Amount
			}
			dataToInsert.Amount = amount
		}

		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo payment order insert")
		}

		for _, item := range input.Items {
			itemToInsert := item.ToPaymentOrderItem()
			itemToInsert.PaymentOrderID = id
			_, err = h.itemsRepo.Insert(tx, *itemToInsert)
			if err != nil {
				return newErrors.Wrap(err, "repo payment order item insert")
			}

			if item.InvoiceID != nil {
				err = updateInvoiceStatus(ctx, *item.InvoiceID, dataToInsert.Amount, len(input.Items), tx, h)

				if err != nil {
					return newErrors.Wrap(err, "update invoice status")
				}
			} else if item.AdditionalExpenseID != nil {
				err = updateAdditionalExpenseStatus(*item.AdditionalExpenseID, dataToInsert.Amount, len(input.Items), tx, h)

				if err != nil {
					return newErrors.Wrap(err, "update additional expense status")
				}
			} else if item.SalaryAdditionalExpenseID != nil {
				err = updateSalaryAdditionalExpenseStatus(*item.SalaryAdditionalExpenseID, dataToInsert.Amount, len(input.Items), tx, h)

				if err != nil {
					return newErrors.Wrap(err, "update salary status")
				}
			}
		}

		for _, item := range input.Items {
			Type := 1

			if input.SourceOfFunding == "Donacija" {
				Type = 2
			}

			currentBudget, _, err := h.currentBudgetService.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
				UnitID:    &input.OrganizationUnitID,
				AccountID: &item.AccountID,
				Type:      &Type,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo current budget get")
			}

			if len(currentBudget) > 0 {
				currentAmount := currentBudget[0].Balance.Sub(decimal.NewFromFloat32(float32(item.Amount)))
				if currentAmount.LessThan(decimal.NewFromInt(0)) {
					return newErrors.Wrap(errors.ErrInsufficientFunds, "repo current budget update balance")
				} else {
					err = h.currentBudgetService.UpdateBalance(ctx, tx, currentBudget[0].ID, currentAmount)
					if err != nil {
						return newErrors.Wrap(err, "repo current budget update balance")
					}
				}
			} else {
				return newErrors.Wrap(errors.ErrInsufficientFunds, "repo current budget update balance")
			}
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo payment order get")
	}

	res := dto.ToPaymentOrderResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *PaymentOrderServiceImpl) UpdatePaymentOrder(ctx context.Context, id int, input dto.PaymentOrderDTO) (*dto.PaymentOrderResponseDTO, error) {
	dataToInsert := input.ToPaymentOrder()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo payment order update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo payment order get")
	}

	response := dto.ToPaymentOrderResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *PaymentOrderServiceImpl) DeletePaymentOrder(ctx context.Context, id int) error {

	input, err := h.GetPaymentOrder(id)

	if err != nil {
		return newErrors.Wrap(err, "repo payment order get")
	}

	for _, item := range input.Items {

		Type := 1

		if input.SourceOfFunding == "Donacija" {
			Type = 2
		}
		currentBudget, _, err := h.currentBudgetService.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
			UnitID:    &input.OrganizationUnitID,
			AccountID: &item.AccountID,
			Type:      &Type,
		})

		if err != nil {
			return newErrors.Wrap(err, "repo current budget get all")
		}

		if len(currentBudget) > 0 {
			currentAmount := currentBudget[0].Balance.Add(decimal.NewFromFloat32(float32(item.Amount)))

			err = h.currentBudgetService.UpdateBalance(ctx, data.Upper, currentBudget[0].ID, currentAmount)
			if err != nil {
				return newErrors.Wrap(err, "repo current budget update balance")
			}

		} else {
			return newErrors.Wrap(errors.ErrNotFound, "repo current budget get all")
		}
	}

	for _, item := range input.Items {
		if item.InvoiceID != nil {
			err = updateInvoiceStatusOnDelete(ctx, *item.InvoiceID, len(input.Items), data.Upper, h)

			if err != nil {
				return newErrors.Wrap(err, "update invoice status on delete")
			}
		} else if item.AdditionalExpenseID != nil {
			err = updateAdditionalExpenseStatusOnDelete(*item.AdditionalExpenseID, len(input.Items), data.Upper, h)

			if err != nil {
				return newErrors.Wrap(err, "update additional expense status on delete")
			}
		} else if item.SalaryAdditionalExpenseID != nil {
			err = updateSalaryAdditionalExpenseStatusOnDelete(*item.SalaryAdditionalExpenseID, len(input.Items), data.Upper, h)

			if err != nil {
				return newErrors.Wrap(err, "update salary additioanl expense status on delete")
			}
		}
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo payment order delete")
	}

	return nil
}

func (h *PaymentOrderServiceImpl) GetPaymentOrder(id int) (*dto.PaymentOrderResponseDTO, error) {
	paymentData, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo payment order get")
	}

	response := dto.ToPaymentOrderResponseDTO(*paymentData)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo payment order item get all")
	}

	for _, item := range items {
		builtItem := dto.PaymentOrderItemResponseDTO{
			ID:                        item.ID,
			PaymentOrderID:            item.PaymentOrderID,
			InvoiceID:                 item.InvoiceID,
			AdditionalExpenseID:       item.AdditionalExpenseID,
			SalaryAdditionalExpenseID: item.SalaryAdditionalExpenseID,
			AccountID:                 item.AccountID,
			Amount:                    item.Amount,
			CreatedAt:                 item.CreatedAt,
			UpdatedAt:                 item.UpdatedAt,
		}

		if item.InvoiceID != nil {
			builtItem.Type = data.TypeInvoice

			item, err := h.invoiceRepo.Get(*item.InvoiceID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo invoice get")
			}

			if item.InvoiceNumber != "" {
				builtItem.Title = item.InvoiceNumber
			} else {
				builtItem.Title = item.ProFormaInvoiceNumber
			}

		} else if item.AdditionalExpenseID != nil {

			additionalItem, err := h.additionalExpensesRepo.Get(*item.AdditionalExpenseID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo additional expense get")
			}

			item, err := h.invoiceRepo.Get(additionalItem.InvoiceID)
			if err != nil {
				return nil, newErrors.Wrap(err, "repo invoice get")
			}

			builtItem.Type = item.Type

			builtItem.Title = item.InvoiceNumber

		} else if item.SalaryAdditionalExpenseID != nil {
			additionalItem, err := h.salaryAdditionalExpensesRepo.Get(*item.SalaryAdditionalExpenseID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo salary additional expense get")
			}

			item, err := h.salariesRepo.Get(additionalItem.SalaryID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo salary get")
			}

			builtItem.Type = data.TypeSalary
			builtItem.Title = "Zarada " + item.Month + " " + string(additionalItem.Title)

		}

		if len(items) == 1 {
			builtItem.Amount = response.Amount
		}

		response.Items = append(response.Items, builtItem)
	}

	return &response, nil
}

func (h *PaymentOrderServiceImpl) GetPaymentOrderByIdOfStatement(id int) (*dto.PaymentOrderResponseDTO, error) {
	paymentData, err := h.repo.GetByIdOfStatement(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo payment order get")
	}

	response := dto.ToPaymentOrderResponseDTO(*paymentData)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": id})
	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo payment order item get all")
	}

	for _, item := range items {
		builtItem := dto.PaymentOrderItemResponseDTO{
			ID:                        item.ID,
			PaymentOrderID:            item.PaymentOrderID,
			InvoiceID:                 item.InvoiceID,
			AdditionalExpenseID:       item.AdditionalExpenseID,
			SalaryAdditionalExpenseID: item.SalaryAdditionalExpenseID,
			AccountID:                 item.AccountID,
			Amount:                    item.Amount,
			CreatedAt:                 item.CreatedAt,
			UpdatedAt:                 item.UpdatedAt,
		}

		if item.InvoiceID != nil {
			builtItem.Type = data.TypeInvoice

			item, err := h.invoiceRepo.Get(*item.InvoiceID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo invoice get")
			}

			if item.InvoiceNumber != "" {
				builtItem.Title = item.InvoiceNumber
			} else {
				builtItem.Title = item.ProFormaInvoiceNumber
			}

		} else if item.AdditionalExpenseID != nil {

			additionalItem, err := h.additionalExpensesRepo.Get(*item.AdditionalExpenseID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo additional expense get")
			}

			item, err := h.invoiceRepo.Get(additionalItem.InvoiceID)
			if err != nil {
				return nil, newErrors.Wrap(err, "repo invoice get")
			}

			builtItem.Type = item.Type

			builtItem.Title = item.InvoiceNumber

		} else if item.SalaryAdditionalExpenseID != nil {
			additionalItem, err := h.salaryAdditionalExpensesRepo.Get(*item.SalaryAdditionalExpenseID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo salary additional expense get")
			}

			item, err := h.salariesRepo.Get(additionalItem.SalaryID)

			if err != nil {
				return nil, newErrors.Wrap(err, "repo salary get")
			}

			builtItem.Type = data.TypeSalary
			builtItem.Title = "Zarada " + item.Month + " " + string(additionalItem.Title)

		}

		if len(items) == 1 {
			builtItem.Amount = response.Amount
		}

		response.Items = append(response.Items, builtItem)
	}

	return &response, nil
}

func (h *PaymentOrderServiceImpl) GetPaymentOrderList(filter dto.PaymentOrderFilterDTO) ([]dto.PaymentOrderResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.Year != nil {
		year := *filter.Year
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_payment": up.Between(startOfYear, endOfYear)})
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
		switch *filter.Status {
		case "Plaćen":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"sap_id is not ": nil})
		case "Na čekanju":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"sap_id is ": nil})
		case "Storniran":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": filter.Status})
		}
	}

	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"sap_id ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	if filter.Status == nil || *filter.Status != "Storniran" {
		cond := up.Or(
			&up.Cond{"status <>": "Storniran"},
			&up.Cond{"status is": nil},
		)
		conditionAndExp = up.And(conditionAndExp, cond)
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
		return nil, nil, newErrors.Wrap(err, "repo payment order get all")
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
		return nil, nil, newErrors.Wrap(err, "repo payment order get all obligations")
	}

	var response []dto.ObligationResponse
	for _, item := range items {
		var invoiceItems []dto.InvoiceItems

		if item.InvoiceID != nil && *item.InvoiceID != 0 {
			conditionAndExp := &up.AndExpr{}
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": *item.InvoiceID})
			articles, _, err := h.invoiceArticlesRepo.GetAll(nil, nil, conditionAndExp, nil)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "repo article get all")
			}

			accountMap := make(map[int]float64)
			totalAccountMap := make(map[int]float64)

			for _, item := range articles {
				if currentAmount, exists := accountMap[item.AccountID]; exists {
					accountMap[item.AccountID] = currentAmount + float64(float64(item.Amount)*(item.NetPrice+item.NetPrice*float64(item.VatPercentage)/100))
					totalAccountMap[item.AccountID] = currentAmount + float64(float64(item.Amount)*(item.NetPrice+item.NetPrice*float64(item.VatPercentage)/100))
				} else {
					accountMap[item.AccountID] = float64(float64(item.Amount) * (item.NetPrice + item.NetPrice*float64(item.VatPercentage)/100))
					totalAccountMap[item.AccountID] = float64(float64(item.Amount) * (item.NetPrice + item.NetPrice*float64(item.VatPercentage)/100))
				}
			}

			conditionAndExp = &up.AndExpr{}
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": *item.InvoiceID})
			paidItems, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "repo payment order items get")
			}

			for _, item := range paidItems {

				paymentOrder, err := h.repo.Get(item.PaymentOrderID)

				if err != nil {
					return nil, nil, newErrors.Wrap(err, "repo payment order get")
				}

				if paymentOrder.Status == nil || *paymentOrder.Status != "Storniran" {
					accountMap[item.SourceAccountID] -= item.Amount
				}
			}

			for account, amount := range accountMap {
				//if amount > 0 {
				invoiceItems = append(invoiceItems, dto.InvoiceItems{
					AccountID:   account,
					RemainPrice: amount,
					TotalPrice:  totalAccountMap[account],
				})
				//}
			}
		}

		response = append(response, dto.ObligationResponse{
			InvoiceID:                 item.InvoiceID,
			AdditionalExpenseID:       item.AdditionalExpenseID,
			SalaryAdditionalExpenseID: item.SalaryAdditionalExpenseID,
			Type:                      item.Type,
			Title:                     item.Title,
			Status:                    item.Status,
			TotalPrice:                item.TotalPrice,
			RemainPrice:               item.RemainPrice,
			InvoiceItems:              invoiceItems,
			CreatedAt:                 item.CreatedAt,
			AccountID:                 item.AccountID,
		})
	}

	return response, total, nil
}

func (h *PaymentOrderServiceImpl) PayPaymentOrder(ctx context.Context, id int, input dto.PaymentOrderDTO) error {
	err := data.Upper.Tx(func(tx up.Session) error {

		err := h.repo.PayPaymentOrder(ctx, tx, id, *input.SAPID, *input.DateOfSAP)
		if err != nil {
			return newErrors.Wrap(err, "repo payment order pay")
		}
		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}

func (h *PaymentOrderServiceImpl) CancelPaymentOrder(ctx context.Context, id int) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		input, err := h.GetPaymentOrder(id)

		if err != nil {
			return newErrors.Wrap(err, "repo payment order get")
		}

		for _, item := range input.Items {
			Type := 1

			if input.SourceOfFunding == "Donacija" {
				Type = 2
			}
			currentBudget, _, err := h.currentBudgetService.GetCurrentBudgetList(dto.CurrentBudgetFilterDTO{
				UnitID:    &input.OrganizationUnitID,
				AccountID: &item.AccountID,
				Type:      &Type,
			})

			if err != nil {
				return newErrors.Wrap(err, "repo current budget get all")
			}

			if len(currentBudget) > 0 {
				currentAmount := currentBudget[0].Balance.Add(decimal.NewFromFloat32(float32(item.Amount)))

				err = h.currentBudgetService.UpdateBalance(ctx, tx, currentBudget[0].ID, currentAmount)
				if err != nil {
					return newErrors.Wrap(err, "repo current budget update balance")
				}

			} else {
				return newErrors.Wrap(errors.ErrNotFound, "repo current budget get all")
			}
		}

		for _, item := range input.Items {
			if item.InvoiceID != nil {
				err = updateInvoiceStatusOnDelete(ctx, *item.InvoiceID, len(input.Items), tx, h)

				if err != nil {
					return newErrors.Wrap(err, "update invoice status on delete")
				}
			} else if item.AdditionalExpenseID != nil {
				err = updateAdditionalExpenseStatusOnDelete(*item.AdditionalExpenseID, len(input.Items), tx, h)

				if err != nil {
					return newErrors.Wrap(err, "update additional expense status on delete")
				}
			} else if item.SalaryAdditionalExpenseID != nil {
				err = updateSalaryAdditionalExpenseStatusOnDelete(*item.SalaryAdditionalExpenseID, len(input.Items), tx, h)

				if err != nil {
					return newErrors.Wrap(err, "update salary additional expense status on delete")
				}
			}
		}

		err = h.repo.CancelPaymentOrder(ctx, tx, id)
		if err != nil {
			return newErrors.Wrap(err, "repo payment order cancel payment order")
		}
		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}

func updateInvoiceStatus(ctx context.Context, id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	invoice, err := h.invoiceRepo.Get(id)

	if err != nil {
		return newErrors.Wrap(err, "repo invoice get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})

	articles, _, err := h.invoiceArticlesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo article get all")
	}

	var price float64
	for _, article := range articles {
		price += float64((article.NetPrice + article.NetPrice*float64(article.VatPercentage)/100) * float64(article.Amount))
	}

	conditionAndExp = &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo payment order items get all")
	}

	statusCanceled := "Storniran"
	for _, item := range items {
		paymentOrder, err := h.repo.Get(item.PaymentOrderID)

		if err != nil {
			return newErrors.Wrap(err, "repo payment order get")
		}

		if paymentOrder.Status == nil || *paymentOrder.Status != statusCanceled {
			amount += paymentOrder.Amount
		}
	}

	if amount+0.09999 >= price || lenOfArray > 1 {
		invoice.Status = data.InvoiceStatusFull
	} else {
		invoice.Status = data.InvoiceStatusPart
	}

	err = h.invoiceRepo.Update(ctx, tx, *invoice)

	if err != nil {
		return newErrors.Wrap(err, "repo invoice update")
	}

	return nil
}

func updateAdditionalExpenseStatus(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.additionalExpensesRepo.Get(id)

	if err != nil {
		return newErrors.Wrap(err, "repo additional expense get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"additional_expense_id": id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo payment order item get all")
	}

	statusCanceled := "Storniran"

	for _, item := range items {
		paymentOrder, err := h.repo.Get(item.PaymentOrderID)

		if err != nil {
			return newErrors.Wrap(err, "repo payment order get")
		}
		if paymentOrder.Status == nil || *paymentOrder.Status != statusCanceled {
			amount += paymentOrder.Amount
		}
	}

	if amount+0.09999 >= float64(item.Price) || lenOfArray > 1 {
		item.Status = data.InvoiceStatusFull
	} else {
		item.Status = data.InvoiceStatusPart
	}

	err = h.additionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return newErrors.Wrap(err, "repo additional expense update")
	}

	return nil
}

func updateSalaryAdditionalExpenseStatus(id int, amount float64, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.salaryAdditionalExpensesRepo.Get(id)

	if err != nil {
		return newErrors.Wrap(err, "repo salary additional expense get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_additional_expense_id": id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo payment order items get all")
	}

	statusCanceled := "Storniran"
	for _, item := range items {
		paymentOrder, err := h.repo.Get(item.PaymentOrderID)

		if err != nil {
			return newErrors.Wrap(err, "repo payment order get")
		}

		if paymentOrder.Status == nil || *paymentOrder.Status != statusCanceled {
			amount += paymentOrder.Amount
		}
	}

	if amount+0.09999 >= item.Amount || lenOfArray > 1 {
		item.Status = data.InvoiceStatusFull
	} else {
		item.Status = data.InvoiceStatusPart
	}

	err = h.salaryAdditionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return newErrors.Wrap(err, "repo salary additional expense update")
	}

	return nil
}

func updateInvoiceStatusOnDelete(ctx context.Context, id int, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	invoice, err := h.invoiceRepo.Get(id)

	if err != nil {
		return newErrors.Wrap(err, "repo invoice get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo payment order item get all")
	}

	//ako smo dobili nazad samo item koji brisemo
	if len(items) == 1 || lenOfArray > 1 {
		invoice.Status = data.InvoiceStatusCreated
	} else {

		numberOfItems := 0
		for _, item := range items {
			paymentOrder, err := h.repo.Get(item.PaymentOrderID)

			if err != nil {
				return newErrors.Wrap(err, "repo payment order get")
			}

			if paymentOrder.Status == nil || *paymentOrder.Status != "Storniran" {
				numberOfItems++
			}
		}

		if numberOfItems > 1 {
			invoice.Status = data.InvoiceStatusPart
		} else {
			invoice.Status = data.InvoiceStatusCreated
		}
	}

	err = h.invoiceRepo.Update(ctx, tx, *invoice)

	if err != nil {
		return newErrors.Wrap(err, "repo invoice update")
	}

	return nil
}

func updateAdditionalExpenseStatusOnDelete(id int, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.additionalExpensesRepo.Get(id)

	if err != nil {
		return newErrors.Wrap(err, "repo additional expense get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"additional_expense_id": id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo payment order item get all")
	}

	if len(items) == 1 || lenOfArray > 1 {
		item.Status = data.InvoiceStatusCreated
	} else {
		numberOfItems := 0
		for _, item := range items {
			paymentOrder, err := h.repo.Get(item.PaymentOrderID)

			if err != nil {
				return newErrors.Wrap(err, "repo payment order get")
			}

			if paymentOrder.Status == nil || *paymentOrder.Status != "Storniran" {
				numberOfItems++
			}
		}

		if numberOfItems > 1 {
			item.Status = data.InvoiceStatusPart
		} else {
			item.Status = data.InvoiceStatusCreated
		}
	}

	err = h.additionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return newErrors.Wrap(err, "repo additional expense update")
	}

	return nil
}

func updateSalaryAdditionalExpenseStatusOnDelete(id int, lenOfArray int, tx up.Session, h *PaymentOrderServiceImpl) error {
	item, err := h.salaryAdditionalExpensesRepo.Get(id)

	if err != nil {
		return newErrors.Wrap(err, "repo salary additional expense get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_additional_expense_id": id})

	items, _, err := h.itemsRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return newErrors.Wrap(err, "repo payment order item get all")
	}

	if len(items) == 1 || lenOfArray > 1 {
		item.Status = data.InvoiceStatusCreated
	} else {
		numberOfItems := 0
		for _, item := range items {
			paymentOrder, err := h.repo.Get(item.PaymentOrderID)

			if err != nil {
				return newErrors.Wrap(err, "repo payment order get")
			}

			if paymentOrder.Status == nil || *paymentOrder.Status != "Storniran" {
				numberOfItems++
			}
		}

		if numberOfItems > 1 {
			item.Status = data.InvoiceStatusPart
		} else {
			item.Status = data.InvoiceStatusCreated
		}
	}

	err = h.salaryAdditionalExpensesRepo.Update(tx, *item)

	if err != nil {
		return newErrors.Wrap(err, "repo salary additional expense update")
	}

	return nil
}
