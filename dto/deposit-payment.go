package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type DepositPaymentDTO struct {
	Payer                     string     `json:"payer"`
	OrganizationUnitID        int        `json:"organization_unit_id"`
	CaseNumber                string     `json:"case_number"`
	PartyName                 string     `json:"party_name"`
	NumberOfBankStatement     string     `json:"number_of_bank_statement"`
	DateOfBankStatement       string     `json:"date_of_bank_statement"`
	AccountID                 int        `json:"account_id"`
	Amount                    float64    `json:"amount"`
	MainBankAccount           bool       `json:"main_bank_account"`
	CurrentBankAccount        string     `json:"current_bank_account"`
	DateOfTransferMainAccount *time.Time `json:"date_of_transfer_main_account"`
}

type DepositPaymentResponseDTO struct {
	ID                        int        `json:"id"`
	Payer                     string     `json:"payer"`
	OrganizationUnitID        int        `json:"organization_unit_id"`
	CaseNumber                string     `json:"case_number"`
	PartyName                 string     `json:"party_name"`
	NumberOfBankStatement     string     `json:"number_of_bank_statement"`
	DateOfBankStatement       string     `json:"date_of_bank_statement"`
	AccountID                 int        `json:"account_id"`
	Amount                    float64    `json:"amount"`
	MainBankAccount           bool       `json:"main_bank_account"`
	CurrentBankAccount        string     `json:"current_bank_account"`
	DateOfTransferMainAccount *time.Time `json:"date_of_transfer_main_account"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

type DepositPaymentFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	Status             *string `json:"status"`
	Search             *string `json:"search"`
	CaseNumber         *string `json:"case_number"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
}

func (dto DepositPaymentDTO) ToDepositPayment() *data.DepositPayment {
	return &data.DepositPayment{
		Payer:                     dto.Payer,
		CaseNumber:                dto.CaseNumber,
		PartyName:                 dto.PartyName,
		NumberOfBankStatement:     dto.NumberOfBankStatement,
		DateOfBankStatement:       dto.DateOfBankStatement,
		AccountID:                 dto.AccountID,
		Amount:                    dto.Amount,
		MainBankAccount:           dto.MainBankAccount,
		DateOfTransferMainAccount: dto.DateOfTransferMainAccount,
		CurrentBankAccount:        dto.CurrentBankAccount,
		OrganizationUnitID:        dto.OrganizationUnitID,
	}
}

func ToDepositPaymentResponseDTO(data data.DepositPayment) DepositPaymentResponseDTO {
	return DepositPaymentResponseDTO{
		ID:                        data.ID,
		Payer:                     data.Payer,
		CaseNumber:                data.CaseNumber,
		PartyName:                 data.PartyName,
		NumberOfBankStatement:     data.NumberOfBankStatement,
		DateOfBankStatement:       data.DateOfBankStatement,
		AccountID:                 data.AccountID,
		OrganizationUnitID:        data.OrganizationUnitID,
		Amount:                    data.Amount,
		MainBankAccount:           data.MainBankAccount,
		DateOfTransferMainAccount: data.DateOfTransferMainAccount,
		CurrentBankAccount:        data.CurrentBankAccount,
		CreatedAt:                 data.CreatedAt,
		UpdatedAt:                 data.UpdatedAt,
	}
}

func ToDepositPaymentListResponseDTO(deposit_payments []*data.DepositPayment) []DepositPaymentResponseDTO {
	dtoList := make([]DepositPaymentResponseDTO, len(deposit_payments))
	for i, x := range deposit_payments {
		dtoList[i] = ToDepositPaymentResponseDTO(*x)
	}
	return dtoList
}
