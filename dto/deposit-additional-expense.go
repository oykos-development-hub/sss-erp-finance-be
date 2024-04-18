package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type DepositAdditionalExpenseDTO struct {
	ID                   int     `json:"id"`
	Title                string  `json:"title"`
	AccountID            int     `json:"account_id"`
	SubjectID            int     `json:"subject_id"`
	BankAccount          string  `json:"bank_account"`
	PaymentOrderID       int     `json:"payment_order_id"`
	PayingPaymentOrderID *int    `json:"paying_payment_order_id"`
	SourceBankAccount    string  `json:"source_bank_account"`
	OrganizationUnitID   int     `json:"organization_unit_id"`
	Price                float32 `json:"price"`
	Status               string  `json:"status"`
}

type DepositAdditionalExpenseResponseDTO struct {
	ID                   int       `json:"id"`
	Title                string    `json:"title"`
	AccountID            int       `json:"account_id"`
	SubjectID            int       `json:"subject_id"`
	BankAccount          string    `json:"bank_account"`
	PaymentOrderID       int       `json:"payment_order_id"`
	PayingPaymentOrderID *int      `json:"paying_payment_order_id"`
	OrganizationUnitID   int       `json:"organization_unit_id"`
	SourceBankAccount    string    `json:"source_bank_account"`
	Price                float32   `json:"price"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type DepositAdditionalExpenseFilterDTO struct {
	Page                 *int    `json:"page"`
	Size                 *int    `json:"size"`
	SortByTitle          *string `json:"sort_by_title"`
	PaymentOrderID       *int    `json:"payment_order_id"`
	PayingPaymentOrderID *int    `json:"paying_payment_order_id"`
	SubjectID            *int    `json:"subject_id"`
	OrganizationUnitID   *int    `json:"organization_unit_id"`
	Year                 *int    `json:"year"`
	Status               *string `json:"status"`
	Search               *string `json:"search"`
}

func (dto DepositAdditionalExpenseDTO) ToDepositAdditionalExpense() *data.DepositAdditionalExpense {
	return &data.DepositAdditionalExpense{
		Title:                dto.Title,
		AccountID:            dto.AccountID,
		SubjectID:            dto.SubjectID,
		BankAccount:          dto.BankAccount,
		Price:                dto.Price,
		SourceBankAccount:    dto.SourceBankAccount,
		PaymentOrderID:       dto.PaymentOrderID,
		PayingPaymentOrderID: dto.PayingPaymentOrderID,
		OrganizationUnitID:   dto.OrganizationUnitID,
		Status:               dto.Status,
	}
}

func ToDepositAdditionalExpenseResponseDTO(data data.DepositAdditionalExpense) DepositAdditionalExpenseResponseDTO {
	return DepositAdditionalExpenseResponseDTO{
		ID:                   data.ID,
		Title:                data.Title,
		AccountID:            data.AccountID,
		SubjectID:            data.SubjectID,
		SourceBankAccount:    data.SourceBankAccount,
		BankAccount:          data.BankAccount,
		Price:                data.Price,
		PaymentOrderID:       data.PaymentOrderID,
		PayingPaymentOrderID: data.PayingPaymentOrderID,
		OrganizationUnitID:   data.OrganizationUnitID,
		Status:               data.Status,
		CreatedAt:            data.CreatedAt,
		UpdatedAt:            data.UpdatedAt,
	}
}

func ToDepositAdditionalExpenseListResponseDTO(deposit_additional_expenses []*data.DepositAdditionalExpense) []DepositAdditionalExpenseResponseDTO {
	dtoList := make([]DepositAdditionalExpenseResponseDTO, len(deposit_additional_expenses))
	for i, x := range deposit_additional_expenses {
		dtoList[i] = ToDepositAdditionalExpenseResponseDTO(*x)
	}
	return dtoList
}
