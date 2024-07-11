package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type AccountingEntryServiceImpl struct {
	App                          *celeritas.Celeritas
	repo                         data.AccountingEntry
	invoiceRepo                  data.Invoice
	articlesRepo                 data.Article
	additionalExpensesRepo       data.AdditionalExpense
	salaryRepo                   data.Salary
	salaryAdditionalExpensesRepo data.SalaryAdditionalExpense
	modelOfAccountingRepo        ModelsOfAccountingService
	items                        data.AccountingEntryItem
	paymentOrderRepo             data.PaymentOrder
	enforcedPaymentRepo          data.EnforcedPayment
}

func NewAccountingEntryServiceImpl(app *celeritas.Celeritas, repo data.AccountingEntry, invoiceRepo data.Invoice, articlesRepo data.Article, additionalExpensesRepo data.AdditionalExpense, salaryRepo data.Salary, salaryAdditionalExpensesRepo data.SalaryAdditionalExpense, modelOfAccountingRepo ModelsOfAccountingService, items data.AccountingEntryItem, paymentOrderRepo data.PaymentOrder, enforcedPaymentRepo data.EnforcedPayment) AccountingEntryService {
	return &AccountingEntryServiceImpl{
		App:                          app,
		repo:                         repo,
		invoiceRepo:                  invoiceRepo,
		articlesRepo:                 articlesRepo,
		additionalExpensesRepo:       additionalExpensesRepo,
		salaryRepo:                   salaryRepo,
		salaryAdditionalExpensesRepo: salaryAdditionalExpensesRepo,
		modelOfAccountingRepo:        modelOfAccountingRepo,
		items:                        items,
		paymentOrderRepo:             paymentOrderRepo,
		enforcedPaymentRepo:          enforcedPaymentRepo,
	}
}

