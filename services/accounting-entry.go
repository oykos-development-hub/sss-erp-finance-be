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
}

func NewAccountingEntryServiceImpl(app *celeritas.Celeritas, repo data.AccountingEntry, invoiceRepo data.Invoice, articlesRepo data.Article, additionalExpensesRepo data.AdditionalExpense, salaryRepo data.Salary, salaryAdditionalExpensesRepo data.SalaryAdditionalExpense, modelOfAccountingRepo ModelsOfAccountingService, items data.AccountingEntryItem) AccountingEntryService {
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
	}
}

func (h *AccountingEntryServiceImpl) CreateAccountingEntry(input dto.AccountingEntryDTO) (*dto.AccountingEntryResponseDTO, error) {
	dataToInsert := input.ToAccountingEntry()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
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
			ID:           item.ID,
			Title:        item.Title,
			EntryID:      item.EntryID,
			AccountID:    item.AccountID,
			CreditAmount: item.CreditAmount,
			DebitAmount:  item.DebitAmount,
			InvoiceID:    item.InvoiceID,
			SalaryID:     item.SalaryID,
			Type:         item.Type,
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
				ID:           item.ID,
				Title:        item.Title,
				EntryID:      item.EntryID,
				AccountID:    item.AccountID,
				CreditAmount: item.CreditAmount,
				DebitAmount:  item.DebitAmount,
				InvoiceID:    item.InvoiceID,
				SalaryID:     item.SalaryID,
				Type:         item.Type,
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
				Title:       string(modelItem.Title),
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
				Title:        string(modelItem.Title),
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
	for _, item := range additionalExpenses {
		price += float64(item.Price)
		switch item.Title {
		case "Neto":
			netPrice += float64(item.Price)
		case "Porez":
			taxPrice += float64(item.Price)
		case "Prirez":
			subTaxPrice += float64(item.Price)
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
				Title:       string(modelItem.Title),
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
				Title:        string(modelItem.Title),
				Type:         data.TypeDecision,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
				SupplierID: invoice.SupplierID,
			})
		case data.TaxTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(taxPrice),
				Title:        string(modelItem.Title),
				Type:         data.TypeDecision,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
			})
		case data.SubTaxTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(subTaxPrice),
				Title:        string(modelItem.Title),
				Type:         data.TypeDecision,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
			})
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
	for _, item := range additionalExpenses {
		price += float64(item.Price)
		switch item.Title {
		case "Neto":
			netPrice += float64(item.Price)
		case "Porez":
			taxPrice += float64(item.Price)
		case "Prirez":
			subTaxPrice += float64(item.Price)
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
				Title:       string(modelItem.Title),
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
				Title:        string(modelItem.Title),
				Type:         data.TypeContract,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
				SupplierID: invoice.SupplierID,
			})
		case data.TaxTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(taxPrice),
				Title:        string(modelItem.Title),
				Type:         data.TypeContract,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
			})
		case data.SubTaxTitle:
			response = append(response, dto.AccountingOrderItemsForObligations{
				AccountID:    modelItem.CreditAccountID,
				CreditAmount: float32(subTaxPrice),
				Title:        string(modelItem.Title),
				Type:         data.TypeContract,
				Invoice: dto.DropdownSimple{
					ID:    invoice.ID,
					Title: invoice.InvoiceNumber,
				},
			})
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

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeSalary,
	})

	if err != nil {
		return nil, err
	}

	//ako ne postoji u bazi odgovarajuci model za taj tip obaveze/naloga vraca se invalid input
	if len(model) != 1 {
		return nil, errors.ErrInvalidInput
	}

	var price float64

	for _, item := range additionalExpenses {
		price += float64(item.Amount)
		bookedItem := &dto.AccountingOrderItemsForObligations{}
		switch item.Type {
		case data.SalaryAdditionalExpenseType(data.ContributionsSalaryExpenseType):
			switch item.Title {
			case data.PIOEmployeeContributionsSalaryTitle:
				bookedItem = buildBookedItemForSalary(item, model[0].Items, data.PIOEmployeeContributionsSalaryTitle)
			case data.PIOEmployerContributionsSalaryTitle:
				bookedItem = buildBookedItemForSalary(item, model[0].Items, data.PIOEmployerContributionsSalaryTitle)
			case data.UnemployementEmployeeContributionsSalaryTitle:
				bookedItem = buildBookedItemForSalary(item, model[0].Items, data.UnemployementEmployeeContributionsSalaryTitle)
			case data.UnemployementEmployerContributionsSalaryTitle:
				bookedItem = buildBookedItemForSalary(item, model[0].Items, data.UnemployementEmployerContributionsSalaryTitle)
			case data.LaborContributionsSalaryTitle:
				bookedItem = buildBookedItemForSalary(item, model[0].Items, data.LaborContributionsSalaryTitle)
			}
		case data.SalaryAdditionalExpenseType(data.TaxesSalaryExpenseType):
			bookedItem = buildBookedItemForSalary(item, model[0].Items, "Porez")
		case data.SalaryAdditionalExpenseType(data.SubTaxesSalaryExpenseType):
			bookedItem = buildBookedItemForSalary(item, model[0].Items, "Prirez")
		case data.SalaryAdditionalExpenseType(data.BanksSalaryExpenseType):
			bookedItem = buildBookedItemForSalaryForBanks(item, model[0].Items, "Banka")
		case data.SalaryAdditionalExpenseType(data.SuspensionsSalaryExpenseType):
			bookedItem = buildBookedItemForSalaryForBanks(item, model[0].Items, "Obustave")
		}
		if bookedItem != nil {
			bookedItem.Salary = dto.DropdownSimple{
				ID:    salary.ID,
				Title: salary.Month,
			}
			response = append(response, *bookedItem)
		}
	}

	for _, model := range model[0].Items {
		if model.Title == data.MainBillTitle {

			newSlice := make([]dto.AccountingOrderItemsForObligations, 0)

			newSlice = append(newSlice, dto.AccountingOrderItemsForObligations{
				AccountID:   model.DebitAccountID,
				DebitAmount: float32(price),
				Title:       string(model.Title),
				Type:        data.TypeSalary,
				Salary: dto.DropdownSimple{
					ID:    salary.ID,
					Title: salary.Month,
				},
			})

			response = append(newSlice, response...)

		}

	}

	return response, nil
}

func buildBookedItemForSalary(item *data.SalaryAdditionalExpense, models []dto.ModelOfAccountingItemResponseDTO, title string) *dto.AccountingOrderItemsForObligations {

	for _, model := range models {
		if string(model.Title) == title {
			response := dto.AccountingOrderItemsForObligations{
				AccountID:    model.CreditAccountID,
				CreditAmount: float32(item.Amount),
				Title:        item.Title,
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
