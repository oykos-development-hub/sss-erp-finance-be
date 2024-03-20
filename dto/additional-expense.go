package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type AdditionalExpenseDTO struct {
	Title              string                       `json:"title"`
	AccountID          int                          `json:"account_id"`
	OrganizationUnitID int                          `json:"organization_unit_id"`
	Price              float32                      `json:"price"`
	SubjectID          int                          `json:"subject_id"`
	BankAccount        string                       `json:"bank_account"`
	InvoiceID          int                          `json:"invoice_id"`
	Status             data.AdditionalExpenseStatus `json:"status"`
}

type AdditionalExpenseResponseDTO struct {
	ID                 int                          `json:"id"`
	Title              string                       `json:"title"`
	AccountID          int                          `json:"account_id"`
	Price              float32                      `json:"price"`
	SubjectID          int                          `json:"subject_id"`
	OrganizationUnitID int                          `json:"organization_unit_id"`
	BankAccount        string                       `json:"bank_account"`
	InvoiceID          int                          `json:"invoice_id"`
	Status             data.AdditionalExpenseStatus `json:"status"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
}

type AdditionalExpenseFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	InvoiceID          *int    `json:"invoice_id"`
	SubjectID          *int    `json:"subject_id"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
	Year               *int    `json:"year"`
	Status             *int    `json:"status"`
	Search             *string `json:"search"`
}

func (dto AdditionalExpenseDTO) ToAdditionalExpense() *data.AdditionalExpense {
	return &data.AdditionalExpense{
		Title:              dto.Title,
		AccountID:          dto.AccountID,
		SubjectID:          dto.SubjectID,
		BankAccount:        dto.BankAccount,
		InvoiceID:          dto.InvoiceID,
		OrganizationUnitID: dto.OrganizationUnitID,
		Status:             dto.Status,
	}
}

func ToAdditionalExpenseResponseDTO(data data.AdditionalExpense) AdditionalExpenseResponseDTO {
	return AdditionalExpenseResponseDTO{
		ID:                 data.ID,
		Title:              data.Title,
		AccountID:          data.AccountID,
		SubjectID:          data.SubjectID,
		OrganizationUnitID: data.OrganizationUnitID,
		BankAccount:        data.BankAccount,
		InvoiceID:          data.InvoiceID,
		Status:             data.Status,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToAdditionalExpenseListResponseDTO(additional_expenses []*data.AdditionalExpense) []AdditionalExpenseResponseDTO {
	dtoList := make([]AdditionalExpenseResponseDTO, len(additional_expenses))
	for i, x := range additional_expenses {
		dtoList[i] = ToAdditionalExpenseResponseDTO(*x)
	}
	return dtoList
}
