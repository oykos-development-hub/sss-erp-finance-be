package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FeePaymentDTO struct {
	FeeID                  int                    `json:"fee_id" validate:"required"`
	PaymentMethod          data.FeePaymentMethod  `json:"payment_method" validate:"required,oneof=1 2 3"`
	Amount                 float64                `json:"amount" validate:"required"`
	PaymentDate            time.Time              `json:"payment_date"`
	PaymentDueDate         time.Time              `json:"payment_due_date"`
	ReceiptNumber          string                 `json:"receipt_number"`
	PaymentReferenceNumber string                 `json:"payment_reference_number"`
	DebitReferenceNumber   string                 `json:"debit_reference_number"`
	Status                 *data.FeePaymentStatus `json:"status"`
}

type FeePaymentResponseDTO struct {
	ID                     int                   `json:"id"`
	FeeID                  int                   `json:"fee_id"`
	PaymentMethod          data.FeePaymentMethod `json:"payment_method"`
	Amount                 float64               `json:"amount"`
	PaymentDate            time.Time             `json:"payment_date"`
	PaymentDueDate         time.Time             `json:"payment_due_date"`
	ReceiptNumber          string                `json:"receipt_number"`
	PaymentReferenceNumber string                `json:"payment_reference_number"`
	DebitReferenceNumber   string                `json:"debit_reference_number"`
	Status                 data.FeePaymentStatus `json:"status"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
}

type FeePaymentFilterDTO struct {
	FeeID int `json:"fee_id"`
}

// ToFeePayment converts FeePaymentDTO to FeePayment
func (dto FeePaymentDTO) ToFeePayment() *data.FeePayment {
	return &data.FeePayment{
		FeeID:                  dto.FeeID,
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

// ToFeePaymentResponseDTO converts FeePayment to FeePaymentResponseDTO
func ToFeePaymentResponseDTO(data data.FeePayment) FeePaymentResponseDTO {
	return FeePaymentResponseDTO{
		ID:                     data.ID,
		FeeID:                  data.FeeID,
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

// ToFeePaymentListResponseDTO converts  []*data.FeePayment to []FeePaymentResponseDTO
func ToFeePaymentListResponseDTO(fees []*data.FeePayment) []FeePaymentResponseDTO {
	dtoList := make([]FeePaymentResponseDTO, len(fees))
	for i, x := range fees {
		dtoList[i] = ToFeePaymentResponseDTO(*x)
	}
	return dtoList
}
