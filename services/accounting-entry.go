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
}

func NewAccountingEntryServiceImpl(app *celeritas.Celeritas, repo data.AccountingEntry, invoiceRepo data.Invoice, articlesRepo data.Article, additionalExpensesRepo data.AdditionalExpense, salaryRepo data.Salary, salaryAdditionalExpensesRepo data.SalaryAdditionalExpense, modelOfAccountingRepo ModelsOfAccountingService) AccountingEntryService {
	return &AccountingEntryServiceImpl{
		App:                          app,
		repo:                         repo,
		invoiceRepo:                  invoiceRepo,
		articlesRepo:                 articlesRepo,
		additionalExpensesRepo:       additionalExpensesRepo,
		salaryRepo:                   salaryRepo,
		salaryAdditionalExpensesRepo: salaryAdditionalExpensesRepo,
		modelOfAccountingRepo:        modelOfAccountingRepo,
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
	err := h.repo.Delete(id)
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

	return &response, nil
}

func (h *AccountingEntryServiceImpl) GetAccountingEntryList(filter dto.AccountingEntryFilterDTO) ([]dto.AccountingEntryResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	// if filter.Year != nil {
	// 	conditionAndExp = up.And(conditionAndExp, &up.Cond{"year": *filter.Year})
	// }

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
			InvoiceID: item.InvoiceID,
			SalaryID:  item.SalaryID,
			Type:      item.Type,
			Title:     item.Title,
			Status:    item.Status,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		})
	}

	return response, total, nil
}

func (h *AccountingEntryServiceImpl) BuildAccountingOrderForObligations(orderData dto.AccountingOrderForObligationsData) (*dto.AccountingOrderForObligations, error) {

	response := dto.AccountingOrderForObligations{}

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
		case data.TypeContract:
		}
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
		price += article.NetPrice + article.VatPrice/100*article.NetPrice
	}

	model, _, err := h.modelOfAccountingRepo.GetModelsOfAccountingList(dto.ModelsOfAccountingFilterDTO{
		Type: &data.TypeInvoice,
	})

	if err != nil {
		return nil, err
	}

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
