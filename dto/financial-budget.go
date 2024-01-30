package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FinancialBudgetDTO struct {
	AccountVersion int `json:"account_version" validate:"required"`
	BudgetID       int `json:"budget_id" validate:"required"`
}

type FinancialBudgetResponseDTO struct {
	ID             int       `json:"id"`
	AccountVersion int       `json:"account_version"`
	BudgetID       int       `json:"budget_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (dto FinancialBudgetDTO) ToFinancialBudget() *data.FinancialBudget {
	return &data.FinancialBudget{
		AccountVersion: dto.AccountVersion,
		BudgetID:       dto.BudgetID,
	}
}

func ToFinancialBudgetResponseDTO(data data.FinancialBudget) FinancialBudgetResponseDTO {
	return FinancialBudgetResponseDTO{
		ID:             data.ID,
		AccountVersion: data.AccountVersion,
		BudgetID:       data.BudgetID,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
}

func ToFinancialBudgetListResponseDTO(financialbudgets []*data.FinancialBudget) []FinancialBudgetResponseDTO {
	dtoList := make([]FinancialBudgetResponseDTO, len(financialbudgets))
	for i, x := range financialbudgets {
		dtoList[i] = ToFinancialBudgetResponseDTO(*x)
	}
	return dtoList
}
