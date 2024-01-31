package handlers

import (
	"net/http"
)

type Handlers struct {
	BudgetHandler                 BudgetHandler
	FinancialBudgetHandler        FinancialBudgetHandler
	FinancialBudgetLimitHandler   FinancialBudgetLimitHandler
	NonFinancialBudgetHandler     NonFinancialBudgetHandler
	NonFinancialBudgetGoalHandler NonFinancialBudgetGoalHandler
	ProgramHandler                ProgramHandler
	ActivityHandler               ActivityHandler
}

type BudgetHandler interface {
	CreateBudget(w http.ResponseWriter, r *http.Request)
	UpdateBudget(w http.ResponseWriter, r *http.Request)
	DeleteBudget(w http.ResponseWriter, r *http.Request)
	GetBudgetById(w http.ResponseWriter, r *http.Request)
	GetBudgetList(w http.ResponseWriter, r *http.Request)
}

type FinancialBudgetHandler interface {
	CreateFinancialBudget(w http.ResponseWriter, r *http.Request)
	UpdateFinancialBudget(w http.ResponseWriter, r *http.Request)
	DeleteFinancialBudget(w http.ResponseWriter, r *http.Request)
	GetFinancialBudgetById(w http.ResponseWriter, r *http.Request)
	GetFinancialBudgetByBudgetID(w http.ResponseWriter, r *http.Request)
	GetFinancialBudgetList(w http.ResponseWriter, r *http.Request)
}

type FinancialBudgetLimitHandler interface {
	CreateFinancialBudgetLimit(w http.ResponseWriter, r *http.Request)
	UpdateFinancialBudgetLimit(w http.ResponseWriter, r *http.Request)
	DeleteFinancialBudgetLimit(w http.ResponseWriter, r *http.Request)
	GetFinancialBudgetLimitById(w http.ResponseWriter, r *http.Request)
	GetFinancialBudgetLimitList(w http.ResponseWriter, r *http.Request)
}

type NonFinancialBudgetHandler interface {
	CreateNonFinancialBudget(w http.ResponseWriter, r *http.Request)
	UpdateNonFinancialBudget(w http.ResponseWriter, r *http.Request)
	DeleteNonFinancialBudget(w http.ResponseWriter, r *http.Request)
	GetNonFinancialBudgetById(w http.ResponseWriter, r *http.Request)
	GetNonFinancialBudgetList(w http.ResponseWriter, r *http.Request)
}

type NonFinancialBudgetGoalHandler interface {
	CreateNonFinancialBudgetGoal(w http.ResponseWriter, r *http.Request)
	UpdateNonFinancialBudgetGoal(w http.ResponseWriter, r *http.Request)
	DeleteNonFinancialBudgetGoal(w http.ResponseWriter, r *http.Request)
	GetNonFinancialBudgetGoalById(w http.ResponseWriter, r *http.Request)
	GetNonFinancialBudgetGoalList(w http.ResponseWriter, r *http.Request)
}

type ProgramHandler interface {
	CreateProgram(w http.ResponseWriter, r *http.Request)
	UpdateProgram(w http.ResponseWriter, r *http.Request)
	DeleteProgram(w http.ResponseWriter, r *http.Request)
	GetProgramById(w http.ResponseWriter, r *http.Request)
	GetProgramList(w http.ResponseWriter, r *http.Request)
}

type ActivityHandler interface {
	CreateActivity(w http.ResponseWriter, r *http.Request)
	UpdateActivity(w http.ResponseWriter, r *http.Request)
	DeleteActivity(w http.ResponseWriter, r *http.Request)
	GetActivityById(w http.ResponseWriter, r *http.Request)
	GetActivityList(w http.ResponseWriter, r *http.Request)
}
