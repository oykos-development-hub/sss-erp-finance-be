package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type BudgetRequestDTO struct {
	OrganizationUnitID int                      `json:"organization_unit_id" validate:"required"`
	BudgetID           int                      `json:"budget_id" validate:"required"`
	RequestType        data.RequestType         `json:"request_type" validate:"required"`
	Status             data.BudgetRequestStatus `json:"status" validate:"required"`
}

type BudgetRequestResponseDTO struct {
	ID                 int                      `json:"id"`
	OrganizationUnitID int                      `json:"organization_unit_id"`
	BudgetID           int                      `json:"budget_id"`
	RequestType        data.RequestType         `json:"request_type"`
	Status             data.BudgetRequestStatus `json:"status"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
}

type BudgetRequestFilterDTO struct {
	Page               *int              `json:"page"`
	Size               *int              `json:"size"`
	OrganizationUnitID *int              `json:"organization_unit_id"`
	BudgetID           int               `json:"budget_id"`
	RequestType        *data.RequestType `json:"request_type"`
}

func (dto BudgetRequestDTO) ToBudgetRequest() *data.BudgetRequest {
	return &data.BudgetRequest{
		OrganizationUnitID: dto.OrganizationUnitID,
		BudgetID:           dto.BudgetID,
		Status:             dto.Status,
		RequestType:        dto.RequestType,
	}
}

func ToBudgetRequestResponseDTO(data data.BudgetRequest) BudgetRequestResponseDTO {
	return BudgetRequestResponseDTO{
		ID:                 data.ID,
		OrganizationUnitID: data.OrganizationUnitID,
		BudgetID:           data.BudgetID,
		RequestType:        data.RequestType,
		Status:             data.Status,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToBudgetRequestListResponseDTO(budgetfinancialrequests []*data.BudgetRequest) []BudgetRequestResponseDTO {
	dtoList := make([]BudgetRequestResponseDTO, len(budgetfinancialrequests))
	for i, x := range budgetfinancialrequests {
		dtoList[i] = ToBudgetRequestResponseDTO(*x)
	}
	return dtoList
}