func (h *AccountingEntryServiceImpl) CreateAccountingEntry(ctx context.Context, input dto.AccountingEntryDTO) (*dto.AccountingEntryResponseDTO, error) {
	dataToInsert := input.ToAccountingEntry()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error

		if len(input.Items) > 0 {
			if (input.Items[0].InvoiceID != nil && *input.Items[0].InvoiceID != 0) || (input.Items[0].SalaryID != nil && *input.Items[0].SalaryID != 0) {
				dataToInsert.Type = data.TypeObligations
			} else if input.Items[0].PaymentOrderID != nil && *input.Items[0].PaymentOrderID != 0 {
				dataToInsert.Type = data.TypePaymentOrder
			} else if input.Items[0].EnforcedPaymentID != nil && *input.Items[0].EnforcedPaymentID != 0 {
				dataToInsert.Type = data.TypeEnforcedPayment
			} else if input.Items[0].ReturnEnforcedPaymentID != nil && *input.Items[0].ReturnEnforcedPaymentID != 0 {
				dataToInsert.Type = data.TypeReturnEnforcedPayment
			}
		}

		if dataToInsert.Type == "" {
			return newErrors.Wrap(errors.ErrInvalidInput, "check input")
		}

		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo accounting entry insert")
		}

		for _, item := range input.Items {
			itemToInsert := item.ToAccountingEntryItem()
			itemToInsert.EntryID = id

			_, err = h.items.Insert(tx, *itemToInsert)
			if err != nil {
				return newErrors.Wrap(err, "repo accounting entry item insert")
			}

			boolTrue := true
			if (item.Type == string(data.TypeInvoice) || item.Type == string(data.TypeContract) || item.Type == string(data.TypeDecision)) &&
				item.InvoiceID != nil && *item.InvoiceID != 0 {

				invoice, err := h.invoiceRepo.Get(*item.InvoiceID)

				if err != nil {
					return newErrors.Wrap(err, "repo invoice get")
				}

				invoice.Registred = &boolTrue

				err = h.invoiceRepo.Update(ctx, tx, *invoice)

				if err != nil {
					return newErrors.Wrap(err, "repo invoice update")
				}
			} else if item.Type == string(data.TypeSalary) && item.SalaryID != nil && *item.SalaryID != 0 {
				salary, err := h.salaryRepo.Get(*item.SalaryID)

				if err != nil {
					return newErrors.Wrap(err, "repo salary get")
				}

				salary.Registred = &boolTrue

				err = h.salaryRepo.Update(ctx, tx, *salary)

				if err != nil {
					return newErrors.Wrap(err, "repo salary update")
				}
			} else if item.Type == string(data.TypePaymentOrder) && item.PaymentOrderID != nil && *item.PaymentOrderID != 0 {
				paymentOrder, err := h.paymentOrderRepo.Get(*item.PaymentOrderID)

				if err != nil {
					return newErrors.Wrap(err, "repo payment order get")
				}

				paymentOrder.Registred = &boolTrue

				err = h.paymentOrderRepo.Update(ctx, tx, *paymentOrder)

				if err != nil {
					return newErrors.Wrap(err, "repo payment order update")
				}
			} else if item.Type == string(data.TypeEnforcedPayment) && item.EnforcedPaymentID != nil && *item.EnforcedPaymentID != 0 {
				enforcedPayment, err := h.enforcedPaymentRepo.Get(*item.EnforcedPaymentID)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment get")
				}

				enforcedPayment.Registred = &boolTrue

				err = h.enforcedPaymentRepo.Update(ctx, tx, *enforcedPayment)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment update")
				}
			} else if item.Type == string(data.TypeReturnEnforcedPayment) && item.ReturnEnforcedPaymentID != nil && *item.ReturnEnforcedPaymentID != 0 {
				enforcedPayment, err := h.enforcedPaymentRepo.Get(*item.ReturnEnforcedPaymentID)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment get")
				}

				enforcedPayment.RegistredReturn = &boolTrue

				err = h.enforcedPaymentRepo.Update(ctx, tx, *enforcedPayment)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment update")
				}
			}

		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo accounting entry get")
	}

	res := dto.ToAccountingEntryResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *AccountingEntryServiceImpl) UpdateAccountingEntry(ctx context.Context, id int, input dto.AccountingEntryDTO) (*dto.AccountingEntryResponseDTO, error) {
	dataToInsert := input.ToAccountingEntry()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo accounting entry update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo accounting entry get")
	}

	response := dto.ToAccountingEntryResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *AccountingEntryServiceImpl) DeleteAccountingEntry(ctx context.Context, id int) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"entry_id": id})
		items, _, err := h.items.GetAll(nil, nil, conditionAndExp, nil)
		if err != nil {
			return newErrors.Wrap(err, "repo accounting entry items get all")
		}

		for _, item := range items {
			boolFalse := false
			if (item.Type == string(data.TypeInvoice) || item.Type == string(data.TypeContract) || item.Type == string(data.TypeDecision)) &&
				item.InvoiceID != nil && *item.InvoiceID != 0 {
				invoice, err := h.invoiceRepo.Get(*item.InvoiceID)

				if err != nil {
					return newErrors.Wrap(err, "repo invoice get")
				}

				invoice.Registred = &boolFalse

				err = h.invoiceRepo.Update(ctx, tx, *invoice)

				if err != nil {
					return newErrors.Wrap(err, "repo invoice update")
				}
			} else if item.Type == string(data.TypeSalary) && item.SalaryID != nil && *item.SalaryID != 0 {
				salary, err := h.salaryRepo.Get(*item.SalaryID)

				if err != nil {
					return newErrors.Wrap(err, "repo salary get")
				}

				salary.Registred = &boolFalse

				err = h.salaryRepo.Update(ctx, tx, *salary)

				if err != nil {
					return newErrors.Wrap(err, "repo salary update")
				}
			} else if item.Type == string(data.TypePaymentOrder) && item.PaymentOrderID != nil && *item.PaymentOrderID != 0 {
				paymentOrder, err := h.paymentOrderRepo.Get(*item.PaymentOrderID)

				if err != nil {
					return newErrors.Wrap(err, "repo payment order get")
				}

				paymentOrder.Registred = &boolFalse

				err = h.paymentOrderRepo.Update(ctx, tx, *paymentOrder)

				if err != nil {
					return newErrors.Wrap(err, "repo payment order update")
				}
			} else if item.Type == string(data.TypeEnforcedPayment) && item.EnforcedPaymentID != nil && *item.EnforcedPaymentID != 0 {
				enforcedPayment, err := h.enforcedPaymentRepo.Get(*item.EnforcedPaymentID)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment get")
				}

				enforcedPayment.Registred = &boolFalse

				err = h.enforcedPaymentRepo.Update(ctx, tx, *enforcedPayment)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment update")
				}
			} else if item.Type == string(data.TypeReturnEnforcedPayment) && item.ReturnEnforcedPaymentID != nil && *item.ReturnEnforcedPaymentID != 0 {
				enforcedPayment, err := h.enforcedPaymentRepo.Get(*item.ReturnEnforcedPaymentID)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment get")
				}

				enforcedPayment.RegistredReturn = &boolFalse

				err = h.enforcedPaymentRepo.Update(ctx, tx, *enforcedPayment)

				if err != nil {
					return newErrors.Wrap(err, "repo enforced payment update")
				}
			}
		}

		err = h.repo.Delete(ctx, id)
		if err != nil {
			return newErrors.Wrap(err, "repo accounting entry delete")
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}

func (h *AccountingEntryServiceImpl) GetAccountingEntry(id int) (*dto.AccountingEntryResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo accounting entry get")
	}

	response := dto.ToAccountingEntryResponseDTO(*data)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"entry_id": id})
	items, _, err := h.items.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo accounting entry item get all")
	}

	var debitAmount float64
	var creditAmount float64
	for _, item := range items {
		responseItem := dto.AccountingEntryItemResponseDTO{
			ID:                      item.ID,
			Title:                   item.Title,
			EntryID:                 item.EntryID,
			AccountID:               item.AccountID,
			CreditAmount:            item.CreditAmount,
			DebitAmount:             item.DebitAmount,
			SupplierID:              item.SupplierID,
			InvoiceID:               item.InvoiceID,
			SalaryID:                item.SalaryID,
			PaymentOrderID:          item.PaymentOrderID,
			EnforcedPaymentID:       item.EnforcedPaymentID,
			ReturnEnforcedPaymentID: item.ReturnEnforcedPaymentID,
			Date:                    item.Date,
			Type:                    item.Type,
		}

		debitAmount += item.DebitAmount
		creditAmount += item.CreditAmount

		response.Items = append(response.Items, responseItem)
	}
	response.CreditAmount = creditAmount
	response.DebitAmount = debitAmount

	return &response, nil
}

