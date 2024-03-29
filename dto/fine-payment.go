package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FinePaymentDTO struct {
	FineID                 int                     `json:"fine_id" validate:"required"`
	PaymentMethod          data.FinePaymentMethod  `json:"payment_method" validate:"required,oneof=1 2 3"`
	Amount                 float64                 `json:"amount" validate:"required"`
	PaymentDate            time.Time               `json:"payment_date"`
	PaymentDueDate         time.Time               `json:"payment_due_date"`
	ReceiptNumber          string                  `json:"receipt_number"`
	PaymentReferenceNumber string                  `json:"payment_reference_number"`
	DebitReferenceNumber   string                  `json:"debit_reference_number"`
	Status                 *data.FinePaymentStatus `json:"status"`
}

type FinePaymentResponseDTO struct {
	ID                     int                    `json:"id"`
	FineID                 int                    `json:"fine_id"`
	PaymentMethod          data.FinePaymentMethod `json:"payment_method"`
	Amount                 float64                `json:"amount"`
	PaymentDate            time.Time              `json:"payment_date"`
	PaymentDueDate         time.Time              `json:"payment_due_date"`
	ReceiptNumber          string                 `json:"receipt_number"`
	PaymentReferenceNumber string                 `json:"payment_reference_number"`
	DebitReferenceNumber   string                 `json:"debit_reference_number"`
	Status                 data.FinePaymentStatus `json:"status"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type FinePaymentFilterDTO struct {
	FineID int `json:"fine_id"`
}

// ToFinePayment converts FinePaymentDTO to FinePayment
func (dto FinePaymentDTO) ToFinePayment() *data.FinePayment {
	return &data.FinePayment{
		FineID:                 dto.FineID,
		PaymentMethod:          dto.PaymentMethod,
		Amount:                 dto.Amount,
		PaymentDate:            dto.PaymentDate,
		PaymentDueDate:         dto.PaymentDueDate,
		ReceiptNumber:          dto.ReceiptNumber,
		PaymentReferenceNumber: dto.PaymentReferenceNumber,
		DebitReferenceNumber:   dto.DebitReferenceNumber,
		Status:                 *dto.Status,
	}
}

// ToFinePaymentResponseDTO converts FinePayment to FinePaymentResponseDTO
func ToFinePaymentResponseDTO(data data.FinePayment) FinePaymentResponseDTO {
	return FinePaymentResponseDTO{
		ID:                     data.ID,
		FineID:                 data.FineID,
		PaymentMethod:          data.PaymentMethod,
		Amount:                 data.Amount,
		PaymentDate:            data.PaymentDate,
		PaymentDueDate:         data.PaymentDueDate,
		ReceiptNumber:          data.ReceiptNumber,
		PaymentReferenceNumber: data.PaymentReferenceNumber,
		DebitReferenceNumber:   data.DebitReferenceNumber,
		Status:                 data.Status,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
	}
}

// ToFinePaymentListResponseDTO converts  []*data.FinePayment to []FinePaymentResponseDTO
func ToFinePaymentListResponseDTO(fines []*data.FinePayment) []FinePaymentResponseDTO {
	dtoList := make([]FinePaymentResponseDTO, len(fines))
	for i, x := range fines {
		dtoList[i] = ToFinePaymentResponseDTO(*x)
	}
	return dtoList
}
