package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FilledFinancialBudgetDTO struct {
	OrganizationUnitID int    `json:"organization_unit_id" validate:"required"`
	BudgetRequestID    int    `json:"budget_request_id" validate:"required"`
	AccountID          int    `json:"account_id" validate:"required"`
	CurrentYear        int    `json:"current_year"`
	NextYear           int    `json:"next_year"`
	YearAfterNext      int    `json:"year_after_next"`
	Description        string `json:"description"`
}

type FilledFinancialBudgetResponseDTO struct {
	ID                 int       `json:"id"`
	OrganizationUnitID int       `json:"organization_unit_id"`
	BudgetRequestID    int       `json:"budget_request_id"`
	AccountID          int       `json:"account_id"`
	CurrentYear        int       `json:"current_year"`
	NextYear           int       `json:"next_year"`
	YearAfterNext      int       `json:"year_after_next"`
	Description        string    `json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type FilledFinancialBudgetFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	OrganizationUnitID int     `json:"organization_unit_id" validate:"required"`
	BudgetRequestID    int     `json:"budget_request_id" validate:"required"`
}

func (dto FilledFinancialBudgetDTO) ToFilledFinancialBudget() *data.FilledFinancialBudget {
	return &data.FilledFinancialBudget{
		OrganizationUnitID: dto.OrganizationUnitID,
		BudgetRequestID:    dto.BudgetRequestID,
		AccountID:          dto.AccountID,
		CurrentYear:        dto.CurrentYear,
		NextYear:           dto.NextYear,
		YearAfterNext:      dto.YearAfterNext,
		Description:        dto.Description,
	}
}

func ToFilledFinancialBudgetResponseDTO(data data.FilledFinancialBudget) FilledFinancialBudgetResponseDTO {
	return FilledFinancialBudgetResponseDTO{
		ID:                 data.ID,
		OrganizationUnitID: data.OrganizationUnitID,
		BudgetRequestID:    data.BudgetRequestID,
		AccountID:          data.AccountID,
		CurrentYear:        data.CurrentYear,
		NextYear:           data.NextYear,
		YearAfterNext:      data.YearAfterNext,
		Description:        data.Description,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToFilledFinancialBudgetListResponseDTO(filledfinancialbudgets []*data.FilledFinancialBudget) []FilledFinancialBudgetResponseDTO {
	dtoList := make([]FilledFinancialBudgetResponseDTO, len(filledfinancialbudgets))
	for i, x := range filledfinancialbudgets {
		dtoList[i] = ToFilledFinancialBudgetResponseDTO(*x)
	}
	return dtoList
}
