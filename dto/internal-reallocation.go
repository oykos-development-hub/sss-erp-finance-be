package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type InternalReallocationDTO struct {
	Title              string                        `json:"title"`
	OrganizationUnitID int                           `json:"organization_unit_id"`
	DateOfRequest      time.Time                     `json:"date_of_request"`
	RequestedBy        int                           `json:"requested_by"`
	FileID             int                           `json:"file_id"`
	BudgetID           int                           `json:"budget_id"`
	Items              []InternalReallocationItemDTO `json:"items"`
}

type InternalReallocationResponseDTO struct {
	ID                 int                                   `json:"id"`
	Title              string                                `json:"title"`
	OrganizationUnitID int                                   `json:"organization_unit_id"`
	DateOfRequest      time.Time                             `json:"date_of_request"`
	RequestedBy        int                                   `json:"requested_by"`
	FileID             int                                   `json:"file_id"`
	BudgetID           int                                   `json:"budget_id"`
	Sum                decimal.Decimal                       `json:"sum"`
	Items              []InternalReallocationItemResponseDTO `json:"items"`
	CreatedAt          time.Time                             `json:"created_at"`
	UpdatedAt          time.Time                             `json:"updated_at"`
}

type InternalReallocationFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	Year               *int    `json:"year"`
	RequestedBy        *int    `json:"requested_by"`
	BudgetID           *int    `json:"budget_id"`
}

func (dto InternalReallocationDTO) ToInternalReallocation() *data.InternalReallocation {
	return &data.InternalReallocation{
		Title:              dto.Title,
		OrganizationUnitID: dto.OrganizationUnitID,
		DateOfRequest:      dto.DateOfRequest,
		RequestedBy:        dto.RequestedBy,
		FileID:             dto.FileID,
		BudgetID:           dto.BudgetID,
	}
}

func ToInternalReallocationResponseDTO(data data.InternalReallocation) InternalReallocationResponseDTO {
	return InternalReallocationResponseDTO{
		ID:                 data.ID,
		Title:              data.Title,
		OrganizationUnitID: data.OrganizationUnitID,
		DateOfRequest:      data.DateOfRequest,
		RequestedBy:        data.RequestedBy,
		FileID:             data.FileID,
		BudgetID:           data.BudgetID,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToInternalReallocationListResponseDTO(internal_reallocations []*data.InternalReallocation) []InternalReallocationResponseDTO {
	dtoList := make([]InternalReallocationResponseDTO, len(internal_reallocations))
	for i, x := range internal_reallocations {
		dtoList[i] = ToInternalReallocationResponseDTO(*x)
	}
	return dtoList
}
