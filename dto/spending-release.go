package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type SpendingReleaseOverviewFilterDTO struct {
	Month    int `json:"month"`
	Year     int `json:"year" validate:"required"`
	BudgetID int `json:"budget_id" validate:"required"`
	UnitID   int `json:"unit_id" validate:"required"`
}

type SpendingReleaseOverview struct {
	Month int             `json:"month"`
	Year  int             `json:"year"`
	Value decimal.Decimal `json:"value"`
}

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

func ToSpendingReleaseOverviewItemDTO(data *data.SpendingReleaseOverview) SpendingReleaseOverview {
	return SpendingReleaseOverview{
		Year:  data.Year,
		Month: data.Month,
		Value: data.Value,
	}
}

func ToSpendingReleaseOverviewDTO(overviewData []data.SpendingReleaseOverview) []SpendingReleaseOverview {
	dtoList := make([]SpendingReleaseOverview, len(overviewData))
	for i, x := range overviewData {
		dtoList[i] = ToSpendingReleaseOverviewItemDTO(&x)
	}
	return dtoList
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