func (h *AccountingEntryServiceImpl) GetAccountingEntryList(filter dto.AccountingEntryFilterDTO) ([]dto.AccountingEntryResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Type != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"type": *filter.Type})
	}

	if filter.DateOfStart != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_booking > ": *filter.DateOfStart})
	}

	if filter.DateOfEnd != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_booking < ": *filter.DateOfEnd})
	}

	if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}

	orders = append(orders, "-created_at")

	var data []*data.AccountingEntry
	var total *uint64
	var err error

	//ako ti bude trebala paginacija za report napravi novi endpoint koji stavlja ovo u ono order case
	if filter.SortForReport != nil && *filter.SortForReport {
		var reportOrder []interface{}
		reportOrder = append(reportOrder, "created_at")
		conditionAndExpInvoice := up.And(conditionAndExp, &up.Cond{"type": "obligations"})
		invoiceData, totalInvoice, err := h.repo.GetAll(nil, nil, conditionAndExpInvoice, reportOrder)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo accounting entry get all")
		}

		conditionAndExpPaymentOrders := up.And(conditionAndExp, &up.Cond{"type": "payment_orders"})
		paymentOrderData, totalPaymentOrder, err := h.repo.GetAll(nil, nil, conditionAndExpPaymentOrders, reportOrder)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo accounting entry get all")
		}

		conditionAndExpEnforcedPayment := up.And(conditionAndExp, &up.Cond{"type": "enforced_payments"})
		enforcedPaymentData, totalEnforcedPayment, err := h.repo.GetAll(nil, nil, conditionAndExpEnforcedPayment, reportOrder)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo accounting entry get all")
		}

		conditionAndExpReturnEnforcedPayment := up.And(conditionAndExp, &up.Cond{"type": "return_enforced_payment"})
		returnEnforcedPaymentData, totalReturnEnforced, err := h.repo.GetAll(nil, nil, conditionAndExpReturnEnforcedPayment, reportOrder)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo accounting entry get all")
		}

		data = append(data, invoiceData...)
		data = append(data, paymentOrderData...)
		data = append(data, enforcedPaymentData...)
		data = append(data, returnEnforcedPaymentData...)
		totalInt := *totalInvoice + *totalPaymentOrder + *totalReturnEnforced + *totalEnforcedPayment
		total = &totalInt
	} else {
		data, total, err = h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo accounting entry get all")
		}
	}

	response := dto.ToAccountingEntryListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"entry_id": response[i].ID})
		items, _, err := h.items.GetAll(nil, nil, conditionAndExp, nil)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo accounting entry item get all")
		}

		var debitAmount float64
		var creditAmount float64

		year := response[i].DateOfBooking.Year()
		yearLastTwoDigits := year % 100

		// Format the string
		formatedIDOfEntry := fmt.Sprintf("%02d-%03d", yearLastTwoDigits, response[i].IDOfEntry)

		for _, item := range items {
			responseItem := dto.AccountingEntryItemResponseDTO{
				ID:                      item.ID,
				Title:                   item.Title,
				EntryID:                 item.EntryID,
				AccountID:               item.AccountID,
				CreditAmount:            item.CreditAmount,
				DebitAmount:             item.DebitAmount,
				InvoiceID:               item.InvoiceID,
				SalaryID:                item.SalaryID,
				PaymentOrderID:          item.PaymentOrderID,
				SupplierID:              item.SupplierID,
				EnforcedPaymentID:       item.EnforcedPaymentID,
				ReturnEnforcedPaymentID: item.ReturnEnforcedPaymentID,
				Date:                    item.Date,
				Type:                    item.Type,
				EntryNumber:             formatedIDOfEntry,
				EntryDate:               response[i].DateOfBooking,
			}

			debitAmount += item.DebitAmount
			creditAmount += item.CreditAmount

			response[i].Items = append(response[i].Items, responseItem)
		}
		response[i].CreditAmount = creditAmount
		response[i].DebitAmount = debitAmount
	}

	return response, total, nil
}

func (h *AccountingEntryServiceImpl) GetAnalyticalCard(filter data.AnalyticalCardFilter) ([]data.AnalyticalCard, error) {

	var response []data.AnalyticalCard

	if filter.SupplierID != nil && *filter.SupplierID != 0 {
		responseItem, err := h.repo.GetAnalyticalCard(filter)

		if err != nil {
			return nil, newErrors.Wrap(err, "repo accounting entry get analytical card")
		}
		responseItem.SupplierID = *filter.SupplierID
		response = append(response, *responseItem)
	} else {
		allSuppliers, err := h.repo.GetAllSuppliers(filter)

		if err != nil {
			return nil, newErrors.Wrap(err, "repo accounting entry get all suppliers")
		}

		for _, supplierID := range allSuppliers {
			filter.SupplierID = &supplierID

			responseItem, err := h.repo.GetAnalyticalCard(filter)

			if err != nil {
				h.App.ErrorLog.Println(err)
				return nil, newErrors.Wrap(err, "repo accounting entry get analytical card")
			}

			responseItem.SupplierID = *filter.SupplierID
			response = append(response, *responseItem)

		}
	}

	return response, nil
}

func (h *AccountingEntryServiceImpl) GetObligationsForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.ObligationForAccounting, *uint64, error) {

	dataFilter := data.ObligationsFilter{
		Page:               filter.Page,
		Size:               filter.Size,
		OrganizationUnitID: filter.OrganizationUnitID,
		Type:               filter.Type,
		Search:             filter.Search,
		DateOfStart:        filter.DateOfStart,
		DateOfEnd:          filter.DateOfEnd,
	}

	items, total, err := h.repo.GetObligationsForAccounting(dataFilter)

	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo accounting entry get obligations for accounting")
	}

	var response []dto.ObligationForAccounting
	for _, item := range items {
		response = append(response, dto.ObligationForAccounting{
			InvoiceID:  item.InvoiceID,
			SalaryID:   item.SalaryID,
			Type:       item.Type,
			Title:      item.Title,
			Status:     item.Status,
			Price:      item.Price,
			SupplierID: item.SupplierID,
			Date:       item.Date,
			CreatedAt:  item.CreatedAt,
		})
	}

	return response, total, nil
}

