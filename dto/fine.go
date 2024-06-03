package dto

import (
	"time"

	"github.com/lib/pq"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type FineDTO struct {
	ActType                data.FineActType `json:"act_type" validate:"required,oneof=1 2"`
	DecisionNumber         string           `json:"decision_number" validate:"required"`
	DecisionDate           time.Time        `json:"decision_date"`
	Subject                string           `json:"subject"`
	JMBG                   string           `json:"jmbg" validate:"required"`
	Residence              string           `json:"residence"`
	Amount                 float64          `json:"amount"`
	PaymentReferenceNumber string           `json:"payment_reference_number"`
	DebitReferenceNumber   string           `json:"debit_reference_number"`
	AccountID              int              `json:"account_id"`
	ExecutionDate          time.Time        `json:"execution_date"`
	PaymentDeadlineDate    time.Time        `json:"payment_deadline_date"`
	Description            string           `json:"description"`
	Status                 *data.FineStatus `json:"status"`
	CourtCosts             *float64         `json:"court_costs"`
	CourtAccountID         *int             `json:"court_account_id"`
	File                   pq.Int64Array    `json:"file"`
}

type FineResponseDTO struct {
	ID                     int                `json:"id"`
	ActType                data.FineActType   `json:"act_type"`
	DecisionNumber         string             `json:"decision_number"`
	DecisionDate           time.Time          `json:"decision_date"`
	Subject                string             `json:"subject"`
	JMBG                   string             `json:"jmbg"`
	Residence              string             `json:"residence"`
	Amount                 float64            `json:"amount"`
	PaymentReferenceNumber string             `json:"payment_reference_number"`
	DebitReferenceNumber   string             `json:"debit_reference_number"`
	AccountID              int                `json:"account_id"`
	ExecutionDate          time.Time          `json:"execution_date"`
	PaymentDeadlineDate    time.Time          `json:"payment_deadline_date"`
	Description            string             `json:"description"`
	Status                 data.FineStatus    `json:"status"`
	CourtCosts             *float64           `json:"court_costs"`
	CourtAccountID         *int               `json:"court_account_id"`
	FineFeeDetailsDTO      *FineFeeDetailsDTO `json:"fine_fee_details"`
	File                   []int              `json:"file"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
}

type FineFeeDetailsDTO struct {
	FeeAllPaymentAmount           float64   `json:"fee_all_payments_amount"`
	FeeAmountGracePeriod          float64   `json:"fee_amount_grace_period"`
	FeeAmountGracePeriodDueDate   time.Time `json:"fee_amount_grace_period_due_date"`
	FeeAmountGracePeriodAvailable bool      `json:"fee_amount_grace_period_available"`
	FeeLeftToPayAmount            float64   `json:"fee_left_to_pay_amount"`

	FeeCourtCostsPaid            float64 `json:"fee_court_costs_paid"`
	FeeCourtCostsLeftToPayAmount float64 `json:"fee_court_costs_left_to_pay_amount"`
}

type FineFilterDTO struct {
	Page              *int    `json:"page"`
	Size              *int    `json:"size"`
	Subject           *string `json:"subject"`
	FilterByActTypeID *int    `json:"act_type_id"`
	Search            *string `json:"search"`
}

// ToFine converts FineDTO to Fine
func (dto FineDTO) ToFine() *data.Fine {
	return &data.Fine{
		ActType:                dto.ActType,
		DecisionNumber:         dto.DecisionNumber,
		DecisionDate:           dto.DecisionDate,
		Subject:                dto.Subject,
		JMBG:                   dto.JMBG,
		Residence:              dto.Residence,
		Amount:                 dto.Amount,
		PaymentReferenceNumber: dto.PaymentReferenceNumber,
		DebitReferenceNumber:   dto.DebitReferenceNumber,
		AccountID:              dto.AccountID,
		ExecutionDate:          dto.ExecutionDate,
		PaymentDeadlineDate:    dto.PaymentDeadlineDate,
		Description:            dto.Description,
		Status:                 *dto.Status,
		CourtCosts:             dto.CourtCosts,
		CourtAccountID:         dto.CourtAccountID,
		File:                   dto.File,
	}
}

// ToFineResponseDTO converts Fine to FineResponseDTO
func ToFineResponseDTO(data data.Fine) FineResponseDTO {
	filesArray := make([]int, len(data.File))
	for i, id := range data.File {
		filesArray[i] = int(id)
	}
	return FineResponseDTO{
		ID:                     data.ID,
		ActType:                data.ActType,
		DecisionNumber:         data.DecisionNumber,
		DecisionDate:           data.DecisionDate,
		Subject:                data.Subject,
		JMBG:                   data.JMBG,
		Residence:              data.Residence,
		Amount:                 data.Amount,
		PaymentReferenceNumber: data.PaymentReferenceNumber,
		DebitReferenceNumber:   data.DebitReferenceNumber,
		AccountID:              data.AccountID,
		ExecutionDate:          data.ExecutionDate,
		PaymentDeadlineDate:    data.PaymentDeadlineDate,
		Description:            data.Description,
		Status:                 data.Status,
		CourtCosts:             data.CourtCosts,
		CourtAccountID:         data.CourtAccountID,
		File:                   filesArray,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
	}
}

// ToFineListResponseDTO converts []*Fine to []FineResponseDTO
func ToFineListResponseDTO(fines []*data.Fine) []FineResponseDTO {
	dtoList := make([]FineResponseDTO, len(fines))
	for i, x := range fines {
		dtoList[i] = ToFineResponseDTO(*x)
	}
	return dtoList
}
