package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type SpendingDynamicFilter struct {
	All *bool `json:"all"`
}

type SpendingDynamicDTO struct {
	BudgetID  int             `json:"budget_id" validate:"required"`
	UnitID    int             `json:"unit_id" validate:"required"`
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

type SpendingDynamicWithEntriesResponseDTO struct {
	ID           int                               `json:"id"`
	BudgetID     int                               `json:"budget_id"`
	UnitID       int                               `json:"unit_id"`
	PlannedTotal decimal.Decimal                   `json:"actual"`
	Entries      []SpendingDynamicEntryResponseDTO `json:"entries"`
	CreatedAt    time.Time                         `json:"created_at"`
	UpdatedAt    time.Time                         `json:"updated_at"`
}

func (dto SpendingDynamicDTO) ToSpendingDynamic() *data.SpendingDynamic {
	return &data.SpendingDynamic{
		BudgetID: dto.BudgetID,
		UnitID:   dto.UnitID,
	}
}

func (dto SpendingDynamicDTO) ToSpendingDynamicEntry() *data.SpendingDynamicEntry {
	return &data.SpendingDynamicEntry{
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

func ToSpendingDynamicWithEntryResponseDTO(data *data.SpendingDynamic, entries []data.SpendingDynamicEntry) *SpendingDynamicWithEntriesResponseDTO {
	return &SpendingDynamicWithEntriesResponseDTO{
		ID:           data.ID,
		BudgetID:     data.BudgetID,
		UnitID:       data.UnitID,
		PlannedTotal: data.PlannedTotal,
		Entries:      ToSpendingDynamicEntryListResponseDTO(entries),
	}
}
