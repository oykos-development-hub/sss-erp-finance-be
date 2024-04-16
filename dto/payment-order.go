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
	IDOfStatement      *string               `json:"id_of_statement"`
	SAPID              *string               `json:"sap_id"`
	DateOfSAP          *time.Time            `json:"date_of_sap"`
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
	IDOfStatement      *string                       `json:"id_of_statement"`
	SAPID              *string                       `json:"sap_id"`
	DateOfSAP          *time.Time                    `json:"date_of_sap"`
	FileID             *int                          `json:"file_id"`
	Items              []PaymentOrderItemResponseDTO `json:"items"`
	Amount             float64                       `json:"amount"`
	Status             string                        `json:"status"`
	CreatedAt          time.Time                     `json:"created_at"`
	UpdatedAt          time.Time                     `json:"updated_at"`
}

type ObligationResponse struct {
	InvoiceID                 *int      `json:"invoice_id"`
	AdditionalExpenseID       *int      `json:"additional_expense_id"`
	SalaryAdditionalExpenseID *int      `json:"salary_additional_expense_id"`
	Type                      string    `json:"type"`
	Title                     string    `json:"title"`
	Status                    string    `json:"status"`
	Price                     float64   `json:"price"`
	CreatedAt                 time.Time `json:"created_at"`
}

type PaymentOrderFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	Status             *string `json:"status"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	SupplierID         *int    `json:"supplier_id"`
	Search             *string `json:"search"`
}

type GetObligationsFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	OrganizationUnitID int     `json:"organization_unit_id"`
	SupplierID         int     `json:"supplier_id"`
	Type               *string `json:"type"`
}

func (dto PaymentOrderDTO) ToPaymentOrder() *data.PaymentOrder {
	return &data.PaymentOrder{
		OrganizationUnitID: dto.OrganizationUnitID,
		SupplierID:         dto.SupplierID,
		BankAccount:        dto.BankAccount,
		DateOfPayment:      dto.DateOfPayment,
		IDOfStatement:      dto.IDOfStatement,
		SAPID:              dto.SAPID,
		DateOfSAP:          dto.DateOfSAP,
		FileID:             dto.FileID,
		Amount:             dto.Amount,
	}
}

func ToPaymentOrderResponseDTO(data data.PaymentOrder) PaymentOrderResponseDTO {

	var status string

	if data.IDOfStatement != nil {
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
		IDOfStatement:      data.IDOfStatement,
		Status:             status,
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
