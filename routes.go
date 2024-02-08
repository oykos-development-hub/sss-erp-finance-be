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

		rt.Post("/invoices", handlers.InvoiceHandler.CreateInvoice)
		rt.Get("/invoices/{id}", handlers.InvoiceHandler.GetInvoiceById)
		rt.Get("/invoices", handlers.InvoiceHandler.GetInvoiceList)
		rt.Put("/invoices/{id}", handlers.InvoiceHandler.UpdateInvoice)
		rt.Delete("/invoices/{id}", handlers.InvoiceHandler.DeleteInvoice)

		rt.Post("/articles", handlers.ArticleHandler.CreateArticle)
		rt.Get("/articles/{id}", handlers.ArticleHandler.GetArticleById)
		rt.Get("/articles", handlers.ArticleHandler.GetArticleList)
		rt.Put("/articles/{id}", handlers.ArticleHandler.UpdateArticle)
		rt.Delete("/articles/{id}", handlers.ArticleHandler.DeleteArticle)

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

		rt.Post("/non-financial-budgets", handlers.NonFinancialBudgetHandler.CreateNonFinancialBudget)
		rt.Get("/non-financial-budgets/{id}", handlers.NonFinancialBudgetHandler.GetNonFinancialBudgetById)
		rt.Get("/non-financial-budgets", handlers.NonFinancialBudgetHandler.GetNonFinancialBudgetList)
		rt.Put("/non-financial-budgets/{id}", handlers.NonFinancialBudgetHandler.UpdateNonFinancialBudget)
		rt.Delete("/non-financial-budgets/{id}", handlers.NonFinancialBudgetHandler.DeleteNonFinancialBudget)

		rt.Post("/non-financial-budget-goals", handlers.NonFinancialBudgetGoalHandler.CreateNonFinancialBudgetGoal)
		rt.Get("/non-financial-budget-goals/{id}", handlers.NonFinancialBudgetGoalHandler.GetNonFinancialBudgetGoalById)
		rt.Get("/non-financial-budget-goals", handlers.NonFinancialBudgetGoalHandler.GetNonFinancialBudgetGoalList)
		rt.Put("/non-financial-budget-goals/{id}", handlers.NonFinancialBudgetGoalHandler.UpdateNonFinancialBudgetGoal)
		rt.Delete("/non-financial-budget-goals/{id}", handlers.NonFinancialBudgetGoalHandler.DeleteNonFinancialBudgetGoal)

		rt.Post("/programs", handlers.ProgramHandler.CreateProgram)
		rt.Get("/programs/{id}", handlers.ProgramHandler.GetProgramById)
		rt.Get("/programs", handlers.ProgramHandler.GetProgramList)
		rt.Put("/programs/{id}", handlers.ProgramHandler.UpdateProgram)
		rt.Delete("/programs/{id}", handlers.ProgramHandler.DeleteProgram)

		rt.Post("/activities", handlers.ActivityHandler.CreateActivity)
		rt.Get("/activities/{id}", handlers.ActivityHandler.GetActivityById)
		rt.Get("/activities", handlers.ActivityHandler.GetActivityList)
		rt.Put("/activities/{id}", handlers.ActivityHandler.UpdateActivity)
		rt.Delete("/activities/{id}", handlers.ActivityHandler.DeleteActivity)

		rt.Post("/goal-indicators", handlers.GoalIndicatorHandler.CreateGoalIndicator)
		rt.Get("/goal-indicators/{id}", handlers.GoalIndicatorHandler.GetGoalIndicatorById)
		rt.Get("/goal-indicators", handlers.GoalIndicatorHandler.GetGoalIndicatorList)
		rt.Put("/goal-indicators/{id}", handlers.GoalIndicatorHandler.UpdateGoalIndicator)
		rt.Delete("/goal-indicators/{id}", handlers.GoalIndicatorHandler.DeleteGoalIndicator)

		rt.Post("/filled-financial-budgets", handlers.FilledFinancialBudgetHandler.CreateFilledFinancialBudget)
		rt.Get("/filled-financial-budgets/{id}", handlers.FilledFinancialBudgetHandler.GetFilledFinancialBudgetById)
		rt.Get("/filled-financial-budgets", handlers.FilledFinancialBudgetHandler.GetFilledFinancialBudgetList)
		rt.Put("/filled-financial-budgets/{id}", handlers.FilledFinancialBudgetHandler.UpdateFilledFinancialBudget)
		rt.Delete("/filled-financial-budgets/{id}", handlers.FilledFinancialBudgetHandler.DeleteFilledFinancialBudget)
	})

	return app.Routes
}
