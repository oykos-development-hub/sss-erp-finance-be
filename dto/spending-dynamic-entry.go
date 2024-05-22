package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type SpendingDynamicEntryResponseDTO struct {
	ID                int             `json:"id"`
	SpendingDynamicID int             `json:"spending_dynamic_id"`
	January           decimal.Decimal `json:"january"`
	February          decimal.Decimal `json:"february"`
	March             decimal.Decimal `json:"march"`
	April             decimal.Decimal `json:"april"`
	May               decimal.Decimal `json:"may"`
	June              decimal.Decimal `json:"june"`
	July              decimal.Decimal `json:"july"`
	August            decimal.Decimal `json:"august"`
	September         decimal.Decimal `json:"september"`
	October           decimal.Decimal `json:"october"`
	November          decimal.Decimal `json:"november"`
	December          decimal.Decimal `json:"december"`
	CreatedAt         time.Time       `json:"created_at"`
}

func ToSpendingDynamicEntryResponseDTO(data *data.SpendingDynamicEntry) *SpendingDynamicEntryResponseDTO {
	return &SpendingDynamicEntryResponseDTO{
		ID:                data.ID,
		SpendingDynamicID: data.SpendingDynamicID,
		January:           data.January,
		February:          data.February,
		March:             data.March,
		April:             data.April,
		May:               data.May,
		June:              data.June,
		July:              data.July,
		August:            data.August,
		September:         data.September,
		October:           data.October,
		November:          data.November,
		December:          data.December,
		CreatedAt:         data.CreatedAt,
	}
}

func ToSpendingDynamicEntryListResponseDTO(spendingdynamics []data.SpendingDynamicEntry) []SpendingDynamicEntryResponseDTO {
	dtoList := make([]SpendingDynamicEntryResponseDTO, len(spendingdynamics))
	for i, x := range spendingdynamics {
		dtoList[i] = *ToSpendingDynamicEntryResponseDTO(&x)
	}
	return dtoList
}
