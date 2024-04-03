package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type SalaryAdditionalExpenseServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.SalaryAdditionalExpense
}

func NewSalaryAdditionalExpenseServiceImpl(app *celeritas.Celeritas, repo data.SalaryAdditionalExpense) SalaryAdditionalExpenseService {
	return &SalaryAdditionalExpenseServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *SalaryAdditionalExpenseServiceImpl) CreateSalaryAdditionalExpense(input dto.SalaryAdditionalExpenseDTO) (*dto.SalaryAdditionalExpenseResponseDTO, error) {
	dataToInsert := input.ToSalaryAdditionalExpense()

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

	res := dto.ToSalaryAdditionalExpenseResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *SalaryAdditionalExpenseServiceImpl) UpdateSalaryAdditionalExpense(id int, input dto.SalaryAdditionalExpenseDTO) (*dto.SalaryAdditionalExpenseResponseDTO, error) {
	dataToInsert := input.ToSalaryAdditionalExpense()
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

	response := dto.ToSalaryAdditionalExpenseResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *SalaryAdditionalExpenseServiceImpl) DeleteSalaryAdditionalExpense(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *SalaryAdditionalExpenseServiceImpl) GetSalaryAdditionalExpense(id int) (*dto.SalaryAdditionalExpenseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToSalaryAdditionalExpenseResponseDTO(*data)

	return &response, nil
}

func (h *SalaryAdditionalExpenseServiceImpl) GetSalaryAdditionalExpenseList(filter dto.SalaryAdditionalExpenseFilterDTO) ([]dto.SalaryAdditionalExpenseResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.SalaryID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"salary_id": *filter.SalaryID})
	}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	/*if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}*/

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToSalaryAdditionalExpenseListResponseDTO(data)

	return response, total, nil
}
