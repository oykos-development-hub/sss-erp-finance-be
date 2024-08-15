package services

import (
	"context"
	"fmt"
	"time"

	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type SpendingReleaseServiceImpl struct {
	App                  *celeritas.Celeritas
	repo                 data.SpendingRelease
	repoBudget           data.Budget
	repoCurrentBudget    data.CurrentBudget
	repoSpendingRequests data.SpendingReleaseRequest
}

func NewSpendingReleaseServiceImpl(app *celeritas.Celeritas, repo data.SpendingRelease, repoCurrentBudget data.CurrentBudget, repoBudget data.Budget, repoSpendingRequests data.SpendingReleaseRequest) SpendingReleaseService {
	return &SpendingReleaseServiceImpl{
		App:                  app,
		repo:                 repo,
		repoBudget:           repoBudget,
		repoCurrentBudget:    repoCurrentBudget,
		repoSpendingRequests: repoSpendingRequests,
	}
}

func (h *SpendingReleaseServiceImpl) CreateSpendingRelease(ctx context.Context, budgetID, unitID int, inputDTOList []dto.SpendingReleaseDTO) ([]dto.SpendingReleaseResponseDTO, error) {
	res := make([]dto.SpendingReleaseResponseDTO, 0, len(inputDTOList))
	currentMonth := int(time.Now().Month())

	currentYear := time.Now().Year()

	existingSpendingRelease, err := h.GetSpendingReleaseList(data.SpendingReleaseFilterDTO{
		BudgetID: &budgetID,
		UnitID:   &unitID,
		Month:    &currentMonth,
		Year:     &currentYear,
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "repo get spending release list")
	}

	if len(existingSpendingRelease) != 0 {
		return nil, newErrors.NewBadRequestError("release already exists")
	}

	for _, inputDTO := range inputDTOList {
		currentBudget, err := h.repoCurrentBudget.GetBy(*up.And(
			up.Cond{"budget_id": budgetID},
			up.Cond{"unit_id": unitID},
			up.Cond{"account_id": inputDTO.AccountID},
			up.Cond{"type": 1},
		))
		if err != nil {
			return nil, newErrors.Wrap(err, "repo current budget get by")
		}

		_, err = h.repo.GetBy(*up.And(up.Cond{"current_budget_id": currentBudget.ID}, up.Cond{"month": currentMonth}))
		if !newErrors.IsErr(err, newErrors.NotFoundCode) {
			return nil, newErrors.NewWithCode(newErrors.SingleMonthSpendingReleaseCode, "only single release is allowed per month")
		}

		budget, err := h.repoBudget.Get(currentBudget.BudgetID)
		if err != nil {
			return nil, newErrors.Wrap(err, "repo budget get")
		}

		inputData := data.SpendingRelease{
			CurrentBudgetID: currentBudget.ID,
			Year:            budget.Year,
			Month:           currentMonth,
			Value:           inputDTO.Value,
			Username:        inputDTO.Username,
		}
		if !inputData.ValidateNewRelease() {
			return nil, newErrors.NewWithCode(newErrors.ReleaseInCurrentMonthCode, "release is possible only in the current month")
		}

		if currentBudget.Vault().Sub(inputData.Value).LessThan(decimal.Zero) {
			return nil, newErrors.NewWithCode(newErrors.NotEnoughFundsCode, "not enough funds")
		}
		conditionAndExp := &up.AndExpr{}
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": unitID})

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"month": currentMonth})

		conditionAndExp = up.And(conditionAndExp, &up.Cond{"year": currentYear})

		spendingReleaseRequest, _, err := h.repoSpendingRequests.GetAll(nil, nil, conditionAndExp, nil)
		if err != nil {
			return nil, newErrors.Wrap(err, "repo spending requests get all")
		}

		if len(spendingReleaseRequest) != 1 {
			return nil, newErrors.Wrap(errors.ErrInvalidInput, "repo spending requests get all")
		}

		spendingReleaseRequest[0].Status = data.SpendingReleaseStatusFilled

		err = h.repoSpendingRequests.Update(data.Upper, *spendingReleaseRequest[0])

		if err != nil {
			return nil, newErrors.Wrap(err, "repo spending requests update")
		}

		id, err := h.repo.Insert(ctx, inputData)
		if err != nil {
			return nil, newErrors.Wrap(err, "repo release insert")
		}

		item, err := h.repo.Get(id)
		if err != nil {
			return nil, newErrors.Wrap(err, "repo release get")
		}

		err = h.repoCurrentBudget.UpdateBalance(currentBudget.ID, currentBudget.Balance.Add(item.Value))
		if err != nil {
			return nil, newErrors.Wrap(err, "repo current budget update balance")
		}

		err = h.repoCurrentBudget.UpdateActual(ctx, currentBudget.ID, currentBudget.Actual.Sub(item.Value))
		if err != nil {
			return nil, newErrors.Wrap(err, "repo current budget update balance")
		}

		resItem := dto.ToSpendingReleaseResponseDTO(item)

		res = append(res, resItem)
	}

	return res, nil
}

