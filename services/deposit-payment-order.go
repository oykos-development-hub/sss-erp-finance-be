package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type DepositPaymentOrderServiceImpl struct {
	App                    *celeritas.Celeritas
	repo                   data.DepositPaymentOrder
	additionalExpenses     DepositAdditionalExpenseService
	additionalExpensesRepo data.DepositAdditionalExpense
}

func NewDepositPaymentOrderServiceImpl(app *celeritas.Celeritas, repo data.DepositPaymentOrder, additionalExpensesRepo data.DepositAdditionalExpense, additionalExpenses DepositAdditionalExpenseService) DepositPaymentOrderService {
	return &DepositPaymentOrderServiceImpl{
		App:                    app,
		repo:                   repo,
		additionalExpensesRepo: additionalExpensesRepo,
		additionalExpenses:     additionalExpenses,
	}
}

func (h *DepositPaymentOrderServiceImpl) CreateDepositPaymentOrder(ctx context.Context, input dto.DepositPaymentOrderDTO) (*dto.DepositPaymentOrderResponseDTO, error) {
	dataToInsert := input.ToDepositPaymentOrder()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo deposit payment order insert")
		}

		for _, item := range input.AdditionalExpenses {
			itemToInsert := item.ToDepositAdditionalExpense()
			itemToInsert.PaymentOrderID = id
			itemToInsert.SourceBankAccount = item.SourceBankAccount
			itemToInsert.Status = "Kreiran"
			if itemToInsert.Price > 0 {
				_, err = h.additionalExpensesRepo.Insert(tx, *itemToInsert)
				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses insert")
				}
			}
		}

		for _, item := range input.AdditionalExpensesForPaying {
			itemToInsert, err := h.additionalExpensesRepo.Get(item.ID)
			if err != nil {
				return newErrors.Wrap(err, "repo deposit payment order additional expenses get")
			}
			itemToInsert.PayingPaymentOrderID = &id
			itemToInsert.Status = "Na čekanju"
			itemToInsert.SourceBankAccount = item.SourceBankAccount
			itemToInsert.ID = item.ID
			err = h.additionalExpensesRepo.Update(tx, *itemToInsert)
			if err != nil {
				return newErrors.Wrap(err, "repo deposit payment order additional expenses update")
			}
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit payment order get")
	}

	res := dto.ToDepositPaymentOrderResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *DepositPaymentOrderServiceImpl) UpdateDepositPaymentOrder(ctx context.Context, id int, input dto.DepositPaymentOrderDTO) (*dto.DepositPaymentOrderResponseDTO, error) {
	dataToInsert := input.ToDepositPaymentOrder()
	dataToInsert.ID = id

	oldData, err := h.GetDepositPaymentOrder(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit payment order get")
	}

	err = data.Upper.Tx(func(tx up.Session) error {
		var err error
		err = h.repo.Update(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo deposit payment order update")
		}

		//update vezanih troskova koji su nastali od tog naloga
		validExpenses := make(map[int]bool)

		for _, item := range oldData.AdditionalExpenses {
			validExpenses[item.ID] = false
		}

		for _, item := range input.AdditionalExpenses {
			_, exists := validExpenses[item.ID]
			if exists {
				validExpenses[item.ID] = true
			} else {
				additionalExpenseData := item.ToDepositAdditionalExpense()
				additionalExpenseData.PaymentOrderID = id
				additionalExpenseData.SourceBankAccount = dataToInsert.SourceBankAccount
				additionalExpenseData.Status = "Kreiran"
				_, err = h.additionalExpensesRepo.Insert(tx, *additionalExpenseData)

				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses insert")
				}
			}
		}

		for itemID, exists := range validExpenses {
			if !exists {
				err := h.additionalExpensesRepo.Delete(itemID)

				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses delete")
				}
			} else {
				for _, item := range input.AdditionalExpenses {
					if item.ID == itemID {
						additionalExpenseData, err := h.additionalExpensesRepo.Get(item.ID)
						if err != nil {
							return newErrors.Wrap(err, "repo deposit payment order additional expenses get")
						}
						additionalExpenseData.ID = itemID
						additionalExpenseData.PaymentOrderID = id
						additionalExpenseData.SourceBankAccount = dataToInsert.SourceBankAccount
						additionalExpenseData.Status = "Kreiran"
						err = h.additionalExpensesRepo.Update(tx, *additionalExpenseData)
						if err != nil {
							return newErrors.Wrap(err, "repo deposit payment order additional expenses update")
						}
					}
				}
			}
		}

		//update vezanih troskova koji se placaju tim nalogom
		validExpensesPaying := make(map[int]bool)

		for _, item := range oldData.AdditionalExpensesForPaying {
			validExpensesPaying[item.ID] = false
		}

		for _, item := range input.AdditionalExpensesForPaying {
			_, exists := validExpensesPaying[item.ID]
			if exists {
				validExpensesPaying[item.ID] = true
			} else {
				additionalExpenseData, err := h.additionalExpensesRepo.Get(item.ID)
				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses get")
				}
				additionalExpenseData.PayingPaymentOrderID = &id
				additionalExpenseData.SourceBankAccount = dataToInsert.SourceBankAccount
				additionalExpenseData.Status = "Na čekanju"
				err = h.additionalExpensesRepo.Update(tx, *additionalExpenseData)

				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses update")
				}
			}
		}

		for itemID, exists := range validExpensesPaying {
			if !exists {
				additionalExpenseData, err := h.additionalExpensesRepo.Get(itemID)

				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses get")
				}

				additionalExpenseData.PayingPaymentOrderID = nil
				additionalExpenseData.SourceBankAccount = dataToInsert.SourceBankAccount
				additionalExpenseData.Status = "Kreiran"

				err = h.additionalExpensesRepo.Update(tx, *additionalExpenseData)

				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses update")
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit payment order delete")
	}

	response := dto.ToDepositPaymentOrderResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *DepositPaymentOrderServiceImpl) PayDepositPaymentOrder(ctx context.Context, id int, input dto.DepositPaymentOrderDTO) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		err = h.repo.PayDepositPaymentOrder(ctx, tx, id, *input.IDOfStatement, *input.DateOfStatement)
		if err != nil {
			return newErrors.Wrap(err, "repo deposit payment order pay")
		}

		additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{
			PaymentOrderID: &id,
		})

		if err != nil {
			return newErrors.Wrap(err, "repo deposit payment order additional expenses get all")
		}

		for _, item := range additionalExpenses {
			if item.Title == "Neto" {
				itemToUpdate, err := h.additionalExpensesRepo.Get(item.ID)
				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses get")
				}
				itemToUpdate.Status = "Plaćen"
				err = h.additionalExpensesRepo.Update(tx, *itemToUpdate)
				if err != nil {
					return newErrors.Wrap(err, "repo deposit payment order additional expenses update")
				}
			}
		}

		additionalExpenses, _, err = h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{
			PayingPaymentOrderID: &id,
		})

		if err != nil {
			return newErrors.Wrap(err, "repo deposit payment order additional expenses get all")
		}

		for _, item := range additionalExpenses {
			itemToUpdate, err := h.additionalExpensesRepo.Get(item.ID)
			if err != nil {
				return newErrors.Wrap(err, "repo deposit payment order additional expenses get")
			}
			itemToUpdate.Status = "Plaćen"
			err = h.additionalExpensesRepo.Update(tx, *itemToUpdate)
			if err != nil {
				return newErrors.Wrap(err, "repo deposit payment order additional expenses update")
			}

		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}

