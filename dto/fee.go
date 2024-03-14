package dto

import (
	"time"

	"github.com/lib/pq"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type FeeDTO struct {
	FeeTypeID              data.FeeType     `json:"fee_type" validate:"required,oneof=1 2"`
	FeeSubcategoryID       data.FeeCategory `json:"fee_subcategory" validate:"required"`
	DecisionNumber         string           `json:"decision_number" validate:"required"`
	DecisionDate           time.Time        `json:"decision_date"`
	Subject                string           `json:"subject"`
	JMBG                   string           `json:"jmbg" validate:"required"`
	Residence              string           `json:"residence"`
	Amount                 float64          `json:"amount"`
	PaymentReferenceNumber string           `json:"payment_reference_number"`
	DebitReferenceNumber   string           `json:"debit_reference_number"`
	ExecutionDate          time.Time        `json:"execution_date"`
	PaymentDeadlineDate    time.Time        `json:"payment_deadline_date"`
	Description            string           `json:"description"`
	Status                 data.FeeStatus   `json:"status" validate:"required,oneof=1 2 3"`
	CourtAccountID         *int             `json:"court_account"`
	File                   pq.Int64Array    `json:"file"`
}

type FeeResponseDTO struct {
	ID                     int              `json:"id"`
	FeeTypeID              data.FeeType     `json:"fee_type_id"`
	FeeSubcategoryID       data.FeeCategory `json:"fee_subcategory_id"`
	DecisionNumber         string           `json:"decision_number"`
	DecisionDate           time.Time        `json:"decision_date"`
	Subject                string           `json:"subject"`
	JMBG                   string           `json:"jmbg"`
	Residence              string           `json:"residence"`
	Amount                 float64          `json:"amount"`
	PaymentReferenceNumber string           `json:"payment_reference_number"`
	DebitReferenceNumber   string           `json:"debit_reference_number"`
	ExecutionDate          time.Time        `json:"execution_date"`
	PaymentDeadlineDate    time.Time        `json:"payment_deadline_date"`
	Description            string           `json:"description"`
	Status                 data.FeeStatus   `json:"status"`
	CourtAccountID         *int             `json:"court_account"`
	FeeDetails             *FeeDetailsDTO   `json:"fee_details"`
	File                   []int            `json:"file"`
	CreatedAt              time.Time        `json:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at"`
}

type FeeDetailsDTO struct {
	FeeLeftToPayAmount  float64 `json:"fee_left_to_pay_amount"`
	FeeAllPaymentAmount float64 `json:"fee_all_payment_amount"`
}

type FeeFilterDTO struct {
	Page                     *int    `json:"page"`
	Size                     *int    `json:"size"`
	FilterByFeeSubcategoryID *int    `json:"fee_subcategory_id"`
	FilterByFeeTypeID        *int    `json:"fee_type_id"`
	Search                   *string `json:"search"`
}

// ToFee converts FeeDTO to Fee
func (dto FeeDTO) ToFee() *data.Fee {
	return &data.Fee{
		FeeTypeID:              dto.FeeTypeID,
		FeeSubcategoryID:       dto.FeeSubcategoryID,
		DecisionNumber:         dto.DecisionNumber,
		DecisionDate:           dto.DecisionDate,
		Subject:                dto.Subject,
		JMBG:                   dto.JMBG,
		Amount:                 dto.Amount,
		PaymentReferenceNumber: dto.PaymentReferenceNumber,
		DebitReferenceNumber:   dto.DebitReferenceNumber,
		ExecutionDate:          dto.ExecutionDate,
		PaymentDeadlineDate:    dto.PaymentDeadlineDate,
		Description:            dto.Description,
		Status:                 dto.Status,
		CourtAccountID:         dto.CourtAccountID,
		File:                   dto.File,
	}
}

// ToFeeResponseDTO converts Fee to FeeResponseDTO
func ToFeeResponseDTO(data data.Fee) FeeResponseDTO {
	filesArray := make([]int, len(data.File))
	for i, id := range data.File {
		filesArray[i] = int(id)
	}
	return FeeResponseDTO{
		ID:                     data.ID,
		FeeTypeID:              data.FeeTypeID,
		FeeSubcategoryID:       data.FeeSubcategoryID,
		DecisionNumber:         data.DecisionNumber,
		DecisionDate:           data.DecisionDate,
		Subject:                data.Subject,
		JMBG:                   data.JMBG,
		Amount:                 data.Amount,
		PaymentReferenceNumber: data.PaymentReferenceNumber,
		DebitReferenceNumber:   data.DebitReferenceNumber,
		ExecutionDate:          data.ExecutionDate,
		PaymentDeadlineDate:    data.PaymentDeadlineDate,
		Description:            data.Description,
		Status:                 data.Status,
		CourtAccountID:         data.CourtAccountID,
		File:                   filesArray,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
	}
}

// ToFeeListResponseDTO converts []*Fee to []FeeResponseDTO
func ToFeeListResponseDTO(fees []*data.Fee) []FeeResponseDTO {
	dtoList := make([]FeeResponseDTO, len(fees))
	for i, x := range fees {
		dtoList[i] = ToFeeResponseDTO(*x)
	}
	return dtoList
}
