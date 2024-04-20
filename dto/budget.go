package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type GetBudgetListInput struct {
	Year         *int `json:"year"`
	BudgetType   *int `json:"budget_type"`
	BudgetStatus *int `json:"budget_status"`
}

type BudgetDTO struct {
	Year         int               `json:"year" validate:"required,gte=2024"`
	BudgetType   data.BudgetType   `json:"budget_type" validate:"required,oneof=1 2"`
	BudgetStatus data.BudgetStatus `json:"budget_status" validate:"required,oneof=1 2 3 4 5"`
}

type BudgetResponseDTO struct {
	ID           int               `json:"id"`
	Year         int               `json:"year"`
	BudgetType   data.BudgetType   `json:"budget_type"`
	BudgetStatus data.BudgetStatus `json:"budget_status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func (dto BudgetDTO) ToBudget() *data.Budget {
	return &data.Budget{
		Year:         dto.Year,
		BudgetType:   dto.BudgetType,
		BudgetStatus: dto.BudgetStatus,
	}
}

func ToBudgetResponseDTO(data data.Budget) BudgetResponseDTO {
	return BudgetResponseDTO{
		ID:           data.ID,
		Year:         data.Year,
		BudgetType:   data.BudgetType,
		BudgetStatus: data.BudgetStatus,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
}

func ToBudgetListResponseDTO(budgets []*data.Budget) []BudgetResponseDTO {
	dtoList := make([]BudgetResponseDTO, len(budgets))
	for i, x := range budgets {
		dtoList[i] = ToBudgetResponseDTO(*x)
	}
	return dtoList
}
