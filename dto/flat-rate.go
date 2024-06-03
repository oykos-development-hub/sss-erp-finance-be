package dto

import (
	"time"

	"github.com/lib/pq"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type FlatRateDTO struct {
	FlatRateType           data.FlatRateType    `json:"flat_rate_type" validate:"required,oneof=1 2"`
	DecisionNumber         string               `json:"decision_number" validate:"required"`
	DecisionDate           time.Time            `json:"decision_date"`
	Subject                string               `json:"subject"`
	JMBG                   string               `json:"jmbg" validate:"required"`
	Residence              string               `json:"residence"`
	Amount                 float64              `json:"amount"`
	PaymentReferenceNumber string               `json:"payment_reference_number"`
	DebitReferenceNumber   string               `json:"debit_reference_number"`
	AccountID              int                  `json:"account_id"`
	ExecutionDate          time.Time            `json:"execution_date"`
	PaymentDeadlineDate    time.Time            `json:"payment_deadline_date"`
	Description            string               `json:"description"`
	Status                 *data.FlatRateStatus `json:"status"`
	CourtCosts             *float64             `json:"court_costs"`
	CourtAccountID         *int                 `json:"court_account_id"`
	File                   pq.Int64Array        `json:"file"`
}

type FlatRateResponseDTO struct {
	ID                     int                 `json:"id"`
	FlatRateType           data.FlatRateType   `json:"flat_rate_type"`
	DecisionNumber         string              `json:"decision_number"`
	DecisionDate           time.Time           `json:"decision_date"`
	Subject                string              `json:"subject"`
	JMBG                   string              `json:"jmbg"`
	Residence              string              `json:"residence"`
	Amount                 float64             `json:"amount"`
	PaymentReferenceNumber string              `json:"payment_reference_number"`
	DebitReferenceNumber   string              `json:"debit_reference_number"`
	AccountID              int                 `json:"account_id"`
	ExecutionDate          time.Time           `json:"execution_date"`
	PaymentDeadlineDate    time.Time           `json:"payment_deadline_date"`
	Description            string              `json:"description"`
	Status                 data.FlatRateStatus `json:"status"`
	CourtCosts             *float64            `json:"court_costs"`
	CourtAccountID         *int                `json:"court_account_id"`
	FlatRateDetails        *FlatRateDetailsDTO `json:"flat_rate_details"`
	File                   []int               `json:"file"`
	CreatedAt              time.Time           `json:"created_at"`
	UpdatedAt              time.Time           `json:"updated_at"`
}

type FlatRateDetailsDTO struct {
	AllPaymentAmount           float64   `json:"all_payments_amount"`
	AmountGracePeriod          float64   `json:"amount_grace_period"`
	AmountGracePeriodDueDate   time.Time `json:"amount_grace_period_due_date"`
	AmountGracePeriodAvailable bool      `json:"amount_grace_period_available"`
	LeftToPayAmount            float64   `json:"left_to_pay_amount"`

	CourtCostsPaid            float64 `json:"court_costs_paid"`
	CourtCostsLeftToPayAmount float64 `json:"court_costs_left_to_pay_amount"`
}

type FlatRateFilterDTO struct {
	Page           *int    `json:"page"`
	Size           *int    `json:"size"`
	Subject        *string `json:"subject"`
	FilterByTypeID *int    `json:"flat_rate_type_id"`
	Search         *string `json:"search"`
}

// ToFlatRate converts FlatRateDTO to FlatRate
func (dto FlatRateDTO) ToFlatRate() *data.FlatRate {
	return &data.FlatRate{
		FlatRateType:           dto.FlatRateType,
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

// ToFlatRateResponseDTO converts FlatRate to FlatRateResponseDTO
func ToFlatRateResponseDTO(data data.FlatRate) FlatRateResponseDTO {
	filesArray := make([]int, len(data.File))
	for i, id := range data.File {
		filesArray[i] = int(id)
	}
	return FlatRateResponseDTO{
		ID:                     data.ID,
		FlatRateType:           data.FlatRateType,
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

// ToFlatRateListResponseDTO converts []*FlatRate to []FlatRateResponseDTO
func ToFlatRateListResponseDTO(flatrates []*data.FlatRate) []FlatRateResponseDTO {
	dtoList := make([]FlatRateResponseDTO, len(flatrates))
	for i, x := range flatrates {
		dtoList[i] = ToFlatRateResponseDTO(*x)
	}
	return dtoList
}
