package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type EnforcedPaymentDTO struct {
	OrganizationUnitID int                        `json:"organization_unit_id"`
	SupplierID         int                        `json:"supplier_id"`
	BankAccount        string                     `json:"bank_account"`
	DateOfPayment      time.Time                  `json:"date_of_payment"`
	DateOfOrder        *time.Time                 `json:"date_of_order"`
	IDOfStatement      *string                    `json:"id_of_statement"`
	SAPID              *string                    `json:"sap_id"`
	Registred          *bool                      `json:"registred"`
	RegistredReturn    *bool                      `json:"registred_return"`
	Status             data.EnforcedPaymentStatus `json:"status"`
	ReturnFileID       *int                       `json:"return_file_id"`
	Description        string                     `json:"description"`
	DateOfSAP          *time.Time                 `json:"date_of_sap"`
	ReturnDate         *time.Time                 `json:"return_date"`
	Items              []EnforcedPaymentItemDTO   `json:"items"`
	Amount             float64                    `json:"amount"`
	AmountForLawyer    float64                    `json:"amount_for_lawyer"`
	AmountForAgent     float64                    `json:"amount_for_agent"`
	FileID             *int                       `json:"file_id"`
}

type EnforcedPaymentResponseDTO struct {
	ID                 int                              `json:"id"`
	OrganizationUnitID int                              `json:"organization_unit_id"`
	SupplierID         int                              `json:"supplier_id"`
	BankAccount        string                           `json:"bank_account"`
	DateOfPayment      time.Time                        `json:"date_of_payment"`
	DateOfOrder        *time.Time                       `json:"date_of_order"`
	IDOfStatement      *string                          `json:"id_of_statement"`
	SAPID              *string                          `json:"sap_id"`
	Description        string                           `json:"description"`
	Status             data.EnforcedPaymentStatus       `json:"status"`
	Registred          *bool                            `json:"registred"`
	RegistredReturn    *bool                            `json:"registred_return"`
	ReturnFileID       *int                             `json:"return_file_id"`
	DateOfSAP          *time.Time                       `json:"date_of_sap"`
	ReturnDate         *time.Time                       `json:"return_date"`
	Items              []EnforcedPaymentItemResponseDTO `json:"items"`
	Amount             float64                          `json:"amount"`
	AmountForLawyer    float64                          `json:"amount_for_lawyer"`
	AmountForAgent     float64                          `json:"amount_for_agent"`
	FileID             *int                             `json:"file_id"`
	CreatedAt          time.Time                        `json:"created_at"`
	UpdatedAt          time.Time                        `json:"updated_at"`
}

type EnforcedPaymentFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	Status             *string `json:"status"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	SupplierID         *int    `json:"supplier_id"`
	Search             *string `json:"search"`
	Year               *int    `json:"year"`
	Registred          *bool   `json:"registred"`
}

func (dto EnforcedPaymentDTO) ToEnforcedPayment() *data.EnforcedPayment {
	return &data.EnforcedPayment{
		OrganizationUnitID: dto.OrganizationUnitID,
		SupplierID:         dto.SupplierID,
		BankAccount:        dto.BankAccount,
		DateOfPayment:      dto.DateOfPayment,
		IDOfStatement:      dto.IDOfStatement,
		SAPID:              dto.SAPID,
		DateOfOrder:        dto.DateOfOrder,
		DateOfSAP:          dto.DateOfSAP,
		FileID:             dto.FileID,
		Registred:          dto.Registred,
		RegistredReturn:    dto.RegistredReturn,
		Description:        dto.Description,
		Amount:             dto.Amount,
		AmountForLawyer:    dto.AmountForLawyer,
		AmountForAgent:     dto.AmountForAgent,
		ReturnDate:         dto.ReturnDate,
		Status:             dto.Status,
		ReturnFileID:       dto.ReturnFileID,
	}
}

func ToEnforcedPaymentResponseDTO(data data.EnforcedPayment) EnforcedPaymentResponseDTO {
	return EnforcedPaymentResponseDTO{
		ID:                 data.ID,
		OrganizationUnitID: data.OrganizationUnitID,
		SupplierID:         data.SupplierID,
		BankAccount:        data.BankAccount,
		DateOfPayment:      data.DateOfPayment,
		IDOfStatement:      data.IDOfStatement,
		SAPID:              data.SAPID,
		Registred:          data.Registred,
		RegistredReturn:    data.RegistredReturn,
		DateOfOrder:        data.DateOfOrder,
		DateOfSAP:          data.DateOfSAP,
		FileID:             data.FileID,
		Description:        data.Description,
		Amount:             data.Amount,
		AmountForLawyer:    data.AmountForLawyer,
		AmountForAgent:     data.AmountForAgent,
		ReturnDate:         data.ReturnDate,
		Status:             data.Status,
		ReturnFileID:       data.ReturnFileID,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToEnforcedPaymentListResponseDTO(enforced_payments []*data.EnforcedPayment) []EnforcedPaymentResponseDTO {
	dtoList := make([]EnforcedPaymentResponseDTO, len(enforced_payments))
	for i, x := range enforced_payments {
		dtoList[i] = ToEnforcedPaymentResponseDTO(*x)
	}
	return dtoList
}