func (h *SpendingReleaseServiceImpl) DeleteSpendingRelease(ctx context.Context, input *dto.DeleteSpendingReleaseInput) error {
	releases, err := h.repo.GetAll(data.SpendingReleaseFilterDTO{
		BudgetID: &input.BudgetID,
		UnitID:   &input.UnitID,
		Month:    &input.Month,
	})
	if err != nil {
		return newErrors.Wrap(err, "repo get all")
	}

	err = data.Upper.Tx(func(tx up.Session) error {
		for _, release := range releases {
			err = h.repo.Delete(ctx, tx, release.ID)
			if err != nil {
				return newErrors.Wrap(err, "repo delete")
			}

			currentBudget, err := h.repoCurrentBudget.Get(release.CurrentBudgetID)
			if err != nil {
				return newErrors.Wrap(err, "repo current budget get")
			}

			if currentBudget.Balance.Sub(release.Value).Cmp(decimal.NewFromInt(0)) < 0 {
				return newErrors.Wrap(errors.ErrInsufficientFunds, "balance")
			}

			err = h.repoCurrentBudget.UpdateBalanceWithTx(ctx, tx, currentBudget.ID, currentBudget.Balance.Sub(release.Value))
			if err != nil {
				return newErrors.Wrap(err, "repo current budget update balance")
			}

			err = h.repoCurrentBudget.UpdateActualWithTx(ctx, tx, currentBudget.ID, currentBudget.Actual.Add(release.Value))
			if err != nil {
				return newErrors.Wrap(err, "repo current budget update balance")
			}
		}
		return nil
	})
	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}

	return nil
}

func (h *SpendingReleaseServiceImpl) GetSpendingRelease(id int) (*dto.SpendingReleaseResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo release get")
	}
	response := dto.ToSpendingReleaseResponseDTO(data)

	return &response, nil
}

func (h *SpendingReleaseServiceImpl) GetSpendingReleaseList(filter data.SpendingReleaseFilterDTO) ([]dto.SpendingReleaseResponseDTO, error) {
	data, err := h.repo.GetAll(filter)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo release get all")
	}
	response := dto.ToSpendingReleaseListResponseDTO(data)

	return response, nil
}

// GetSpendingReleaseOverview implements SpendingReleaseService.
func (h *SpendingReleaseServiceImpl) GetSpendingReleaseOverview(filter dto.SpendingReleaseOverviewFilterDTO) ([]dto.SpendingReleaseOverview, error) {
	data, err := h.repo.GetAllSum(filter.Month, filter.Year, filter.BudgetID, filter.UnitID)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo release get all sum")
	}

	response := dto.ToSpendingReleaseOverviewDTO(data)

	return response, nil
}

func (h *SpendingReleaseServiceImpl) StartMonthlyTaskForSpendingReleases() {
	go h.monthlyScheduler()
}

func (h *SpendingReleaseServiceImpl) monthlyScheduler() {
	for {
		now := time.Now()
		year, month, _ := now.Date()
		location := now.Location()
		//var nextMonth time.Time

		/*if month == 12 {
			nextMonth = time.Date(year, month+2, 1, 0, 0, 0, 0, location)
		} else {
			nextMonth = time.Date(year, month+1, 1, 0, 0, 0, 0, location)
		}*/

		nextMonth := time.Date(year, month, 15, 15, 0, 0, 0, location)

		waitDuration := time.Until(nextMonth)
		fmt.Printf("Sleeping until %v for release trigger\n", nextMonth)

		time.Sleep(waitDuration)

		h.executeMonthlyTask()
	}
}

func (h *SpendingReleaseServiceImpl) executeMonthlyTask() {
	fmt.Printf("Task for releases started")

	conditionAndExp := &up.AndExpr{}

	now := time.Now()
	year, monthTime, _ := now.Date()
	formattedString := fmt.Sprintf("%d-01-01T00:00:00Z", year)

	conditionAndExp = up.And(conditionAndExp, &up.Cond{"created_at >= ": formattedString})

	currentBudget, _, _ := h.repoCurrentBudget.GetAll(nil, nil, conditionAndExp, nil)

	month := int(monthTime) - 1

	for _, item := range currentBudget {
		spendingReleaseItem, _ := h.repo.GetAll(data.SpendingReleaseFilterDTO{
			CurrentBudgetID: &item.ID,
			Year:            &year,
			Month:           &month,
		})

		ctx := context.Background()
		ctx = contextutil.SetUserIDInContext(ctx, 0)

		if spendingReleaseItem == nil && len(spendingReleaseItem) == 0 {
			_, _ = h.repo.Insert(ctx, data.SpendingRelease{
				CurrentBudgetID: item.ID,
				Year:            year,
				Month:           month,
				Value:           decimal.NewFromInt(0),
			})
		}
	}

	fmt.Printf("Task for releases ended")

}
