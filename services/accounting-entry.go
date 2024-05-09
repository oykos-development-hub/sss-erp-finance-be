package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

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

func (h *AccountingEntryServiceImpl) CreateAccountingEntry(input dto.AccountingEntryDTO) (*dto.AccountingEntryResponseDTO, error) {
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
			}
		}

		if dataToInsert.Type == "" {
			return errors.ErrInvalidInput
		}

		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range input.Items {
			itemToInsert := item.ToAccountingEntryItem()
			itemToInsert.EntryID = id

			_, err = h.items.Insert(tx, *itemToInsert)
			if err != nil {
				return err
			}

			boolTrue := true
			if item.Title == string(data.MainBillTitle) && item.InvoiceID != nil && *item.InvoiceID != 0 {
				invoice, err := h.invoiceRepo.Get(*item.InvoiceID)

				if err != nil {
					return err
				}

				invoice.Registred = &boolTrue

				err = h.invoiceRepo.Update(tx, *invoice)

				if err != nil {
					return err
				}
			} else if item.Title == string(data.MainBillTitle) && item.SalaryID != nil && *item.SalaryID != 0 {
				salary, err := h.salaryRepo.Get(*item.SalaryID)

				if err != nil {
					return err
				}

				salary.Registred = &boolTrue

				err = h.salaryRepo.Update(tx, *salary)

				if err != nil {
					return err
				}
			} else if item.Title == string(data.MainBillTitle) && item.PaymentOrderID != nil && *item.PaymentOrderID != 0 {
				paymentOrder, err := h.paymentOrderRepo.Get(*item.PaymentOrderID)

				if err != nil {
					return err
				}

				paymentOrder.Registred = &boolTrue

				err = h.paymentOrderRepo.Update(tx, *paymentOrder)

				if err != nil {
					return err
				}
			} else if item.Title == string(data.LawyerCostTitle) && item.EnforcedPaymentID != nil && *item.EnforcedPaymentID != 0 {
				enforcedPayment, err := h.enforcedPaymentRepo.Get(*item.EnforcedPaymentID)

				if err != nil {
					return err
				}

				enforcedPayment.Registred = &boolTrue

				err = h.enforcedPaymentRepo.Update(tx, *enforcedPayment)

				if err != nil {
					return err
				}
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

	res := dto.ToAccountingEntryResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *AccountingEntryServiceImpl) UpdateAccountingEntry(id int, input dto.AccountingEntryDTO) (*dto.AccountingEntryResponseDTO, error) {
	dataToInsert := input.ToAccountingEntry()
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

	response := dto.ToAccountingEntryResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *AccountingEntryServiceImpl) DeleteAccountingEntry(id int) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"entry_id": id})
		items, _, err := h.items.GetAll(nil, nil, conditionAndExp, nil)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, item := range items {
			boolFalse := false
			if item.Title == string(data.MainBillTitle) && item.InvoiceID != nil && *item.InvoiceID != 0 {
				invoice, err := h.invoiceRepo.Get(*item.InvoiceID)

				if err != nil {
					return err
				}

				invoice.Registred = &boolFalse

				err = h.invoiceRepo.Update(tx, *invoice)

				if err != nil {
					return err
				}
			} else if item.Title == string(data.MainBillTitle) && item.SalaryID != nil && *item.SalaryID != 0 {
				salary, err := h.salaryRepo.Get(*item.SalaryID)

				if err != nil {
					return err
				}

				salary.Registred = &boolFalse

				err = h.salaryRepo.Update(tx, *salary)

				if err != nil {
					return err
				}
			} else if item.Title == string(data.MainBillTitle) && item.PaymentOrderID != nil && *item.PaymentOrderID != 0 {
				paymentOrder, err := h.paymentOrderRepo.Get(*item.PaymentOrderID)

				if err != nil {
					return err
				}

				paymentOrder.Registred = &boolFalse

				err = h.paymentOrderRepo.Update(tx, *paymentOrder)

				if err != nil {
					return err
				}
			} else if item.Title == string(data.LawyerCostTitle) && item.EnforcedPaymentID != nil && *item.EnforcedPaymentID != 0 {
				enforcedPayment, err := h.enforcedPaymentRepo.Get(*item.EnforcedPaymentID)

				if err != nil {
					return err
				}

				enforcedPayment.Registred = &boolFalse

				err = h.enforcedPaymentRepo.Update(tx, *enforcedPayment)

				if err != nil {
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
	})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *AccountingEntryServiceImpl) GetAccountingEntry(id int) (*dto.AccountingEntryResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToAccountingEntryResponseDTO(*data)

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"entry_id": id})
	items, _, err := h.items.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}

	var debitAmount float64
	var creditAmount float64
	for _, item := range items {
		responseItem := dto.AccountingEntryItemResponseDTO{
			ID:             item.ID,
			Title:          item.Title,
			EntryID:        item.EntryID,
			AccountID:      item.AccountID,
			CreditAmount:   item.CreditAmount,
			DebitAmount:    item.DebitAmount,
			InvoiceID:      item.InvoiceID,
			SalaryID:       item.SalaryID,
			PaymentOrderID: item.PaymentOrderID,
			Type:           item.Type,
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
	response := dto.ToAccountingEntryListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"entry_id": response[i].ID})
		items, _, err := h.items.GetAll(nil, nil, conditionAndExp, nil)

		if err != nil {
			h.App.ErrorLog.Println(err)
			return nil, nil, errors.ErrNotFound
		}

		var debitAmount float64
		var creditAmount float64
		for _, item := range items {
			responseItem := dto.AccountingEntryItemResponseDTO{
				ID:             item.ID,
				Title:          item.Title,
				EntryID:        item.EntryID,
				AccountID:      item.AccountID,
				CreditAmount:   item.CreditAmount,
				DebitAmount:    item.DebitAmount,
				InvoiceID:      item.InvoiceID,
				SalaryID:       item.SalaryID,
				PaymentOrderID: item.PaymentOrderID,
				Type:           item.Type,
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

func (h *AccountingEntryServiceImpl) GetObligationsForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.ObligationForAccounting, *uint64, error) {

	dataFilter := data.ObligationsFilter{
		Page:               filter.Page,
		Size:               filter.Size,
		OrganizationUnitID: filter.OrganizationUnitID,
		Type:               filter.Type,
		Search:             filter.Search,
	}

	items, total, err := h.repo.GetObligationsForAccounting(dataFilter)

	if err != nil {
		return nil, nil, err
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
	}

	items, total, err := h.repo.GetPaymentOrdersForAccounting(dataFilter)

	if err != nil {
		return nil, nil, err
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
	}

	items, total, err := h.repo.GetEnforcedPaymentsForAccounting(dataFilter)

	if err != nil {
		return nil, nil, err
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
	}

	items, total, err := h.repo.GetReturnedEnforcedPaymentsForAccounting(dataFilter)

	if err != nil {
		return nil, nil, err
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
			return nil, err
		}

		switch invoice.Type {
		case data.TypeInvoice:
			item, err := buildAccountingOrderForInvoice(id, h)

			if err != nil {
				return nil, err
			}

			response.Items = append(response.Items, item...)
		case data.TypeDecision:
			item, err := buildAccountingOrderForDecisions(id, h)

			if err != nil {
				return nil, err
			}

			response.Items = append(response.Items, item...)
		case data.TypeContract:
			item, err := buildAccountingOrderForContracts(id, h)

			if err != nil {
				return nil, err
			}
			response.Items = append(response.Items, item...)
		}
	}

	for _, id := range orderData.SalaryID {
		item, err := buildAccountingOrderForSalaries(id, h)

		if err != nil {
			return nil, err
		}
		response.Items = append(response.Items, item...)
	}

	for _, id := range orderData.PaymentOrderID {
		item, err := buildAccountingOrderForPaymentOrder(id, h)

		if err != nil {
			return nil, err
		}
		response.Items = append(response.Items, item...)
	}

	for _, id := range orderData.EnforcedPaymentID {
		item, err := buildAccountingOrderForEnforcedPayment(id, h)

		if err != nil {
			return nil, err
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
		return nil, err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	articles, _, err := h.articlesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, err
	}

	var price float64
	for _, article := range articles {
		price += (article.NetPrice + float64(article.VatPercentage)/100*article.NetPrice) * float64(article.Amount)
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeInvoice,
	})

	if err != nil {
		return nil, err
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, errors.ErrInvalidInput
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.MainBillTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(price),
				Title:       modelItem.Title,
				Type:        data.TypeInvoice,
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
		return nil, err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	additionalExpenses, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, err
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
	for _, item := range additionalExpenses {
		price += float64(item.Price)
		switch item.Title {
		case data.NetTitle:
			netPrice += float64(item.Price)
		case data.ObligationTaxTitle:
			taxPrice += float64(item.Price)
		case data.ObligationSubTaxTitle:
			subTaxPrice += float64(item.Price)
		case data.ContributionForPIOTitle:
			contributionForPIO += float64(item.Price)
		case data.ContributionForUnemploymentTitle:
			contributionForUnemployment += float64(item.Price)
		case data.ContributionForPIOEmployeeTitle:
			contributionForPIOEmployee += float64(item.Price)
		case data.ContributionForPIOEmployerTitle:
			contributionForPIOEmployer += float64(item.Price)
		case data.ContributionForUnemploymentEmployeeTitle:
			contributionForUnemploymentEmployee += float64(item.Price)
		case data.ContributionForUnemploymentEmployerTitle:
			contributionForUnemploymentEmployer += float64(item.Price)
		case data.LaborFundTitle:
			contributionForLaborFund += float64(item.Price)
		}
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeDecision,
	})

	if err != nil {
		return nil, err
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, errors.ErrInvalidInput
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.MainBillTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(price),
				Title:       modelItem.Title,
				Type:        data.TypeDecision,
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
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.SubTaxTitle:
			if subTaxPrice > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(subTaxPrice),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.LaborContributionsTitle:
			if contributionForLaborFund > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForLaborFund),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.PIOContributionsTitle:
			if contributionForPIO > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIO),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.UnemployementContributionsTitle:
			if contributionForUnemploymentEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemployment),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.PIOEmployeeContributionsTitle:
			if contributionForPIOEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.PIOEmployerContributionsTitle:
			if contributionForPIOEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.UnemployementEmployeeContributionsTitle:
			if contributionForUnemploymentEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.UnemployementEmployerContributionsTitle:
			if contributionForUnemploymentEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeDecision,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
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
		return nil, err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": id})
	additionalExpenses, _, err := h.additionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, err
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
	for _, item := range additionalExpenses {
		price += float64(item.Price)
		switch item.Title {
		case data.NetTitle:
			netPrice += float64(item.Price)
		case data.ObligationTaxTitle:
			taxPrice += float64(item.Price)
		case data.ObligationSubTaxTitle:
			subTaxPrice += float64(item.Price)
		case data.ContributionForPIOTitle:
			contributionForPIO += float64(item.Price)
		case data.ContributionForUnemploymentTitle:
			contributionForUnemployment += float64(item.Price)
		case data.ContributionForPIOEmployeeTitle:
			contributionForPIOEmployee += float64(item.Price)
		case data.ContributionForPIOEmployerTitle:
			contributionForPIOEmployer += float64(item.Price)
		case data.ContributionForUnemploymentEmployeeTitle:
			contributionForUnemploymentEmployee += float64(item.Price)
		case data.ContributionForUnemploymentEmployerTitle:
			contributionForUnemploymentEmployer += float64(item.Price)
		case data.LaborFundTitle:
			contributionForLaborFund += float64(item.Price)
		}
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeContract,
	})

	if err != nil {
		return nil, err
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, errors.ErrInvalidInput
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.MainBillTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(price),
				Title:       modelItem.Title,
				Type:        data.TypeContract,
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
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.SubTaxTitle:
			if subTaxPrice > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(subTaxPrice),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.LaborContributionsTitle:
			if contributionForLaborFund > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForLaborFund),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.PIOContributionsTitle:
			if contributionForPIO > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIO),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.UnemployementContributionsTitle:
			if contributionForUnemploymentEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemployment),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.PIOEmployeeContributionsTitle:
			if contributionForPIOEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.PIOEmployerContributionsTitle:
			if contributionForPIOEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForPIOEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.UnemployementEmployeeContributionsTitle:
			if contributionForUnemploymentEmployee > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployee),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
				})
			}
		case data.UnemployementEmployerContributionsTitle:
			if contributionForUnemploymentEmployer > 0 {
				response = append(response, dto.AccountingOrderItemsForObligations{
					AccountID:    modelItem.CreditAccountID,
					CreditAmount: float32(contributionForUnemploymentEmployer),
					Title:        modelItem.Title,
					Type:         data.TypeContract,
					Invoice: dto.DropdownSimple{
						ID:    invoice.ID,
						Title: invoice.InvoiceNumber,
					},
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
		return nil, err
	}

	conditionAndExp := &up.AndExpr{}
	conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_id": id})
	additionalExpenses, _, err := h.salaryAdditionalExpensesRepo.GetAll(nil, nil, conditionAndExp, nil)

	if err != nil {
		return nil, err
	}

	models, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeSalary,
	})

	if err != nil {
		return nil, err
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(models) != 1 {
		return nil, errors.ErrInvalidInput
	}

	var price float64
	var taxPrice float64

	for _, item := range additionalExpenses {
		price += float64(item.Amount)
		bookedItem := &dto.AccountingOrderItemsForObligations{}
		switch item.Type {
		case data.SalaryAdditionalExpenseType(data.ContributionsSalaryExpenseType):
			switch item.Title {
			case data.PIOEmployeeContributionsTitle:
				bookedItem = buildBookedItemForSalary(item, models[0].Items, data.PIOEmployeeContributionsTitle)
			case data.PIOEmployerContributionsTitle:
				bookedItem = buildBookedItemForSalary(item, models[0].Items, data.PIOEmployerContributionsTitle)
			case data.UnemployementEmployeeContributionsTitle:
				bookedItem = buildBookedItemForSalary(item, models[0].Items, data.UnemployementEmployeeContributionsTitle)
			case data.UnemployementEmployerContributionsTitle:
				bookedItem = buildBookedItemForSalary(item, models[0].Items, data.UnemployementEmployerContributionsTitle)
			case data.LaborContributionsTitle:
				bookedItem = buildBookedItemForSalary(item, models[0].Items, data.LaborContributionsTitle)
			}
		case data.SalaryAdditionalExpenseType(data.TaxesSalaryExpenseType):
			taxPrice += item.Amount
		case data.SalaryAdditionalExpenseType(data.SubTaxesSalaryExpenseType):
			bookedItem = buildBookedItemForSalary(item, models[0].Items, "Prirez")
		case data.SalaryAdditionalExpenseType(data.BanksSalaryExpenseType):
			bookedItem = buildBookedItemForSalaryForBanks(item, models[0].Items, "Banka")
		case data.SalaryAdditionalExpenseType(data.SuspensionsSalaryExpenseType):
			bookedItem = buildBookedItemForSalaryForBanks(item, models[0].Items, "Obustave")
		}
		if bookedItem != nil {
			bookedItem.Salary = dto.DropdownSimple{
				ID:    salary.ID,
				Title: salary.Month,
			}
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
				Salary: dto.DropdownSimple{
					ID:    salary.ID,
					Title: salary.Month,
				},
			})

			response = append(newSlice, response...)

		} else if model.Title == data.TaxTitle {
			if taxPrice > 0 {
				bookedItem := buildBookedItemForSalary(&data.SalaryAdditionalExpense{Amount: taxPrice}, models[0].Items, "Porez")
				response = append(response, *bookedItem)
			}
		}

	}

	return response, nil
}

