package main

import (
	"log"
	"os"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/handlers"
	"gitlab.sudovi.me/erp/finance-api/middleware"
	"gitlab.sudovi.me/erp/finance-api/services"

	"github.com/oykos-development-hub/celeritas"
)

func initApplication() *celeritas.Celeritas {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init celeritas
	cel := &celeritas.Celeritas{}
	err = cel.New(path)
	if err != nil {
		log.Fatal(err)
	}

	cel.AppName = "gitlab.sudovi.me/erp/finance-api"

	models := data.New(cel.DB.Pool)

	ArticleService := services.NewArticleServiceImpl(cel, models.Article)
	ArticleHandler := handlers.NewArticleHandler(cel, ArticleService)

	InvoiceService := services.NewInvoiceServiceImpl(cel, models.Invoice, ArticleService)
	InvoiceHandler := handlers.NewInvoiceHandler(cel, InvoiceService, ArticleService)

	BudgetService := services.NewBudgetServiceImpl(cel, models.Budget)
	BudgetHandler := handlers.NewBudgetHandler(cel, BudgetService)

	FinancialBudgetService := services.NewFinancialBudgetServiceImpl(cel, models.FinancialBudget)
	FinancialBudgetHandler := handlers.NewFinancialBudgetHandler(cel, FinancialBudgetService)

	FinancialBudgetLimitService := services.NewFinancialBudgetLimitServiceImpl(cel, models.FinancialBudgetLimit)
	FinancialBudgetLimitHandler := handlers.NewFinancialBudgetLimitHandler(cel, FinancialBudgetLimitService)

	NonFinancialBudgetService := services.NewNonFinancialBudgetServiceImpl(cel, models.NonFinancialBudget)
	NonFinancialBudgetHandler := handlers.NewNonFinancialBudgetHandler(cel, NonFinancialBudgetService)

	NonFinancialBudgetGoalService := services.NewNonFinancialBudgetGoalServiceImpl(cel, models.NonFinancialBudgetGoal)
	NonFinancialBudgetGoalHandler := handlers.NewNonFinancialBudgetGoalHandler(cel, NonFinancialBudgetGoalService)

	ProgramService := services.NewProgramServiceImpl(cel, models.Program)
	ProgramHandler := handlers.NewProgramHandler(cel, ProgramService)

	ActivityService := services.NewActivityServiceImpl(cel, models.Activity)
	ActivityHandler := handlers.NewActivityHandler(cel, ActivityService)

	GoalIndicatorService := services.NewGoalIndicatorServiceImpl(cel, models.GoalIndicator)
	GoalIndicatorHandler := handlers.NewGoalIndicatorHandler(cel, GoalIndicatorService)

	FilledFinancialBudgetService := services.NewFilledFinancialBudgetServiceImpl(cel, models.FilledFinancialBudget)
	FilledFinancialBudgetHandler := handlers.NewFilledFinancialBudgetHandler(cel, FilledFinancialBudgetService)

	BudgetRequestService := services.NewBudgetRequestServiceImpl(cel, models.BudgetRequest)
	BudgetRequestHandler := handlers.NewBudgetRequestHandler(cel, BudgetRequestService)

	myHandlers := &handlers.Handlers{
		InvoiceHandler: InvoiceHandler,
		ArticleHandler: ArticleHandler,

		BudgetHandler:                 BudgetHandler,
		FinancialBudgetHandler:        FinancialBudgetHandler,
		FinancialBudgetLimitHandler:   FinancialBudgetLimitHandler,
		NonFinancialBudgetHandler:     NonFinancialBudgetHandler,
		NonFinancialBudgetGoalHandler: NonFinancialBudgetGoalHandler,
		ProgramHandler:                ProgramHandler,
		ActivityHandler:               ActivityHandler,
		GoalIndicatorHandler:          GoalIndicatorHandler,
		FilledFinancialBudgetHandler:  FilledFinancialBudgetHandler,
		BudgetRequestHandler:          BudgetRequestHandler,
	}

	myMiddleware := &middleware.Middleware{
		App: cel,
	}

	cel.Routes = routes(cel, myMiddleware, myHandlers)

	return cel
}
