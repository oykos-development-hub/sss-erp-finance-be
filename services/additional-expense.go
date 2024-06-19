package services

import (
	"fmt"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type AdditionalExpenseServiceImpl struct {
	App         *celeritas.Celeritas
	repo        data.AdditionalExpense
	invoiceRepo data.Invoice
}

func NewAdditionalExpenseServiceImpl(app *celeritas.Celeritas, repo data.AdditionalExpense, invoiceRepo data.Invoice) AdditionalExpenseService {
	return &AdditionalExpenseServiceImpl{
		App:         app,
		repo:        repo,
		invoiceRepo: invoiceRepo,
	}
}

func (h *AdditionalExpenseServiceImpl) DeleteAdditionalExpense(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		return newErrors.Wrap(err, "repo additional expense delete")
	}

	return nil
}

func (h *AdditionalExpenseServiceImpl) GetAdditionalExpense(id int) (*dto.AdditionalExpenseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo additional expense get")
	}

	response := dto.ToAdditionalExpenseResponseDTO(*data)

	return &response, nil
}

func (h *AdditionalExpenseServiceImpl) GetAdditionalExpenseList(filter dto.AdditionalExpenseFilterDTO) ([]dto.AdditionalExpenseResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.InvoiceID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": *filter.InvoiceID})
	} else {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"title <> ": "Neto"})
	}

	if filter.Status != nil && *filter.Status != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	if filter.SubjectID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"subject_id": *filter.SubjectID})
	}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"title ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	if filter.Year != nil {
		year := *filter.Year
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"created_at": up.Between(startOfYear, endOfYear)})
	}

	if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}

	orders = append(orders, "-created_at")

	items, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo additional expenses get all")
	}
	response := dto.ToAdditionalExpenseListResponseDTO(items)

	for i := 0; i < len(response); i++ {
		invoice, err := h.invoiceRepo.Get(response[i].InvoiceID)

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo invoice get all")
		}

		response[i].ObligationType = invoice.Type
		response[i].ObligationNumber = invoice.InvoiceNumber
		response[i].ObligationSupplierID = invoice.SupplierID

	}

	return response, total, nil
}
