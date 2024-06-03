package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ProcedureCostPaymentDTO struct {
	ProcedureCostID        int                              `json:"procedure_cost_id" validate:"required"`
	PaymentMethod          data.ProcedureCostPaymentMethod  `json:"payment_method" validate:"required,oneof=1 2 3"`
	Amount                 float64                          `json:"amount" validate:"required"`
	PaymentDate            time.Time                        `json:"payment_date"`
	PaymentDueDate         time.Time                        `json:"payment_due_date"`
	ReceiptNumber          string                           `json:"receipt_number"`
	PaymentReferenceNumber string                           `json:"payment_reference_number"`
	DebitReferenceNumber   string                           `json:"debit_reference_number"`
	Status                 *data.ProcedureCostPaymentStatus `json:"status"`
}

type ProcedureCostPaymentResponseDTO struct {
	ID                     int                             `json:"id"`
	ProcedureCostID        int                             `json:"procedure_cost_id"`
	PaymentMethod          data.ProcedureCostPaymentMethod `json:"payment_method"`
	Amount                 float64                         `json:"amount"`
	PaymentDate            time.Time                       `json:"payment_date"`
	PaymentDueDate         time.Time                       `json:"payment_due_date"`
	ReceiptNumber          string                          `json:"receipt_number"`
	PaymentReferenceNumber string                          `json:"payment_reference_number"`
	DebitReferenceNumber   string                          `json:"debit_reference_number"`
	Status                 data.ProcedureCostPaymentStatus `json:"status"`
	CreatedAt              time.Time                       `json:"created_at"`
	UpdatedAt              time.Time                       `json:"updated_at"`
}

type ProcedureCostPaymentFilterDTO struct {
	ProcedureCostID int `json:"procedure_cost_id"`
}

// ToProcedureCostPayment converts ProcedureCostPaymentDTO to ProcedureCostPayment
func (dto ProcedureCostPaymentDTO) ToProcedureCostPayment() *data.ProcedureCostPayment {
	return &data.ProcedureCostPayment{
		ProcedureCostID:        dto.ProcedureCostID,
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

// ToProcedureCostPaymentResponseDTO converts ProcedureCostPayment to ProcedureCostPaymentResponseDTO
func ToProcedureCostPaymentResponseDTO(data data.ProcedureCostPayment) ProcedureCostPaymentResponseDTO {
	return ProcedureCostPaymentResponseDTO{
		ID:                     data.ID,
		ProcedureCostID:        data.ProcedureCostID,
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

// ToProcedureCostPaymentListResponseDTO converts  []*data.ProcedureCostPayment to []ProcedureCostPaymentResponseDTO
func ToProcedureCostPaymentListResponseDTO(procedurecosts []*data.ProcedureCostPayment) []ProcedureCostPaymentResponseDTO {
	dtoList := make([]ProcedureCostPaymentResponseDTO, len(procedurecosts))
	for i, x := range procedurecosts {
		dtoList[i] = ToProcedureCostPaymentResponseDTO(*x)
	}
	return dtoList
}