func buildBookedItemForSalary(item *data.SalaryAdditionalExpense, models []dto.ModelOfAccountingItemResponseDTO, title data.AccountingOrderItemsTitle) *dto.AccountingOrderItemsForObligations {

	for _, model := range models {
		if string(model.Title) == string(title) {
			response := dto.AccountingOrderItemsForObligations{
				AccountID:    model.CreditAccountID,
				CreditAmount: float32(item.Amount),
				Title:        model.Title,
				Type:         data.TypeSalary,
			}
			return &response
		}
	}

	return nil
}

func buildBookedItemForSalaryForBanks(item *data.SalaryAdditionalExpense, models []dto.ModelOfAccountingItemResponseDTO, title string) *dto.AccountingOrderItemsForObligations {

	for _, model := range models {
		if string(model.Title) == title {
			response := dto.AccountingOrderItemsForObligations{
				AccountID:    model.CreditAccountID,
				CreditAmount: float32(item.Amount),
				Title:        item.Title,
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
		return nil, err
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypePaymentOrder,
	})

	if err != nil {
		return nil, err
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, errors.ErrInvalidInput
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(paymentOrder.Amount),
				Title:       modelItem.Title,
				Type:        data.TypePaymentOrder,
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
		return nil, err
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeEnforcedPayment,
	})

	if err != nil {
		return nil, err
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, errors.ErrInvalidInput
	}

	for _, modelItem := range model[0].Items {
		switch modelItem.Title {
		case data.SupplierTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:   modelItem.DebitAccountID,
				DebitAmount: float32(enforcedPayment.Amount),
				Title:       modelItem.Title,
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
					EnforcedPayment: dto.DropdownSimple{
						ID:    enforcedPayment.ID,
						Title: *enforcedPayment.SAPID,
					},
				})
			}
		case data.EnforcedPaymentTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(enforcedPayment.Amount + enforcedPayment.AmountForAgent + enforcedPayment.AmountForLawyer),
				Title:        modelItem.Title,
				Type:         data.TypeEnforcedPayment,
				EnforcedPayment: dto.DropdownSimple{
					ID:    enforcedPayment.ID,
					Title: *enforcedPayment.SAPID,
				},
			})
		}

	}

	return response, nil
}
