package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type DepositPaymentOrderDTO struct {
	OrganizationUnitID          int                           `json:"organization_unit_id"`
	CaseNumber                  string                        `json:"case_number"`
	SupplierID                  int                           `json:"supplier_id"`
	NetAmount                   float64                       `json:"net_amount"`
	BankAccount                 string                        `json:"bank_account"`
	DateOfPayment               time.Time                     `json:"date_of_payment"`
	DateOfStatement             *time.Time                    `json:"date_of_statement"`
	IDOfStatement               *string                       `json:"id_of_statement"`
	AdditionalExpenses          []DepositAdditionalExpenseDTO `json:"additional_expenses"`
	AdditionalExpensesForPaying []DepositAdditionalExpenseDTO `json:"additional_expenses_for_paying"`
}

type DepositPaymentOrderResponseDTO struct {
	ID                          int                                   `json:"id"`
	OrganizationUnitID          int                                   `json:"organization_unit_id"`
	CaseNumber                  string                                `json:"case_number"`
	SupplierID                  int                                   `json:"supplier_id"`
	NetAmount                   float64                               `json:"net_amount"`
	BankAccount                 string                                `json:"bank_account"`
	DateOfPayment               time.Time                             `json:"date_of_payment"`
	DateOfStatement             *time.Time                            `json:"date_of_statement"`
	IDOfStatement               *string                               `json:"id_of_statement"`
	Status                      string                                `json:"status"`
	AdditionalExpenses          []DepositAdditionalExpenseResponseDTO `json:"additional_expenses"`
	AdditionalExpensesForPaying []DepositAdditionalExpenseResponseDTO `json:"additional_expenses_for_paying"`
	CreatedAt                   time.Time                             `json:"created_at"`
	UpdatedAt                   time.Time                             `json:"updated_at"`
}

type DepositPaymentOrderFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	Status             *string `json:"status"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	SupplierID         *int    `json:"supplier_id"`
	Search             *string `json:"search"`
	CaseNumber         *string `json:"case_number"`
}

func (dto DepositPaymentOrderDTO) ToDepositPaymentOrder() *data.DepositPaymentOrder {
	return &data.DepositPaymentOrder{
		OrganizationUnitID: dto.OrganizationUnitID,
		CaseNumber:         dto.CaseNumber,
		SupplierID:         dto.SupplierID,
		NetAmount:          dto.NetAmount,
		BankAccount:        dto.BankAccount,
		DateOfPayment:      dto.DateOfPayment,
		DateOfStatement:    dto.DateOfStatement,
		IDOfStatement:      dto.IDOfStatement,
	}
}

func ToDepositPaymentOrderResponseDTO(data data.DepositPaymentOrder) DepositPaymentOrderResponseDTO {
	response := DepositPaymentOrderResponseDTO{
		ID:                 data.ID,
		OrganizationUnitID: data.OrganizationUnitID,
		CaseNumber:         data.CaseNumber,
		SupplierID:         data.SupplierID,
		NetAmount:          data.NetAmount,
		BankAccount:        data.BankAccount,
		DateOfPayment:      data.DateOfPayment,
		DateOfStatement:    data.DateOfStatement,
		IDOfStatement:      data.IDOfStatement,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}

	if data.IDOfStatement != nil {
		response.Status = "Plaćen"
	} else {
		response.Status = "Na čekanju"
	}

	return response
}

func ToDepositPaymentOrderListResponseDTO(deposit_payment_orders []*data.DepositPaymentOrder) []DepositPaymentOrderResponseDTO {
	dtoList := make([]DepositPaymentOrderResponseDTO, len(deposit_payment_orders))
	for i, x := range deposit_payment_orders {
		dtoList[i] = ToDepositPaymentOrderResponseDTO(*x)
	}
	return dtoList
}
