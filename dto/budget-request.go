package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type BudgetRequestDTO struct {
	ParentID           *int                     `json:"parent_id"`
	OrganizationUnitID int                      `json:"organization_unit_id" validate:"required"`
	BudgetID           int                      `json:"budget_id" validate:"required"`
	RequestType        data.RequestType         `json:"request_type" validate:"required,oneof=1 2 3 4 5"`
	Status             data.BudgetRequestStatus `json:"status" validate:"required,oneof=1 2 3 4 5 6 7"`
	Comment            string                   `json:"comment"`
}

type BudgetRequestResponseDTO struct {
	ID                 int                      `json:"id"`
	ParentID           *int                     `json:"parent_id"`
	OrganizationUnitID int                      `json:"organization_unit_id"`
	BudgetID           int                      `json:"budget_id"`
	RequestType        data.RequestType         `json:"request_type"`
	Status             data.BudgetRequestStatus `json:"status"`
	Comment            string                   `json:"comment"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
}

type BudgetRequestFilterDTO struct {
	Page               *int              `json:"page"`
	Size               *int              `json:"size"`
	ParentID           *int              `json:"parent_id"`
	OrganizationUnitID *int              `json:"organization_unit_id"`
	BudgetID           *int              `json:"budget_id"`
	RequestType        *data.RequestType `json:"request_type"`
	RequestTypes       *[]interface{}    `json:"request_types"`
	Statuses           []interface{}     `json:"statuses"`
}

func (dto BudgetRequestDTO) ToBudgetRequest() *data.BudgetRequest {
	return &data.BudgetRequest{
		ParentID:           dto.ParentID,
		OrganizationUnitID: dto.OrganizationUnitID,
		BudgetID:           dto.BudgetID,
		Status:             dto.Status,
		RequestType:        dto.RequestType,
		Comment:            dto.Comment,
	}
}

func ToBudgetRequestResponseDTO(data data.BudgetRequest) BudgetRequestResponseDTO {
	return BudgetRequestResponseDTO{
		ID:                 data.ID,
		ParentID:           data.ParentID,
		OrganizationUnitID: data.OrganizationUnitID,
		BudgetID:           data.BudgetID,
		RequestType:        data.RequestType,
		Status:             data.Status,
		Comment:            data.Comment,
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
