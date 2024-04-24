package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type AccountingEntryDTO struct {
	Title string `json:"title" validate:"required,min=2"`
}

type AccountingEntryResponseDTO struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountingEntryFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
}

type ObligationForAccounting struct {
	InvoiceID *int      `json:"invoice_id"`
	SalaryID  *int      `json:"salary_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (dto AccountingEntryDTO) ToAccountingEntry() *data.AccountingEntry {
	return &data.AccountingEntry{
		Title: dto.Title,
	}
}

func ToAccountingEntryResponseDTO(data data.AccountingEntry) AccountingEntryResponseDTO {
	return AccountingEntryResponseDTO{
		ID:        data.ID,
		Title:     data.Title,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func ToAccountingEntryListResponseDTO(accounting_entries []*data.AccountingEntry) []AccountingEntryResponseDTO {
	dtoList := make([]AccountingEntryResponseDTO, len(accounting_entries))
	for i, x := range accounting_entries {
		dtoList[i] = ToAccountingEntryResponseDTO(*x)
	}
	return dtoList
}
