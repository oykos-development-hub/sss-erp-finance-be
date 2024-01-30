package handlers

import (
	"net/http"
)

type Handlers struct {
	BudgetHandler          BudgetHandler
	FinancialBudgetHandler FinancialBudgetHandler
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
