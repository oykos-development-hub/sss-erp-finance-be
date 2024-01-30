package services

import (
	"gitlab.sudovi.me/erp/finance-api/dto"
)

type BaseService interface {
	RandomString(n int) string
	Encrypt(text string) (string, error)
	Decrypt(crypto string) (string, error)
}

type BudgetService interface {
	CreateBudget(input dto.BudgetDTO) (*dto.BudgetResponseDTO, error)
	UpdateBudget(id int, input dto.BudgetDTO) (*dto.BudgetResponseDTO, error)
	DeleteBudget(id int) error
	GetBudget(id int) (*dto.BudgetResponseDTO, error)
	GetBudgetList(input dto.GetBudgetListInput) ([]dto.BudgetResponseDTO, error)
}

type FinancialBudgetService interface {
	CreateFinancialBudget(input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error)
	UpdateFinancialBudget(id int, input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error)
	DeleteFinancialBudget(id int) error
	GetFinancialBudget(id int) (*dto.FinancialBudgetResponseDTO, error)
	GetFinancialBudgetList() ([]dto.FinancialBudgetResponseDTO, error)
	GetFinancialBudgetByBudgetID(id int) (*dto.FinancialBudgetResponseDTO, error)
}
