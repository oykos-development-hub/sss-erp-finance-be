package services

import (
	"context"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type SalaryServiceImpl struct {
	App                         *celeritas.Celeritas
	repo                        data.Salary
	salaryAdditionalExpenseRepo data.SalaryAdditionalExpense
	salaryAdditionalService     SalaryAdditionalExpenseService
}

func NewSalaryServiceImpl(app *celeritas.Celeritas, repo data.Salary, salaryAdditionalExpenseRepo data.SalaryAdditionalExpense, salaryAdditionalService SalaryAdditionalExpenseService) SalaryService {
	return &SalaryServiceImpl{
		App:                         app,
		repo:                        repo,
		salaryAdditionalExpenseRepo: salaryAdditionalExpenseRepo,
		salaryAdditionalService:     salaryAdditionalService,
	}
}

func (h *SalaryServiceImpl) CreateSalary(ctx context.Context, input dto.SalaryDTO) (*dto.SalaryResponseDTO, error) {
	dataToInsert := input.ToSalary()

	var id int
	var err error
	id, err = h.repo.Insert(ctx, data.Upper, *dataToInsert)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo salary insert")
	}

	for _, additionalExpense := range input.SalaryAdditionalExpenses {
		additionalExpenseData := additionalExpense.ToSalaryAdditionalExpense()
		additionalExpenseData.SalaryID = id
		additionalExpenseData.Status = "Kreiran"
		if additionalExpenseData.Amount > 0 {
			_, err = h.salaryAdditionalExpenseRepo.Insert(data.Upper, *additionalExpenseData)
			if err != nil {
				return nil, newErrors.Wrap(err, "repo salary additional expenses insert")
			}
		}
	}

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo salary get")
	}

	res := dto.ToSalaryResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *SalaryServiceImpl) UpdateSalary(ctx context.Context, id int, input dto.SalaryDTO) (*dto.SalaryResponseDTO, error) {
	dataToInsert := input.ToSalary()
	dataToInsert.ID = id

	oldData, err := h.GetSalary(id)

	if err != nil {
		return nil, newErrors.Wrap(err, "repo salary get")
	}

	err = data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo salary update")
		}

		validExpenses := make(map[int]bool)

		for _, item := range oldData.SalaryAdditionalExpenses {
			validExpenses[item.ID] = false
		}

		for _, item := range input.SalaryAdditionalExpenses {
			_, exists := validExpenses[item.ID]
			if exists {
				validExpenses[item.ID] = true
			} else {
				additionalExpenseData := item.ToSalaryAdditionalExpense()
				additionalExpenseData.SalaryID = id
				additionalExpenseData.Status = "Kreiran"
				if additionalExpenseData.Amount > 0 {
					_, err = h.salaryAdditionalExpenseRepo.Insert(tx, *additionalExpenseData)

					if err != nil {
						return newErrors.Wrap(err, "repo salary additional expenses insert")
					}
				}
			}
		}

		for itemID, exists := range validExpenses {
			if !exists {
				err := h.salaryAdditionalExpenseRepo.Delete(itemID)

				if err != nil {
					return newErrors.Wrap(err, "repo salary additional expenses delete")
				}
			} else {
				for _, item := range input.SalaryAdditionalExpenses {
					if item.ID == itemID {
						additionalExpenseData := item.ToSalaryAdditionalExpense()
						additionalExpenseData.ID = id
						additionalExpenseData.SalaryID = id
						if additionalExpenseData.Amount > 0 {
							err := h.salaryAdditionalExpenseRepo.Update(tx, *additionalExpenseData)
							if err != nil {
								return newErrors.Wrap(err, "repo salary additional expenses update")
							}
						} else {
							err := h.salaryAdditionalExpenseRepo.Delete(additionalExpenseData.ID)
							if err != nil {
								return newErrors.Wrap(err, "repo salary additional expenses delete")
							}
						}
					}
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
		return nil, newErrors.Wrap(err, "repo salary get")
	}

	response := dto.ToSalaryResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *SalaryServiceImpl) DeleteSalary(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo salary delete")
	}

	return nil
}

func (h *SalaryServiceImpl) GetSalary(id int) (*dto.SalaryResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo salary get")
	}

	response := dto.ToSalaryResponseDTO(*data)

	additionalExpenses, _, err := h.salaryAdditionalService.GetSalaryAdditionalExpenseList(dto.SalaryAdditionalExpenseFilterDTO{
		SalaryID: &id,
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "repo salary additional expenses get")
	}

	response.SalaryAdditionalExpenses = additionalExpenses

	for _, additionalExpense := range additionalExpenses {
		if additionalExpense.Type == "banks" {
			response.NetPrice += additionalExpense.Amount
		} else if additionalExpense.Type == "suspensions" {
			response.ObligationsPrice += additionalExpense.Amount
		} else {
			response.VatPrice += additionalExpense.Amount
		}
	}

	response.GrossPrice = response.VatPrice + response.NetPrice

	return &response, nil
}

func (h *SalaryServiceImpl) GetSalaryList(filter dto.SalaryFilterDTO) ([]dto.SalaryResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.ActivityID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"activity_id": *filter.ActivityID})
	}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Status != nil && *filter.Status != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	if filter.Registred != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"registred": *filter.Registred})
	}

	if filter.Year != nil {
		year := *filter.Year
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_calculation": up.Between(startOfYear, endOfYear)})
	}

	if filter.Month != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"month": *filter.Month})
	}

	/*if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}*/

	orders = append(orders, "-created_at")

	salaryData, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo salary get all")
	}
	response := dto.ToSalaryListResponseDTO(salaryData)

	for i := 0; i < len(response); i++ {
		additionalExpenses, _, err := h.salaryAdditionalService.GetSalaryAdditionalExpenseList(dto.SalaryAdditionalExpenseFilterDTO{
			SalaryID: &response[i].ID,
		})

		if err != nil {
			return nil, nil, newErrors.Wrap(err, "repo salary additional expenses get all")
		}

		response[i].SalaryAdditionalExpenses = additionalExpenses
		response[i].Deletable = true
		for _, additionalExpense := range additionalExpenses {

			if additionalExpense.Status != data.InvoiceStatusCreated {
				response[i].Deletable = false
			}

			if additionalExpense.Type == "banks" {
				response[i].NetPrice += additionalExpense.Amount
			} else if additionalExpense.Type == "suspensions" {
				response[i].ObligationsPrice += additionalExpense.Amount
			} else {
				response[i].VatPrice += additionalExpense.Amount
			}
		}

		if response[i].Registred != nil && *response[i].Registred {
			response[i].Deletable = false
		}

		response[i].GrossPrice = response[i].VatPrice + response[i].NetPrice
	}

	return response, total, nil
}
