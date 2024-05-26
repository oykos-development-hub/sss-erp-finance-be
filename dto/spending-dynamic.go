package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type SpendingDynamicDTO struct {
	BudgetID  int             `json:"budget_id" validate:"required"`
	UnitID    int             `json:"unit_id" validate:"required"`
	AccountID int             `json:"account_id" validate:"required"`
	Username  string          `json:"username" validate:"required"`
	January   decimal.Decimal `json:"january" validate:"required"`
	February  decimal.Decimal `json:"february" validate:"required"`
	March     decimal.Decimal `json:"march" validate:"required"`
	April     decimal.Decimal `json:"april" validate:"required"`
	May       decimal.Decimal `json:"may" validate:"required"`
	June      decimal.Decimal `json:"june" validate:"required"`
	July      decimal.Decimal `json:"july" validate:"required"`
	August    decimal.Decimal `json:"august" validate:"required"`
	September decimal.Decimal `json:"september" validate:"required"`
	October   decimal.Decimal `json:"october" validate:"required"`
	November  decimal.Decimal `json:"november" validate:"required"`
	December  decimal.Decimal `json:"december" validate:"required"`
}

type SpendingDynamicHistoryResponseDTO struct {
	BudgetID  int       `json:"budget_id"`
	UnitID    int       `json:"unit_id"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
}

type SpendingDynamicFilter struct {
	Version *int `json:"version"`
}

type SpendingDynamicWithEntryResponseDTO struct {
	ID              int `json:"id"`
	CurrentBudgetID int `json:"current_budget_id"`
	AccountID       int `json:"account_id"`
	BudgetID        int `json:"budget_id"`
	UnitID          int `json:"unit_id"`
	SpendingDynamicEntryResponseDTO
	Actual decimal.Decimal `json:"actual"`
}

func (dto SpendingDynamicDTO) ToSpendingDynamicEntry() *data.SpendingDynamicEntry {
	return &data.SpendingDynamicEntry{
		Username:  dto.Username,
		January:   dto.January,
		February:  dto.February,
		March:     dto.March,
		April:     dto.April,
		May:       dto.May,
		June:      dto.June,
		July:      dto.July,
		August:    dto.August,
		September: dto.September,
		October:   dto.October,
		November:  dto.November,
		December:  dto.December,
	}
}

func ToSpendingDynamicWithEntryResponseDTO(data *data.SpendingDynamic, entry *data.SpendingDynamicEntry, currentBudget *data.CurrentBudget) *SpendingDynamicWithEntryResponseDTO {
	return &SpendingDynamicWithEntryResponseDTO{
		ID:                              data.ID,
		CurrentBudgetID:                 data.CurrentBudgetID,
		AccountID:                       currentBudget.AccountID,
		BudgetID:                        currentBudget.BudgetID,
		UnitID:                          currentBudget.UnitID,
		Actual:                          currentBudget.Actual,
		SpendingDynamicEntryResponseDTO: *ToSpendingDynamicEntryResponseDTO(entry),
	}
}
