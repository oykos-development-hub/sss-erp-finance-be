package handlers

import "net/http"

type Handlers struct {
	InvoiceHandler                  InvoiceHandler
	ArticleHandler                  ArticleHandler
	BudgetHandler                   BudgetHandler
	FinancialBudgetHandler          FinancialBudgetHandler
	FinancialBudgetLimitHandler     FinancialBudgetLimitHandler
	NonFinancialBudgetHandler       NonFinancialBudgetHandler
	NonFinancialBudgetGoalHandler   NonFinancialBudgetGoalHandler
	ProgramHandler                  ProgramHandler
	ActivityHandler                 ActivityHandler
	GoalIndicatorHandler            GoalIndicatorHandler
	FilledFinancialBudgetHandler    FilledFinancialBudgetHandler
	BudgetRequestHandler            BudgetRequestHandler
	FeeHandler                      FeeHandler
	FeePaymentHandler               FeePaymentHandler
	FineHandler                     FineHandler
	FinePaymentHandler              FinePaymentHandler
	ProcedureCostHandler            ProcedureCostHandler
	ProcedureCostPaymentHandler     ProcedureCostPaymentHandler
	AdditionalExpenseHandler        AdditionalExpenseHandler
	FlatRateHandler                 FlatRateHandler
	FlatRatePaymentHandler          FlatRatePaymentHandler
	PropBenConfHandler              PropBenConfHandler
	PropBenConfPaymentHandler       PropBenConfPaymentHandler
	TaxAuthorityCodebookHandler     TaxAuthorityCodebookHandler
	SalaryHandler                   SalaryHandler
	SalaryAdditionalExpenseHandler  SalaryAdditionalExpenseHandler
	FixedDepositHandler             FixedDepositHandler
	FixedDepositItemHandler         FixedDepositItemHandler
	FixedDepositDispatchHandler     FixedDepositDispatchHandler
	FixedDepositJudgeHandler        FixedDepositJudgeHandler
	FixedDepositWillHandler         FixedDepositWillHandler
	FixedDepositWillDispatchHandler FixedDepositWillDispatchHandler
	DepositPaymentHandler           DepositPaymentHandler
	DepositPaymentOrderHandler      DepositPaymentOrderHandler
	DepositAdditionalExpenseHandler DepositAdditionalExpenseHandler
	PaymentOrderHandler             PaymentOrderHandler
	PaymentOrderItemHandler         PaymentOrderItemHandler
}

type InvoiceHandler interface {
	CreateInvoice(w http.ResponseWriter, r *http.Request)
	UpdateInvoice(w http.ResponseWriter, r *http.Request)
	DeleteInvoice(w http.ResponseWriter, r *http.Request)
	GetInvoiceById(w http.ResponseWriter, r *http.Request)
	GetInvoiceList(w http.ResponseWriter, r *http.Request)
}

