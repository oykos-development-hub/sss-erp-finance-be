package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type CurrentBudgetDTO struct {
	BudgetID      int             `json:"budget_id"`
	UnitID        int             `json:"unit_id"`
	AccountID     int             `json:"account_id"`
	InitialActual decimal.Decimal `json:"initial_actual"`
	Actual        decimal.Decimal `json:"actual"`
	Balance       decimal.Decimal `json:"balance"`
}

type CurrentBudgetResponseDTO struct {
	ID            int             `json:"id"`
	BudgetID      int             `json:"budget_id"`
	UnitID        int             `json:"unit_id"`
	AccountID     int             `json:"account_id"`
	InitialActual decimal.Decimal `json:"initial_actual"`
	Actual        decimal.Decimal `json:"actual"`
	Balance       decimal.Decimal `json:"balance"`
	CreatedAt     time.Time       `json:"created_at,omitempty"`
}

type CurrentBudgetFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	UnitID      *int    `json:"unit_id"`
	AccountID   *int    `json:"account_id"`
}

func (dto CurrentBudgetDTO) ToCurrentBudget() *data.CurrentBudget {
	return &data.CurrentBudget{
		BudgetID:      dto.BudgetID,
		UnitID:        dto.UnitID,
		AccountID:     dto.AccountID,
		InitialActual: dto.InitialActual,
		Actual:        dto.Actual,
		Balance:       dto.Balance,
	}
}

func ToCurrentBudgetResponseDTO(data *data.CurrentBudget) CurrentBudgetResponseDTO {
	return CurrentBudgetResponseDTO{
		ID:            data.ID,
		BudgetID:      data.BudgetID,
		UnitID:        data.UnitID,
		AccountID:     data.AccountID,
		InitialActual: data.InitialActual,
		Actual:        data.Actual,
		Balance:       data.Balance,
		CreatedAt:     data.CreatedAt,
	}
}

func ToCurrentBudgetListResponseDTO(currentbudgets []*data.CurrentBudget) []CurrentBudgetResponseDTO {
	dtoList := make([]CurrentBudgetResponseDTO, len(currentbudgets))
	for i, x := range currentbudgets {
		dtoList[i] = ToCurrentBudgetResponseDTO(x)
	}
	return dtoList
}
