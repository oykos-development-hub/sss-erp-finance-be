package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type SpendingReleaseDTO struct {
	BudgetID  int             `json:"budget_id"`
	UnitID    int             `json:"unit_id"`
	AccountID int             `json:"account_id"`
	Month     int             `json:"month"`
	Value     decimal.Decimal `json:"value"`
}

type SpendingReleaseResponseDTO struct {
	ID              int             `json:"id"`
	CurrentBudgetID int             `json:"current_budget_id"`
	BudgetID        int             `json:"budget_id"`
	UnitID          int             `json:"unit_id"`
	AccountID       int             `json:"account_id"`
	Year            int             `json:"year"`
	Month           int             `json:"month"`
	Value           decimal.Decimal `json:"value"`
	CreatedAt       time.Time       `json:"created_at"`
}

func ToSpendingReleaseResponseDTO(data *data.SpendingReleaseWithCurrentBudget) SpendingReleaseResponseDTO {
	return SpendingReleaseResponseDTO{
		ID:              data.ID,
		BudgetID:        data.BudgetID,
		UnitID:          data.UnitID,
		AccountID:       data.AccountID,
		CurrentBudgetID: data.CurrentBudgetID,
		Year:            data.Year,
		Month:           data.Month,
		Value:           data.Value,
		CreatedAt:       data.CreatedAt,
	}
}

func ToSpendingReleaseListResponseDTO(spendingreleases []data.SpendingReleaseWithCurrentBudget) []SpendingReleaseResponseDTO {
	dtoList := make([]SpendingReleaseResponseDTO, len(spendingreleases))
	for i, x := range spendingreleases {
		dtoList[i] = ToSpendingReleaseResponseDTO(&x)
	}
	return dtoList
}
