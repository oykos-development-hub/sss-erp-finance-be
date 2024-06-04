package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type InternalReallocationItemDTO struct {
	ReallocationID       int             `json:"reallocation_id"`
	SourceAccountID      int             `json:"source_account_id"`
	DestinationAccountID int             `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}

type InternalReallocationItemResponseDTO struct {
	ID                   int             `json:"id"`
	ReallocationID       int             `json:"reallocation_id"`
	SourceAccountID      int             `json:"source_account_id"`
	DestinationAccountID int             `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type InternalReallocationItemFilterDTO struct {
	Page           *int    `json:"page"`
	Size           *int    `json:"size"`
	SortByTitle    *string `json:"sort_by_title"`
	ReallocationID *int    `json:"reallocation_id"`
}

func (dto InternalReallocationItemDTO) ToInternalReallocationItem() *data.InternalReallocationItem {
	return &data.InternalReallocationItem{
		ReallocationID:       dto.ReallocationID,
		SourceAccountID:      dto.SourceAccountID,
		DestinationAccountID: dto.DestinationAccountID,
		Amount:               dto.Amount,
	}
}

func ToInternalReallocationItemResponseDTO(data data.InternalReallocationItem) InternalReallocationItemResponseDTO {
	return InternalReallocationItemResponseDTO{
		ID:                   data.ID,
		ReallocationID:       data.ReallocationID,
		SourceAccountID:      data.SourceAccountID,
		DestinationAccountID: data.DestinationAccountID,
		Amount:               data.Amount,
		CreatedAt:            data.CreatedAt,
		UpdatedAt:            data.UpdatedAt,
	}
}

func ToInternalReallocationItemListResponseDTO(internal_reallocation_items []*data.InternalReallocationItem) []InternalReallocationItemResponseDTO {
	dtoList := make([]InternalReallocationItemResponseDTO, len(internal_reallocation_items))
	for i, x := range internal_reallocation_items {
		dtoList[i] = ToInternalReallocationItemResponseDTO(*x)
	}
	return dtoList
}
