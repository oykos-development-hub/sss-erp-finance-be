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

type FinancialBudgetLimitService interface {
	CreateFinancialBudgetLimit(input dto.FinancialBudgetLimitDTO) (*dto.FinancialBudgetLimitResponseDTO, error)
	UpdateFinancialBudgetLimit(id int, input dto.FinancialBudgetLimitDTO) (*dto.FinancialBudgetLimitResponseDTO, error)
	DeleteFinancialBudgetLimit(id int) error
	GetFinancialBudgetLimit(id int) (*dto.FinancialBudgetLimitResponseDTO, error)
	GetFinancialBudgetLimitList(filter dto.FinancialBudgetLimitFilterDTO) ([]dto.FinancialBudgetLimitResponseDTO, *uint64, error)
}

type NonFinancialBudgetService interface {
	CreateNonFinancialBudget(input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error)
	UpdateNonFinancialBudget(id int, input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error)
	DeleteNonFinancialBudget(id int) error
	GetNonFinancialBudget(id int) (*dto.NonFinancialBudgetResponseDTO, error)
	GetNonFinancialBudgetList(filter dto.NonFinancialBudgetFilterDTO) ([]dto.NonFinancialBudgetResponseDTO, *uint64, error)
}

type NonFinancialBudgetGoalService interface {
	CreateNonFinancialBudgetGoal(input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	UpdateNonFinancialBudgetGoal(id int, input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	DeleteNonFinancialBudgetGoal(id int) error
	GetNonFinancialBudgetGoal(id int) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	GetNonFinancialBudgetGoalList(filter dto.NonFinancialBudgetGoalFilterDTO) ([]dto.NonFinancialBudgetGoalResponseDTO, *uint64, error)
}

type ProgramService interface {
	CreateProgram(input dto.ProgramDTO) (*dto.ProgramResponseDTO, error)
	UpdateProgram(id int, input dto.ProgramDTO) (*dto.ProgramResponseDTO, error)
	DeleteProgram(id int) error
	GetProgram(id int) (*dto.ProgramResponseDTO, error)
	GetProgramList(filter dto.ProgramFilterDTO) ([]dto.ProgramResponseDTO, *uint64, error)
}

type ActivityService interface {
	CreateActivity(input dto.ActivityDTO) (*dto.ActivityResponseDTO, error)
	UpdateActivity(id int, input dto.ActivityDTO) (*dto.ActivityResponseDTO, error)
	DeleteActivity(id int) error
	GetActivity(id int) (*dto.ActivityResponseDTO, error)
	GetActivityList(filter dto.ActivityFilterDTO) ([]dto.ActivityResponseDTO, *uint64, error)
}
