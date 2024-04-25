package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ModelOfAccountingItemDTO struct {
	Title           string `json:"title"`
	ModelID         int    `json:"model_id"`
	DebitAccountID  int    `json:"debit_account_id"`
	CreditAccountID int    `json:"credit_account_id"`
}

type ModelOfAccountingItemResponseDTO struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	ModelID         int       `json:"model_id"`
	DebitAccountID  int       `json:"debit_account_id"`
	CreditAccountID int       `json:"credit_account_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ModelOfAccountingItemFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
}

func (dto ModelOfAccountingItemDTO) ToModelOfAccountingItem() *data.ModelOfAccountingItem {
	return &data.ModelOfAccountingItem{
		Title: dto.Title,
	}
}

func ToModelOfAccountingItemResponseDTO(data data.ModelOfAccountingItem) ModelOfAccountingItemResponseDTO {
	return ModelOfAccountingItemResponseDTO{
		ID:        data.ID,
		Title:     data.Title,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func ToModelOfAccountingItemListResponseDTO(model_of_accounting_items []*data.ModelOfAccountingItem) []ModelOfAccountingItemResponseDTO {
	dtoList := make([]ModelOfAccountingItemResponseDTO, len(model_of_accounting_items))
	for i, x := range model_of_accounting_items {
		dtoList[i] = ToModelOfAccountingItemResponseDTO(*x)
	}
	return dtoList
}
