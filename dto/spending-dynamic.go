package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type SpendingDynamicDTO struct {
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
	ID              int             `json:"id"`
	AccountID       int             `json:"account_id"`
	BudgetID        int             `json:"budget_id"`
	UnitID          int             `json:"unit_id"`
	CurrentBudgetID int             `json:"current_budget_id"`
	Actual          decimal.Decimal `json:"actual"`
	TotalSavings    decimal.Decimal `json:"total_savings"`
	Username        string          `json:"username"`
	January         MonthEntry      `json:"january"`
	February        MonthEntry      `json:"february"`
	March           MonthEntry      `json:"march"`
	April           MonthEntry      `json:"april"`
	May             MonthEntry      `json:"may"`
	June            MonthEntry      `json:"june"`
	July            MonthEntry      `json:"july"`
	August          MonthEntry      `json:"august"`
	September       MonthEntry      `json:"september"`
	October         MonthEntry      `json:"october"`
	November        MonthEntry      `json:"november"`
	December        MonthEntry      `json:"december"`
	CreatedAt       time.Time       `json:"created_at"`
}

func (t *SpendingDynamicWithEntryResponseDTO) GetTotalSavings() decimal.Decimal {
	return decimal.Sum(
		t.January.Savings,
		t.February.Savings,
		t.March.Savings,
		t.April.Savings,
		t.May.Savings,
		t.June.Savings,
		t.July.Savings,
		t.August.Savings,
		t.September.Savings,
		t.October.Savings,
		t.November.Savings,
		t.December.Savings,
	)
}

type MonthEntry struct {
	Value   decimal.Decimal `json:"value"`
	Savings decimal.Decimal `json:"savings"`
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