func (h *AccountingEntryServiceImpl) GetPaymentOrdersForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.PaymentOrdersForAccounting, *uint64, error) {

	dataFilter := data.ObligationsFilter{
		Page:               filter.Page,
		Size:               filter.Size,
		OrganizationUnitID: filter.OrganizationUnitID,
		Search:             filter.Search,
		DateOfStart:        filter.DateOfStart,
		DateOfEnd:          filter.DateOfEnd,
	}

	items, total, err := h.repo.GetPaymentOrdersForAccounting(dataFilter)

	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo accounting entry get payment orders for accounting")
	}

	var response []dto.PaymentOrdersForAccounting
	for _, item := range items {
		response = append(response, dto.PaymentOrdersForAccounting{
			PaymentOrderID: item.PaymentOrderID,
			Title:          item.Title,
			Price:          item.Price,
			SupplierID:     item.SupplierID,
			Date:           item.Date,
			CreatedAt:      item.CreatedAt,
		})
	}

	return response, total, nil
}

func (h *AccountingEntryServiceImpl) GetEnforcedPaymentsForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.PaymentOrdersForAccounting, *uint64, error) {

	dataFilter := data.ObligationsFilter{
		Page:               filter.Page,
		Size:               filter.Size,
		OrganizationUnitID: filter.OrganizationUnitID,
		Search:             filter.Search,
		DateOfStart:        filter.DateOfStart,
		DateOfEnd:          filter.DateOfEnd,
	}

	items, total, err := h.repo.GetEnforcedPaymentsForAccounting(dataFilter)

	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo accounting entry get enforced payments for accounting")
	}

	var response []dto.PaymentOrdersForAccounting
	for _, item := range items {
		response = append(response, dto.PaymentOrdersForAccounting{
			PaymentOrderID: item.PaymentOrderID,
			Title:          item.Title,
			Price:          item.Price,
			SupplierID:     item.SupplierID,
			Date:           item.Date,
			CreatedAt:      item.CreatedAt,
		})
	}

	return response, total, nil
}

func (h *AccountingEntryServiceImpl) GetReturnedEnforcedPaymentsForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.PaymentOrdersForAccounting, *uint64, error) {

	dataFilter := data.ObligationsFilter{
		Page:               filter.Page,
		Size:               filter.Size,
		OrganizationUnitID: filter.OrganizationUnitID,
		Search:             filter.Search,
		DateOfStart:        filter.DateOfStart,
		DateOfEnd:          filter.DateOfEnd,
	}

	items, total, err := h.repo.GetReturnedEnforcedPaymentsForAccounting(dataFilter)

	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo accounting entry get returned enforced payments for accounting")
	}

	var response []dto.PaymentOrdersForAccounting
	for _, item := range items {
		response = append(response, dto.PaymentOrdersForAccounting{
			PaymentOrderID: item.PaymentOrderID,
			Title:          item.Title,
			Price:          item.Price,
			SupplierID:     item.SupplierID,
			Date:           item.Date,
			CreatedAt:      item.CreatedAt,
		})
	}

	return response, total, nil
}

func (h *AccountingEntryServiceImpl) BuildAccountingOrderForObligations(orderData dto.AccountingOrderForObligationsData) (*dto.AccountingOrderForObligations, error) {

	response := dto.AccountingOrderForObligations{}
	response.DateOfBooking = orderData.DateOfBooking

	for _, id := range orderData.InvoiceID {
		invoice, err := h.invoiceRepo.Get(id)

		if err != nil {
			return nil, newErrors.Wrap(err, "repo invoice get")
		}

		switch invoice.Type {
		case data.TypeInvoice:
			item, err := buildAccountingOrderForInvoice(id, h)

			if err != nil {
				return nil, newErrors.Wrap(err, "build accounting order for invoice")
			}

			response.Items = append(response.Items, item...)
		case data.TypeDecision:
			item, err := buildAccountingOrderForDecisions(id, h)

			if err != nil {
				return nil, newErrors.Wrap(err, "build accounting order for decisions")
			}

			response.Items = append(response.Items, item...)
		case data.TypeContract:
			item, err := buildAccountingOrderForContracts(id, h)

			if err != nil {
				return nil, newErrors.Wrap(err, "build accounting order for contracts")
			}
			response.Items = append(response.Items, item...)
		}
	}

	for _, id := range orderData.SalaryID {
		item, err := buildAccountingOrderForSalaries(id, h)

		if err != nil {
			return nil, newErrors.Wrap(err, "build accounting order for salaries")
		}
		response.Items = append(response.Items, item...)
	}

	for _, id := range orderData.PaymentOrderID {
		item, err := buildAccountingOrderForPaymentOrder(id, h)

		if err != nil {
			return nil, newErrors.Wrap(err, "build accounting order for payment order")
		}
		response.Items = append(response.Items, item...)
	}

	for _, id := range orderData.EnforcedPaymentID {
		item, err := buildAccountingOrderForEnforcedPayment(id, h)

		if err != nil {
			return nil, newErrors.Wrap(err, "build accounting order for enforced payment")
		}
		response.Items = append(response.Items, item...)
	}

	for _, id := range orderData.ReturnEnforcedPaymentID {
		item, err := buildAccountingOrderForReturnEnforcedPayment(id, h)

		if err != nil {
			return nil, newErrors.Wrap(err, "build accounting order for return enforced payment")
		}
		response.Items = append(response.Items, item...)

	}

	for _, item := range response.Items {
		response.CreditAmount += item.CreditAmount
		response.DebitAmount += item.DebitAmount
	}

	return &response, nil
}

