package handlers

import "net/http"

type Handlers struct {
	InvoiceHandler                InvoiceHandler
	ArticleHandler                ArticleHandler
	BudgetHandler                 BudgetHandler
	FinancialBudgetHandler        FinancialBudgetHandler
	FinancialBudgetLimitHandler   FinancialBudgetLimitHandler
	NonFinancialBudgetHandler     NonFinancialBudgetHandler
	NonFinancialBudgetGoalHandler NonFinancialBudgetGoalHandler
	ProgramHandler                ProgramHandler
	ActivityHandler               ActivityHandler
	GoalIndicatorHandler          GoalIndicatorHandler
	FilledFinancialBudgetHandler  FilledFinancialBudgetHandler
	BudgetRequestHandler          BudgetRequestHandler

	FeeHandler        FeeHandler
	FeePaymentHandler FeePaymentHandler

	FineHandler        FineHandler
	FinePaymentHandler FinePaymentHandler
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
