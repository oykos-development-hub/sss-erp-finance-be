package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ExternalReallocationDTO struct {
	ID                            int                           `json:"id"`
	Title                         string                        `json:"title"`
	Status                        data.ReallocationStatus       `json:"status"`
	SourceOrganizationUnitID      int                           `json:"source_organization_unit_id"`
	DestinationOrganizationUnitID int                           `json:"destination_organization_unit_id"`
	DateOfRequest                 time.Time                     `json:"date_of_request"`
	DateOfActionDestOrgUnit       time.Time                     `json:"date_of_action_dest_org_unit"`
	DateOfActionSSS               time.Time                     `json:"date_of_action_sss"`
	RequestedBy                   int                           `json:"requested_by"`
	AcceptedBy                    int                           `json:"accepted_by"`
	FileID                        int                           `json:"file_id"`
	DestinationOrgUnitFileID      int                           `json:"destination_org_unit_file_id"`
	SSSFileID                     int                           `json:"sss_file_id"`
	BudgetID                      int                           `json:"budget_id"`
	Items                         []ExternalReallocationItemDTO `json:"items"`
}

type ExternalReallocationResponseDTO struct {
	ID                            int                                   `json:"id"`
	Title                         string                                `json:"title"`
	Status                        data.ReallocationStatus               `json:"status"`
	SourceOrganizationUnitID      int                                   `json:"source_organization_unit_id"`
	DestinationOrganizationUnitID int                                   `json:"destination_organization_unit_id"`
	DateOfRequest                 time.Time                             `json:"date_of_request"`
	DateOfActionDestOrgUnit       time.Time                             `json:"date_of_action_dest_org_unit"`
	DateOfActionSSS               time.Time                             `json:"date_of_action_sss"`
	RequestedBy                   int                                   `json:"requested_by"`
	AcceptedBy                    int                                   `json:"accepted_by"`
	FileID                        int                                   `json:"file_id"`
	DestinationOrgUnitFileID      int                                   `json:"destination_org_unit_file_id"`
	SSSFileID                     int                                   `json:"sss_file_id"`
	BudgetID                      int                                   `json:"budget_id"`
	Items                         []ExternalReallocationItemResponseDTO `json:"items"`
	CreatedAt                     time.Time                             `json:"created_at"`
	UpdatedAt                     time.Time                             `json:"updated_at"`
}

type ExternalReallocationFilterDTO struct {
	Page                          *int    `json:"page"`
	Size                          *int    `json:"size"`
	SortByTitle                   *string `json:"sort_by_title"`
	Status                        *string `json:"status"`
	SourceOrganizationUnitID      *int    `json:"source_organization_unit_id"`
	DestinationOrganizationUnitID *int    `json:"destination_organization_unit_id"`
	OrganizationUnitID            *int    `json:"organization_unit_id"`
	RequestedBy                   *int    `json:"requested_by"`
	BudgetID                      *int    `json:"budget_id"`
}

func (dto ExternalReallocationDTO) ToExternalReallocation() *data.ExternalReallocation {
	return &data.ExternalReallocation{
		Title:                         dto.Title,
		Status:                        dto.Status,
		SourceOrganizationUnitID:      dto.SourceOrganizationUnitID,
		DestinationOrganizationUnitID: dto.DestinationOrganizationUnitID,
		DateOfRequest:                 dto.DateOfRequest,
		DateOfActionDestOrgUnit:       dto.DateOfActionDestOrgUnit,
		DateOfActionSSS:               dto.DateOfActionSSS,
		RequestedBy:                   dto.RequestedBy,
		AcceptedBy:                    dto.AcceptedBy,
		FileID:                        dto.FileID,
		DestinationOrgUnitFileID:      dto.DestinationOrgUnitFileID,
		SSSFileID:                     dto.SSSFileID,
		BudgetID:                      dto.BudgetID,
	}
}

func ToExternalReallocationResponseDTO(data data.ExternalReallocation) ExternalReallocationResponseDTO {
	return ExternalReallocationResponseDTO{
		ID:                            data.ID,
		Title:                         data.Title,
		Status:                        data.Status,
		SourceOrganizationUnitID:      data.SourceOrganizationUnitID,
		DestinationOrganizationUnitID: data.DestinationOrganizationUnitID,
		DateOfRequest:                 data.DateOfRequest,
		DateOfActionDestOrgUnit:       data.DateOfActionDestOrgUnit,
		DateOfActionSSS:               data.DateOfActionSSS,
		RequestedBy:                   data.RequestedBy,
		AcceptedBy:                    data.AcceptedBy,
		FileID:                        data.FileID,
		DestinationOrgUnitFileID:      data.DestinationOrgUnitFileID,
		SSSFileID:                     data.SSSFileID,
		BudgetID:                      data.BudgetID,
		CreatedAt:                     data.CreatedAt,
		UpdatedAt:                     data.UpdatedAt,
	}
}

func ToExternalReallocationListResponseDTO(external_reallocations []*data.ExternalReallocation) []ExternalReallocationResponseDTO {
	dtoList := make([]ExternalReallocationResponseDTO, len(external_reallocations))
	for i, x := range external_reallocations {
		dtoList[i] = ToExternalReallocationResponseDTO(*x)
	}
	return dtoList
}
