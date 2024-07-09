package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type SpendingReleaseRequestDTO struct {
	Year                   int                        `json:"year"`
	Month                  int                        `json:"month"`
	OrganizationUnitID     int                        `json:"organization_unit_id"`
	OrganizationUnitFileID int                        `json:"organization_unit_file_id"`
	SSSFileID              int                        `json:"sss_file_id"`
	Status                 data.SpendingReleaseStatus `json:"status"`
}

type SpendingReleaseRequestResponseDTO struct {
	ID                     int                        `json:"id"`
	Year                   int                        `json:"year"`
	Month                  int                        `json:"month"`
	OrganizationUnitID     int                        `json:"organization_unit_id"`
	OrganizationUnitFileID int                        `json:"organization_unit_file_id"`
	SSSFileID              int                        `json:"sss_file_id"`
	Status                 data.SpendingReleaseStatus `json:"status"`
	CreatedAt              time.Time                  `json:"created_at"`
}

type SpendingReleaseRequestFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	Status             *string `json:"status"`
}

func (dto SpendingReleaseRequestDTO) ToSpendingReleaseRequest() *data.SpendingReleaseRequest {
	return &data.SpendingReleaseRequest{
		Year:                   dto.Year,
		Month:                  dto.Month,
		OrganizationUnitID:     dto.OrganizationUnitID,
		OrganizationUnitFileID: dto.OrganizationUnitFileID,
		SSSFileID:              dto.SSSFileID,
		Status:                 dto.Status,
	}
}

func ToSpendingReleaseRequestResponseDTO(data data.SpendingReleaseRequest) SpendingReleaseRequestResponseDTO {
	return SpendingReleaseRequestResponseDTO{
		ID:                     data.ID,
		Year:                   data.Year,
		Month:                  data.Month,
		OrganizationUnitID:     data.OrganizationUnitID,
		OrganizationUnitFileID: data.OrganizationUnitFileID,
		SSSFileID:              data.SSSFileID,
		Status:                 data.Status,
		CreatedAt:              data.CreatedAt,
	}
}

func ToSpendingReleaseRequestListResponseDTO(spending_release_requests []*data.SpendingReleaseRequest) []SpendingReleaseRequestResponseDTO {
	dtoList := make([]SpendingReleaseRequestResponseDTO, len(spending_release_requests))
	for i, x := range spending_release_requests {
		dtoList[i] = ToSpendingReleaseRequestResponseDTO(*x)
	}
	return dtoList
}
