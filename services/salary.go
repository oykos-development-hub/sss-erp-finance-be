package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type SalaryServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.Salary
}

func NewSalaryServiceImpl(app *celeritas.Celeritas, repo data.Salary) SalaryService {
	return &SalaryServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *SalaryServiceImpl) CreateSalary(input dto.SalaryDTO) (*dto.SalaryResponseDTO, error) {
	dataToInsert := input.ToSalary()

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

	res := dto.ToSalaryResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *SalaryServiceImpl) UpdateSalary(id int, input dto.SalaryDTO) (*dto.SalaryResponseDTO, error) {
	dataToInsert := input.ToSalary()
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

	response := dto.ToSalaryResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *SalaryServiceImpl) DeleteSalary(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *SalaryServiceImpl) GetSalary(id int) (*dto.SalaryResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToSalaryResponseDTO(*data)

	return &response, nil
}

func (h *SalaryServiceImpl) GetSalaryList(filter dto.SalaryFilterDTO) ([]dto.SalaryResponseDTO, *uint64, error) {
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
	response := dto.ToSalaryListResponseDTO(data)

	return response, total, nil
}
