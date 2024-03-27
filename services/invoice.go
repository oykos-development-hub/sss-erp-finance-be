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

type InvoiceServiceImpl struct {
	App                    *celeritas.Celeritas
	repo                   data.Invoice
	additionalExpensesRepo data.AdditionalExpense
	articlesServices       ArticleService
	additionalExpenses     AdditionalExpenseService
}

func NewInvoiceServiceImpl(app *celeritas.Celeritas, repo data.Invoice, additionalExpensesRepo data.AdditionalExpense, articles ArticleService, additionalExpenses AdditionalExpenseService) InvoiceService {
	return &InvoiceServiceImpl{
		App:                    app,
		repo:                   repo,
		additionalExpensesRepo: additionalExpensesRepo,
		articlesServices:       articles,
		additionalExpenses:     additionalExpenses,
	}
}

func (h *InvoiceServiceImpl) CreateInvoice(input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error) {
	invoice := input.ToInvoice()

	if invoice.SSSInvoiceReceiptDate != nil {
		invoice.Status = "waiting"
	} else {
		invoice.Status = "created"
	}

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *invoice)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, additionalExpense := range input.AdditionalExpenses {
			additionalExpenseData := additionalExpense.ToAdditionalExpense()
			additionalExpenseData.InvoiceID = id
			if _, err = h.additionalExpensesRepo.Insert(tx, *additionalExpenseData); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response, err := h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToInvoiceResponseDTO(*response)

	return &res, nil
}

func (h *InvoiceServiceImpl) UpdateInvoice(id int, input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error) {
	invoice := input.ToInvoice()
	invoice.ID = id

	if invoice.SSSInvoiceReceiptDate != nil {
		invoice.Status = "waiting"
	} else {
		invoice.Status = "created"
	}

	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		err = h.repo.Update(tx, *invoice)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, additionalExpense := range input.AdditionalExpenses {
			additionalExpenseData := additionalExpense.ToAdditionalExpense()
			if err = h.additionalExpensesRepo.Update(tx, *additionalExpenseData); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToInvoiceResponseDTO(*invoice)

	return &response, nil
}

func (h *InvoiceServiceImpl) DeleteInvoice(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *InvoiceServiceImpl) GetInvoice(id int) (*dto.InvoiceResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToInvoiceResponseDTO(*data)

	articles, _, err := h.articlesServices.GetArticleList(dto.ArticleFilterDTO{InvoiceID: &id})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	response.Articles = articles

	additionaExpenses, _, err := h.additionalExpenses.GetAdditionalExpenseList(dto.AdditionalExpenseFilterDTO{InvoiceID: &id})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	response.AdditionalExpenses = additionaExpenses

	return &response, nil
}

func (h *InvoiceServiceImpl) GetInvoiceList(input dto.InvoicesFilter) ([]dto.InvoiceResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}

	if input.Year != nil {
		year := *input.Year
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_invoice": up.Between(startOfYear, endOfYear)})
	}

	if input.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *input.Status})
	}

	if input.SupplierID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"supplier_id": *input.SupplierID})
	}

	if input.ActivityID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"activity_id": *input.ActivityID})
	}

	if input.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *input.OrganizationUnitID})
	}

	if input.Type != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"type": *input.Type})
	}

	if input.Search != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Search)
		search := up.Or(
			up.Cond{"invoice_number ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	data, total, err := h.repo.GetAll(input.Page, input.Size, conditionAndExp)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, err
	}
	response := dto.ToInvoiceListResponseDTO(data)

	return response, total, nil
}
