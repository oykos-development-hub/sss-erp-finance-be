package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type AccountingEntryItemDTO struct {
	Title        string  `json:"title"`
	EntryID      int     `json:"entry_id"`
	AccountID    int     `json:"account_id"`
	CreditAmount float64 `json:"credit_amount"`
	DebitAmount  float64 `json:"debit_amount"`
	InvoiceID    int     `json:"invoice_id"`
	SalaryID     int     `json:"salary_id"`
	Type         string  `json:"type"`
}

type AccountingEntryItemResponseDTO struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	EntryID      int       `json:"entry_id"`
	AccountID    int       `json:"account_id"`
	CreditAmount float64   `json:"credit_amount"`
	DebitAmount  float64   `json:"debit_amount"`
	InvoiceID    int       `json:"invoice_id"`
	SalaryID     int       `json:"salary_id"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AccountingEntryItemFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	EntryID     *int    `json:"entry_id"`
}

func (dto AccountingEntryItemDTO) ToAccountingEntryItem() *data.AccountingEntryItem {
	return &data.AccountingEntryItem{
		Title: dto.Title,
	}
}

func ToAccountingEntryItemResponseDTO(data data.AccountingEntryItem) AccountingEntryItemResponseDTO {
	return AccountingEntryItemResponseDTO{
		ID:        data.ID,
		Title:     data.Title,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func ToAccountingEntryItemListResponseDTO(accounting_entry_items []*data.AccountingEntryItem) []AccountingEntryItemResponseDTO {
	dtoList := make([]AccountingEntryItemResponseDTO, len(accounting_entry_items))
	for i, x := range accounting_entry_items {
		dtoList[i] = ToAccountingEntryItemResponseDTO(*x)
	}
	return dtoList
}
