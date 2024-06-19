package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type FilledFinancialBudgetServiceImpl struct {
	App               *celeritas.Celeritas
	repo              data.FilledFinancialBudget
	reqRepo           data.BudgetRequest
	currentBudgetRepo data.CurrentBudget
}

func NewFilledFinancialBudgetServiceImpl(
	app *celeritas.Celeritas,
	repo data.FilledFinancialBudget,
	reqRepo data.BudgetRequest,
	currentBudgetRepo data.CurrentBudget,
) FilledFinancialBudgetService {
	return &FilledFinancialBudgetServiceImpl{
		App:               app,
		repo:              repo,
		reqRepo:           reqRepo,
		currentBudgetRepo: currentBudgetRepo,
	}
}

func (h *FilledFinancialBudgetServiceImpl) CreateFilledFinancialBudget(ctx context.Context, input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error) {
	data := input.ToFilledFinancialBudget()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget get")
	}

	res := dto.ToFilledFinancialBudgetResponseDTO(data)

	return &res, nil
}

func (h *FilledFinancialBudgetServiceImpl) UpdateFilledFinancialBudget(ctx context.Context, id int, input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error) {
	inputData := input.ToFilledFinancialBudget()
	inputData.ID = id

	err := h.repo.Update(ctx, *inputData)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget update")
	}

	resData, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget update")
	}

	response := dto.ToFilledFinancialBudgetResponseDTO(resData)

	return &response, nil
}

func (h *FilledFinancialBudgetServiceImpl) UpdateActualFilledFinancialBudget(ctx context.Context, id int, actual decimal.Decimal) (*dto.FilledFinancialBudgetResponseDTO, error) {
	err := h.repo.UpdateActual(ctx, id, actual)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget update actual")
	}

	item, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget get")
	}

	budgetRequest, err := h.reqRepo.Get(item.BudgetRequestID)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo budget request get")
	}

	// TODO: check if there is only one insert. If we allow official to update actual, then we need to update it here too.
	_, err = h.currentBudgetRepo.Insert(ctx, data.CurrentBudget{
		BudgetID:      budgetRequest.BudgetID,
		UnitID:        budgetRequest.OrganizationUnitID,
		AccountID:     item.AccountID,
		InitialActual: item.Actual.Decimal,
		Actual:        item.Actual.Decimal,
		Balance:       decimal.Zero,
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "repo current budget insert")
	}

	response := dto.ToFilledFinancialBudgetResponseDTO(item)

	return &response, nil
}

func (h *FilledFinancialBudgetServiceImpl) DeleteFilledFinancialBudget(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo filled financial budget delete")
	}

	return nil
}

func (h *FilledFinancialBudgetServiceImpl) GetFilledFinancialBudget(id int) (*dto.FilledFinancialBudgetResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget get")
	}

	response := dto.ToFilledFinancialBudgetResponseDTO(data)

	return &response, nil
}

func (h *FilledFinancialBudgetServiceImpl) GetFilledFinancialBudgetList(filter dto.FilledFinancialBudgetFilterDTO) ([]dto.FilledFinancialBudgetResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_request_id": filter.BudgetRequestID})

	if len(filter.AccountIdList) > 0 {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"account_id": up.In(filter.AccountIdList...)})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo filled financial budget get all")
	}

	response := dto.ToFilledFinancialBudgetListResponseDTO(data)

	return response, total, nil
}

func (h *FilledFinancialBudgetServiceImpl) GetSummaryFilledFinancialRequests(budgetID int, requestType data.RequestType) ([]dto.FilledFinancialBudgetResponseDTO, error) {
	data, err := h.repo.GetSummaryFilledFinancialRequests(budgetID, requestType)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo filled financial budget get summary filled financial requests")
	}
	response := dto.ToFilledFinancialBudgetListResponseDTO(data)

	return response, nil
}
