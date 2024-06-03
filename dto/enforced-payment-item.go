package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type EnforcedPaymentItemDTO struct {
	PaymentOrderID int     `json:"payment_order_id"`
	Title          string  `json:"title"`
	Amount         float32 `json:"amount"`
	InvoiceID      *int    `json:"invoice_id"`
	AccountID      int     `json:"account_id"`
}

type EnforcedPaymentItemResponseDTO struct {
	ID             int       `json:"id"`
	PaymentOrderID int       `json:"payment_order_id"`
	Title          string    `json:"title"`
	Amount         float32   `json:"amount"`
	InvoiceID      *int      `json:"invoice_id"`
	AccountID      int       `json:"account_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type EnforcedPaymentItemFilterDTO struct {
	Page           *int    `json:"page"`
	Size           *int    `json:"size"`
	SortByTitle    *string `json:"sort_by_title"`
	PaymentOrderID *int    `json:"payment_order_id"`
}

func (dto EnforcedPaymentItemDTO) ToEnforcedPaymentItem() *data.EnforcedPaymentItem {
	return &data.EnforcedPaymentItem{
		PaymentOrderID: dto.PaymentOrderID,
		Title:          dto.Title,
		Amount:         dto.Amount,
		InvoiceID:      dto.InvoiceID,
		AccountID:      dto.AccountID,
	}
}

func ToEnforcedPaymentItemResponseDTO(data data.EnforcedPaymentItem) EnforcedPaymentItemResponseDTO {
	return EnforcedPaymentItemResponseDTO{
		ID:             data.ID,
		PaymentOrderID: data.PaymentOrderID,
		Title:          data.Title,
		Amount:         data.Amount,
		InvoiceID:      data.InvoiceID,
		AccountID:      data.AccountID,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
}

func ToEnforcedPaymentItemListResponseDTO(enforced_payment_items []*data.EnforcedPaymentItem) []EnforcedPaymentItemResponseDTO {
	dtoList := make([]EnforcedPaymentItemResponseDTO, len(enforced_payment_items))
	for i, x := range enforced_payment_items {
		dtoList[i] = ToEnforcedPaymentItemResponseDTO(*x)
	}
	return dtoList
}
