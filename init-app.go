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

	ErrorLogService := services.NewErrorLogServiceImpl(cel, models.ErrorLog)
	ErrorLogHandler := handlers.NewErrorLogHandler(cel, ErrorLogService)

	ArticleService := services.NewArticleServiceImpl(cel, models.Article)
	ArticleHandler := handlers.NewArticleHandler(cel, ArticleService, ErrorLogService)

	//obligations and demands
	AdditionalExpenseService := services.NewAdditionalExpenseServiceImpl(cel, models.AdditionalExpense, models.Invoice)
	AdditionalExpenseHandler := handlers.NewAdditionalExpenseHandler(cel, AdditionalExpenseService, ErrorLogService)

	InvoiceService := services.NewInvoiceServiceImpl(cel, models.Invoice, models.AdditionalExpense, ArticleService, AdditionalExpenseService)
	InvoiceHandler := handlers.NewInvoiceHandler(cel, InvoiceService, ArticleService, ErrorLogService)

	BudgetService := services.NewBudgetServiceImpl(cel, models.Budget)
	BudgetHandler := handlers.NewBudgetHandler(cel, BudgetService, ErrorLogService)

	SpendingDynamicService := services.NewSpendingDynamicServiceImpl(cel, models.SpendingDynamicEntry, models.CurrentBudget)
	SpendingDynamicHandler := handlers.NewSpendingDynamicHandler(cel, SpendingDynamicService, ErrorLogService)

	SpendingReleaseService := services.NewSpendingReleaseServiceImpl(cel, models.SpendingRelease, models.CurrentBudget, models.Budget, models.SpendingReleaseRequest)
	SpendingReleaseHandler := handlers.NewSpendingReleaseHandler(cel, SpendingReleaseService, ErrorLogService)

	CurrentBudgetService := services.NewCurrentBudgetServiceImpl(cel, models.CurrentBudget, SpendingDynamicService)
	CurrentBudgetHandler := handlers.NewCurrentBudgetHandler(cel, CurrentBudgetService, ErrorLogService)

	FinancialBudgetService := services.NewFinancialBudgetServiceImpl(cel, models.FinancialBudget)
	FinancialBudgetHandler := handlers.NewFinancialBudgetHandler(cel, FinancialBudgetService, ErrorLogService)

	FinancialBudgetLimitService := services.NewFinancialBudgetLimitServiceImpl(cel, models.FinancialBudgetLimit)
	FinancialBudgetLimitHandler := handlers.NewFinancialBudgetLimitHandler(cel, FinancialBudgetLimitService, ErrorLogService)

	NonFinancialBudgetService := services.NewNonFinancialBudgetServiceImpl(cel, models.NonFinancialBudget)
	NonFinancialBudgetHandler := handlers.NewNonFinancialBudgetHandler(cel, NonFinancialBudgetService, ErrorLogService)

	NonFinancialBudgetGoalService := services.NewNonFinancialBudgetGoalServiceImpl(cel, models.NonFinancialBudgetGoal)
	NonFinancialBudgetGoalHandler := handlers.NewNonFinancialBudgetGoalHandler(cel, NonFinancialBudgetGoalService, ErrorLogService)

	ProgramService := services.NewProgramServiceImpl(cel, models.Program)
	ProgramHandler := handlers.NewProgramHandler(cel, ProgramService, ErrorLogService)

	ActivityService := services.NewActivityServiceImpl(cel, models.Activity)
	ActivityHandler := handlers.NewActivityHandler(cel, ActivityService, ErrorLogService)

	GoalIndicatorService := services.NewGoalIndicatorServiceImpl(cel, models.GoalIndicator)
	GoalIndicatorHandler := handlers.NewGoalIndicatorHandler(cel, GoalIndicatorService, ErrorLogService)

	FilledFinancialBudgetService := services.NewFilledFinancialBudgetServiceImpl(cel, models.FilledFinancialBudget, models.BudgetRequest, models.CurrentBudget, CurrentBudgetService)
	FilledFinancialBudgetHandler := handlers.NewFilledFinancialBudgetHandler(cel, FilledFinancialBudgetService, ErrorLogService)

	BudgetRequestService := services.NewBudgetRequestServiceImpl(cel, models.BudgetRequest)
	BudgetRequestHandler := handlers.NewBudgetRequestHandler(cel, BudgetRequestService, ErrorLogService)

	// fees

	FeeSharedLogicService := services.NewFeeSharedLogicServiceImpl(cel, models.Fee, models.FeePayment)

	FeePaymentService := services.NewFeePaymentServiceImpl(cel, models.FeePayment, FeeSharedLogicService)
	FeePaymentHandler := handlers.NewFeePaymentHandler(cel, FeePaymentService, ErrorLogService)

	FeeService := services.NewFeeServiceImpl(cel, models.Fee, FeeSharedLogicService)
	FeeHandler := handlers.NewFeeHandler(cel, FeeService, ErrorLogService)

	// fines
	FineSharedLogicService := services.NewFineSharedLogicServiceImpl(cel, models.Fine, models.FinePayment)

	FinePaymentService := services.NewFinePaymentServiceImpl(cel, models.FinePayment, FineSharedLogicService)
	FinePaymentHandler := handlers.NewFinePaymentHandler(cel, FinePaymentService, ErrorLogService)

	FineService := services.NewFineServiceImpl(cel, models.Fine, FineSharedLogicService)
	FineHandler := handlers.NewFineHandler(cel, FineService, ErrorLogService)

	// procedure cost
	ProcedureCostSharedLogicService := services.NewProcedureCostSharedLogicServiceImpl(cel, models.ProcedureCost, models.ProcedureCostPayment)

	ProcedureCostPaymentService := services.NewProcedureCostPaymentServiceImpl(cel, models.ProcedureCostPayment, ProcedureCostSharedLogicService)
	ProcedureCostPaymentHandler := handlers.NewProcedureCostPaymentHandler(cel, ProcedureCostPaymentService, ErrorLogService)

	ProcedureCostService := services.NewProcedureCostServiceImpl(cel, models.ProcedureCost, ProcedureCostSharedLogicService)
	ProcedureCostHandler := handlers.NewProcedureCostHandler(cel, ProcedureCostService, ErrorLogService)

	FlatRateSharedLogicService := services.NewFlatRateSharedLogicServiceImpl(cel, models.FlatRate, models.FlatRatePayment)

	FlatRatePaymentService := services.NewFlatRatePaymentServiceImpl(cel, models.FlatRatePayment, FlatRateSharedLogicService)
	FlatRatePaymentHandler := handlers.NewFlatRatePaymentHandler(cel, FlatRatePaymentService, ErrorLogService)

	FlatRateService := services.NewFlatRateServiceImpl(cel, models.FlatRate, FlatRateSharedLogicService)
	FlatRateHandler := handlers.NewFlatRateHandler(cel, FlatRateService, ErrorLogService)

	// property benefits confiscation
	PropBenConfSharedLogicService := services.NewPropBenConfSharedLogicServiceImpl(cel, models.PropBenConf, models.PropBenConfPayment)

	PropBenConfService := services.NewPropBenConfServiceImpl(cel, models.PropBenConf, PropBenConfSharedLogicService)
	PropBenConfHandler := handlers.NewPropBenConfHandler(cel, PropBenConfService, ErrorLogService)

	PropBenConfPaymentService := services.NewPropBenConfPaymentServiceImpl(cel, models.PropBenConfPayment, PropBenConfSharedLogicService)
	PropBenConfPaymentHandler := handlers.NewPropBenConfPaymentHandler(cel, PropBenConfPaymentService, ErrorLogService)

	TaxAuthorityCodebookService := services.NewTaxAuthorityCodebookServiceImpl(cel, models.TaxAuthorityCodebook)
	TaxAuthorityCodebookHandler := handlers.NewTaxAuthorityCodebookHandler(cel, TaxAuthorityCodebookService, ErrorLogService)

	SalaryAdditionalExpenseService := services.NewSalaryAdditionalExpenseServiceImpl(cel, models.SalaryAdditionalExpense)
	SalaryAdditionalExpenseHandler := handlers.NewSalaryAdditionalExpenseHandler(cel, SalaryAdditionalExpenseService, ErrorLogService)

	SalaryService := services.NewSalaryServiceImpl(cel, models.Salary, models.SalaryAdditionalExpense, SalaryAdditionalExpenseService)
	SalaryHandler := handlers.NewSalaryHandler(cel, SalaryService, ErrorLogService)

	FixedDepositItemService := services.NewFixedDepositItemServiceImpl(cel, models.FixedDepositItem)
	FixedDepositItemHandler := handlers.NewFixedDepositItemHandler(cel, FixedDepositItemService, ErrorLogService)

	FixedDepositDispatchService := services.NewFixedDepositDispatchServiceImpl(cel, models.FixedDepositDispatch)
	FixedDepositDispatchHandler := handlers.NewFixedDepositDispatchHandler(cel, FixedDepositDispatchService, ErrorLogService)

	FixedDepositJudgeService := services.NewFixedDepositJudgeServiceImpl(cel, models.FixedDepositJudge)
	FixedDepositJudgeHandler := handlers.NewFixedDepositJudgeHandler(cel, FixedDepositJudgeService, ErrorLogService)

	FixedDepositService := services.NewFixedDepositServiceImpl(cel, models.FixedDeposit, FixedDepositItemService, FixedDepositDispatchService, FixedDepositJudgeService)
	FixedDepositHandler := handlers.NewFixedDepositHandler(cel, FixedDepositService, ErrorLogService)

	FixedDepositWillDispatchService := services.NewFixedDepositWillDispatchServiceImpl(cel, models.FixedDepositWillDispatch)
	FixedDepositWillDispatchHandler := handlers.NewFixedDepositWillDispatchHandler(cel, FixedDepositWillDispatchService, ErrorLogService)

	FixedDepositWillService := services.NewFixedDepositWillServiceImpl(cel, models.FixedDepositWill, FixedDepositJudgeService, FixedDepositWillDispatchService)
	FixedDepositWillHandler := handlers.NewFixedDepositWillHandler(cel, FixedDepositWillService, ErrorLogService)

	DepositPaymentService := services.NewDepositPaymentServiceImpl(cel, models.DepositPayment)
	DepositPaymentHandler := handlers.NewDepositPaymentHandler(cel, DepositPaymentService, ErrorLogService)

	DepositAdditionalExpenseService := services.NewDepositAdditionalExpenseServiceImpl(cel, models.DepositAdditionalExpense, models.DepositPaymentOrder)
	DepositAdditionalExpenseHandler := handlers.NewDepositAdditionalExpenseHandler(cel, DepositAdditionalExpenseService, ErrorLogService)

	DepositPaymentOrderService := services.NewDepositPaymentOrderServiceImpl(cel, models.DepositPaymentOrder, models.DepositAdditionalExpense, DepositAdditionalExpenseService)
	DepositPaymentOrderHandler := handlers.NewDepositPaymentOrderHandler(cel, DepositPaymentOrderService, ErrorLogService)

	PaymentOrderItemService := services.NewPaymentOrderItemServiceImpl(cel, models.PaymentOrderItem)
	PaymentOrderItemHandler := handlers.NewPaymentOrderItemHandler(cel, PaymentOrderItemService, ErrorLogService)

	PaymentOrderService := services.NewPaymentOrderServiceImpl(cel, CurrentBudgetService, models.PaymentOrder, models.PaymentOrderItem, models.Invoice, models.Article, models.AdditionalExpense, models.SalaryAdditionalExpense, models.Salary)
	PaymentOrderHandler := handlers.NewPaymentOrderHandler(cel, PaymentOrderService, ErrorLogService)

	EnforcedPaymentService := services.NewEnforcedPaymentServiceImpl(cel, models.EnforcedPayment, CurrentBudgetService, models.EnforcedPaymentItem, models.Invoice, models.Article, models.AdditionalExpense, models.PaymentOrderItem, models.PaymentOrder)
	EnforcedPaymentHandler := handlers.NewEnforcedPaymentHandler(cel, EnforcedPaymentService, ErrorLogService)

	EnforcedPaymentItemService := services.NewEnforcedPaymentItemServiceImpl(cel, models.EnforcedPaymentItem)
	EnforcedPaymentItemHandler := handlers.NewEnforcedPaymentItemHandler(cel, EnforcedPaymentItemService, ErrorLogService)

	ModelsOfAccountingService := services.NewModelsOfAccountingServiceImpl(cel, models.ModelsOfAccounting, models.ModelOfAccountingItem)
	ModelsOfAccountingHandler := handlers.NewModelsOfAccountingHandler(cel, ModelsOfAccountingService, ErrorLogService)

	ModelOfAccountingItemService := services.NewModelOfAccountingItemServiceImpl(cel, models.ModelOfAccountingItem)
	ModelOfAccountingItemHandler := handlers.NewModelOfAccountingItemHandler(cel, ModelOfAccountingItemService, ErrorLogService)

	AccountingEntryService := services.NewAccountingEntryServiceImpl(cel, models.AccountingEntry, models.Invoice, models.Article, models.AdditionalExpense, models.Salary, models.SalaryAdditionalExpense, ModelsOfAccountingService, models.AccountingEntryItem, models.PaymentOrder, models.EnforcedPayment)
	AccountingEntryHandler := handlers.NewAccountingEntryHandler(cel, AccountingEntryService, ErrorLogService)

	AccountingEntryItemService := services.NewAccountingEntryItemServiceImpl(cel, models.AccountingEntryItem)
	AccountingEntryItemHandler := handlers.NewAccountingEntryItemHandler(cel, AccountingEntryItemService, ErrorLogService)

	InternalReallocationService := services.NewInternalReallocationServiceImpl(cel, models.InternalReallocation, models.InternalReallocationItem, models.CurrentBudget, models.SpendingDynamicEntry)
	InternalReallocationHandler := handlers.NewInternalReallocationHandler(cel, InternalReallocationService, ErrorLogService)

	InternalReallocationItemService := services.NewInternalReallocationItemServiceImpl(cel, models.InternalReallocationItem)
	InternalReallocationItemHandler := handlers.NewInternalReallocationItemHandler(cel, InternalReallocationItemService, ErrorLogService)

	ExternalReallocationService := services.NewExternalReallocationServiceImpl(cel, models.ExternalReallocation, models.ExternalReallocationItem, models.CurrentBudget, models.SpendingDynamicEntry)
	ExternalReallocationHandler := handlers.NewExternalReallocationHandler(cel, ExternalReallocationService, ErrorLogService)

	ExternalReallocationItemService := services.NewExternalReallocationItemServiceImpl(cel, models.ExternalReallocationItem)
	ExternalReallocationItemHandler := handlers.NewExternalReallocationItemHandler(cel, ExternalReallocationItemService, ErrorLogService)

	LogService := services.NewLogServiceImpl(cel, models.Log)
	LogHandler := handlers.NewLogHandler(cel, LogService, ErrorLogService)

	SpendingReleaseRequestService := services.NewSpendingReleaseRequestServiceImpl(cel, models.SpendingReleaseRequest)
	SpendingReleaseRequestHandler := handlers.NewSpendingReleaseRequestHandler(cel, SpendingReleaseRequestService, ErrorLogService)

	myHandlers := &handlers.Handlers{
		InvoiceHandler:                  InvoiceHandler,
		ArticleHandler:                  ArticleHandler,
		BudgetHandler:                   BudgetHandler,
		FinancialBudgetHandler:          FinancialBudgetHandler,
		FinancialBudgetLimitHandler:     FinancialBudgetLimitHandler,
		NonFinancialBudgetHandler:       NonFinancialBudgetHandler,
		NonFinancialBudgetGoalHandler:   NonFinancialBudgetGoalHandler,
		ProgramHandler:                  ProgramHandler,
		ActivityHandler:                 ActivityHandler,
		GoalIndicatorHandler:            GoalIndicatorHandler,
		FilledFinancialBudgetHandler:    FilledFinancialBudgetHandler,
		BudgetRequestHandler:            BudgetRequestHandler,
		FeeHandler:                      FeeHandler,
		FeePaymentHandler:               FeePaymentHandler,
		FineHandler:                     FineHandler,
		FinePaymentHandler:              FinePaymentHandler,
		ProcedureCostHandler:            ProcedureCostHandler,
		ProcedureCostPaymentHandler:     ProcedureCostPaymentHandler,
		FlatRateHandler:                 FlatRateHandler,
		FlatRatePaymentHandler:          FlatRatePaymentHandler,
		AdditionalExpenseHandler:        AdditionalExpenseHandler,
		PropBenConfHandler:              PropBenConfHandler,
		PropBenConfPaymentHandler:       PropBenConfPaymentHandler,
		TaxAuthorityCodebookHandler:     TaxAuthorityCodebookHandler,
		SalaryHandler:                   SalaryHandler,
		SalaryAdditionalExpenseHandler:  SalaryAdditionalExpenseHandler,
		FixedDepositHandler:             FixedDepositHandler,
		FixedDepositItemHandler:         FixedDepositItemHandler,
		FixedDepositDispatchHandler:     FixedDepositDispatchHandler,
		FixedDepositJudgeHandler:        FixedDepositJudgeHandler,
		FixedDepositWillHandler:         FixedDepositWillHandler,
		FixedDepositWillDispatchHandler: FixedDepositWillDispatchHandler,
		DepositPaymentHandler:           DepositPaymentHandler,
		DepositPaymentOrderHandler:      DepositPaymentOrderHandler,
		DepositAdditionalExpenseHandler: DepositAdditionalExpenseHandler,
		PaymentOrderHandler:             PaymentOrderHandler,
		PaymentOrderItemHandler:         PaymentOrderItemHandler,
		EnforcedPaymentHandler:          EnforcedPaymentHandler,
		EnforcedPaymentItemHandler:      EnforcedPaymentItemHandler,
		AccountingEntryHandler:          AccountingEntryHandler,
		ModelsOfAccountingHandler:       ModelsOfAccountingHandler,
		ModelOfAccountingItemHandler:    ModelOfAccountingItemHandler,
		AccountingEntryItemHandler:      AccountingEntryItemHandler,
		SpendingDynamicHandler:          SpendingDynamicHandler,
		CurrentBudgetHandler:            CurrentBudgetHandler,
		SpendingReleaseHandler:          SpendingReleaseHandler,
		InternalReallocationHandler:     InternalReallocationHandler,
		InternalReallocationItemHandler: InternalReallocationItemHandler,
		ExternalReallocationHandler:     ExternalReallocationHandler,
		ExternalReallocationItemHandler: ExternalReallocationItemHandler,
		LogHandler:                      LogHandler,
		SpendingReleaseRequestHandler:   SpendingReleaseRequestHandler,
		ErrorLogHandler:                 ErrorLogHandler,
	}

	myMiddleware := &middleware.Middleware{
		App: cel,
	}

	cel.Routes = routes(cel, myMiddleware, myHandlers)

	SpendingReleaseService.StartMonthlyTaskForSpendingReleases()

	return cel
}
