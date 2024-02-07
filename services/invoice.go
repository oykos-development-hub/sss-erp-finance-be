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
	App              *celeritas.Celeritas
	repo             data.Invoice
	articlesServices ArticleService
}

func NewInvoiceServiceImpl(app *celeritas.Celeritas, repo data.Invoice, articles ArticleService) InvoiceService {
	return &InvoiceServiceImpl{
		App:              app,
		repo:             repo,
		articlesServices: articles,
	}
}

func (h *InvoiceServiceImpl) CreateInvoice(input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error) {
	data := input.ToInvoice()

	if data.SSSInvoiceReceiptDate != nil {
		data.Status = "waiting"
	} else {
		data.Status = "created"
	}

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToInvoiceResponseDTO(*data)

	return &res, nil
}

func (h *InvoiceServiceImpl) UpdateInvoice(id int, input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error) {
	data := input.ToInvoice()
	data.ID = id

	if data.SSSInvoiceReceiptDate != nil {
		data.Status = "waiting"
	} else {
		data.Status = "created"
	}

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToInvoiceResponseDTO(*data)

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

	if input.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *input.OrganizationUnitID})
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
