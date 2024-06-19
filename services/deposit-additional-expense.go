package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type DepositAdditionalExpenseServiceImpl struct {
	App    *celeritas.Celeritas
	repo   data.DepositAdditionalExpense
	orders data.DepositPaymentOrder
}

func NewDepositAdditionalExpenseServiceImpl(app *celeritas.Celeritas, repo data.DepositAdditionalExpense, orders data.DepositPaymentOrder) DepositAdditionalExpenseService {
	return &DepositAdditionalExpenseServiceImpl{
		App:    app,
		repo:   repo,
		orders: orders,
	}
}

func (h *DepositAdditionalExpenseServiceImpl) CreateDepositAdditionalExpense(input dto.DepositAdditionalExpenseDTO) (*dto.DepositAdditionalExpenseResponseDTO, error) {
	dataToInsert := input.ToDepositAdditionalExpense()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo deposit additional expenses insert")
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit additional expenses delete")
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
			return newErrors.Wrap(err, "repo deposit additional expenses update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit additional expenses delete")
	}

	response := dto.ToDepositAdditionalExpenseResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *DepositAdditionalExpenseServiceImpl) DeleteDepositAdditionalExpense(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		return newErrors.Wrap(err, "repo deposit additional expenses delete")
	}

	return nil
}

func (h *DepositAdditionalExpenseServiceImpl) GetDepositAdditionalExpense(id int) (*dto.DepositAdditionalExpenseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit additional expenses get")
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

	if filter.PayingPaymentOrderID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"paying_payment_order_id": *filter.PayingPaymentOrderID})
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
		return nil, nil, newErrors.Wrap(err, "repo deposit additional expenses get all")
	}
	response := dto.ToDepositAdditionalExpenseListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		item, err := h.orders.Get(response[i].PaymentOrderID)
		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo payment orders get")
		}

		response[i].CaseNumber = item.CaseNumber
	}

	return response, total, nil
}
