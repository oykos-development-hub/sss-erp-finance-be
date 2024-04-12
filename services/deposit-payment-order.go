package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type DepositPaymentOrderServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.DepositPaymentOrder
}

func NewDepositPaymentOrderServiceImpl(app *celeritas.Celeritas, repo data.DepositPaymentOrder) DepositPaymentOrderService {
	return &DepositPaymentOrderServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *DepositPaymentOrderServiceImpl) CreateDepositPaymentOrder(input dto.DepositPaymentOrderDTO) (*dto.DepositPaymentOrderResponseDTO, error) {
	dataToInsert := input.ToDepositPaymentOrder()

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

	res := dto.ToDepositPaymentOrderResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *DepositPaymentOrderServiceImpl) UpdateDepositPaymentOrder(id int, input dto.DepositPaymentOrderDTO) (*dto.DepositPaymentOrderResponseDTO, error) {
	dataToInsert := input.ToDepositPaymentOrder()
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

	response := dto.ToDepositPaymentOrderResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *DepositPaymentOrderServiceImpl) DeleteDepositPaymentOrder(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *DepositPaymentOrderServiceImpl) GetDepositPaymentOrder(id int) (*dto.DepositPaymentOrderResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToDepositPaymentOrderResponseDTO(*data)

	return &response, nil
}

func (h *DepositPaymentOrderServiceImpl) GetDepositPaymentOrderList(filter dto.DepositPaymentOrderFilterDTO) ([]dto.DepositPaymentOrderResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.CaseNumber != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"case_number": *filter.CaseNumber})
	}

	if filter.SupplierID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"supplier_id": *filter.SupplierID})
	}

	if filter.Status != nil {
		switch *filter.Status {
		case "Plaćen":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"id_of_statement is not ": nil})
		case "Na čekanju":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"id_of_statement is ": nil})
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
	response := dto.ToDepositPaymentOrderListResponseDTO(data)

	return response, total, nil
}
