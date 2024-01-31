package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FinancialBudgetLimitDTO struct {
	FinancialBudgetID  int `json:"financial_budget_id" validate:"required"`
	OrganizationUnitID int `json:"organization_unit_id" validate:"required"`
	Limit              int `json:"limit" validate:"required"`
}

type FinancialBudgetLimitResponseDTO struct {
	ID                 int       `json:"id"`
	FinancialBudgetID  int       `json:"financial_budget_id"`
	OrganizationUnitID int       `json:"organization_unit_id"`
	Limit              int       `json:"limit"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type FinancialBudgetLimitFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
}

func (dto FinancialBudgetLimitDTO) ToFinancialBudgetLimit() *data.FinancialBudgetLimit {
	return &data.FinancialBudgetLimit{
		FinancialBudgetID:  dto.FinancialBudgetID,
		OrganizationUnitID: dto.OrganizationUnitID,
		Limit:              dto.Limit,
	}
}

func ToFinancialBudgetLimitResponseDTO(data data.FinancialBudgetLimit) FinancialBudgetLimitResponseDTO {
	return FinancialBudgetLimitResponseDTO{
		ID:                 data.ID,
		FinancialBudgetID:  data.FinancialBudgetID,
		OrganizationUnitID: data.OrganizationUnitID,
		Limit:              data.Limit,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToFinancialBudgetLimitListResponseDTO(financialbudgetlimits []*data.FinancialBudgetLimit) []FinancialBudgetLimitResponseDTO {
	dtoList := make([]FinancialBudgetLimitResponseDTO, len(financialbudgetlimits))
	for i, x := range financialbudgetlimits {
		dtoList[i] = ToFinancialBudgetLimitResponseDTO(*x)
	}
	return dtoList
}