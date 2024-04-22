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

func (h *DepositPaymentOrderServiceImpl) CreateDepositPaymentOrder(input dto.DepositPaymentOrderDTO) (*dto.DepositPaymentOrderResponseDTO, error) {
	dataToInsert := input.ToDepositPaymentOrder()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return err
		}

		for _, item := range input.AdditionalExpenses {
			itemToInsert := item.ToDepositAdditionalExpense()
			itemToInsert.PaymentOrderID = id
			itemToInsert.SourceBankAccount = item.SourceBankAccount
			itemToInsert.Status = "Kreiran"
			if itemToInsert.Price > 0 {
				_, err = h.additionalExpensesRepo.Insert(tx, *itemToInsert)
				if err != nil {
					return err
				}
			}
		}

		for _, item := range input.AdditionalExpensesForPaying {
			itemToInsert, err := h.additionalExpensesRepo.Get(item.ID)
			if err != nil {
				return err
			}
			itemToInsert.PayingPaymentOrderID = &id
			itemToInsert.Status = "Na čekanju"
			itemToInsert.SourceBankAccount = item.SourceBankAccount
			itemToInsert.ID = item.ID
			err = h.additionalExpensesRepo.Update(tx, *itemToInsert)
			if err != nil {
				return err
			}
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

	oldData, err := h.GetDepositPaymentOrder(id)

	if err != nil {
		return nil, err
	}

	err = data.Upper.Tx(func(tx up.Session) error {
		var err error
		err = h.repo.Update(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
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
					return err
				}
			}
		}

		for itemID, exists := range validExpenses {
			if !exists {
				err := h.additionalExpensesRepo.Delete(itemID)

				if err != nil {
					return err
				}
			} else {
				for _, item := range input.AdditionalExpenses {
					if item.ID == itemID {
						additionalExpenseData, err := h.additionalExpensesRepo.Get(item.ID)
						if err != nil {
							return err
						}
						additionalExpenseData.ID = itemID
						additionalExpenseData.PaymentOrderID = id
						additionalExpenseData.SourceBankAccount = dataToInsert.SourceBankAccount
						additionalExpenseData.Status = "Kreiran"
						err = h.additionalExpensesRepo.Update(tx, *additionalExpenseData)
						if err != nil {
							return err
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
					return err
				}
				additionalExpenseData.PayingPaymentOrderID = &id
				additionalExpenseData.SourceBankAccount = dataToInsert.SourceBankAccount
				additionalExpenseData.Status = "Na čekanju"
				err = h.additionalExpensesRepo.Update(tx, *additionalExpenseData)

				if err != nil {
					return err
				}
			}
		}

		for itemID, exists := range validExpensesPaying {
			if !exists {
				additionalExpenseData, err := h.additionalExpensesRepo.Get(itemID)

				if err != nil {
					return err
				}

				additionalExpenseData.PayingPaymentOrderID = nil
				additionalExpenseData.SourceBankAccount = dataToInsert.SourceBankAccount
				additionalExpenseData.Status = "Kreiran"

				err = h.additionalExpensesRepo.Update(tx, *additionalExpenseData)

				if err != nil {
					return err
				}
			}
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

func (h *DepositPaymentOrderServiceImpl) PayDepositPaymentOrder(id int, input dto.DepositPaymentOrderDTO) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		err = h.repo.PayDepositPaymentOrder(tx, id, *input.IDOfStatement, *input.DateOfStatement)
		if err != nil {
			return errors.ErrInternalServer
		}

		additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{
			PaymentOrderID: &id,
		})

		if err != nil {
			return err
		}

		for _, item := range additionalExpenses {
			if item.Title == "Neto" {
				itemToUpdate, err := h.additionalExpensesRepo.Get(item.ID)
				if err != nil {
					return err
				}
				itemToUpdate.Status = "Plaćen"
				err = h.additionalExpensesRepo.Update(tx, *itemToUpdate)
				if err != nil {
					return err
				}
			}
		}

		additionalExpenses, _, err = h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{
			PayingPaymentOrderID: &id,
		})

		if err != nil {
			return err
		}

		for _, item := range additionalExpenses {
			itemToUpdate, err := h.additionalExpensesRepo.Get(item.ID)
			if err != nil {
				return err
			}
			itemToUpdate.Status = "Plaćen"
			err = h.additionalExpensesRepo.Update(tx, *itemToUpdate)
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (h *DepositPaymentOrderServiceImpl) DeleteDepositPaymentOrder(id int) error {
	err := data.Upper.Tx(func(tx up.Session) error {
		additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{
			PayingPaymentOrderID: &id,
		})

		if err != nil {
			return err
		}

		for _, item := range additionalExpenses {
			itemToUpdate, err := h.additionalExpensesRepo.Get(item.ID)
			if err != nil {
				return err
			}
			itemToUpdate.Status = "Kreiran"
			err = h.additionalExpensesRepo.Update(tx, *itemToUpdate)
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		return err
	}

	err = h.repo.Delete(id)
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

	additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PaymentOrderID: &id})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	response.AdditionalExpenses = additionalExpenses

	additionalExpenses, _, err = h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PayingPaymentOrderID: &id})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	for i := 0; i < len(additionalExpenses); i++ {
		order, err := h.repo.Get(additionalExpenses[i].PaymentOrderID)

		if err != nil {
			h.App.ErrorLog.Println(err)
			return nil, err
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
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToDepositPaymentOrderListResponseDTO(data)

	for i := 0; i < len(response); i++ {
		additionalExpenses, _, err := h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PaymentOrderID: &response[i].ID})

		if err != nil {
			h.App.ErrorLog.Println(err)
			return nil, nil, errors.ErrInternalServer
		}

		response[i].AdditionalExpenses = additionalExpenses

		additionalExpenses, _, err = h.additionalExpenses.GetDepositAdditionalExpenseList(dto.DepositAdditionalExpenseFilterDTO{PayingPaymentOrderID: &response[i].ID})

		if err != nil {
			h.App.ErrorLog.Println(err)
			return nil, nil, errors.ErrInternalServer
		}

		response[i].AdditionalExpensesForPaying = additionalExpenses
	}

	return response, total, nil
}
