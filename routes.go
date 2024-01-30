package main

import (
	"gitlab.sudovi.me/erp/finance-api/handlers"
	"gitlab.sudovi.me/erp/finance-api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/oykos-development-hub/celeritas"
)

func routes(app *celeritas.Celeritas, middleware *middleware.Middleware, handlers *handlers.Handlers) *chi.Mux {
	// middleware must come before any routes

	//api
	app.Routes.Route("/api", func(rt chi.Router) {

		rt.Post("/budgets", handlers.BudgetHandler.CreateBudget)
		rt.Get("/budgets/{id}", handlers.BudgetHandler.GetBudgetById)
		rt.Get("/budgets", handlers.BudgetHandler.GetBudgetList)
		rt.Get("/budgets/{id}/financial", handlers.FinancialBudgetHandler.GetFinancialBudgetByBudgetID)
		rt.Put("/budgets/{id}", handlers.BudgetHandler.UpdateBudget)
		rt.Delete("/budgets/{id}", handlers.BudgetHandler.DeleteBudget)

		rt.Post("/financial-budgets", handlers.FinancialBudgetHandler.CreateFinancialBudget)
		rt.Get("/financial-budgets/{id}", handlers.FinancialBudgetHandler.GetFinancialBudgetById)
		rt.Get("/financial-budgets", handlers.FinancialBudgetHandler.GetFinancialBudgetList)
		rt.Put("/financial-budgets/{id}", handlers.FinancialBudgetHandler.UpdateFinancialBudget)
		rt.Delete("/financial-budgets/{id}", handlers.FinancialBudgetHandler.DeleteFinancialBudget)

		rt.Post("/financial-budget-limits", handlers.FinancialBudgetLimitHandler.CreateFinancialBudgetLimit)
		rt.Get("/financial-budget-limits/{id}", handlers.FinancialBudgetLimitHandler.GetFinancialBudgetLimitById)
		rt.Get("/financial-budget-limits", handlers.FinancialBudgetLimitHandler.GetFinancialBudgetLimitList)
		rt.Put("/financial-budget-limits/{id}", handlers.FinancialBudgetLimitHandler.UpdateFinancialBudgetLimit)
		rt.Delete("/financial-budget-limits/{id}", handlers.FinancialBudgetLimitHandler.DeleteFinancialBudgetLimit)
	})

	return app.Routes
}
