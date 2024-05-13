package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type AccountingEntryItemDTO struct {
	Title                   string  `json:"title"`
	EntryID                 int     `json:"entry_id"`
	SupplierID              int     `json:"supplier_id"`
	AccountID               int     `json:"account_id"`
	CreditAmount            float64 `json:"credit_amount"`
	DebitAmount             float64 `json:"debit_amount"`
	InvoiceID               *int    `json:"invoice_id"`
	SalaryID                *int    `json:"salary_id"`
	PaymentOrderID          *int    `json:"payment_order_id"`
	EnforcedPaymentID       *int    `json:"enforced_payment_id"`
	ReturnEnforcedPaymentID *int    `json:"return_enforced_payment_id"`
	Type                    string  `json:"type"`
	Date                    string  `json:"date"`
}

type AccountingEntryItemResponseDTO struct {
	ID                      int       `json:"id"`
	Title                   string    `json:"title"`
	SupplierID              int       `json:"supplier_id"`
	EntryID                 int       `json:"entry_id"`
	AccountID               int       `json:"account_id"`
	CreditAmount            float64   `json:"credit_amount"`
	DebitAmount             float64   `json:"debit_amount"`
	InvoiceID               *int      `json:"invoice_id"`
	SalaryID                *int      `json:"salary_id"`
	PaymentOrderID          *int      `json:"payment_order_id"`
	EnforcedPaymentID       *int      `json:"enforced_payment_id"`
	ReturnEnforcedPaymentID *int      `json:"return_enforced_payment_id"`
	Type                    string    `json:"type"`
	Date                    string    `json:"date"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type AccountingEntryItemFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	EntryID     *int    `json:"entry_id"`
}

func (dto AccountingEntryItemDTO) ToAccountingEntryItem() *data.AccountingEntryItem {
	return &data.AccountingEntryItem{
		Title:                   dto.Title,
		EntryID:                 dto.EntryID,
		SupplierID:              dto.SupplierID,
		AccountID:               dto.AccountID,
		CreditAmount:            dto.CreditAmount,
		DebitAmount:             dto.DebitAmount,
		InvoiceID:               dto.InvoiceID,
		SalaryID:                dto.SalaryID,
		PaymentOrderID:          dto.PaymentOrderID,
		EnforcedPaymentID:       dto.EnforcedPaymentID,
		ReturnEnforcedPaymentID: dto.ReturnEnforcedPaymentID,
		Type:                    dto.Type,
		Date:                    dto.Date,
	}
}

func ToAccountingEntryItemResponseDTO(data data.AccountingEntryItem) AccountingEntryItemResponseDTO {
	return AccountingEntryItemResponseDTO{
		ID:                      data.ID,
		Title:                   data.Title,
		EntryID:                 data.EntryID,
		SupplierID:              data.SupplierID,
		AccountID:               data.AccountID,
		CreditAmount:            data.CreditAmount,
		DebitAmount:             data.DebitAmount,
		InvoiceID:               data.InvoiceID,
		SalaryID:                data.SalaryID,
		EnforcedPaymentID:       data.EnforcedPaymentID,
		ReturnEnforcedPaymentID: data.ReturnEnforcedPaymentID,
		Type:                    data.Type,
		Date:                    data.Date,
		CreatedAt:               data.CreatedAt,
		UpdatedAt:               data.UpdatedAt,
	}
}

func ToAccountingEntryItemListResponseDTO(accounting_entry_items []*data.AccountingEntryItem) []AccountingEntryItemResponseDTO {
	dtoList := make([]AccountingEntryItemResponseDTO, len(accounting_entry_items))
	for i, x := range accounting_entry_items {
		dtoList[i] = ToAccountingEntryItemResponseDTO(*x)
	}
	return dtoList
}