func buildAccountingOrderForInvoice(id int, h *AccountingEntryServiceImpl) ([]dto.AccountingOrderItemsForObligations, error) {
	response := []dto.AccountingOrderItemsForObligations{}

	invoice, err := h.invoiceRepo.Get(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo invoice get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	articles, _, err := h.articlesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo article get all")
	}

	var price float64
	for _, article := range articles {
		price += (article.NetPrice + float64(article.VatPercentage)/100*article.NetPrice) * float64(article.Amount)
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeInvoice,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting get all")
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo model of accounting get all")
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.MainBillTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(price),
				Title:       modelItem.Title,
				Type:        data.TypeInvoice,
				Date:        invoice.DateOfInvoice,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
			})
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(price),
				Title:        modelItem.Title,
				Type:         data.TypeInvoice,
				Date:         invoice.DateOfInvoice,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
				SupplierID: invoice.SupplierID,
			})
		}

	}

	return response, nil
}

func buildAccountingOrderForDecisions(id int, h *AccountingEntryServiceImpl) ([]dto.AccountingOrderItemsForObligations, error) {
	response := []dto.AccountingOrderItemsForObligations{}

	invoice, err := h.invoiceRepo.Get(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo invoice get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	additionalExpenses, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo additional expenses get all")
	}

	var price float64
	var netPrice float64
	var taxPrice float64
	var subTaxPrice float64
	var contributionForPIO float64
	var contributionForUnemployment float64
	var contributionForPIOEmployee float64
	var contributionForPIOEmployer float64
	var contributionForUnemploymentEmployee float64
	var contributionForUnemploymentEmployer float64
	var contributionForLaborFund float64
	var taxSupplierID int
	var subTaxSupplierID int
	var PIOSupplierID int
	var UnemploymentSupplierID int
	var PIOEmployeeSupplierID int
	var PIOEmployerSupplierID int
	var UnemploymentEmployeeSupplierID int
	var UnemploymentEmployerSupplierID int
	var LaborFundSupplierID int

	for _, item := range additionalExpenses {
		price += float64(item.Price)
		switch item.Title {
		case data.NetTitle:
			netPrice += float64(item.Price)
		case data.ObligationTaxTitle:
			taxPrice += float64(item.Price)
			taxSupplierID = item.SubjectID
		case data.ObligationSubTaxTitle:
			subTaxPrice += float64(item.Price)
			subTaxSupplierID = item.SubjectID
		case data.ContributionForPIOTitle:
			contributionForPIO += float64(item.Price)
			PIOSupplierID = item.SubjectID
		case data.ContributionForUnemploymentTitle:
			contributionForUnemployment += float64(item.Price)
			UnemploymentSupplierID = item.SubjectID
		case data.ContributionForPIOEmployeeTitle:
			contributionForPIOEmployee += float64(item.Price)
			PIOEmployeeSupplierID = item.SubjectID
		case data.ContributionForPIOEmployerTitle:
			contributionForPIOEmployer += float64(item.Price)
			PIOEmployerSupplierID = item.SubjectID
		case data.ContributionForUnemploymentEmployeeTitle:
			contributionForUnemploymentEmployee += float64(item.Price)
			UnemploymentEmployeeSupplierID = item.SubjectID
		case data.ContributionForUnemploymentEmployerTitle:
			contributionForUnemploymentEmployer += float64(item.Price)
			UnemploymentEmployerSupplierID = item.SubjectID
		case data.LaborFundTitle:
			contributionForLaborFund += float64(item.Price)
			LaborFundSupplierID = item.SubjectID
		}
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeDecision,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting get all")
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo model of accounting get all")
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.MainBillTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(price),
				Title:       modelItem.Title,
				Type:        data.TypeDecision,
				Date:        invoice.DateOfInvoice,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
			})
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(netPrice),
				Title:        modelItem.Title,
				Type:         data.TypeDecision,
				Date:         invoice.DateOfInvoice,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
				SupplierID: invoice.SupplierID,
			})
		case data.TaxTitle:
			if taxPrice > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(taxPrice),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: taxSupplierID,
				})
			}
		case data.SubTaxTitle:
			if subTaxPrice > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(subTaxPrice),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: subTaxSupplierID,
				})
			}
		case data.LaborContributionsTitle:
			if contributionForLaborFund > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForLaborFund),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: LaborFundSupplierID,
				})
			}
		case data.PIOContributionsTitle:
			if contributionForPIO > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIO),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: PIOSupplierID,
				})
			}
		case data.UnemployementContributionsTitle:
			if contributionForUnemployment > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemployment),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: UnemploymentSupplierID,
				})
			}
		case data.PIOEmployeeContributionsTitle:
			if contributionForPIOEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: PIOEmployeeSupplierID,
				})
			}
		case data.PIOEmployerContributionsTitle:
			if contributionForPIOEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: PIOEmployerSupplierID,
				})
			}
		case data.UnemployementEmployeeContributionsTitle:
			if contributionForUnemploymentEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: UnemploymentEmployeeSupplierID,
				})
			}
		case data.UnemployementEmployerContributionsTitle:
			if contributionForUnemploymentEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: UnemploymentEmployerSupplierID,
				})
			}
		}

	}

	return response, nil
}

