package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type SalaryAdditionalExpenseDTO struct {
	ID                 int     `json:"id"`
	SalaryID           int     `json:"salary_id"`
	AccountID          int     `json:"account_id"`
	Amount             float64 `json:"amount"`
	SubjectID          int     `json:"subject_id"`
	BankAccount        string  `json:"bank_account"`
	Status             string  `json:"status"`
	OrganizationUnitID int     `json:"organization_unit_id"`
	Type               string  `json:"type"`
}

type SalaryAdditionalExpenseResponseDTO struct {
	ID                 int       `json:"id"`
	SalaryID           int       `json:"salary_id"`
	AccountID          int       `json:"account_id"`
	Amount             float64   `json:"amount"`
	SubjectID          int       `json:"subject_id"`
	BankAccount        string    `json:"bank_account"`
	Status             string    `json:"status"`
	OrganizationUnitID int       `json:"organization_unit_id"`
	Type               string    `json:"type"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type SalaryAdditionalExpenseFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	SalaryID           *int    `json:"salary_id"`
	Status             *string `json:"status"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
}

func (dto SalaryAdditionalExpenseDTO) ToSalaryAdditionalExpense() *data.SalaryAdditionalExpense {
	return &data.SalaryAdditionalExpense{
		SalaryID:           dto.SalaryID,
		AccountID:          dto.AccountID,
		Amount:             dto.Amount,
		SubjectID:          dto.SubjectID,
		BankAccount:        dto.BankAccount,
		Status:             dto.Status,
		OrganizationUnitID: dto.OrganizationUnitID,
		Type:               dto.Type,
	}
}

func ToSalaryAdditionalExpenseResponseDTO(data data.SalaryAdditionalExpense) SalaryAdditionalExpenseResponseDTO {
	return SalaryAdditionalExpenseResponseDTO{
		ID:                 data.ID,
		SalaryID:           data.SalaryID,
		AccountID:          data.AccountID,
		Amount:             data.Amount,
		SubjectID:          data.SubjectID,
		BankAccount:        data.BankAccount,
		Status:             data.Status,
		OrganizationUnitID: data.OrganizationUnitID,
		Type:               data.Type,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToSalaryAdditionalExpenseListResponseDTO(salaryAdditionalExpenses []*data.SalaryAdditionalExpense) []SalaryAdditionalExpenseResponseDTO {
	dtoList := make([]SalaryAdditionalExpenseResponseDTO, len(salaryAdditionalExpenses))
	for i, x := range salaryAdditionalExpenses {
		dtoList[i] = ToSalaryAdditionalExpenseResponseDTO(*x)
	}
	return dtoList
}
