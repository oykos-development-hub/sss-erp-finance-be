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

	//obligations and demands
	AdditionalExpenseService := services.NewAdditionalExpenseServiceImpl(cel, models.AdditionalExpense)
	AdditionalExpenseHandler := handlers.NewAdditionalExpenseHandler(cel, AdditionalExpenseService)

	InvoiceService := services.NewInvoiceServiceImpl(cel, models.Invoice, ArticleService, AdditionalExpenseService)
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

	// fees

	FeeSharedLogicService := services.NewFeeSharedLogicServiceImpl(cel, models.Fee, models.FeePayment)

	FeePaymentService := services.NewFeePaymentServiceImpl(cel, models.FeePayment, FeeSharedLogicService)
	FeePaymentHandler := handlers.NewFeePaymentHandler(cel, FeePaymentService)

	FeeService := services.NewFeeServiceImpl(cel, models.Fee, FeeSharedLogicService)
	FeeHandler := handlers.NewFeeHandler(cel, FeeService)

	// fines
	FineSharedLogicService := services.NewFineSharedLogicServiceImpl(cel, models.Fine, models.FinePayment)

	FinePaymentService := services.NewFinePaymentServiceImpl(cel, models.FinePayment, FineSharedLogicService)
	FinePaymentHandler := handlers.NewFinePaymentHandler(cel, FinePaymentService)

	FineService := services.NewFineServiceImpl(cel, models.Fine, FineSharedLogicService)
	FineHandler := handlers.NewFineHandler(cel, FineService)

	// procedure cost
	ProcedureCostSharedLogicService := services.NewProcedureCostSharedLogicServiceImpl(cel, models.ProcedureCost, models.ProcedureCostPayment)

	ProcedureCostPaymentService := services.NewProcedureCostPaymentServiceImpl(cel, models.ProcedureCostPayment, ProcedureCostSharedLogicService)
	ProcedureCostPaymentHandler := handlers.NewProcedureCostPaymentHandler(cel, ProcedureCostPaymentService)

	ProcedureCostService := services.NewProcedureCostServiceImpl(cel, models.ProcedureCost, ProcedureCostSharedLogicService)
	ProcedureCostHandler := handlers.NewProcedureCostHandler(cel, ProcedureCostService)

	FlatRateSharedLogicService := services.NewFlatRateSharedLogicServiceImpl(cel, models.FlatRate, models.FlatRatePayment)

	FlatRatePaymentService := services.NewFlatRatePaymentServiceImpl(cel, models.FlatRatePayment, FlatRateSharedLogicService)
	FlatRatePaymentHandler := handlers.NewFlatRatePaymentHandler(cel, FlatRatePaymentService)

	FlatRateService := services.NewFlatRateServiceImpl(cel, models.FlatRate, FlatRateSharedLogicService)
	FlatRateHandler := handlers.NewFlatRateHandler(cel, FlatRateService)

	// property benefits confiscation
	PropBenConfSharedLogicService := services.NewPropBenConfSharedLogicServiceImpl(cel, models.PropBenConf, models.PropBenConfPayment)

	PropBenConfService := services.NewPropBenConfServiceImpl(cel, models.PropBenConf, PropBenConfSharedLogicService)
	PropBenConfHandler := handlers.NewPropBenConfHandler(cel, PropBenConfService)

	PropBenConfPaymentService := services.NewPropBenConfPaymentServiceImpl(cel, models.PropBenConfPayment, PropBenConfSharedLogicService)
	PropBenConfPaymentHandler := handlers.NewPropBenConfPaymentHandler(cel, PropBenConfPaymentService)

	TaxAuthorityCodebookService := services.NewTaxAuthorityCodebookServiceImpl(cel, models.TaxAuthorityCodebook)
	TaxAuthorityCodebookHandler := handlers.NewTaxAuthorityCodebookHandler(cel, TaxAuthorityCodebookService)

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

		FeeHandler:        FeeHandler,
		FeePaymentHandler: FeePaymentHandler,

		FineHandler:        FineHandler,
		FinePaymentHandler: FinePaymentHandler,

		ProcedureCostHandler:        ProcedureCostHandler,
		ProcedureCostPaymentHandler: ProcedureCostPaymentHandler,
		FlatRateHandler:             FlatRateHandler,
		FlatRatePaymentHandler:      FlatRatePaymentHandler,
		AdditionalExpenseHandler:    AdditionalExpenseHandler,

		PropBenConfHandler:          PropBenConfHandler,
		PropBenConfPaymentHandler:   PropBenConfPaymentHandler,
		TaxAuthorityCodebookHandler: TaxAuthorityCodebookHandler,
	}

	myMiddleware := &middleware.Middleware{
		App: cel,
	}

	cel.Routes = routes(cel, myMiddleware, myHandlers)

	return cel
}
