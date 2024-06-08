package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type FilledFinancialBudgetDTO struct {
	BudgetRequestID int             `json:"budget_request_id" validate:"required"`
	AccountID       int             `json:"account_id" validate:"required"`
	CurrentYear     decimal.Decimal `json:"current_year"`
	NextYear        decimal.Decimal `json:"next_year"`
	YearAfterNext   decimal.Decimal `json:"year_after_next"`
	Description     string          `json:"description"`
}

type FilledActualFinancialBudgetDTO struct {
	Actual decimal.Decimal `json:"actual" validate:"required"`
}

type FilledFinancialBudgetResponseDTO struct {
	ID              int                 `json:"id,omitempty"`
	BudgetRequestID int                 `json:"budget_request_id,omitempty"`
	AccountID       int                 `json:"account_id"`
	CurrentYear     decimal.Decimal     `json:"current_year"`
	NextYear        decimal.Decimal     `json:"next_year"`
	YearAfterNext   decimal.Decimal     `json:"year_after_next"`
	Actual          decimal.NullDecimal `json:"actual"`
	Description     string              `json:"description,omitempty"`
	CreatedAt       time.Time           `json:"created_at,omitempty"`
	UpdatedAt       time.Time           `json:"updated_at,omitempty"`
}

type FilledFinancialBudgetFilterDTO struct {
	Page            *int    `json:"page"`
	Size            *int    `json:"size"`
	SortByTitle     *string `json:"sort_by_title"`
	BudgetRequestID int     `json:"budget_request_id" validate:"required"`
	AccountIdList   []any   `json:"accounts"`
}

func (dto FilledFinancialBudgetDTO) ToFilledFinancialBudget() *data.FilledFinancialBudget {
	return &data.FilledFinancialBudget{
		BudgetRequestID: dto.BudgetRequestID,
		AccountID:       dto.AccountID,
		CurrentYear:     dto.CurrentYear,
		NextYear:        dto.NextYear,
		YearAfterNext:   dto.YearAfterNext,
		Description:     dto.Description,
	}
}

func ToFilledFinancialBudgetResponseDTO(data *data.FilledFinancialBudget) FilledFinancialBudgetResponseDTO {
	return FilledFinancialBudgetResponseDTO{
		ID:              data.ID,
		BudgetRequestID: data.BudgetRequestID,
		AccountID:       data.AccountID,
		CurrentYear:     data.CurrentYear,
		NextYear:        data.NextYear,
		YearAfterNext:   data.YearAfterNext,
		Actual:          data.Actual,
		Description:     data.Description,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
	}
}

func ToFilledFinancialBudgetListResponseDTO(filledfinancialbudgets []data.FilledFinancialBudget) []FilledFinancialBudgetResponseDTO {
	dtoList := make([]FilledFinancialBudgetResponseDTO, len(filledfinancialbudgets))
	for i, x := range filledfinancialbudgets {
		dtoList[i] = ToFilledFinancialBudgetResponseDTO(&x)
	}
	return dtoList
}
