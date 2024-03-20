package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type AdditionalExpenseServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.AdditionalExpense
}

func NewAdditionalExpenseServiceImpl(app *celeritas.Celeritas, repo data.AdditionalExpense) AdditionalExpenseService {
	return &AdditionalExpenseServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *AdditionalExpenseServiceImpl) CreateAdditionalExpense(input dto.AdditionalExpenseDTO) (*dto.AdditionalExpenseResponseDTO, error) {
	data := input.ToAdditionalExpense()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToAdditionalExpenseResponseDTO(*data)

	return &res, nil
}

func (h *AdditionalExpenseServiceImpl) UpdateAdditionalExpense(id int, input dto.AdditionalExpenseDTO) (*dto.AdditionalExpenseResponseDTO, error) {
	data := input.ToAdditionalExpense()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToAdditionalExpenseResponseDTO(*data)

	return &response, nil
}

func (h *AdditionalExpenseServiceImpl) DeleteAdditionalExpense(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *AdditionalExpenseServiceImpl) GetAdditionalExpense(id int) (*dto.AdditionalExpenseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToAdditionalExpenseResponseDTO(*data)

	return &response, nil
}

func (h *AdditionalExpenseServiceImpl) GetAdditionalExpenseList(filter dto.AdditionalExpenseFilterDTO) ([]dto.AdditionalExpenseResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.InvoiceID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": *filter.InvoiceID})
	}

	if filter.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	if filter.SubjectID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"subject_id": *filter.SubjectID})
	}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Search != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"title ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	/*if filter.Year != nil {
		year := *filter.Year
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_invoice": up.Between(startOfYear, endOfYear)})
	}*/

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
	response := dto.ToAdditionalExpenseListResponseDTO(data)

	return response, total, nil
}
