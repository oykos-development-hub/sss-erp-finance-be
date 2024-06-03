package dto

import (
	"time"

	"github.com/lib/pq"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type ProcedureCostDTO struct {
	ProcedureCostType      data.ProcedureCostType    `json:"procedure_cost_type" validate:"required,oneof=1 2"`
	DecisionNumber         string                    `json:"decision_number" validate:"required"`
	DecisionDate           time.Time                 `json:"decision_date"`
	Subject                string                    `json:"subject"`
	JMBG                   string                    `json:"jmbg" validate:"required"`
	Residence              string                    `json:"residence"`
	Amount                 float64                   `json:"amount"`
	PaymentReferenceNumber string                    `json:"payment_reference_number"`
	DebitReferenceNumber   string                    `json:"debit_reference_number"`
	AccountID              int                       `json:"account_id"`
	ExecutionDate          time.Time                 `json:"execution_date"`
	PaymentDeadlineDate    time.Time                 `json:"payment_deadline_date"`
	Description            string                    `json:"description"`
	Status                 *data.ProcedureCostStatus `json:"status"`
	CourtCosts             *float64                  `json:"court_costs"`
	CourtAccountID         *int                      `json:"court_account_id"`
	File                   pq.Int64Array             `json:"file"`
}

type ProcedureCostResponseDTO struct {
	ID                     int                      `json:"id"`
	ProcedureCostType      data.ProcedureCostType   `json:"procedure_cost_type"`
	DecisionNumber         string                   `json:"decision_number"`
	DecisionDate           time.Time                `json:"decision_date"`
	Subject                string                   `json:"subject"`
	JMBG                   string                   `json:"jmbg"`
	Residence              string                   `json:"residence"`
	Amount                 float64                  `json:"amount"`
	PaymentReferenceNumber string                   `json:"payment_reference_number"`
	DebitReferenceNumber   string                   `json:"debit_reference_number"`
	AccountID              int                      `json:"account_id"`
	ExecutionDate          time.Time                `json:"execution_date"`
	PaymentDeadlineDate    time.Time                `json:"payment_deadline_date"`
	Description            string                   `json:"description"`
	Status                 data.ProcedureCostStatus `json:"status"`
	CourtCosts             *float64                 `json:"court_costs"`
	CourtAccountID         *int                     `json:"court_account_id"`
	ProcedureCostDetails   *ProcedureCostDetailsDTO `json:"procedure_cost_details"`
	File                   []int                    `json:"file"`
	CreatedAt              time.Time                `json:"created_at"`
	UpdatedAt              time.Time                `json:"updated_at"`
}

type ProcedureCostDetailsDTO struct {
	AllPaymentAmount           float64   `json:"all_payments_amount"`
	AmountGracePeriod          float64   `json:"amount_grace_period"`
	AmountGracePeriodDueDate   time.Time `json:"amount_grace_period_due_date"`
	AmountGracePeriodAvailable bool      `json:"amount_grace_period_available"`
	LeftToPayAmount            float64   `json:"left_to_pay_amount"`

	CourtCostsPaid            float64 `json:"court_costs_paid"`
	CourtCostsLeftToPayAmount float64 `json:"court_costs_left_to_pay_amount"`
}

type ProcedureCostFilterDTO struct {
	Page                        *int    `json:"page"`
	Size                        *int    `json:"size"`
	Subject                     *string `json:"subject"`
	FilterByProcedureCostTypeID *int    `json:"procedure_cost_type_id"`
	Search                      *string `json:"search"`
}

// ToProcedureCost converts ProcedureCostDTO to ProcedureCost
func (dto ProcedureCostDTO) ToProcedureCost() *data.ProcedureCost {
	return &data.ProcedureCost{
		ProcedureCostType:      dto.ProcedureCostType,
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

// ToProcedureCostResponseDTO converts ProcedureCost to ProcedureCostResponseDTO
func ToProcedureCostResponseDTO(data data.ProcedureCost) ProcedureCostResponseDTO {
	filesArray := make([]int, len(data.File))
	for i, id := range data.File {
		filesArray[i] = int(id)
	}
	return ProcedureCostResponseDTO{
		ID:                     data.ID,
		ProcedureCostType:      data.ProcedureCostType,
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

// ToProcedureCostListResponseDTO converts []*ProcedureCost to []ProcedureCostResponseDTO
func ToProcedureCostListResponseDTO(procedurecosts []*data.ProcedureCost) []ProcedureCostResponseDTO {
	dtoList := make([]ProcedureCostResponseDTO, len(procedurecosts))
	for i, x := range procedurecosts {
		dtoList[i] = ToProcedureCostResponseDTO(*x)
	}
	return dtoList
}
