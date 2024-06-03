package dto

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type PropBenConfDTO struct {
	PropBenConfType        data.PropBenConfType    `json:"property_benefits_confiscation_type" validate:"required,oneof=1 2"`
	DecisionNumber         string                  `json:"decision_number" validate:"required"`
	DecisionDate           time.Time               `json:"decision_date"`
	Subject                string                  `json:"subject"`
	JMBG                   string                  `json:"jmbg" validate:"required"`
	Residence              string                  `json:"residence"`
	Amount                 decimal.Decimal         `json:"amount"`
	PaymentReferenceNumber string                  `json:"payment_reference_number"`
	DebitReferenceNumber   string                  `json:"debit_reference_number"`
	AccountID              int                     `json:"account_id"`
	ExecutionDate          time.Time               `json:"execution_date"`
	PaymentDeadlineDate    time.Time               `json:"payment_deadline_date"`
	Description            string                  `json:"description"`
	Status                 *data.PropBenConfStatus `json:"status"`
	CourtCosts             *decimal.Decimal        `json:"court_costs"`
	CourtAccountID         *int                    `json:"court_account_id"`
	File                   pq.Int64Array           `json:"file"`
}

type PropBenConfResponseDTO struct {
	ID                     int                    `json:"id"`
	PropBenConfType        data.PropBenConfType   `json:"property_benefits_confiscation_type"`
	DecisionNumber         string                 `json:"decision_number"`
	DecisionDate           time.Time              `json:"decision_date"`
	Subject                string                 `json:"subject"`
	JMBG                   string                 `json:"jmbg"`
	Residence              string                 `json:"residence"`
	Amount                 decimal.Decimal        `json:"amount"`
	PaymentReferenceNumber string                 `json:"payment_reference_number"`
	DebitReferenceNumber   string                 `json:"debit_reference_number"`
	AccountID              int                    `json:"account_id"`
	ExecutionDate          time.Time              `json:"execution_date"`
	PaymentDeadlineDate    time.Time              `json:"payment_deadline_date"`
	Description            string                 `json:"description"`
	Status                 data.PropBenConfStatus `json:"status"`
	CourtCosts             *decimal.Decimal       `json:"court_costs"`
	CourtAccountID         *int                   `json:"court_account_id"`
	PropBenConfDetailsDTO  *PropBenConfDetailsDTO `json:"property_benefits_confiscation_details"`
	File                   []int                  `json:"file"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type PropBenConfDetailsDTO struct {
	AllPaymentAmount           decimal.Decimal `json:"all_payments_amount"`
	AmountGracePeriod          decimal.Decimal `json:"amount_grace_period"`
	AmountGracePeriodDueDate   time.Time       `json:"amount_grace_period_due_date"`
	AmountGracePeriodAvailable bool            `json:"amount_grace_period_available"`
	LeftToPayAmount            decimal.Decimal `json:"left_to_pay_amount"`

	CourtCostsPaid            decimal.Decimal `json:"court_costs_paid"`
	CourtCostsLeftToPayAmount decimal.Decimal `json:"court_costs_left_to_pay_amount"`
}

type PropBenConfFilterDTO struct {
	Page                      *int    `json:"page"`
	Size                      *int    `json:"size"`
	Subject                   *string `json:"subject"`
	FilterByPropBenConfTypeID *int    `json:"property_benefits_confiscation_type_id"`
	Search                    *string `json:"search"`
}

// ToPropBenConf converts PropBenConfDTO to PropBenConf
func (dto PropBenConfDTO) ToPropBenConf() *data.PropBenConf {
	return &data.PropBenConf{
		PropBenConfType:        dto.PropBenConfType,
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

// ToPropBenConfResponseDTO converts PropBenConf to PropBenConfResponseDTO
func ToPropBenConfResponseDTO(data data.PropBenConf) PropBenConfResponseDTO {
	filesArray := make([]int, len(data.File))
	for i, id := range data.File {
		filesArray[i] = int(id)
	}
	return PropBenConfResponseDTO{
		ID:                     data.ID,
		PropBenConfType:        data.PropBenConfType,
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

// ToPropBenConfListResponseDTO converts []*PropBenConf to []PropBenConfResponseDTO
func ToPropBenConfListResponseDTO(propbenconfs []*data.PropBenConf) []PropBenConfResponseDTO {
	dtoList := make([]PropBenConfResponseDTO, len(propbenconfs))
	for i, x := range propbenconfs {
		dtoList[i] = ToPropBenConfResponseDTO(*x)
	}
	return dtoList
}
