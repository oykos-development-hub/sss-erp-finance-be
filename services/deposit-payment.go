package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type DepositPaymentServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.DepositPayment
}

func NewDepositPaymentServiceImpl(app *celeritas.Celeritas, repo data.DepositPayment) DepositPaymentService {
	return &DepositPaymentServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *DepositPaymentServiceImpl) CreateDepositPayment(input dto.DepositPaymentDTO) (*dto.DepositPaymentResponseDTO, error) {
	dataToInsert := input.ToDepositPayment()

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

	res := dto.ToDepositPaymentResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *DepositPaymentServiceImpl) UpdateDepositPayment(id int, input dto.DepositPaymentDTO) (*dto.DepositPaymentResponseDTO, error) {
	dataToInsert := input.ToDepositPayment()
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

	response := dto.ToDepositPaymentResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *DepositPaymentServiceImpl) DeleteDepositPayment(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *DepositPaymentServiceImpl) GetDepositPayment(id int) (*dto.DepositPaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToDepositPaymentResponseDTO(*data)

	return &response, nil
}

func (h *DepositPaymentServiceImpl) GetDepositPaymentList(filter dto.DepositPaymentFilterDTO) ([]dto.DepositPaymentResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Status != nil {
		switch *filter.Status {
		case "Prolazni račun":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"main_bank_account": true})
		case "Prelazni račun":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"main_bank_account": false})
		}
	}

	if filter.Search != nil {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"payer ILIKE": likeCondition},
			up.Cond{"case_number ILIKE": likeCondition},
			up.Cond{"party_name ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	//if filter.SortByTitle != nil {
	//	if *filter.SortByTitle == "asc" {
	//		orders = append(orders, "-title")
	//	} else {
	//		orders = append(orders, "title")
	//	}
	//}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToDepositPaymentListResponseDTO(data)

	return response, total, nil
}
