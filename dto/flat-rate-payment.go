package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FlatRatePaymentDTO struct {
	FlatRateID             int                         `json:"flat_rate_id" validate:"required"`
	PaymentMethod          data.FlatRatePaymentMethod  `json:"payment_method" validate:"required,oneof=1 2 3"`
	Amount                 float64                     `json:"amount" validate:"required"`
	PaymentDate            time.Time                   `json:"payment_date"`
	PaymentDueDate         time.Time                   `json:"payment_due_date"`
	ReceiptNumber          string                      `json:"receipt_number"`
	PaymentReferenceNumber string                      `json:"payment_reference_number"`
	DebitReferenceNumber   string                      `json:"debit_reference_number"`
	Status                 *data.FlatRatePaymentStatus `json:"status"`
}

type FlatRatePaymentResponseDTO struct {
	ID                     int                        `json:"id"`
	FlatRateID             int                        `json:"flat_rate_id"`
	PaymentMethod          data.FlatRatePaymentMethod `json:"payment_method"`
	Amount                 float64                    `json:"amount"`
	PaymentDate            time.Time                  `json:"payment_date"`
	PaymentDueDate         time.Time                  `json:"payment_due_date"`
	ReceiptNumber          string                     `json:"receipt_number"`
	PaymentReferenceNumber string                     `json:"payment_reference_number"`
	DebitReferenceNumber   string                     `json:"debit_reference_number"`
	Status                 data.FlatRatePaymentStatus `json:"status"`
	CreatedAt              time.Time                  `json:"created_at"`
	UpdatedAt              time.Time                  `json:"updated_at"`
}

type FlatRatePaymentFilterDTO struct {
	FlatRateID int `json:"flat_rate_id"`
}

// ToFlatRatePayment converts FlatRatePaymentDTO to FlatRatePayment
func (dto FlatRatePaymentDTO) ToFlatRatePayment() *data.FlatRatePayment {
	return &data.FlatRatePayment{
		FlatRateID:             dto.FlatRateID,
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

// ToFlatRatePaymentResponseDTO converts FlatRatePayment to FlatRatePaymentResponseDTO
func ToFlatRatePaymentResponseDTO(data data.FlatRatePayment) FlatRatePaymentResponseDTO {
	return FlatRatePaymentResponseDTO{
		ID:                     data.ID,
		FlatRateID:             data.FlatRateID,
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

// ToFlatRatePaymentListResponseDTO converts  []*data.FlatRatePayment to []FlatRatePaymentResponseDTO
func ToFlatRatePaymentListResponseDTO(flatrates []*data.FlatRatePayment) []FlatRatePaymentResponseDTO {
	dtoList := make([]FlatRatePaymentResponseDTO, len(flatrates))
	for i, x := range flatrates {
		dtoList[i] = ToFlatRatePaymentResponseDTO(*x)
	}
	return dtoList
}
