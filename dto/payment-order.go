package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type PaymentOrderDTO struct {
	OrganizationUnitID int                   `json:"organization_unit_id"`
	SupplierID         int                   `json:"supplier_id"`
	BankAccount        string                `json:"bank_account"`
	DateOfPayment      time.Time             `json:"date_of_payment"`
	DateOfOrder        *time.Time            `json:"date_of_order"`
	IDOfStatement      *int                  `json:"id_of_statement"`
	SAPID              *string               `json:"sap_id"`
	Registred          *bool                 `json:"registred"`
	SourceOfFunding    string                `json:"source_of_funding"`
	Description        string                `json:"description"`
	DateOfSAP          *time.Time            `json:"date_of_sap"`
	Status             string                `json:"status"`
	Items              []PaymentOrderItemDTO `json:"items"`
	Amount             float64               `json:"amount"`
	FileID             *int                  `json:"file_id"`
}

type PaymentOrderResponseDTO struct {
	ID                 int                           `json:"id"`
	OrganizationUnitID int                           `json:"organization_unit_id"`
	SupplierID         int                           `json:"supplier_id"`
	BankAccount        string                        `json:"bank_account"`
	DateOfPayment      time.Time                     `json:"date_of_payment"`
	DateOfOrder        *time.Time                    `json:"date_of_order"`
	IDOfStatement      *int                          `json:"id_of_statement"`
	SAPID              *string                       `json:"sap_id"`
	Registred          *bool                         `json:"registred"`
	SourceOfFunding    string                        `json:"source_of_funding"`
	Description        string                        `json:"description"`
	DateOfSAP          *time.Time                    `json:"date_of_sap"`
	FileID             *int                          `json:"file_id"`
	Items              []PaymentOrderItemResponseDTO `json:"items"`
	Amount             float64                       `json:"amount"`
	Status             string                        `json:"status"`
	CreatedAt          time.Time                     `json:"created_at"`
	UpdatedAt          time.Time                     `json:"updated_at"`
}

type ObligationResponse struct {
	InvoiceID                 *int                   `json:"invoice_id"`
	AdditionalExpenseID       *int                   `json:"additional_expense_id"`
	SalaryAdditionalExpenseID *int                   `json:"salary_additional_expense_id"`
	Type                      data.TypesOfObligation `json:"type"`
	Title                     string                 `json:"title"`
	Status                    string                 `json:"status"`
	TotalPrice                float64                `json:"total_price"`
	RemainPrice               float64                `json:"remain_price"`
	CreatedAt                 time.Time              `json:"created_at"`
}

type PaymentOrderFilterDTO struct {
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

type GetObligationsFilterDTO struct {
	Page               *int                    `json:"page"`
	Size               *int                    `json:"size"`
	OrganizationUnitID int                     `json:"organization_unit_id"`
	SupplierID         int                     `json:"supplier_id"`
	Type               *data.TypesOfObligation `json:"type"`
	Search             *string                 `json:"search"`
}

func (dto PaymentOrderDTO) ToPaymentOrder() *data.PaymentOrder {
	return &data.PaymentOrder{
		OrganizationUnitID: dto.OrganizationUnitID,
		SupplierID:         dto.SupplierID,
		BankAccount:        dto.BankAccount,
		DateOfPayment:      dto.DateOfPayment,
		IDOfStatement:      dto.IDOfStatement,
		SAPID:              dto.SAPID,
		SourceOfFunding:    dto.SourceOfFunding,
		DateOfOrder:        dto.DateOfOrder,
		DateOfSAP:          dto.DateOfSAP,
		FileID:             dto.FileID,
		Description:        dto.Description,
		Amount:             dto.Amount,
		Registred:          dto.Registred,
	}
}

func ToPaymentOrderResponseDTO(data data.PaymentOrder) PaymentOrderResponseDTO {

	var status string
	if data.Status == "Storniran" {
		status = "Storniran"
	} else if data.SAPID != nil {
		status = "Plaćen"
	} else {
		status = "Na čekanju"
	}

	return PaymentOrderResponseDTO{
		ID:                 data.ID,
		OrganizationUnitID: data.OrganizationUnitID,
		SupplierID:         data.SupplierID,
		BankAccount:        data.BankAccount,
		DateOfPayment:      data.DateOfPayment,
		DateOfOrder:        data.DateOfOrder,
		IDOfStatement:      data.IDOfStatement,
		Status:             status,
		Registred:          data.Registred,
		SourceOfFunding:    data.SourceOfFunding,
		Description:        data.Description,
		Amount:             data.Amount,
		SAPID:              data.SAPID,
		DateOfSAP:          data.DateOfSAP,
		FileID:             data.FileID,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToPaymentOrderListResponseDTO(payment_orders []*data.PaymentOrder) []PaymentOrderResponseDTO {
	dtoList := make([]PaymentOrderResponseDTO, len(payment_orders))
	for i, x := range payment_orders {
		dtoList[i] = ToPaymentOrderResponseDTO(*x)
	}
	return dtoList
}
