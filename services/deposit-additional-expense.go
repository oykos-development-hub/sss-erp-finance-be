package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type DepositAdditionalExpenseServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.DepositAdditionalExpense
}

func NewDepositAdditionalExpenseServiceImpl(app *celeritas.Celeritas, repo data.DepositAdditionalExpense) DepositAdditionalExpenseService {
	return &DepositAdditionalExpenseServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *DepositAdditionalExpenseServiceImpl) CreateDepositAdditionalExpense(input dto.DepositAdditionalExpenseDTO) (*dto.DepositAdditionalExpenseResponseDTO, error) {
	dataToInsert := input.ToDepositAdditionalExpense()

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

	res := dto.ToDepositAdditionalExpenseResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *DepositAdditionalExpenseServiceImpl) UpdateDepositAdditionalExpense(id int, input dto.DepositAdditionalExpenseDTO) (*dto.DepositAdditionalExpenseResponseDTO, error) {
	dataToInsert := input.ToDepositAdditionalExpense()
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

	response := dto.ToDepositAdditionalExpenseResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *DepositAdditionalExpenseServiceImpl) DeleteDepositAdditionalExpense(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *DepositAdditionalExpenseServiceImpl) GetDepositAdditionalExpense(id int) (*dto.DepositAdditionalExpenseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToDepositAdditionalExpenseResponseDTO(*data)

	return &response, nil
}

func (h *DepositAdditionalExpenseServiceImpl) GetDepositAdditionalExpenseList(filter dto.DepositAdditionalExpenseFilterDTO) ([]dto.DepositAdditionalExpenseResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.PaymentOrderID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"payment_order_id": *filter.PaymentOrderID})
	} else {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"title <> ": "Neto"})
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
	response := dto.ToDepositAdditionalExpenseListResponseDTO(data)

	return response, total, nil
}
