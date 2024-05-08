package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ModelsOfAccountingDTO struct {
	Items []ModelOfAccountingItemDTO `json:"items"`
	Title string                     `json:"title"`
	Type  data.TypesOfObligation     `json:"type"`
}

type ModelsOfAccountingResponseDTO struct {
	ID        int                                `json:"id"`
	Title     string                             `json:"title"`
	Type      data.TypesOfObligation             `json:"type"`
	Items     []ModelOfAccountingItemResponseDTO `json:"items"`
	CreatedAt time.Time                          `json:"created_at"`
	UpdatedAt time.Time                          `json:"updated_at"`
}

type ModelsOfAccountingFilterDTO struct {
	Page        *int                    `json:"page"`
	Size        *int                    `json:"size"`
	Type        *data.TypesOfObligation `json:"type"`
	Search      *string                 `json:"search"`
	SortByTitle *string                 `json:"sort_by_title"`
}

func (dto ModelsOfAccountingDTO) ToModelsOfAccounting() *data.ModelsOfAccounting {
	return &data.ModelsOfAccounting{
		Title: dto.Title,
		Type:  dto.Type,
	}
}

func ToModelsOfAccountingResponseDTO(data data.ModelsOfAccounting) ModelsOfAccountingResponseDTO {
	return ModelsOfAccountingResponseDTO{
		ID:        data.ID,
		Title:     data.Title,
		Type:      data.Type,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func ToModelsOfAccountingListResponseDTO(models_of_accountings []*data.ModelsOfAccounting) []ModelsOfAccountingResponseDTO {
	dtoList := make([]ModelsOfAccountingResponseDTO, len(models_of_accountings))
	for i, x := range models_of_accountings {
		dtoList[i] = ToModelsOfAccountingResponseDTO(*x)
	}
	return dtoList
}