type ArticleHandler interface {
	CreateArticle(w http.ResponseWriter, r *http.Request)
	UpdateArticle(w http.ResponseWriter, r *http.Request)
	DeleteArticle(w http.ResponseWriter, r *http.Request)
	GetArticleById(w http.ResponseWriter, r *http.Request)
	GetArticleList(w http.ResponseWriter, r *http.Request)
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

type GoalIndicatorHandler interface {
	CreateGoalIndicator(w http.ResponseWriter, r *http.Request)
	UpdateGoalIndicator(w http.ResponseWriter, r *http.Request)
	DeleteGoalIndicator(w http.ResponseWriter, r *http.Request)
	GetGoalIndicatorById(w http.ResponseWriter, r *http.Request)
	GetGoalIndicatorList(w http.ResponseWriter, r *http.Request)
}

type FilledFinancialBudgetHandler interface {
	CreateFilledFinancialBudget(w http.ResponseWriter, r *http.Request)
	UpdateFilledFinancialBudget(w http.ResponseWriter, r *http.Request)
	DeleteFilledFinancialBudget(w http.ResponseWriter, r *http.Request)
	GetFilledFinancialBudgetById(w http.ResponseWriter, r *http.Request)
	GetFilledFinancialBudgetList(w http.ResponseWriter, r *http.Request)
}

type BudgetRequestHandler interface {
	CreateBudgetRequest(w http.ResponseWriter, r *http.Request)
	UpdateBudgetRequest(w http.ResponseWriter, r *http.Request)
	DeleteBudgetRequest(w http.ResponseWriter, r *http.Request)
	GetBudgetRequestById(w http.ResponseWriter, r *http.Request)
	GetBudgetRequestList(w http.ResponseWriter, r *http.Request)
}

// fees
type FeeHandler interface {
	CreateFee(w http.ResponseWriter, r *http.Request)
	GetFeeById(w http.ResponseWriter, r *http.Request)
	UpdateFee(w http.ResponseWriter, r *http.Request)
	DeleteFee(w http.ResponseWriter, r *http.Request)
	GetFeeList(w http.ResponseWriter, r *http.Request)
}

type FeePaymentHandler interface {
	CreateFeePayment(w http.ResponseWriter, r *http.Request)
	DeleteFeePayment(w http.ResponseWriter, r *http.Request)
	UpdateFeePayment(w http.ResponseWriter, r *http.Request)
	GetFeePaymentById(w http.ResponseWriter, r *http.Request)
	GetFeePaymentList(w http.ResponseWriter, r *http.Request)
}

// fines
type FineHandler interface {
	CreateFine(w http.ResponseWriter, r *http.Request)
	GetFineById(w http.ResponseWriter, r *http.Request)
	UpdateFine(w http.ResponseWriter, r *http.Request)
	DeleteFine(w http.ResponseWriter, r *http.Request)
	GetFineList(w http.ResponseWriter, r *http.Request)
}

type FinePaymentHandler interface {
	CreateFinePayment(w http.ResponseWriter, r *http.Request)
	DeleteFinePayment(w http.ResponseWriter, r *http.Request)
	UpdateFinePayment(w http.ResponseWriter, r *http.Request)
	GetFinePaymentById(w http.ResponseWriter, r *http.Request)
	GetFinePaymentList(w http.ResponseWriter, r *http.Request)
}

// procedure costs
type ProcedureCostHandler interface {
	CreateProcedureCost(w http.ResponseWriter, r *http.Request)
	GetProcedureCostById(w http.ResponseWriter, r *http.Request)
	UpdateProcedureCost(w http.ResponseWriter, r *http.Request)
	DeleteProcedureCost(w http.ResponseWriter, r *http.Request)
	GetProcedureCostList(w http.ResponseWriter, r *http.Request)
}

type ProcedureCostPaymentHandler interface {
	CreateProcedureCostPayment(w http.ResponseWriter, r *http.Request)
	DeleteProcedureCostPayment(w http.ResponseWriter, r *http.Request)
	UpdateProcedureCostPayment(w http.ResponseWriter, r *http.Request)
	GetProcedureCostPaymentById(w http.ResponseWriter, r *http.Request)
	GetProcedureCostPaymentList(w http.ResponseWriter, r *http.Request)
}

type AdditionalExpenseHandler interface {
	DeleteAdditionalExpense(w http.ResponseWriter, r *http.Request)
	GetAdditionalExpenseById(w http.ResponseWriter, r *http.Request)
	GetAdditionalExpenseList(w http.ResponseWriter, r *http.Request)
}

// flat rate
type FlatRateHandler interface {
	CreateFlatRate(w http.ResponseWriter, r *http.Request)
	GetFlatRateById(w http.ResponseWriter, r *http.Request)
	UpdateFlatRate(w http.ResponseWriter, r *http.Request)
	DeleteFlatRate(w http.ResponseWriter, r *http.Request)
	GetFlatRateList(w http.ResponseWriter, r *http.Request)
}

type FlatRatePaymentHandler interface {
	CreateFlatRatePayment(w http.ResponseWriter, r *http.Request)
	GetFlatRatePaymentById(w http.ResponseWriter, r *http.Request)
	UpdateFlatRatePayment(w http.ResponseWriter, r *http.Request)
	DeleteFlatRatePayment(w http.ResponseWriter, r *http.Request)
	GetFlatRatePaymentList(w http.ResponseWriter, r *http.Request)
}

// property benefits confiscation
type PropBenConfHandler interface {
	CreatePropBenConf(w http.ResponseWriter, r *http.Request)
	GetPropBenConfById(w http.ResponseWriter, r *http.Request)
	UpdatePropBenConf(w http.ResponseWriter, r *http.Request)
	DeletePropBenConf(w http.ResponseWriter, r *http.Request)
	GetPropBenConfList(w http.ResponseWriter, r *http.Request)
}

type PropBenConfPaymentHandler interface {
	CreatePropBenConfPayment(w http.ResponseWriter, r *http.Request)
	DeletePropBenConfPayment(w http.ResponseWriter, r *http.Request)
	UpdatePropBenConfPayment(w http.ResponseWriter, r *http.Request)
	GetPropBenConfPaymentById(w http.ResponseWriter, r *http.Request)
	GetPropBenConfPaymentList(w http.ResponseWriter, r *http.Request)
}

type TaxAuthorityCodebookHandler interface {
	CreateTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request)
	UpdateTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request)
	DeactivateTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request)
	DeleteTaxAuthorityCodebook(w http.ResponseWriter, r *http.Request)
	GetTaxAuthorityCodebookById(w http.ResponseWriter, r *http.Request)
	GetTaxAuthorityCodebookList(w http.ResponseWriter, r *http.Request)
}

