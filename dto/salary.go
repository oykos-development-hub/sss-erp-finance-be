package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type SalaryDTO struct {
	ActivityID               int                          `json:"activity_id"`
	Month                    string                       `json:"month"`
	DateOfCalculation        time.Time                    `json:"date_of_calculation"`
	Description              string                       `json:"description"`
	OrganizationUnitID       int                          `json:"organization_unit_id"`
	Status                   string                       `json:"status"`
	NumberOfEmployees        int                          `json:"number_of_employees"`
	SalaryAdditionalExpenses []SalaryAdditionalExpenseDTO `json:"salary_additional_expenses"`
	Registred                *bool                        `json:"registred"`
}

type SalaryResponseDTO struct {
	ID                       int                                  `json:"id"`
	ActivityID               int                                  `json:"activity_id"`
	Month                    string                               `json:"month"`
	DateOfCalculation        time.Time                            `json:"date_of_calculation"`
	Description              string                               `json:"description"`
	Status                   string                               `json:"status"`
	OrganizationUnitID       int                                  `json:"organization_unit_id"`
	SalaryAdditionalExpenses []SalaryAdditionalExpenseResponseDTO `json:"salary_additional_expenses"`
	GrossPrice               decimal.Decimal                      `json:"gross_price"`
	Registred                *bool                                `json:"registred"`
	VatPrice                 decimal.Decimal                      `json:"vat_price"`
	ObligationsPrice         decimal.Decimal                      `json:"obligations_price"`
	NetPrice                 decimal.Decimal                      `json:"net_price"`
	NumberOfEmployees        int                                  `json:"number_of_employees"`
	Deletable                bool                                 `json:"is_deletable"`
	CreatedAt                time.Time                            `json:"created_at"`
	UpdatedAt                time.Time                            `json:"updated_at"`
}

type SalaryFilterDTO struct {
	Page               *int    `json:"page"`
	Size               *int    `json:"size"`
	SortByTitle        *string `json:"sort_by_title"`
	Month              *string `json:"month"`
	Status             *string `json:"status"`
	Year               *int    `json:"year"`
	ActivityID         *int    `json:"activity_id"`
	Registred          *bool   `json:"registred"`
	OrganizationUnitID *int    `json:"organization_unit_id"`
}

func (dto SalaryDTO) ToSalary() *data.Salary {
	return &data.Salary{
		ActivityID:         dto.ActivityID,
		Month:              dto.Month,
		DateOfCalculation:  dto.DateOfCalculation,
		Description:        dto.Description,
		Status:             dto.Status,
		Registred:          dto.Registred,
		OrganizationUnitID: dto.OrganizationUnitID,
		NumberOfEmployees:  dto.NumberOfEmployees,
	}
}

func ToSalaryResponseDTO(data data.Salary) SalaryResponseDTO {
	return SalaryResponseDTO{
		ID:                 data.ID,
		ActivityID:         data.ActivityID,
		Month:              data.Month,
		DateOfCalculation:  data.DateOfCalculation,
		Description:        data.Description,
		Status:             data.Status,
		OrganizationUnitID: data.OrganizationUnitID,
		NumberOfEmployees:  data.NumberOfEmployees,
		Registred:          data.Registred,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToSalaryListResponseDTO(salaries []*data.Salary) []SalaryResponseDTO {
	dtoList := make([]SalaryResponseDTO, len(salaries))
	for i, x := range salaries {
		dtoList[i] = ToSalaryResponseDTO(*x)
	}
	return dtoList
}