func buildAccountingOrderForContracts(id int, h *AccountingEntryServiceImpl) ([]dto.AccountingOrderItemsForObligations, error) {
	response := []dto.AccountingOrderItemsForObligations{}

	invoice, err := h.invoiceRepo.Get(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo invoice get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	additionalExpenses, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo additional expenses get all")
	}

	var price float64
	var netPrice float64
	var taxPrice float64
	var subTaxPrice float64
	var contributionForPIO float64
	var contributionForUnemployment float64
	var contributionForPIOEmployee float64
	var contributionForPIOEmployer float64
	var contributionForUnemploymentEmployee float64
	var contributionForUnemploymentEmployer float64
	var contributionForLaborFund float64

	var taxSupplierID int
	var subTaxSupplierID int
	var PIOSupplierID int
	var UnemploymentSupplierID int
	var PIOEmployeeSupplierID int
	var PIOEmployerSupplierID int
	var UnemploymentEmployeeSupplierID int
	var UnemploymentEmployerSupplierID int
	var LaborFundSupplierID int

	for _, item := range additionalExpenses {
		price += float64(item.Price)
		switch item.Title {
		case data.NetTitle:
			netPrice += float64(item.Price)
		case data.ObligationTaxTitle:
			taxPrice += float64(item.Price)
			taxSupplierID = item.SubjectID
		case data.ObligationSubTaxTitle:
			subTaxPrice += float64(item.Price)
			subTaxSupplierID = item.SubjectID
		case data.ContributionForPIOTitle:
			contributionForPIO += float64(item.Price)
			PIOSupplierID = item.SubjectID
		case data.ContributionForUnemploymentTitle:
			contributionForUnemployment += float64(item.Price)
			UnemploymentSupplierID = item.SubjectID
		case data.ContributionForPIOEmployeeTitle:
			contributionForPIOEmployee += float64(item.Price)
			PIOEmployeeSupplierID = item.SubjectID
		case data.ContributionForPIOEmployerTitle:
			contributionForPIOEmployer += float64(item.Price)
			PIOEmployerSupplierID = item.SubjectID
		case data.ContributionForUnemploymentEmployeeTitle:
			contributionForUnemploymentEmployee += float64(item.Price)
			UnemploymentEmployeeSupplierID = item.SubjectID
		case data.ContributionForUnemploymentEmployerTitle:
			contributionForUnemploymentEmployer += float64(item.Price)
			UnemploymentEmployerSupplierID = item.SubjectID
		case data.LaborFundTitle:
			contributionForLaborFund += float64(item.Price)
			LaborFundSupplierID = item.SubjectID
		}
	}
	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeContract,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting get all")
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo model of accounting get all")
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.MainBillTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(price),
				Title:       modelItem.Title,
				Type:        data.TypeContract,
				Date:        invoice.DateOfInvoice,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
			})
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(netPrice),
				Title:        modelItem.Title,
				Type:         data.TypeContract,
				Date:         invoice.DateOfInvoice,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
				SupplierID: invoice.SupplierID,
			})
		case data.TaxTitle:
			if taxPrice > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(taxPrice),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: taxSupplierID,
				})
			}
		case data.SubTaxTitle:
			if subTaxPrice > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(subTaxPrice),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: subTaxSupplierID,
				})
			}
		case data.LaborContributionsTitle:
			if contributionForLaborFund > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForLaborFund),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: LaborFundSupplierID,
				})
			}
		case data.PIOContributionsTitle:
			if contributionForPIO > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIO),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: PIOSupplierID,
				})
			}
		case data.UnemployementContributionsTitle:
			if contributionForUnemployment > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemployment),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: UnemploymentSupplierID,
				})
			}
		case data.PIOEmployeeContributionsTitle:
			if contributionForPIOEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: PIOEmployeeSupplierID,
				})
			}
		case data.PIOEmployerContributionsTitle:
			if contributionForPIOEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: PIOEmployerSupplierID,
				})
			}
		case data.UnemployementEmployeeContributionsTitle:
			if contributionForUnemploymentEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: UnemploymentEmployeeSupplierID,
				})
			}
		case data.UnemployementEmployerContributionsTitle:
			if contributionForUnemploymentEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Date:         invoice.DateOfInvoice,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
					SupplierID: UnemploymentEmployerSupplierID,
				})
			}
		}

	}

	return response, nil
}

