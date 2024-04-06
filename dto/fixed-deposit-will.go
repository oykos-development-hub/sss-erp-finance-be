package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FixedDepositWillDTO struct {
	OrganizationUnitID int        `json:"organization_unit_id"`
	Subject            string     `json:"subject"`
	FatherName         string     `json:"father_name"`
	DateOfBirth        time.Time  `json:"date_of_birth"`
	JMBG               string     `json:"jmbg"`
	CaseNumberSI       string     `json:"case_number_si"`
	CaseNumberRS       string     `json:"case_number_rs"`
	DateOfReceiptSI    *time.Time `json:"date_of_receipt_si"`
	DateOfReceiptRS    *time.Time `json:"date_of_receipt_rs"`
	DateOfEnd          *time.Time `json:"date_of_end"`
	Status             string     `json:"status"`
	FileID             int        `json:"file_id"`
}

type FixedDepositWillResponseDTO struct {
	ID                 int                                   `json:"id"`
	OrganizationUnitID int                                   `json:"organization_unit_id"`
	Subject            string                                `json:"subject"`
	FatherName         string                                `json:"father_name"`
	DateOfBirth        time.Time                             `json:"date_of_birth"`
	JMBG               string                                `json:"jmbg"`
	CaseNumberSI       string                                `json:"case_number_si"`
	CaseNumberRS       string                                `json:"case_number_rs"`
	DateOfReceiptSI    *time.Time                            `json:"date_of_receipt_si"`
	DateOfReceiptRS    *time.Time                            `json:"date_of_receipt_rs"`
	DateOfEnd          *time.Time                            `json:"date_of_end"`
	Status             string                                `json:"status"`
	FileID             int                                   `json:"file_id"`
	Judges             []FixedDepositJudgeResponseDTO        `json:"judges"`
	Dispatches         []FixedDepositWillDispatchResponseDTO `json:"dispatches"`
	CreatedAt          time.Time                             `json:"created_at"`
	UpdatedAt          time.Time                             `json:"updated_at"`
}

type FixedDepositWillFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	Search             *string `json:"search"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	Status             *string `json:"status"`
}

func (dto FixedDepositWillDTO) ToFixedDepositWill() *data.FixedDepositWill {
	return &data.FixedDepositWill{
		OrganizationUnitID: dto.OrganizationUnitID,
		Subject:            dto.Subject,
		FatherName:         dto.FatherName,
		DateOfBirth:        dto.DateOfBirth,
		JMBG:               dto.JMBG,
		CaseNumberSI:       dto.CaseNumberSI,
		CaseNumberRS:       dto.CaseNumberRS,
		DateOfReceiptSI:    dto.DateOfReceiptSI,
		DateOfReceiptRS:    dto.DateOfReceiptRS,
		DateOfEnd:          dto.DateOfEnd,
		Status:             dto.Status,
		FileID:             dto.FileID,
	}
}

func ToFixedDepositWillResponseDTO(data data.FixedDepositWill) FixedDepositWillResponseDTO {
	return FixedDepositWillResponseDTO{
		ID:                 data.ID,
		OrganizationUnitID: data.OrganizationUnitID,
		Subject:            data.Subject,
		FatherName:         data.FatherName,
		DateOfBirth:        data.DateOfBirth,
		JMBG:               data.JMBG,
		CaseNumberSI:       data.CaseNumberSI,
		CaseNumberRS:       data.CaseNumberRS,
		DateOfReceiptSI:    data.DateOfReceiptSI,
		DateOfReceiptRS:    data.DateOfReceiptRS,
		DateOfEnd:          data.DateOfEnd,
		Status:             data.Status,
		FileID:             data.FileID,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToFixedDepositWillListResponseDTO(fixed_deposit_wills []*data.FixedDepositWill) []FixedDepositWillResponseDTO {
	dtoList := make([]FixedDepositWillResponseDTO, len(fixed_deposit_wills))
	for i, x := range fixed_deposit_wills {
		dtoList[i] = ToFixedDepositWillResponseDTO(*x)
	}
	return dtoList
}
