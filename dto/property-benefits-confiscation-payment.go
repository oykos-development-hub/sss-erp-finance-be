package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type PropBenConfPaymentDTO struct {
	PropBenConfID          int                           `json:"property_benefits_confiscation_id" validate:"required"`
	PaymentMethod          data.PropBenConfPaymentMethod `json:"payment_method" validate:"required,oneof=1 2 3"`
	Amount                 float64                       `json:"amount" validate:"required"`
	PaymentDate            time.Time                     `json:"payment_date"`
	PaymentDueDate         time.Time                     `json:"payment_due_date"`
	ReceiptNumber          string                        `json:"receipt_number"`
	PaymentReferenceNumber string                        `json:"payment_reference_number"`
	DebitReferenceNumber   string                        `json:"debit_reference_number"`
	Status                 data.PropBenConfPaymentStatus `json:"status"`
}

type PropBenConfPaymentResponseDTO struct {
	ID                     int                           `json:"id"`
	PropBenConfID          int                           `json:"property_benefits_confiscation_id"`
	PaymentMethod          data.PropBenConfPaymentMethod `json:"payment_method"`
	Amount                 float64                       `json:"amount"`
	PaymentDate            time.Time                     `json:"payment_date"`
	PaymentDueDate         time.Time                     `json:"payment_due_date"`
	ReceiptNumber          string                        `json:"receipt_number"`
	PaymentReferenceNumber string                        `json:"payment_reference_number"`
	DebitReferenceNumber   string                        `json:"debit_reference_number"`
	Status                 data.PropBenConfPaymentStatus `json:"status"`
	CreatedAt              time.Time                     `json:"created_at"`
	UpdatedAt              time.Time                     `json:"updated_at"`
}

type PropBenConfPaymentFilterDTO struct {
	PropBenConfID int `json:"property_benefits_confiscation_id"`
}

// ToPropBenConfPayment converts PropBenConfPaymentDTO to PropBenConfPayment
func (dto PropBenConfPaymentDTO) ToPropBenConfPayment() *data.PropBenConfPayment {
	return &data.PropBenConfPayment{
		PropBenConfID:          dto.PropBenConfID,
		PaymentMethod:          dto.PaymentMethod,
		Amount:                 dto.Amount,
		PaymentDate:            dto.PaymentDate,
		PaymentDueDate:         dto.PaymentDueDate,
		ReceiptNumber:          dto.ReceiptNumber,
		PaymentReferenceNumber: dto.PaymentReferenceNumber,
		DebitReferenceNumber:   dto.DebitReferenceNumber,
		Status:                 dto.Status,
	}
}

// ToPropBenConfPaymentResponseDTO converts PropBenConfPayment to PropBenConfPaymentResponseDTO
func ToPropBenConfPaymentResponseDTO(data data.PropBenConfPayment) PropBenConfPaymentResponseDTO {
	return PropBenConfPaymentResponseDTO{
		ID:                     data.ID,
		PropBenConfID:          data.PropBenConfID,
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

// ToPropBenConfPaymentListResponseDTO converts  []*data.PropBenConfPayment to []PropBenConfPaymentResponseDTO
func ToPropBenConfPaymentListResponseDTO(propbenconfs []*data.PropBenConfPayment) []PropBenConfPaymentResponseDTO {
	dtoList := make([]PropBenConfPaymentResponseDTO, len(propbenconfs))
	for i, x := range propbenconfs {
		dtoList[i] = ToPropBenConfPaymentResponseDTO(*x)
	}
	return dtoList
}
