package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FixedDepositDTO struct {
	OrganizationUnitID   int        `json:"organization_unit_id"`
	Subject              string     `json:"subject"`
	CaseNumber           string     `json:"case_number"`
	JudgeID              int        `json:"judge_id"`
	DateOfRecipiet       *time.Time `json:"date_of_recipiet"`
	DateOfCase           *time.Time `json:"date_of_case"`
	DateOfFinality       *time.Time `json:"date_of_finality"`
	DateOfEnforceability *time.Time `json:"date_of_enforceability"`
	AccountID            int        `json:"account_id"`
	Type                 string     `json:"type"`
	FileID               int        `json:"file_id"`
	Status               string     `json:"status"`
}

type FixedDepositResponseDTO struct {
	ID                   int                               `json:"id"`
	OrganizationUnitID   int                               `json:"organization_unit_id"`
	Subject              string                            `json:"subject"`
	JudgeID              int                               `json:"judge_id"`
	CaseNumber           string                            `json:"case_number"`
	DateOfRecipiet       *time.Time                        `json:"date_of_recipiet"`
	DateOfCase           *time.Time                        `json:"date_of_case"`
	DateOfFinality       *time.Time                        `json:"date_of_finality"`
	DateOfEnforceability *time.Time                        `json:"date_of_enforceability"`
	AccountID            int                               `json:"account_id"`
	FileID               int                               `json:"file_id"`
	Status               string                            `json:"status"`
	Type                 string                            `json:"type"`
	Items                []FixedDepositItemResponseDTO     `json:"items"`
	Dispatches           []FixedDepositDispatchResponseDTO `json:"dispatches"`
	CreatedAt            time.Time                         `json:"created_at"`
	UpdatedAt            time.Time                         `json:"updated_at"`
}

type FixedDepositFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	Type               *string `json:"type"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	Search             *string `json:"search"`
	JudgeID            *int    `json:"judge_id"`
	Status             *string `json:"status"`
	Subject            *string `json:"subject"`
}

func (dto FixedDepositDTO) ToFixedDeposit() *data.FixedDeposit {
	return &data.FixedDeposit{
		OrganizationUnitID:   dto.OrganizationUnitID,
		Subject:              dto.Subject,
		JudgeID:              dto.JudgeID,
		CaseNumber:           dto.CaseNumber,
		DateOfRecipiet:       dto.DateOfRecipiet,
		DateOfCase:           dto.DateOfCase,
		DateOfFinality:       dto.DateOfFinality,
		DateOfEnforceability: dto.DateOfEnforceability,
		AccountID:            dto.AccountID,
		Type:                 dto.Type,
		FileID:               dto.FileID,
		Status:               dto.Status,
	}
}

func ToFixedDepositResponseDTO(data data.FixedDeposit) FixedDepositResponseDTO {
	return FixedDepositResponseDTO{
		ID:                 data.ID,
		OrganizationUnitID: data.OrganizationUnitID,
		Subject:            data.Subject,
		JudgeID:            data.JudgeID,
		CaseNumber:         data.CaseNumber,
		DateOfRecipiet:     data.DateOfRecipiet,
		DateOfCase:         data.DateOfCase,
		DateOfFinality:     data.DateOfFinality,
		AccountID:          data.AccountID,
		FileID:             data.FileID,
		Type:               data.Type,
		Status:             data.Status,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToFixedDepositListResponseDTO(fixed_deposits []*data.FixedDeposit) []FixedDepositResponseDTO {
	dtoList := make([]FixedDepositResponseDTO, len(fixed_deposits))
	for i, x := range fixed_deposits {
		dtoList[i] = ToFixedDepositResponseDTO(*x)
	}
	return dtoList
}
