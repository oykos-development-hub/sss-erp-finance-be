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
		rt.Patch("/filled-financial-budgets/{id}/actual", handlers.FilledFinancialBudgetHandler.UpdateActualFilledFinancialBudget)
		rt.Delete("/filled-financial-budgets/{id}", handlers.FilledFinancialBudgetHandler.DeleteFilledFinancialBudget)
		rt.Get("/budgets/{id}/filled-financial-summary/{req_type}", handlers.FilledFinancialBudgetHandler.GetSummaryFilledFinancialRequests)

		rt.Post("/budget-requests", handlers.BudgetRequestHandler.CreateBudgetRequest)
		rt.Get("/budget-requests/{id}", handlers.BudgetRequestHandler.GetBudgetRequestById)
		rt.Get("/budget-requests", handlers.BudgetRequestHandler.GetBudgetRequestList)
		rt.Put("/budget-requests/{id}", handlers.BudgetRequestHandler.UpdateBudgetRequest)
		rt.Delete("/budget-requests/{id}", handlers.BudgetRequestHandler.DeleteBudgetRequest)

		rt.Post("/fees", handlers.FeeHandler.CreateFee)
		rt.Get("/fees/{id}", handlers.FeeHandler.GetFeeById)
		rt.Get("/fees", handlers.FeeHandler.GetFeeList)
		rt.Put("/fees/{id}", handlers.FeeHandler.UpdateFee)
		rt.Delete("/fees/{id}", handlers.FeeHandler.DeleteFee)

		rt.Post("/fee-payments", handlers.FeePaymentHandler.CreateFeePayment)
		rt.Get("/fee-payments/{id}", handlers.FeePaymentHandler.GetFeePaymentById)
		rt.Get("/fee-payments", handlers.FeePaymentHandler.GetFeePaymentList)
		rt.Put("/fee-payments/{id}", handlers.FeePaymentHandler.UpdateFeePayment)
		rt.Delete("/fee-payments/{id}", handlers.FeePaymentHandler.DeleteFeePayment)

		rt.Post("/fines", handlers.FineHandler.CreateFine)
		rt.Get("/fines/{id}", handlers.FineHandler.GetFineById)
		rt.Get("/fines", handlers.FineHandler.GetFineList)
		rt.Put("/fines/{id}", handlers.FineHandler.UpdateFine)
		rt.Delete("/fines/{id}", handlers.FineHandler.DeleteFine)

		rt.Post("/fine-payments", handlers.FinePaymentHandler.CreateFinePayment)
		rt.Get("/fine-payments/{id}", handlers.FinePaymentHandler.GetFinePaymentById)
		rt.Get("/fine-payments", handlers.FinePaymentHandler.GetFinePaymentList)
		rt.Put("/fine-payments/{id}", handlers.FinePaymentHandler.UpdateFinePayment)
		rt.Delete("/fine-payments/{id}", handlers.FinePaymentHandler.DeleteFinePayment)

		rt.Post("/procedure-costs", handlers.ProcedureCostHandler.CreateProcedureCost)
		rt.Get("/procedure-costs/{id}", handlers.ProcedureCostHandler.GetProcedureCostById)
		rt.Get("/procedure-costs", handlers.ProcedureCostHandler.GetProcedureCostList)
		rt.Put("/procedure-costs/{id}", handlers.ProcedureCostHandler.UpdateProcedureCost)
		rt.Delete("/procedure-costs/{id}", handlers.ProcedureCostHandler.DeleteProcedureCost)

		rt.Post("/procedure-cost-payments", handlers.ProcedureCostPaymentHandler.CreateProcedureCostPayment)
		rt.Get("/procedure-cost-payments/{id}", handlers.ProcedureCostPaymentHandler.GetProcedureCostPaymentById)
		rt.Get("/procedure-cost-payments", handlers.ProcedureCostPaymentHandler.GetProcedureCostPaymentList)
		rt.Put("/procedure-cost-payments/{id}", handlers.ProcedureCostPaymentHandler.UpdateProcedureCostPayment)
		rt.Delete("/procedure-cost-payments/{id}", handlers.ProcedureCostPaymentHandler.DeleteProcedureCostPayment)

		rt.Get("/additional-expenses/{id}", handlers.AdditionalExpenseHandler.GetAdditionalExpenseById)
		rt.Get("/additional-expenses", handlers.AdditionalExpenseHandler.GetAdditionalExpenseList)
		rt.Delete("/additional-expenses/{id}", handlers.AdditionalExpenseHandler.DeleteAdditionalExpense)

		rt.Post("/flat-rates", handlers.FlatRateHandler.CreateFlatRate)
		rt.Get("/flat-rates/{id}", handlers.FlatRateHandler.GetFlatRateById)
		rt.Get("/flat-rates", handlers.FlatRateHandler.GetFlatRateList)
		rt.Put("/flat-rates/{id}", handlers.FlatRateHandler.UpdateFlatRate)
		rt.Delete("/flat-rates/{id}", handlers.FlatRateHandler.DeleteFlatRate)

		rt.Post("/flat-rate-payments", handlers.FlatRatePaymentHandler.CreateFlatRatePayment)
		rt.Get("/flat-rate-payments/{id}", handlers.FlatRatePaymentHandler.GetFlatRatePaymentById)
		rt.Get("/flat-rate-payments", handlers.FlatRatePaymentHandler.GetFlatRatePaymentList)
		rt.Put("/flat-rate-payments/{id}", handlers.FlatRatePaymentHandler.UpdateFlatRatePayment)
		rt.Delete("/flat-rate-payments/{id}", handlers.FlatRatePaymentHandler.DeleteFlatRatePayment)

		rt.Post("/property-benefits-confiscations", handlers.PropBenConfHandler.CreatePropBenConf)
		rt.Get("/property-benefits-confiscations/{id}", handlers.PropBenConfHandler.GetPropBenConfById)
		rt.Get("/property-benefits-confiscations", handlers.PropBenConfHandler.GetPropBenConfList)
		rt.Put("/property-benefits-confiscations/{id}", handlers.PropBenConfHandler.UpdatePropBenConf)
		rt.Delete("/property-benefits-confiscations/{id}", handlers.PropBenConfHandler.DeletePropBenConf)

		rt.Post("/property-benefits-confiscation-payments", handlers.PropBenConfPaymentHandler.CreatePropBenConfPayment)
		rt.Get("/property-benefits-confiscation-payments/{id}", handlers.PropBenConfPaymentHandler.GetPropBenConfPaymentById)
		rt.Get("/property-benefits-confiscation-payments", handlers.PropBenConfPaymentHandler.GetPropBenConfPaymentList)
		rt.Put("/property-benefits-confiscation-payments/{id}", handlers.PropBenConfPaymentHandler.UpdatePropBenConfPayment)
		rt.Delete("/property-benefits-confiscation-payments/{id}", handlers.PropBenConfPaymentHandler.DeletePropBenConfPayment)

		rt.Post("/tax-authority-codebooks", handlers.TaxAuthorityCodebookHandler.CreateTaxAuthorityCodebook)
		rt.Get("/tax-authority-codebooks/{id}", handlers.TaxAuthorityCodebookHandler.GetTaxAuthorityCodebookById)
		rt.Get("/tax-authority-codebooks", handlers.TaxAuthorityCodebookHandler.GetTaxAuthorityCodebookList)
		rt.Put("/tax-authority-codebooks/{id}", handlers.TaxAuthorityCodebookHandler.UpdateTaxAuthorityCodebook)
		rt.Put("/tax-authority-codebook-deactivate/{id}", handlers.TaxAuthorityCodebookHandler.DeactivateTaxAuthorityCodebook)
		rt.Delete("/tax-authority-codebooks/{id}", handlers.TaxAuthorityCodebookHandler.DeleteTaxAuthorityCodebook)

		rt.Post("/salaries", handlers.SalaryHandler.CreateSalary)
		rt.Get("/salaries/{id}", handlers.SalaryHandler.GetSalaryById)
		rt.Get("/salaries", handlers.SalaryHandler.GetSalaryList)
		rt.Put("/salaries/{id}", handlers.SalaryHandler.UpdateSalary)
		rt.Delete("/salaries/{id}", handlers.SalaryHandler.DeleteSalary)

		rt.Post("/salary-additional-expenses", handlers.SalaryAdditionalExpenseHandler.CreateSalaryAdditionalExpense)
		rt.Get("/salary-additional-expenses/{id}", handlers.SalaryAdditionalExpenseHandler.GetSalaryAdditionalExpenseById)
		rt.Get("/salary-additional-expenses", handlers.SalaryAdditionalExpenseHandler.GetSalaryAdditionalExpenseList)
		rt.Put("/salary-additional-expenses/{id}", handlers.SalaryAdditionalExpenseHandler.UpdateSalaryAdditionalExpense)
		rt.Delete("/salary-additional-expenses/{id}", handlers.SalaryAdditionalExpenseHandler.DeleteSalaryAdditionalExpense)

		rt.Post("/fixed-deposits", handlers.FixedDepositHandler.CreateFixedDeposit)
		rt.Get("/fixed-deposits/{id}", handlers.FixedDepositHandler.GetFixedDepositById)
		rt.Get("/fixed-deposits", handlers.FixedDepositHandler.GetFixedDepositList)
		rt.Put("/fixed-deposits/{id}", handlers.FixedDepositHandler.UpdateFixedDeposit)
		rt.Delete("/fixed-deposits/{id}", handlers.FixedDepositHandler.DeleteFixedDeposit)

		rt.Post("/fixed-deposit-items", handlers.FixedDepositItemHandler.CreateFixedDepositItem)
		rt.Get("/fixed-deposit-items/{id}", handlers.FixedDepositItemHandler.GetFixedDepositItemById)
		rt.Get("/fixed-deposit-items", handlers.FixedDepositItemHandler.GetFixedDepositItemList)
		rt.Put("/fixed-deposit-items/{id}", handlers.FixedDepositItemHandler.UpdateFixedDepositItem)
		rt.Delete("/fixed-deposit-items/{id}", handlers.FixedDepositItemHandler.DeleteFixedDepositItem)

		rt.Post("/fixed-deposit-dispatches", handlers.FixedDepositDispatchHandler.CreateFixedDepositDispatch)
		rt.Get("/fixed-deposit-dispatches/{id}", handlers.FixedDepositDispatchHandler.GetFixedDepositDispatchById)
		rt.Get("/fixed-deposit-dispatches", handlers.FixedDepositDispatchHandler.GetFixedDepositDispatchList)
		rt.Put("/fixed-deposit-dispatches/{id}", handlers.FixedDepositDispatchHandler.UpdateFixedDepositDispatch)
		rt.Delete("/fixed-deposit-dispatches/{id}", handlers.FixedDepositDispatchHandler.DeleteFixedDepositDispatch)

		rt.Post("/fixed-deposit-judges", handlers.FixedDepositJudgeHandler.CreateFixedDepositJudge)
		rt.Get("/fixed-deposit-judges/{id}", handlers.FixedDepositJudgeHandler.GetFixedDepositJudgeById)
		rt.Get("/fixed-deposit-judges", handlers.FixedDepositJudgeHandler.GetFixedDepositJudgeList)
		rt.Put("/fixed-deposit-judges/{id}", handlers.FixedDepositJudgeHandler.UpdateFixedDepositJudge)
		rt.Delete("/fixed-deposit-judges/{id}", handlers.FixedDepositJudgeHandler.DeleteFixedDepositJudge)

		rt.Post("/fixed-deposit-wills", handlers.FixedDepositWillHandler.CreateFixedDepositWill)
		rt.Get("/fixed-deposit-wills/{id}", handlers.FixedDepositWillHandler.GetFixedDepositWillById)
		rt.Get("/fixed-deposit-wills", handlers.FixedDepositWillHandler.GetFixedDepositWillList)
		rt.Put("/fixed-deposit-wills/{id}", handlers.FixedDepositWillHandler.UpdateFixedDepositWill)
		rt.Delete("/fixed-deposit-wills/{id}", handlers.FixedDepositWillHandler.DeleteFixedDepositWill)

		rt.Post("/fixed-deposit-will-dispatches", handlers.FixedDepositWillDispatchHandler.CreateFixedDepositWillDispatch)
		rt.Get("/fixed-deposit-will-dispatches/{id}", handlers.FixedDepositWillDispatchHandler.GetFixedDepositWillDispatchById)
		rt.Get("/fixed-deposit-will-dispatches", handlers.FixedDepositWillDispatchHandler.GetFixedDepositWillDispatchList)
		rt.Put("/fixed-deposit-will-dispatches/{id}", handlers.FixedDepositWillDispatchHandler.UpdateFixedDepositWillDispatch)
		rt.Delete("/fixed-deposit-will-dispatches/{id}", handlers.FixedDepositWillDispatchHandler.DeleteFixedDepositWillDispatch)

		rt.Post("/deposit-payments", handlers.DepositPaymentHandler.CreateDepositPayment)
		rt.Get("/deposit-payments/{id}", handlers.DepositPaymentHandler.GetDepositPaymentById)
		rt.Get("/deposit-payments", handlers.DepositPaymentHandler.GetDepositPaymentList)
		rt.Get("/deposit-payments-case-number", handlers.DepositPaymentHandler.GetDepositPaymentsByCaseNumber)
		rt.Get("/get-case-number", handlers.DepositPaymentHandler.GetCaseNumber)
		rt.Get("/get-initial-state", handlers.DepositPaymentHandler.GetInitialState)
		rt.Put("/deposit-payments/{id}", handlers.DepositPaymentHandler.UpdateDepositPayment)
		rt.Delete("/deposit-payments/{id}", handlers.DepositPaymentHandler.DeleteDepositPayment)

		rt.Post("/deposit-payment-orders", handlers.DepositPaymentOrderHandler.CreateDepositPaymentOrder)
		rt.Get("/deposit-payment-orders/{id}", handlers.DepositPaymentOrderHandler.GetDepositPaymentOrderById)
		rt.Get("/deposit-payment-orders", handlers.DepositPaymentOrderHandler.GetDepositPaymentOrderList)
		rt.Put("/deposit-payment-orders/{id}", handlers.DepositPaymentOrderHandler.UpdateDepositPaymentOrder)
		rt.Put("/pay-deposit-payment-order/{id}", handlers.DepositPaymentOrderHandler.PayDepositPaymentOrder)
		rt.Delete("/deposit-payment-orders/{id}", handlers.DepositPaymentOrderHandler.DeleteDepositPaymentOrder)

		rt.Post("/deposit-additional-expenses", handlers.DepositAdditionalExpenseHandler.CreateDepositAdditionalExpense)
		rt.Get("/deposit-additional-expenses/{id}", handlers.DepositAdditionalExpenseHandler.GetDepositAdditionalExpenseById)
		rt.Get("/deposit-additional-expenses", handlers.DepositAdditionalExpenseHandler.GetDepositAdditionalExpenseList)
		rt.Put("/deposit-additional-expenses/{id}", handlers.DepositAdditionalExpenseHandler.UpdateDepositAdditionalExpense)
		rt.Delete("/deposit-additional-expenses/{id}", handlers.DepositAdditionalExpenseHandler.DeleteDepositAdditionalExpense)

		rt.Post("/payment-orders", handlers.PaymentOrderHandler.CreatePaymentOrder)
		rt.Get("/payment-orders/{id}", handlers.PaymentOrderHandler.GetPaymentOrderById)
		rt.Get("/payment-orders", handlers.PaymentOrderHandler.GetPaymentOrderList)
		rt.Get("/get-all-obligations", handlers.PaymentOrderHandler.GetAllObligations)
		rt.Put("/payment-orders/{id}", handlers.PaymentOrderHandler.UpdatePaymentOrder)
		rt.Put("/pay-payment-order/{id}", handlers.PaymentOrderHandler.PayPaymentOrder)
		rt.Delete("/payment-orders/{id}", handlers.PaymentOrderHandler.DeletePaymentOrder)

		rt.Post("/payment-order-items", handlers.PaymentOrderItemHandler.CreatePaymentOrderItem)
		rt.Get("/payment-order-items/{id}", handlers.PaymentOrderItemHandler.GetPaymentOrderItemById)
		rt.Get("/payment-order-items", handlers.PaymentOrderItemHandler.GetPaymentOrderItemList)
		rt.Put("/payment-order-items/{id}", handlers.PaymentOrderItemHandler.UpdatePaymentOrderItem)
		rt.Delete("/payment-order-items/{id}", handlers.PaymentOrderItemHandler.DeletePaymentOrderItem)

		rt.Post("/enforced-payments", handlers.EnforcedPaymentHandler.CreateEnforcedPayment)
		rt.Get("/enforced-payments/{id}", handlers.EnforcedPaymentHandler.GetEnforcedPaymentById)
		rt.Get("/enforced-payments", handlers.EnforcedPaymentHandler.GetEnforcedPaymentList)
		rt.Put("/enforced-payments/{id}", handlers.EnforcedPaymentHandler.UpdateEnforcedPayment)
		rt.Put("/return-enforced-payment/{id}", handlers.EnforcedPaymentHandler.ReturnEnforcedPayment)
		rt.Delete("/enforced-payments/{id}", handlers.EnforcedPaymentHandler.DeleteEnforcedPayment)

		rt.Post("/enforced-payment-items", handlers.EnforcedPaymentItemHandler.CreateEnforcedPaymentItem)
		rt.Get("/enforced-payment-items/{id}", handlers.EnforcedPaymentItemHandler.GetEnforcedPaymentItemById)
		rt.Get("/enforced-payment-items", handlers.EnforcedPaymentItemHandler.GetEnforcedPaymentItemList)
		rt.Put("/enforced-payment-items/{id}", handlers.EnforcedPaymentItemHandler.UpdateEnforcedPaymentItem)
		rt.Delete("/enforced-payment-items/{id}", handlers.EnforcedPaymentItemHandler.DeleteEnforcedPaymentItem)

		rt.Post("/accounting-entries", handlers.AccountingEntryHandler.CreateAccountingEntry)
		rt.Get("/accounting-entries/{id}", handlers.AccountingEntryHandler.GetAccountingEntryById)
		rt.Get("/accounting-entries", handlers.AccountingEntryHandler.GetAccountingEntryList)
		rt.Put("/accounting-entries/{id}", handlers.AccountingEntryHandler.UpdateAccountingEntry)
		rt.Delete("/accounting-entries/{id}", handlers.AccountingEntryHandler.DeleteAccountingEntry)

		rt.Get("/get-obligations-for-accounting", handlers.AccountingEntryHandler.GetObligationsForAccounting)
		rt.Get("/build-accounting-order-for-obligations", handlers.AccountingEntryHandler.BuildAccountingOrderForObligations)

		rt.Post("/models-of-accountings", handlers.ModelsOfAccountingHandler.CreateModelsOfAccounting)
		rt.Put("/models-of-accountings/{id}", handlers.ModelsOfAccountingHandler.UpdateModelsOfAccounting)
		rt.Get("/models-of-accountings/{id}", handlers.ModelsOfAccountingHandler.GetModelsOfAccountingById)
		rt.Get("/models-of-accountings", handlers.ModelsOfAccountingHandler.GetModelsOfAccountingList)

		rt.Post("/model-of-accounting-items", handlers.ModelOfAccountingItemHandler.CreateModelOfAccountingItem)
		rt.Get("/model-of-accounting-items/{id}", handlers.ModelOfAccountingItemHandler.GetModelOfAccountingItemById)
		rt.Get("/model-of-accounting-items", handlers.ModelOfAccountingItemHandler.GetModelOfAccountingItemList)
		rt.Put("/model-of-accounting-items/{id}", handlers.ModelOfAccountingItemHandler.UpdateModelOfAccountingItem)
	})

	return app.Routes
}