func (h *DepositPaymentOrderServiceImpl) DeleteDepositPaymentOrder(ctx context.Context, id int) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{
			PayingPaymentOrderID: &id,
		})

		if err != nil {
			return newErrors.Wrap(err, "repo deposit payment order additional expenses get all")
		}

		for _, item := range additionalExpenses {
			itemToUpdate, err := h.additionalExpensesRepo.Get(item.ID)
			if err != nil {
				return newErrors.Wrap(err, "repo deposit payment order additional expenses get")
			}
			itemToUpdate.Status = "Kreiran"
			err = h.additionalExpensesRepo.Update(tx, *itemToUpdate)
			if err != nil {
				return newErrors.Wrap(err, "repo deposit payment order additional expenses update")
			}

		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo deposit payment order delete")
	}

	return nil
}

func (h *DepositPaymentOrderServiceImpl) GetDepositPaymentOrder(id int) (*dto.DepositPaymentOrderResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit payment order get")
	}

	response := dto.ToDepositPaymentOrderResponseDTO(*data)

	additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PaymentOrderID: &id})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit payment order additional expenses get all")
	}

	response.AdditionalExpenses = additionalExpenses

	additionalExpenses, _, err = h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PayingPaymentOrderID: &id})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo deposit payment order additional expenses get all")
	}

	for i := 0; i < len(additionalExpenses); i++ {
		order, err := h.repo.Get(additionalExpenses[i].PaymentOrderID)

		if err != nil {
			return nil, newErrors.Wrap(err, "repo deposit payment order get")
		}

		additionalExpenses[i].CaseNumber = order.CaseNumber

	}

	response.AdditionalExpensesForPaying = additionalExpenses

	return &response, nil
}

func (h *DepositPaymentOrderServiceImpl) GetDepositPaymentOrderList(filter dto.DepositPaymentOrderFilterDTO) ([]dto.DepositPaymentOrderResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.CaseNumber != nil && *filter.CaseNumber != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"case_number": *filter.CaseNumber})
	}

	if filter.SourceBankAccount != nil && *filter.SourceBankAccount != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"source_bank_account": *filter.SourceBankAccount})
	}

	if filter.SupplierID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"supplier_id": *filter.SupplierID})
	}

	if filter.SourceBankAccount != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"source_bank_account": *filter.SourceBankAccount})
	}

	if filter.Status != nil {
		switch *filter.Status {
		case "Plaćen":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"id_of_statement is not ": nil})
		case "Na čekanju":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"id_of_statement is ": nil})
		}
	}

	if filter.Search != nil && *filter.Search != "" {
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
		return nil, nil, newErrors.Wrap(err, "repo deposit payment order get all")
	}
	response := dto.ToDepositPaymentOrderListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PaymentOrderID: &response[i].ID})

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo deposit payment order additional expenses get all")
		}

		response[i].AdditionalExpenses = additionalExpenses

		additionalExpenses, _, err = h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PayingPaymentOrderID: &response[i].ID})

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo deposit payment order additional expenses get all")
		}

		response[i].AdditionalExpensesForPaying = additionalExpenses
	}

	return response, total, nil
}
