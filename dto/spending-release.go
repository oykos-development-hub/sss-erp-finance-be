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
	Month           int             `json:"month"`
	Value           decimal.Decimal `json:"value"`
	CreatedAt       time.Time       `json:"created_at"`
}

type SpendingReleaseFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
}

func ToSpendingReleaseResponseDTO(data data.SpendingRelease) SpendingReleaseResponseDTO {
	return SpendingReleaseResponseDTO{
		ID:              data.ID,
		CurrentBudgetID: data.CurrentBudgetID,
		Month:           data.Month,
		Value:           data.Value,
		CreatedAt:       data.CreatedAt,
	}
}

func ToSpendingReleaseListResponseDTO(spendingreleases []*data.SpendingRelease) []SpendingReleaseResponseDTO {
	dtoList := make([]SpendingReleaseResponseDTO, len(spendingreleases))
	for i, x := range spendingreleases {
		dtoList[i] = ToSpendingReleaseResponseDTO(*x)
	}
	return dtoList
}