func buildAccountingOrderForSalaries(id int, h *AccountingEntryServiceImpl) ([]dto.AccountingOrderItemsForObligations, error) {
	response := []dto.AccountingOrderItemsForObligations{}

	salary, err := h.salaryRepo.Get(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo salary get")
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_id": id})
	additionalExpenses, _, err := h.salaryAdditionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo salary additional expenses get all")
	}

	models, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeSalary,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting get all")
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(models) != 1 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo model of accounting get all")
	}

	var price float64
	var taxPrice float64
	var taxSupplierID int

	for _, item := range additionalExpenses {
		price += float64(item.Amount)
		bookedItem := &dto.AccountingOrderItemsForObligations{}
		switch item.Type {
		case data.SalaryAdditionalExpenseType(data.ContributionsSalaryExpenseType):
			switch item.Title {
			case data.PIOEmployeeContributionsTitle:
				if item.Amount > 0 {
					bookedItem = buildBookedItemForSalary(item, models[0].Items, data.PIOEmployeeContributionsTitle)
				}
			case data.PIOEmployerContributionsTitle:
				if item.Amount > 0 {
					bookedItem = buildBookedItemForSalary(item, models[0].Items, data.PIOEmployerContributionsTitle)
				}
			case data.UnemployementEmployeeContributionsTitle:
				if item.Amount > 0 {
					bookedItem = buildBookedItemForSalary(item, models[0].Items, data.UnemployementEmployeeContributionsTitle)
				}
			case data.UnemployementEmployerContributionsTitle:
				if item.Amount > 0 {
					bookedItem = buildBookedItemForSalary(item, models[0].Items, data.UnemployementEmployerContributionsTitle)
				}
			case data.LaborContributionsTitle:
				if item.Amount > 0 {
					bookedItem = buildBookedItemForSalary(item, models[0].Items, data.LaborContributionsTitle)
				}
			}
		case data.SalaryAdditionalExpenseType(data.TaxesSalaryExpenseType):
			taxPrice += item.Amount
			taxSupplierID = item.SubjectID
		case data.SalaryAdditionalExpenseType(data.SubTaxesSalaryExpenseType):
			if item.Amount > 0 {
				bookedItem = buildBookedItemForSalary(item, models[0].Items, data.SubTaxTitle)
			}
		case data.SalaryAdditionalExpenseType(data.BanksSalaryExpenseType):
			if item.Amount > 0 {
				bookedItem = buildBookedItemForSalaryForBanks(item, models[0].Items, data.BankTitle)
			}
		case data.SalaryAdditionalExpenseType(data.SuspensionsSalaryExpenseType):
			if item.Amount > 0 {
				bookedItem = buildBookedItemForSalaryForBanks(item, models[0].Items, data.SuspensionsTitle)
			}
		}
		if bookedItem != nil && bookedItem.Title != "" {
			bookedItem.Salary = dto.DropdownSimple{
				ID:    salary.ID,
				Title: salary.Month,
			}
			bookedItem.Date = salary.DateOfCalculation
			response = append(response, *bookedItem)
		}
	}

	for _, model := range models[0].Items {
		if model.Title == data.MainBillTitle {

			newSlice := make([]dto.AccountingOrderItemsForObligations, 0)

			newSlice = append(newSlice, dto.AccountingOrderItemsForObligations{
				AccountID:   model.DebitAccountID,
				DebitAmount: float32(price),
				Title:       model.Title,
				Type:        data.TypeSalary,
				Date:        salary.DateOfCalculation,
				Salary: dto.DropdownSimple{
					ID:    salary.ID,
					Title: salary.Month,
				},
			})

			response = append(newSlice, response...)

		}

		if model.Title == data.TaxTitle {
			if taxPrice > 0 {
				bookedItem := buildBookedItemForSalary(&data.SalaryAdditionalExpense{Amount: taxPrice, SubjectID: taxSupplierID}, models[0].Items, data.TaxTitle)

				index := 0
				if len(response) > 0 {
					if response[0].Title == data.MainBillTitle {
						index++
					}
				}

				bookedItem.Salary.ID = salary.ID
				bookedItem.Salary.Title = salary.Month
				bookedItem.Date = salary.DateOfCalculation

				response = append(response[:index], append([]dto.AccountingOrderItemsForObligations{*bookedItem}, response[index:]...)...)
			}
		}

	}

	return response, nil
}

func buildBookedItemForSalary(item *data.SalaryAdditionalExpense, models []dto.ModelOfAccountingItemResponseDTO, title data.AccountingOrderItemsTitle) *dto.AccountingOrderItemsForObligations {

	for _, model := range models {
		if string(model.Title) == string(title) {
			title := model.Title
			if model.Title == data.SubTaxTitle {
				title = title + " - " + item.Title
			}
			response := dto.AccountingOrderItemsForObligations{
				AccountID:    model.CreditAccountID,
				CreditAmount: float32(item.Amount),
				Title:        title,
				Type:         data.TypeSalary,
				SupplierID:   item.SubjectID,
			}
			return &response
		}
	}

	return nil
}

func buildBookedItemForSalaryForBanks(item *data.SalaryAdditionalExpense, models []dto.ModelOfAccountingItemResponseDTO, title data.AccountingOrderItemsTitle) *dto.AccountingOrderItemsForObligations {

	for _, model := range models {
		if string(model.Title) == string(title) {
			response := dto.AccountingOrderItemsForObligations{
				AccountID:    model.CreditAccountID,
				CreditAmount: float32(item.Amount),
				Title:        item.Title + " - " + title,
				Type:         data.TypeSalary,
				SupplierID:   item.SubjectID,
			}
			return &response
		}
	}

	return nil
}

