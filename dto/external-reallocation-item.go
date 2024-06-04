package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type ExternalReallocationItemDTO struct {
	ReallocationID       int             `json:"reallocation_id"`
	SourceAccountID      int             `json:"source_account_id"`
	DestinationAccountID int             `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}

type ExternalReallocationItemResponseDTO struct {
	ID                   int             `json:"id"`
	ReallocationID       int             `json:"reallocation_id"`
	SourceAccountID      int             `json:"source_account_id"`
	DestinationAccountID int             `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type ExternalReallocationItemFilterDTO struct {
	Page           *int    `json:"page"`
	Size           *int    `json:"size"`
	SortByTitle    *string `json:"sort_by_title"`
	ReallocationID *int    `json:"reallocation_id"`
}

func (dto ExternalReallocationItemDTO) ToExternalReallocationItem() *data.ExternalReallocationItem {
	return &data.ExternalReallocationItem{
		ReallocationID:       dto.ReallocationID,
		SourceAccountID:      dto.SourceAccountID,
		DestinationAccountID: dto.DestinationAccountID,
		Amount:               dto.Amount,
	}
}

func ToExternalReallocationItemResponseDTO(data data.ExternalReallocationItem) ExternalReallocationItemResponseDTO {
	return ExternalReallocationItemResponseDTO{
		ID:                   data.ID,
		ReallocationID:       data.ReallocationID,
		SourceAccountID:      data.SourceAccountID,
		DestinationAccountID: data.DestinationAccountID,
		Amount:               data.Amount,
		CreatedAt:            data.CreatedAt,
		UpdatedAt:            data.UpdatedAt,
	}
}

func ToExternalReallocationItemListResponseDTO(external_reallocation_items []*data.ExternalReallocationItem) []ExternalReallocationItemResponseDTO {
	dtoList := make([]ExternalReallocationItemResponseDTO, len(external_reallocation_items))
	for i, x := range external_reallocation_items {
		dtoList[i] = ToExternalReallocationItemResponseDTO(*x)
	}
	return dtoList
}
