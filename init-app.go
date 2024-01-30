package main

import (
	"log"
	"os"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/handlers"
	"gitlab.sudovi.me/erp/finance-api/middleware"

	"github.com/oykos-development-hub/celeritas"
	"gitlab.sudovi.me/erp/finance-api/services"
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

	BudgetService := services.NewBudgetServiceImpl(cel, models.Budget)
	BudgetHandler := handlers.NewBudgetHandler(cel, BudgetService)

	FinancialBudgetService := services.NewFinancialBudgetServiceImpl(cel, models.FinancialBudget)
	FinancialBudgetHandler := handlers.NewFinancialBudgetHandler(cel, FinancialBudgetService)

	myHandlers := &handlers.Handlers{
		BudgetHandler:          BudgetHandler,
		FinancialBudgetHandler: FinancialBudgetHandler,
	}

	myMiddleware := &middleware.Middleware{
		App: cel,
	}

	cel.Routes = routes(cel, myMiddleware, myHandlers)

	return cel
}