type SalaryHandler interface {
	CreateSalary(w http.ResponseWriter, r *http.Request)
	UpdateSalary(w http.ResponseWriter, r *http.Request)
	DeleteSalary(w http.ResponseWriter, r *http.Request)
	GetSalaryById(w http.ResponseWriter, r *http.Request)
	GetSalaryList(w http.ResponseWriter, r *http.Request)
}

type SalaryAdditionalExpenseHandler interface {
	CreateSalaryAdditionalExpense(w http.ResponseWriter, r *http.Request)
	UpdateSalaryAdditionalExpense(w http.ResponseWriter, r *http.Request)
	DeleteSalaryAdditionalExpense(w http.ResponseWriter, r *http.Request)
	GetSalaryAdditionalExpenseById(w http.ResponseWriter, r *http.Request)
	GetSalaryAdditionalExpenseList(w http.ResponseWriter, r *http.Request)
}

type FixedDepositHandler interface {
	CreateFixedDeposit(w http.ResponseWriter, r *http.Request)
	UpdateFixedDeposit(w http.ResponseWriter, r *http.Request)
	DeleteFixedDeposit(w http.ResponseWriter, r *http.Request)
	GetFixedDepositById(w http.ResponseWriter, r *http.Request)
	GetFixedDepositList(w http.ResponseWriter, r *http.Request)
}

type FixedDepositItemHandler interface {
	CreateFixedDepositItem(w http.ResponseWriter, r *http.Request)
	UpdateFixedDepositItem(w http.ResponseWriter, r *http.Request)
	DeleteFixedDepositItem(w http.ResponseWriter, r *http.Request)
	GetFixedDepositItemById(w http.ResponseWriter, r *http.Request)
	GetFixedDepositItemList(w http.ResponseWriter, r *http.Request)
}

type FixedDepositDispatchHandler interface {
	CreateFixedDepositDispatch(w http.ResponseWriter, r *http.Request)
	UpdateFixedDepositDispatch(w http.ResponseWriter, r *http.Request)
	DeleteFixedDepositDispatch(w http.ResponseWriter, r *http.Request)
	GetFixedDepositDispatchById(w http.ResponseWriter, r *http.Request)
	GetFixedDepositDispatchList(w http.ResponseWriter, r *http.Request)
}

