package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type InvoicesFilter struct {
	Search             *string `json:"search"`
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	Year               *string `json:"year"`
	Status             *string `json:"status"`
	SupplierID         *int    `json:"supplier_id"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
}

type InvoiceDTO struct {
	InvoiceNumber         string    `json:"invoice_number"`
	Status                string    `json:"status"`
	GrossPrice            float64   `json:"gross_price"`
	VATPrice              float64   `json:"vat_price"`
	SupplierID            int       `json:"supplier_id"`
	OrderID               int       `json:"order_id"`
	OrganizationUnitID    int       `json:"organization_unit_id"`
	DateOfInvoice         time.Time `json:"date_of_invoice"`
	ReceiptDate           time.Time `json:"receipt_date"`
	DateOfPayment         time.Time `json:"date_of_payment"`
	SSSInvoiceReceiptDate time.Time `json:"sss_invoice_receipt_date"`
	FileID                int       `json:"file_id"`
	BankAccount           string    `json:"bank_account"`
	Description           string    `json:"description"`
}

type InvoiceResponseDTO struct {
	ID                    int                  `json:"id"`
	InvoiceNumber         string               `json:"invoice_number"`
	Status                string               `json:"status"`
	GrossPrice            float64              `json:"gross_price"`
	VATPrice              float64              `json:"vat_price"`
	SupplierID            int                  `json:"supplier_id"`
	OrderID               int                  `json:"order_id"`
	OrganizationUnitID    int                  `json:"organization_unit_id"`
	DateOfInvoice         time.Time            `json:"date_of_invoice"`
	ReceiptDate           time.Time            `json:"receipt_date"`
	DateOfPayment         time.Time            `json:"date_of_payment"`
	SSSInvoiceReceiptDate time.Time            `json:"sss_invoice_receipt_date"`
	FileID                int                  `json:"file_id"`
	BankAccount           string               `json:"bank_account"`
	Description           string               `json:"description"`
	Articles              []ArticleResponseDTO `json:"articles"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

func (dto InvoiceDTO) ToInvoice() *data.Invoice {
	return &data.Invoice{
		InvoiceNumber:         dto.InvoiceNumber,
		Status:                dto.Status,
		GrossPrice:            dto.GrossPrice,
		VATPrice:              dto.VATPrice,
		SupplierID:            dto.SupplierID,
		OrderID:               dto.OrderID,
		OrganizationUnitID:    dto.OrganizationUnitID,
		DateOfInvoice:         dto.DateOfInvoice,
		ReceiptDate:           dto.ReceiptDate,
		DateOfPayment:         dto.DateOfPayment,
		SSSInvoiceReceiptDate: dto.SSSInvoiceReceiptDate,
		FileID:                dto.FileID,
		BankAccount:           dto.BankAccount,
		Description:           dto.Description,
	}
}

func ToInvoiceResponseDTO(data data.Invoice) InvoiceResponseDTO {
	return InvoiceResponseDTO{
		ID:                    data.ID,
		InvoiceNumber:         data.InvoiceNumber,
		Status:                data.Status,
		GrossPrice:            data.GrossPrice,
		VATPrice:              data.VATPrice,
		SupplierID:            data.SupplierID,
		OrderID:               data.OrderID,
		OrganizationUnitID:    data.OrganizationUnitID,
		DateOfInvoice:         data.DateOfInvoice,
		ReceiptDate:           data.ReceiptDate,
		DateOfPayment:         data.DateOfPayment,
		SSSInvoiceReceiptDate: data.SSSInvoiceReceiptDate,
		FileID:                data.FileID,
		BankAccount:           data.BankAccount,
		Description:           data.Description,
		CreatedAt:             data.CreatedAt,
		UpdatedAt:             data.UpdatedAt,
	}
}

func ToInvoiceListResponseDTO(invoices []*data.Invoice) []InvoiceResponseDTO {
	dtoList := make([]InvoiceResponseDTO, len(invoices))
	for i, x := range invoices {
		dtoList[i] = ToInvoiceResponseDTO(*x)
	}
	return dtoList
}