func buildAccountingOrderForPaymentOrder(id int, h *AccountingEntryServiceImpl) ([]dto.AccountingOrderItemsForObligations, error) {
	response := []dto.AccountingOrderItemsForObligations{}

	paymentOrder, err := h.paymentOrderRepo.Get(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo payment order get")
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypePaymentOrder,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting get all")
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo model of accounting get all")
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(paymentOrder.Amount),
				Title:       modelItem.Title,
				Type:        data.TypePaymentOrder,
				Date:        *paymentOrder.DateOfSAP,
				PaymentOrder: dto.DropdownSimple{
					ID:    paymentOrder.ID,
					Title: *paymentOrder.SAPID,
				},
				SupplierID: paymentOrder.SupplierID,
			})
		case data.MainBillTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(paymentOrder.Amount),
				Title:        modelItem.Title,
				Type:         data.TypePaymentOrder,
				Date:         *paymentOrder.DateOfSAP,
				PaymentOrder: dto.DropdownSimple{
					ID:    paymentOrder.ID,
					Title: *paymentOrder.SAPID,
				},
			})
		case data.CostTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(paymentOrder.Amount),
				Title:       modelItem.Title,
				Date:        *paymentOrder.DateOfSAP,
				Type:        data.TypePaymentOrder,
				PaymentOrder: dto.DropdownSimple{
					ID:    paymentOrder.ID,
					Title: *paymentOrder.SAPID,
				},
			})
		case data.AllocatedAmountTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(paymentOrder.Amount),
				Title:        modelItem.Title,
				Date:         *paymentOrder.DateOfSAP,
				Type:         data.TypePaymentOrder,
				PaymentOrder: dto.DropdownSimple{
					ID:    paymentOrder.ID,
					Title: *paymentOrder.SAPID,
				},
			})
		}

	}

	return response, nil
}

func buildAccountingOrderForEnforcedPayment(id int, h *AccountingEntryServiceImpl) ([]dto.AccountingOrderItemsForObligations, error) {
	response := []dto.AccountingOrderItemsForObligations{}

	enforcedPayment, err := h.enforcedPaymentRepo.Get(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment get")
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeEnforcedPayment,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting get all")
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo model of accounting get all")
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(enforcedPayment.Amount),
				Title:       modelItem.Title,
				Date:        *enforcedPayment.DateOfSAP,
				Type:        data.TypeEnforcedPayment,
				EnforcedPayment: dto.DropdownSimple{
					ID:    enforcedPayment.ID,
					Title: *enforcedPayment.SAPID,
				},
				SupplierID: enforcedPayment.SupplierID,
			})
		case data.ProcessCostTitle:
			if enforcedPayment.AmountForAgent > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:   modelItem.DebitAccountID,
					DebitAmount: float32(enforcedPayment.AmountForAgent),
					Title:       modelItem.Title,
					Type:        data.TypeEnforcedPayment,
					Date:        *enforcedPayment.DateOfSAP,
					EnforcedPayment: dto.DropdownSimple{
						ID:    enforcedPayment.ID,
						Title: *enforcedPayment.SAPID,
					},
				})
			}
		case data.LawyerCostTitle:
			if enforcedPayment.AmountForLawyer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:   modelItem.DebitAccountID,
					DebitAmount: float32(enforcedPayment.AmountForLawyer),
					Title:       modelItem.Title,
					Type:        data.TypeEnforcedPayment,
					Date:        *enforcedPayment.DateOfSAP,
					EnforcedPayment: dto.DropdownSimple{
						ID:    enforcedPayment.ID,
						Title: *enforcedPayment.SAPID,
					},
				})
			}
		case data.BankCostTitle:
			if enforcedPayment.AmountForBank > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:   modelItem.DebitAccountID,
					DebitAmount: float32(enforcedPayment.AmountForBank),
					Title:       modelItem.Title,
					Type:        data.TypeEnforcedPayment,
					Date:        *enforcedPayment.DateOfSAP,
					EnforcedPayment: dto.DropdownSimple{
						ID:    enforcedPayment.ID,
						Title: *enforcedPayment.SAPID,
					},
				})
			}
		case data.EnforcedPaymentTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(enforcedPayment.Amount + enforcedPayment.AmountForAgent + enforcedPayment.AmountForLawyer + enforcedPayment.AmountForBank),
				Title:        modelItem.Title,
				Type:         data.TypeEnforcedPayment,
				Date:         *enforcedPayment.DateOfSAP,
				EnforcedPayment: dto.DropdownSimple{
					ID:    enforcedPayment.ID,
					Title: *enforcedPayment.SAPID,
				},
			})
		}

	}

	return response, nil
}

func buildAccountingOrderForReturnEnforcedPayment(id int, h *AccountingEntryServiceImpl) ([]dto.AccountingOrderItemsForObligations, error) {
	response := []dto.AccountingOrderItemsForObligations{}

	enforcedPayment, err := h.enforcedPaymentRepo.Get(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo enforced payment get")
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeReturnEnforcedPayment,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo model of accounting get all")
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, newErrors.Wrap(errors.ErrNotFound, "repo model of accounting get all")
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.EnforcedPaymentTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(*enforcedPayment.ReturnAmount),
				Title:       modelItem.Title,
				Type:        data.TypeReturnEnforcedPayment,
				Date:        *enforcedPayment.ReturnDate,
				ReturnEnforcedPayment: dto.DropdownSimple{
					ID:    enforcedPayment.ID,
					Title: *enforcedPayment.SAPID,
				},
			})
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(*enforcedPayment.ReturnAmount),
				Title:        modelItem.Title,
				Date:         *enforcedPayment.DateOfSAP,
				Type:         data.TypeReturnEnforcedPayment,
				ReturnEnforcedPayment: dto.DropdownSimple{
					ID:    enforcedPayment.ID,
					Title: *enforcedPayment.SAPID,
				},
				SupplierID: enforcedPayment.SupplierID,
			})
		}

	}

	return response, nil
}