type FixedDepositJudgeHandler interface {
	CreateFixedDepositJudge(w http.ResponseWriter, r *http.Request)
	UpdateFixedDepositJudge(w http.ResponseWriter, r *http.Request)
	DeleteFixedDepositJudge(w http.ResponseWriter, r *http.Request)
	GetFixedDepositJudgeById(w http.ResponseWriter, r *http.Request)
	GetFixedDepositJudgeList(w http.ResponseWriter, r *http.Request)
}

type FixedDepositWillHandler interface {
	CreateFixedDepositWill(w http.ResponseWriter, r *http.Request)
	UpdateFixedDepositWill(w http.ResponseWriter, r *http.Request)
	DeleteFixedDepositWill(w http.ResponseWriter, r *http.Request)
	GetFixedDepositWillById(w http.ResponseWriter, r *http.Request)
	GetFixedDepositWillList(w http.ResponseWriter, r *http.Request)
}

type FixedDepositWillDispatchHandler interface {
	CreateFixedDepositWillDispatch(w http.ResponseWriter, r *http.Request)
	UpdateFixedDepositWillDispatch(w http.ResponseWriter, r *http.Request)
	DeleteFixedDepositWillDispatch(w http.ResponseWriter, r *http.Request)
	GetFixedDepositWillDispatchById(w http.ResponseWriter, r *http.Request)
	GetFixedDepositWillDispatchList(w http.ResponseWriter, r *http.Request)
}

type DepositPaymentHandler interface {
	CreateDepositPayment(w http.ResponseWriter, r *http.Request)
	UpdateDepositPayment(w http.ResponseWriter, r *http.Request)
	DeleteDepositPayment(w http.ResponseWriter, r *http.Request)
	GetDepositPaymentById(w http.ResponseWriter, r *http.Request)
	GetDepositPaymentList(w http.ResponseWriter, r *http.Request)
	GetDepositPaymentsByCaseNumber(w http.ResponseWriter, r *http.Request)
	GetCaseNumber(w http.ResponseWriter, r *http.Request)
}

type DepositPaymentOrderHandler interface {
	CreateDepositPaymentOrder(w http.ResponseWriter, r *http.Request)
	UpdateDepositPaymentOrder(w http.ResponseWriter, r *http.Request)
	PayDepositPaymentOrder(w http.ResponseWriter, r *http.Request)
	DeleteDepositPaymentOrder(w http.ResponseWriter, r *http.Request)
	GetDepositPaymentOrderById(w http.ResponseWriter, r *http.Request)
	GetDepositPaymentOrderList(w http.ResponseWriter, r *http.Request)
}

type DepositAdditionalExpenseHandler interface {
	CreateDepositAdditionalExpense(w http.ResponseWriter, r *http.Request)
	UpdateDepositAdditionalExpense(w http.ResponseWriter, r *http.Request)
	DeleteDepositAdditionalExpense(w http.ResponseWriter, r *http.Request)
	GetDepositAdditionalExpenseById(w http.ResponseWriter, r *http.Request)
	GetDepositAdditionalExpenseList(w http.ResponseWriter, r *http.Request)
}

type PaymentOrderHandler interface {
	CreatePaymentOrder(w http.ResponseWriter, r *http.Request)
	UpdatePaymentOrder(w http.ResponseWriter, r *http.Request)
	DeletePaymentOrder(w http.ResponseWriter, r *http.Request)
	GetPaymentOrderById(w http.ResponseWriter, r *http.Request)
	GetPaymentOrderList(w http.ResponseWriter, r *http.Request)
	GetAllObligations(w http.ResponseWriter, r *http.Request)
}

type PaymentOrderItemHandler interface {
	CreatePaymentOrderItem(w http.ResponseWriter, r *http.Request)
	UpdatePaymentOrderItem(w http.ResponseWriter, r *http.Request)
	DeletePaymentOrderItem(w http.ResponseWriter, r *http.Request)
	GetPaymentOrderItemById(w http.ResponseWriter, r *http.Request)
	GetPaymentOrderItemList(w http.ResponseWriter, r *http.Request)
}
