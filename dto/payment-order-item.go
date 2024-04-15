package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type PaymentOrderItemDTO struct {
	PaymentOrderID            int  `json:"payment_order_id"`
	InvoiceID                 *int `json:"invoice_id"`
	AdditionalExpenseID       *int `json:"additional_expense"`
	SalaryAdditionalExpenseID *int `json:"salary_additional_expense"`
	AccountID                 int  `json:"account_id"`
}

type PaymentOrderItemResponseDTO struct {
	ID                        int       `json:"id"`
	PaymentOrderID            int       `json:"payment_order_id"`
	InvoiceID                 *int      `json:"invoice_id"`
	AdditionalExpenseID       *int      `json:"additional_expense_id"`
	SalaryAdditionalExpenseID *int      `json:"salary_additional_expense_id"`
	Type                      string    `json:"type"`
	AccountID                 int       `json:"account_id"`
	Amount                    float64   `json:"amount"`
	Title                     string    `json:"title"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type PaymentOrderItemFilterDTO struct {
	Page                      *int    `json:"page"`
	Size                      *int    `json:"size"`
	SortByTitle               *string `json:"sort_by_title"`
	PaymentOrderID            *int    `json:"payment_order_id"`
	InvoiceID                 *int    `json:"invoice_id"`
	AdditionalExpenseID       *int    `json:"additional_expense"`
	SalaryAdditionalExpenseID *int    `json:"salary_additional_expense"`
}

func (dto PaymentOrderItemDTO) ToPaymentOrderItem() *data.PaymentOrderItem {
	return &data.PaymentOrderItem{
		PaymentOrderID:            dto.PaymentOrderID,
		InvoiceID:                 dto.InvoiceID,
		AccountID:                 dto.AccountID,
		AdditionalExpenseID:       dto.AdditionalExpenseID,
		SalaryAdditionalExpenseID: dto.SalaryAdditionalExpenseID,
	}
}

func ToPaymentOrderItemResponseDTO(data data.PaymentOrderItem) PaymentOrderItemResponseDTO {
	return PaymentOrderItemResponseDTO{
		ID:                        data.ID,
		PaymentOrderID:            data.PaymentOrderID,
		InvoiceID:                 data.InvoiceID,
		AdditionalExpenseID:       data.AdditionalExpenseID,
		SalaryAdditionalExpenseID: data.SalaryAdditionalExpenseID,
		AccountID:                 data.AccountID,
		CreatedAt:                 data.CreatedAt,
		UpdatedAt:                 data.UpdatedAt,
	}
}

func ToPaymentOrderItemListResponseDTO(payment_order_items []*data.PaymentOrderItem) []PaymentOrderItemResponseDTO {
	dtoList := make([]PaymentOrderItemResponseDTO, len(payment_order_items))
	for i, x := range payment_order_items {
		dtoList[i] = ToPaymentOrderItemResponseDTO(*x)
	}
	return dtoList
}
