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
	AdditionalExpenseService := services.NewAdditionalExpenseServiceImpl(cel, models.AdditionalExpense, models.Invoice)
	AdditionalExpenseHandler := handlers.NewAdditionalExpenseHandler(cel, AdditionalExpenseService)

	InvoiceService := services.NewInvoiceServiceImpl(cel, models.Invoice, models.AdditionalExpense, ArticleService, AdditionalExpenseService)
	InvoiceHandler := handlers.NewInvoiceHandler(cel, InvoiceService, ArticleService)

	BudgetService := services.NewBudgetServiceImpl(cel, models.Budget)
	BudgetHandler := handlers.NewBudgetHandler(cel, BudgetService)

	SpendingDynamicService := services.NewSpendingDynamicServiceImpl(cel, models.SpendingDynamicEntry, models.CurrentBudget)
	SpendingDynamicHandler := handlers.NewSpendingDynamicHandler(cel, SpendingDynamicService)

	SpendingReleaseService := services.NewSpendingReleaseServiceImpl(cel, models.SpendingRelease, models.CurrentBudget, models.Budget)
	SpendingReleaseHandler := handlers.NewSpendingReleaseHandler(cel, SpendingReleaseService)

	CurrentBudgetService := services.NewCurrentBudgetServiceImpl(cel, models.CurrentBudget, SpendingDynamicService)
	CurrentBudgetHandler := handlers.NewCurrentBudgetHandler(cel, CurrentBudgetService)

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

	FilledFinancialBudgetService := services.NewFilledFinancialBudgetServiceImpl(cel, models.FilledFinancialBudget, models.BudgetRequest, models.CurrentBudget, CurrentBudgetService)
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

	SalaryAdditionalExpenseService := services.NewSalaryAdditionalExpenseServiceImpl(cel, models.SalaryAdditionalExpense)
	SalaryAdditionalExpenseHandler := handlers.NewSalaryAdditionalExpenseHandler(cel, SalaryAdditionalExpenseService)

	SalaryService := services.NewSalaryServiceImpl(cel, models.Salary, models.SalaryAdditionalExpense, SalaryAdditionalExpenseService)
	SalaryHandler := handlers.NewSalaryHandler(cel, SalaryService)

	FixedDepositItemService := services.NewFixedDepositItemServiceImpl(cel, models.FixedDepositItem)
	FixedDepositItemHandler := handlers.NewFixedDepositItemHandler(cel, FixedDepositItemService)

	FixedDepositDispatchService := services.NewFixedDepositDispatchServiceImpl(cel, models.FixedDepositDispatch)
	FixedDepositDispatchHandler := handlers.NewFixedDepositDispatchHandler(cel, FixedDepositDispatchService)

	FixedDepositJudgeService := services.NewFixedDepositJudgeServiceImpl(cel, models.FixedDepositJudge)
	FixedDepositJudgeHandler := handlers.NewFixedDepositJudgeHandler(cel, FixedDepositJudgeService)

	FixedDepositService := services.NewFixedDepositServiceImpl(cel, models.FixedDeposit, FixedDepositItemService, FixedDepositDispatchService, FixedDepositJudgeService)
	FixedDepositHandler := handlers.NewFixedDepositHandler(cel, FixedDepositService)

	FixedDepositWillDispatchService := services.NewFixedDepositWillDispatchServiceImpl(cel, models.FixedDepositWillDispatch)
	FixedDepositWillDispatchHandler := handlers.NewFixedDepositWillDispatchHandler(cel, FixedDepositWillDispatchService)

	FixedDepositWillService := services.NewFixedDepositWillServiceImpl(cel, models.FixedDepositWill, FixedDepositJudgeService, FixedDepositWillDispatchService)
	FixedDepositWillHandler := handlers.NewFixedDepositWillHandler(cel, FixedDepositWillService)

	DepositPaymentService := services.NewDepositPaymentServiceImpl(cel, models.DepositPayment)
	DepositPaymentHandler := handlers.NewDepositPaymentHandler(cel, DepositPaymentService)

	DepositAdditionalExpenseService := services.NewDepositAdditionalExpenseServiceImpl(cel, models.DepositAdditionalExpense, models.DepositPaymentOrder)
	DepositAdditionalExpenseHandler := handlers.NewDepositAdditionalExpenseHandler(cel, DepositAdditionalExpenseService)

	DepositPaymentOrderService := services.NewDepositPaymentOrderServiceImpl(cel, models.DepositPaymentOrder, models.DepositAdditionalExpense, DepositAdditionalExpenseService)
	DepositPaymentOrderHandler := handlers.NewDepositPaymentOrderHandler(cel, DepositPaymentOrderService)

	PaymentOrderItemService := services.NewPaymentOrderItemServiceImpl(cel, models.PaymentOrderItem)
	PaymentOrderItemHandler := handlers.NewPaymentOrderItemHandler(cel, PaymentOrderItemService)

	PaymentOrderService := services.NewPaymentOrderServiceImpl(cel, CurrentBudgetService, models.PaymentOrder, models.PaymentOrderItem, models.Invoice, models.Article, models.AdditionalExpense, models.SalaryAdditionalExpense, models.Salary)
	PaymentOrderHandler := handlers.NewPaymentOrderHandler(cel, PaymentOrderService)

	EnforcedPaymentService := services.NewEnforcedPaymentServiceImpl(cel, models.EnforcedPayment, CurrentBudgetService, models.EnforcedPaymentItem, models.Invoice, models.Article, models.AdditionalExpense, models.PaymentOrderItem, models.PaymentOrder)
	EnforcedPaymentHandler := handlers.NewEnforcedPaymentHandler(cel, EnforcedPaymentService)

	EnforcedPaymentItemService := services.NewEnforcedPaymentItemServiceImpl(cel, models.EnforcedPaymentItem)
	EnforcedPaymentItemHandler := handlers.NewEnforcedPaymentItemHandler(cel, EnforcedPaymentItemService)

	ModelsOfAccountingService := services.NewModelsOfAccountingServiceImpl(cel, models.ModelsOfAccounting, models.ModelOfAccountingItem)
	ModelsOfAccountingHandler := handlers.NewModelsOfAccountingHandler(cel, ModelsOfAccountingService)

	ModelOfAccountingItemService := services.NewModelOfAccountingItemServiceImpl(cel, models.ModelOfAccountingItem)
	ModelOfAccountingItemHandler := handlers.NewModelOfAccountingItemHandler(cel, ModelOfAccountingItemService)

	AccountingEntryService := services.NewAccountingEntryServiceImpl(cel, models.AccountingEntry, models.Invoice, models.Article, models.AdditionalExpense, models.Salary, models.SalaryAdditionalExpense, ModelsOfAccountingService, models.AccountingEntryItem, models.PaymentOrder, models.EnforcedPayment)
	AccountingEntryHandler := handlers.NewAccountingEntryHandler(cel, AccountingEntryService)

	AccountingEntryItemService := services.NewAccountingEntryItemServiceImpl(cel, models.AccountingEntryItem)
	AccountingEntryItemHandler := handlers.NewAccountingEntryItemHandler(cel, AccountingEntryItemService)

	InternalReallocationService := services.NewInternalReallocationServiceImpl(cel, models.InternalReallocation, models.InternalReallocationItem, models.CurrentBudget, models.SpendingDynamicEntry)
	InternalReallocationHandler := handlers.NewInternalReallocationHandler(cel, InternalReallocationService)

	InternalReallocationItemService := services.NewInternalReallocationItemServiceImpl(cel, models.InternalReallocationItem)
	InternalReallocationItemHandler := handlers.NewInternalReallocationItemHandler(cel, InternalReallocationItemService)

	ExternalReallocationService := services.NewExternalReallocationServiceImpl(cel, models.ExternalReallocation, models.ExternalReallocationItem, models.CurrentBudget, models.SpendingDynamicEntry)
	ExternalReallocationHandler := handlers.NewExternalReallocationHandler(cel, ExternalReallocationService)

	ExternalReallocationItemService := services.NewExternalReallocationItemServiceImpl(cel, models.ExternalReallocationItem)
	ExternalReallocationItemHandler := handlers.NewExternalReallocationItemHandler(cel, ExternalReallocationItemService)

	LogService := services.NewLogServiceImpl(cel, models.Log)
	LogHandler := handlers.NewLogHandler(cel, LogService)

		
	SpendingReleaseRequestService := services.NewSpendingReleaseRequestServiceImpl(cel, models.SpendingReleaseRequest)
	SpendingReleaseRequestHandler := handlers.NewSpendingReleaseRequestHandler(cel, SpendingReleaseRequestService)

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
		SpendingReleaseRequestHandler: SpendingReleaseRequestHandler,
	}

	myMiddleware := &middleware.Middleware{
		App: cel,
	}

	cel.Routes = routes(cel, myMiddleware, myHandlers)

	return cel
}
