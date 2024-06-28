package services

import (
	"context"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type CurrentBudgetServiceImpl struct {
	App             *celeritas.Celeritas
	repo            data.CurrentBudget
	spendingService SpendingDynamicService
}

func NewCurrentBudgetServiceImpl(
	app *celeritas.Celeritas,
	repo data.CurrentBudget,
	spendingService SpendingDynamicService,
) CurrentBudgetService {
	return &CurrentBudgetServiceImpl{
		App:             app,
		repo:            repo,
		spendingService: spendingService,
	}
}

func (h *CurrentBudgetServiceImpl) CreateCurrentBudget(ctx context.Context, input dto.CurrentBudgetDTO) (*dto.CurrentBudgetResponseDTO, error) {
	data := input.ToCurrentBudget()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo current budget insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo current budget get")
	}

	res := dto.ToCurrentBudgetResponseDTO(data)

	err = h.spendingService.CreateInititalSpendingDynamicFromCurrentBudget(ctx, data)
	if err != nil {
		return nil, newErrors.Wrap(err, "svc spending create initial spending dynamic from current budget")
	}

	return &res, nil
}

func (h *CurrentBudgetServiceImpl) UpdateActual(ctx context.Context, unitID, budgetID, accountID int, actual decimal.Decimal) (*dto.CurrentBudgetResponseDTO, error) {
	currentBudget, err := h.repo.GetBy(*up.And(
		up.Cond{"budget_id": budgetID},
		up.Cond{"unit_id": unitID},
		up.Cond{"account_id": accountID},
	))
	if err != nil {
		return nil, newErrors.Wrap(err, "repo current budget get by")
	}

	err = h.repo.UpdateActual(ctx, currentBudget.ID, actual)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo current budget update actual")
	}

	response := dto.ToCurrentBudgetResponseDTO(currentBudget)

	return &response, nil
}

func (h *CurrentBudgetServiceImpl) UpdateBalance(ctx context.Context, tx up.Session, id int, balance decimal.Decimal) error {
	err := h.repo.UpdateBalanceWithTx(ctx, tx, id, balance)
	if err != nil {
		return newErrors.Wrap(err, "repo current budget update balance with tx")
	}

	return nil
}

func (h *CurrentBudgetServiceImpl) GetCurrentBudget(id int) (*dto.CurrentBudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.WrapNotFoundError(err, "repo current budget get by")
	}
	response := dto.ToCurrentBudgetResponseDTO(data)

	return &response, nil
}

func (h *CurrentBudgetServiceImpl) GetCurrentBudgetList(filter dto.CurrentBudgetFilterDTO) ([]dto.CurrentBudgetResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.AccountID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"account_id": *filter.AccountID})
	}

	if filter.UnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"unit_id": *filter.UnitID})
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
		return nil, nil, newErrors.Wrap(err, "repo current budget get all")
	}
	response := dto.ToCurrentBudgetListResponseDTO(data)

	return response, total, nil
}

func (h *CurrentBudgetServiceImpl) GetCurrentBudgetUnitList() ([]int, error) {
	data, err := h.repo.GetCurrentBudgetUnits(int(time.Now().Year()))
	if err != nil {
		return nil, newErrors.Wrap(err, "repo get current budget units")
	}

	return data, nil
}

func (h *CurrentBudgetServiceImpl) GetAcctualCurrentBudget(organizationUnitID int) ([]dto.CurrentBudgetResponseDTO, error) {
	data, err := h.repo.GetActualCurrentBudget(organizationUnitID)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo current budget get actual current budget")
	}
	response := dto.ToCurrentBudgetListResponseDTO(data)

	return response, nil
}
