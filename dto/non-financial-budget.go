package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type NonFinancialBudgetDTO struct {
	BudetID            int `json:"budget_id" validate:"required"`
	OrganizationUnitID int `json:"organization_unit_id" validate:"required"`

	ImplContactFullName     string `json:"impl_contact_fullname" validate:"required"`
	ImplContactWorkingPlace string `json:"impl_contact_working_place" validate:"required"`
	ImplContactPhone        string `json:"impl_contact_phone" validate:"required"`
	ImplContactEmail        string `json:"impl_contact_email" validate:"required"`

	ContactFullName     string `json:"contact_fullname" validate:"required"`
	ContactWorkingPlace string `json:"contact_working_place" validate:"required"`
	ContactPhone        string `json:"contact_phone" validate:"required"`
	ContactEmail        string `json:"contact_email" validate:"required"`
}

type NonFinancialBudgetResponseDTO struct {
	ID                 int `json:"id"`
	BudetID            int `json:"budget_id"`
	OrganizationUnitID int `json:"organization_unit_id"`

	ImplContactFullName     string `json:"impl_contact_fullname"`
	ImplContactWorkingPlace string `json:"impl_contact_working_place"`
	ImplContactPhone        string `json:"impl_contact_phone"`
	ImplContactEmail        string `json:"impl_contact_email"`

	ContactFullName     string `json:"contact_fullname"`
	ContactWorkingPlace string `json:"contact_working_place"`
	ContactPhone        string `json:"contact_phone"`
	ContactEmail        string `json:"contact_email"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NonFinancialBudgetFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	BudgetID    *int    `json:"budget_id"`
	SortByTitle *string `json:"sort_by_title"`
}

func (dto NonFinancialBudgetDTO) ToNonFinancialBudget() *data.NonFinancialBudget {
	return &data.NonFinancialBudget{
		BudetID:                 dto.BudetID,
		OrganizationUnitID:      dto.OrganizationUnitID,
		ImplContactFullName:     dto.ImplContactFullName,
		ImplContactWorkingPlace: dto.ImplContactWorkingPlace,
		ImplContactPhone:        dto.ImplContactPhone,
		ImplContactEmail:        dto.ImplContactEmail,
		ContactFullName:         dto.ContactFullName,
		ContactWorkingPlace:     dto.ContactWorkingPlace,
		ContactPhone:            dto.ContactPhone,
		ContactEmail:            dto.ContactEmail,
	}
}

func ToNonFinancialBudgetResponseDTO(data data.NonFinancialBudget) NonFinancialBudgetResponseDTO {
	return NonFinancialBudgetResponseDTO{
		ID:                      data.ID,
		BudetID:                 data.BudetID,
		OrganizationUnitID:      data.OrganizationUnitID,
		ImplContactFullName:     data.ImplContactFullName,
		ImplContactWorkingPlace: data.ImplContactWorkingPlace,
		ImplContactPhone:        data.ImplContactPhone,
		ImplContactEmail:        data.ImplContactEmail,
		ContactFullName:         data.ContactFullName,
		ContactWorkingPlace:     data.ContactWorkingPlace,
		ContactPhone:            data.ContactPhone,
		ContactEmail:            data.ContactEmail,
		CreatedAt:               data.CreatedAt,
		UpdatedAt:               data.UpdatedAt,
	}
}

func ToNonFinancialBudgetListResponseDTO(nonfinancialbudgets []*data.NonFinancialBudget) []NonFinancialBudgetResponseDTO {
	dtoList := make([]NonFinancialBudgetResponseDTO, len(nonfinancialbudgets))
	for i, x := range nonfinancialbudgets {
		dtoList[i] = ToNonFinancialBudgetResponseDTO(*x)
	}
	return dtoList
}
