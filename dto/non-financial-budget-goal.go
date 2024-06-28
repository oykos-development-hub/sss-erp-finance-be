package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type NonFinancialBudgetGoalDTO struct {
	NonFinancialBudgetID int    `json:"non_financial_budget_id" validate:"required"`
	Title                string `json:"title" validate:"required"`
	Description          string `json:"description"`
}

type NonFinancialBudgetGoalResponseDTO struct {
	ID                   int       `json:"id"`
	NonFinancialBudgetID int       `json:"non_financial_budget_id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type NonFinancialBudgetGoalFilterDTO struct {
	NonFinancialBudgetID *int `json:"non_financial_budget_id"`
}

func (dto NonFinancialBudgetGoalDTO) ToNonFinancialBudgetGoal() *data.NonFinancialBudgetGoal {
	return &data.NonFinancialBudgetGoal{
		NonFinancialBudgetID: dto.NonFinancialBudgetID,
		Title:                dto.Title,
		Description:          dto.Description,
	}
}

func ToNonFinancialBudgetGoalResponseDTO(data data.NonFinancialBudgetGoal) NonFinancialBudgetGoalResponseDTO {
	return NonFinancialBudgetGoalResponseDTO{
		ID:                   data.ID,
		NonFinancialBudgetID: data.NonFinancialBudgetID,
		Title:                data.Title,
		Description:          data.Description,
		CreatedAt:            data.CreatedAt,
		UpdatedAt:            data.UpdatedAt,
	}
}

func ToNonFinancialBudgetGoalListResponseDTO(nonfinancialbudgetgoals []*data.NonFinancialBudgetGoal) []NonFinancialBudgetGoalResponseDTO {
	dtoList := make([]NonFinancialBudgetGoalResponseDTO, len(nonfinancialbudgetgoals))
	for i, x := range nonfinancialbudgetgoals {
		dtoList[i] = ToNonFinancialBudgetGoalResponseDTO(*x)
	}
	return dtoList
}
