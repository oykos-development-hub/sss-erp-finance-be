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

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		if invoice.Status != data.InvoiceStatusIncomplete {
			invoice.Status = data.InvoiceStatusCreated
		}
		id, err = h.repo.Insert(tx, *invoice)
		if err != nil {
			return errors.ErrInternalServer
		}

		for _, additionalExpense := range input.AdditionalExpenses {
			additionalExpenseData := additionalExpense.ToAdditionalExpense()
			additionalExpenseData.InvoiceID = id
			additionalExpenseData.OrganizationUnitID = input.OrganizationUnitID
			additionalExpenseData.Status = data.InvoiceStatusCreated
			if additionalExpenseData.Price > 0 {
				if _, err = h.additionalExpensesRepo.Insert(tx, *additionalExpenseData); err != nil {
					return err
				}
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

	oldData, err := h.GetInvoice(id)

	if err != nil {
		return nil, err
	}

	if input.Type == data.TypeInvoice && oldData.Status == data.InvoiceStatusIncomplete {
		statusCreated := true
		for _, article := range oldData.Articles {
			if article.AccountID == 0 {
				statusCreated = false
			}
		}

		if statusCreated {
			input.Status = data.InvoiceStatusCreated
		}
	}

	err = data.Upper.Tx(func(tx up.Session) error {
		var err error
		err = h.repo.Update(tx, *invoice)
		if err != nil {
			return errors.ErrInternalServer
		}

		validExpenses := make(map[int]bool)

		for _, item := range oldData.AdditionalExpenses {
			validExpenses[item.ID] = false
		}

		for _, item := range input.AdditionalExpenses {
			_, exists := validExpenses[item.ID]
			if exists {
				validExpenses[item.ID] = true
			} else {
				additionalExpenseData := item.ToAdditionalExpense()
				additionalExpenseData.InvoiceID = id
				additionalExpenseData.OrganizationUnitID = input.OrganizationUnitID
				additionalExpenseData.Status = data.InvoiceStatusCreated
				if additionalExpenseData.Price > 0 {
					_, err = h.additionalExpensesRepo.Insert(tx, *additionalExpenseData)

					if err != nil {
						return err
					}
				}
			}
		}

		for itemID, exists := range validExpenses {
			if !exists {
				err := h.additionalExpensesRepo.Delete(itemID)

				if err != nil {
					return err
				}
			} else {
				for _, item := range input.AdditionalExpenses {
					if item.ID == itemID {
						additionalExpenseData := item.ToAdditionalExpense()
						additionalExpenseData.ID = item.ID
						additionalExpenseData.InvoiceID = id
						additionalExpenseData.OrganizationUnitID = input.OrganizationUnitID
						additionalExpenseData.Status = data.InvoiceStatusCreated
						err := h.additionalExpensesRepo.Update(tx, *additionalExpenseData)
						if err != nil {
							return err
						}
					}
				}
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

	additionalExpenses, _, err := h.additionalExpenses.GetAdditionalExpenseList(dto.AdditionalExpenseFilterDTO{InvoiceID: &id})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	response.AdditionalExpenses = additionalExpenses

	if len(additionalExpenses) > 0 {
		response.Status = additionalExpenses[len(additionalExpenses)-1].Status
		response.NetPrice = float64(additionalExpenses[len(additionalExpenses)-1].Price)
	}

	for j := 0; j < len(additionalExpenses)-1; j++ {
		response.VATPrice += float64(additionalExpenses[j].Price)
	}

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

	if input.Status != nil && *input.Status != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *input.Status})
	}

	if input.SupplierID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"supplier_id": *input.SupplierID})
	}

	if input.OrderID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"order_id": *input.OrderID})
	}

	if input.ActivityID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"activity_id": *input.ActivityID})
	}

	if input.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *input.OrganizationUnitID})
	}

	if input.Registred != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"registred": *input.Registred})
	}

	if input.Type != nil && *input.Type != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *input.Type)
		search := up.Or(
			up.Cond{"type ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	if input.PassedToInventory != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"passed_to_inventory": *input.PassedToInventory})
	}

	if input.Search != nil && *input.Search != "" {
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

	for i := 0; i < len(response); i++ {
		additionalExpenses, _, err := h.additionalExpenses.GetAdditionalExpenseList(dto.AdditionalExpenseFilterDTO{
			InvoiceID: &response[i].ID,
		})

		if err != nil {
			h.App.ErrorLog.Println(err)
			return nil, nil, err
		}

		if len(additionalExpenses) > 0 {
			response[i].NetPrice = float64(additionalExpenses[len(additionalExpenses)-1].Price)
			response[i].Status = additionalExpenses[len(additionalExpenses)-1].Status
		}

		for j := 0; j < len(additionalExpenses)-1; j++ {
			response[i].VATPrice += float64(additionalExpenses[j].Price)
		}
	}

	return response, total, nil
}
