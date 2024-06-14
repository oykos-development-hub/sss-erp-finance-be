package services

import (
	"context"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
)

type BaseService interface {
	RandomString(n int) string
	Encrypt(text string) (string, error)
	Decrypt(crypto string) (string, error)
}

type InvoiceService interface {
	CreateInvoice(ctx context.Context, input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error)
	UpdateInvoice(ctx context.Context, id int, input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error)
	DeleteInvoice(ctx context.Context, id int) error
	GetInvoice(id int) (*dto.InvoiceResponseDTO, error)
	GetInvoiceList(input dto.InvoicesFilter) ([]dto.InvoiceResponseDTO, *uint64, error)
}

type ArticleService interface {
	CreateArticle(input dto.ArticleDTO) (*dto.ArticleResponseDTO, error)
	UpdateArticle(id int, input dto.ArticleDTO) (*dto.ArticleResponseDTO, error)
	DeleteArticle(id int) error
	GetArticle(id int) (*dto.ArticleResponseDTO, error)
	GetArticleList(filter dto.ArticleFilterDTO) ([]dto.ArticleResponseDTO, *uint64, error)
}

type BudgetService interface {
	CreateBudget(ctx context.Context, input dto.BudgetDTO) (*dto.BudgetResponseDTO, error)
	UpdateBudget(ctx context.Context, id int, input dto.BudgetDTO) (*dto.BudgetResponseDTO, error)
	DeleteBudget(ctx context.Context, id int) error
	GetBudget(id int) (*dto.BudgetResponseDTO, error)
	GetBudgetList(input dto.GetBudgetListInput) ([]dto.BudgetResponseDTO, error)
}

type FinancialBudgetService interface {
	CreateFinancialBudget(ctx context.Context, input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error)
	UpdateFinancialBudget(ctx context.Context, id int, input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error)
	DeleteFinancialBudget(ctx context.Context, id int) error
	GetFinancialBudget(id int) (*dto.FinancialBudgetResponseDTO, error)
	GetFinancialBudgetList() ([]dto.FinancialBudgetResponseDTO, error)
	GetFinancialBudgetByBudgetID(id int) (*dto.FinancialBudgetResponseDTO, error)
}

type FinancialBudgetLimitService interface {
	CreateFinancialBudgetLimit(input dto.FinancialBudgetLimitDTO) (*dto.FinancialBudgetLimitResponseDTO, error)
	UpdateFinancialBudgetLimit(id int, input dto.FinancialBudgetLimitDTO) (*dto.FinancialBudgetLimitResponseDTO, error)
	DeleteFinancialBudgetLimit(id int) error
	GetFinancialBudgetLimit(id int) (*dto.FinancialBudgetLimitResponseDTO, error)
	GetFinancialBudgetLimitList(filter dto.FinancialBudgetLimitFilterDTO) ([]dto.FinancialBudgetLimitResponseDTO, *uint64, error)
}

type NonFinancialBudgetService interface {
	CreateNonFinancialBudget(ctx context.Context, input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error)
	UpdateNonFinancialBudget(ctx context.Context, id int, input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error)
	DeleteNonFinancialBudget(ctx context.Context, id int) error
	GetNonFinancialBudget(id int) (*dto.NonFinancialBudgetResponseDTO, error)
	GetNonFinancialBudgetList(filter dto.NonFinancialBudgetFilterDTO) ([]dto.NonFinancialBudgetResponseDTO, *uint64, error)
}

type NonFinancialBudgetGoalService interface {
	CreateNonFinancialBudgetGoal(ctx context.Context, input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	UpdateNonFinancialBudgetGoal(ctx context.Context, id int, input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	DeleteNonFinancialBudgetGoal(ctx context.Context, id int) error
	GetNonFinancialBudgetGoal(id int) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	GetNonFinancialBudgetGoalList(filter dto.NonFinancialBudgetGoalFilterDTO) ([]dto.NonFinancialBudgetGoalResponseDTO, *uint64, error)
}

type ProgramService interface {
	CreateProgram(ctx context.Context, input dto.ProgramDTO) (*dto.ProgramResponseDTO, error)
	UpdateProgram(ctx context.Context, id int, input dto.ProgramDTO) (*dto.ProgramResponseDTO, error)
	DeleteProgram(ctx context.Context, id int) error
	GetProgram(id int) (*dto.ProgramResponseDTO, error)
	GetProgramList(filter dto.ProgramFilterDTO) ([]dto.ProgramResponseDTO, *uint64, error)
}

type ActivityService interface {
	CreateActivity(ctx context.Context, input dto.ActivityDTO) (*dto.ActivityResponseDTO, error)
	UpdateActivity(ctx context.Context, id int, input dto.ActivityDTO) (*dto.ActivityResponseDTO, error)
	DeleteActivity(ctx context.Context, id int) error
	GetActivity(id int) (*dto.ActivityResponseDTO, error)
	GetActivityList(filter dto.ActivityFilterDTO) ([]dto.ActivityResponseDTO, *uint64, error)
}

type GoalIndicatorService interface {
	CreateGoalIndicator(ctx context.Context, input dto.GoalIndicatorDTO) (*dto.GoalIndicatorResponseDTO, error)
	UpdateGoalIndicator(ctx context.Context, id int, input dto.GoalIndicatorDTO) (*dto.GoalIndicatorResponseDTO, error)
	DeleteGoalIndicator(ctx context.Context, id int) error
	GetGoalIndicator(id int) (*dto.GoalIndicatorResponseDTO, error)
	GetGoalIndicatorList(filter dto.GoalIndicatorFilterDTO) ([]dto.GoalIndicatorResponseDTO, *uint64, error)
}

type FilledFinancialBudgetService interface {
	CreateFilledFinancialBudget(ctx context.Context, input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error)
	UpdateFilledFinancialBudget(ctx context.Context, id int, input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error)
	UpdateActualFilledFinancialBudget(ctx context.Context, id int, amount decimal.Decimal) (*dto.FilledFinancialBudgetResponseDTO, error)
	DeleteFilledFinancialBudget(ctx context.Context, id int) error
	GetFilledFinancialBudget(id int) (*dto.FilledFinancialBudgetResponseDTO, error)
	GetFilledFinancialBudgetList(filter dto.FilledFinancialBudgetFilterDTO) ([]dto.FilledFinancialBudgetResponseDTO, *uint64, error)
	GetSummaryFilledFinancialRequests(budgetID int, requestType data.RequestType) ([]dto.FilledFinancialBudgetResponseDTO, error)
}

type BudgetRequestService interface {
	CreateBudgetRequest(ctx context.Context, input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error)
	UpdateBudgetRequest(ctx context.Context, id int, input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error)
	DeleteBudgetRequest(ctx context.Context, id int) error
	GetBudgetRequest(id int) (*dto.BudgetRequestResponseDTO, error)
	GetBudgetRequestList(filter dto.BudgetRequestFilterDTO) ([]dto.BudgetRequestResponseDTO, *uint64, error)
}

type FeeSharedLogicService interface {
	CalculateFeeDetailsAndUpdateStatus(ctx context.Context, pfineId int) (*dto.FeeDetailsDTO, data.FeeStatus, error)
}

type FeeService interface {
	CreateFee(ctx context.Context, input dto.FeeDTO) (*dto.FeeResponseDTO, error)
	GetFee(id int) (*dto.FeeResponseDTO, error)
	GetFeeList(filter dto.FeeFilterDTO) ([]dto.FeeResponseDTO, *uint64, error)
	UpdateFee(ctx context.Context, id int, input dto.FeeDTO) (*dto.FeeResponseDTO, error)
	DeleteFee(ctx context.Context, id int) error
}

type FeePaymentService interface {
	CreateFeePayment(ctx context.Context, input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error)
	DeleteFeePayment(ctx context.Context, id int) error
	UpdateFeePayment(ctx context.Context, id int, input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error)
	GetFeePayment(id int) (*dto.FeePaymentResponseDTO, error)
	GetFeePaymentList(filter dto.FeePaymentFilterDTO) ([]dto.FeePaymentResponseDTO, *uint64, error)
}

type FineSharedLogicService interface {
	CalculateFineDetailsAndUpdateStatus(ctx context.Context, fineId int) (*dto.FineFeeDetailsDTO, data.FineStatus, error)
}

type FineService interface {
	CreateFine(ctx context.Context, input dto.FineDTO) (*dto.FineResponseDTO, error)
	GetFine(id int) (*dto.FineResponseDTO, error)
	GetFineList(filter dto.FineFilterDTO) ([]dto.FineResponseDTO, *uint64, error)
	UpdateFine(ctx context.Context, id int, input dto.FineDTO) (*dto.FineResponseDTO, error)
	DeleteFine(ctx context.Context, id int) error
}

type FinePaymentService interface {
	CreateFinePayment(ctx context.Context, input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error)
	DeleteFinePayment(ctx context.Context, id int) error
	UpdateFinePayment(ctx context.Context, id int, input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error)
	GetFinePayment(id int) (*dto.FinePaymentResponseDTO, error)
	GetFinePaymentList(filter dto.FinePaymentFilterDTO) ([]dto.FinePaymentResponseDTO, *uint64, error)
}

type ProcedureCostSharedLogicService interface {
	CalculateProcedureCostDetailsAndUpdateStatus(ctx context.Context, procedureCostId int) (*dto.ProcedureCostDetailsDTO, data.ProcedureCostStatus, error)
}

type ProcedureCostService interface {
	CreateProcedureCost(ctx context.Context, input dto.ProcedureCostDTO) (*dto.ProcedureCostResponseDTO, error)
	GetProcedureCost(id int) (*dto.ProcedureCostResponseDTO, error)
	GetProcedureCostList(filter dto.ProcedureCostFilterDTO) ([]dto.ProcedureCostResponseDTO, *uint64, error)
	UpdateProcedureCost(ctx context.Context, id int, input dto.ProcedureCostDTO) (*dto.ProcedureCostResponseDTO, error)
	DeleteProcedureCost(ctx context.Context, id int) error
}

type ProcedureCostPaymentService interface {
	CreateProcedureCostPayment(ctx context.Context, input dto.ProcedureCostPaymentDTO) (*dto.ProcedureCostPaymentResponseDTO, error)
	DeleteProcedureCostPayment(ctx context.Context, id int) error
	UpdateProcedureCostPayment(ctx context.Context, id int, input dto.ProcedureCostPaymentDTO) (*dto.ProcedureCostPaymentResponseDTO, error)
	GetProcedureCostPayment(id int) (*dto.ProcedureCostPaymentResponseDTO, error)
	GetProcedureCostPaymentList(filter dto.ProcedureCostPaymentFilterDTO) ([]dto.ProcedureCostPaymentResponseDTO, *uint64, error)
}

type AdditionalExpenseService interface {
	DeleteAdditionalExpense(id int) error
	GetAdditionalExpense(id int) (*dto.AdditionalExpenseResponseDTO, error)
	GetAdditionalExpenseList(filter dto.AdditionalExpenseFilterDTO) ([]dto.AdditionalExpenseResponseDTO, *uint64, error)
}

type FlatRateSharedLogicService interface {
	CalculateFlatRateDetailsAndUpdateStatus(ctx context.Context, flatRateId int) (*dto.FlatRateDetailsDTO, data.FlatRateStatus, error)
}

type FlatRateService interface {
	CreateFlatRate(ctx context.Context, input dto.FlatRateDTO) (*dto.FlatRateResponseDTO, error)
	GetFlatRate(id int) (*dto.FlatRateResponseDTO, error)
	GetFlatRateList(filter dto.FlatRateFilterDTO) ([]dto.FlatRateResponseDTO, *uint64, error)
	UpdateFlatRate(ctx context.Context, id int, input dto.FlatRateDTO) (*dto.FlatRateResponseDTO, error)
	DeleteFlatRate(ctx context.Context, id int) error
}

type FlatRatePaymentService interface {
	CreateFlatRatePayment(ctx context.Context, input dto.FlatRatePaymentDTO) (*dto.FlatRatePaymentResponseDTO, error)
	DeleteFlatRatePayment(ctx context.Context, id int) error
	UpdateFlatRatePayment(ctx context.Context, id int, input dto.FlatRatePaymentDTO) (*dto.FlatRatePaymentResponseDTO, error)
	GetFlatRatePayment(id int) (*dto.FlatRatePaymentResponseDTO, error)
	GetFlatRatePaymentList(filter dto.FlatRatePaymentFilterDTO) ([]dto.FlatRatePaymentResponseDTO, *uint64, error)
}

type PropBenConfSharedLogicService interface {
	CalculatePropBenConfDetailsAndUpdateStatus(ctx context.Context, propbenconfId int) (*dto.PropBenConfDetailsDTO, data.PropBenConfStatus, error)
}

type PropBenConfService interface {
	CreatePropBenConf(ctx context.Context, input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error)
	GetPropBenConf(id int) (*dto.PropBenConfResponseDTO, error)
	GetPropBenConfList(filter dto.PropBenConfFilterDTO) ([]dto.PropBenConfResponseDTO, *uint64, error)
	UpdatePropBenConf(ctx context.Context, id int, input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error)
	DeletePropBenConf(ctx context.Context, id int) error
}

type PropBenConfPaymentService interface {
	CreatePropBenConfPayment(ctx context.Context, input dto.PropBenConfPaymentDTO) (*dto.PropBenConfPaymentResponseDTO, error)
	DeletePropBenConfPayment(ctx context.Context, id int) error
	UpdatePropBenConfPayment(ctx context.Context, id int, input dto.PropBenConfPaymentDTO) (*dto.PropBenConfPaymentResponseDTO, error)
	GetPropBenConfPayment(id int) (*dto.PropBenConfPaymentResponseDTO, error)
	GetPropBenConfPaymentList(filter dto.PropBenConfPaymentFilterDTO) ([]dto.PropBenConfPaymentResponseDTO, *uint64, error)
}

type TaxAuthorityCodebookService interface {
	CreateTaxAuthorityCodebook(ctx context.Context, input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error)
	UpdateTaxAuthorityCodebook(ctx context.Context, id int, input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error)
	DeleteTaxAuthorityCodebook(ctx context.Context, id int) error
	DeactivateTaxAuthorityCodebook(ctx context.Context, id int, active bool) error
	GetTaxAuthorityCodebook(id int) (*dto.TaxAuthorityCodebookResponseDTO, error)
	GetTaxAuthorityCodebookList(filter dto.TaxAuthorityCodebookFilterDTO) ([]dto.TaxAuthorityCodebookResponseDTO, *uint64, error)
}

type SalaryService interface {
	CreateSalary(ctx context.Context, input dto.SalaryDTO) (*dto.SalaryResponseDTO, error)
	UpdateSalary(ctx context.Context, id int, input dto.SalaryDTO) (*dto.SalaryResponseDTO, error)
	DeleteSalary(ctx context.Context, id int) error
	GetSalary(id int) (*dto.SalaryResponseDTO, error)
	GetSalaryList(filter dto.SalaryFilterDTO) ([]dto.SalaryResponseDTO, *uint64, error)
}

type SalaryAdditionalExpenseService interface {
	CreateSalaryAdditionalExpense(input dto.SalaryAdditionalExpenseDTO) (*dto.SalaryAdditionalExpenseResponseDTO, error)
	UpdateSalaryAdditionalExpense(id int, input dto.SalaryAdditionalExpenseDTO) (*dto.SalaryAdditionalExpenseResponseDTO, error)
	DeleteSalaryAdditionalExpense(id int) error
	GetSalaryAdditionalExpense(id int) (*dto.SalaryAdditionalExpenseResponseDTO, error)
	GetSalaryAdditionalExpenseList(filter dto.SalaryAdditionalExpenseFilterDTO) ([]dto.SalaryAdditionalExpenseResponseDTO, *uint64, error)
}

type FixedDepositService interface {
	CreateFixedDeposit(ctx context.Context, input dto.FixedDepositDTO) (*dto.FixedDepositResponseDTO, error)
	UpdateFixedDeposit(ctx context.Context, id int, input dto.FixedDepositDTO) (*dto.FixedDepositResponseDTO, error)
	DeleteFixedDeposit(ctx context.Context, id int) error
	GetFixedDeposit(id int) (*dto.FixedDepositResponseDTO, error)
	GetFixedDepositList(filter dto.FixedDepositFilterDTO) ([]dto.FixedDepositResponseDTO, *uint64, error)
}

type FixedDepositItemService interface {
	CreateFixedDepositItem(ctx context.Context, input dto.FixedDepositItemDTO) (*dto.FixedDepositItemResponseDTO, error)
	UpdateFixedDepositItem(ctx context.Context, id int, input dto.FixedDepositItemDTO) (*dto.FixedDepositItemResponseDTO, error)
	DeleteFixedDepositItem(ctx context.Context, id int) error
	GetFixedDepositItem(id int) (*dto.FixedDepositItemResponseDTO, error)
	GetFixedDepositItemList(filter dto.FixedDepositItemFilterDTO) ([]dto.FixedDepositItemResponseDTO, *uint64, error)
}

type FixedDepositDispatchService interface {
	CreateFixedDepositDispatch(ctx context.Context, input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error)
	UpdateFixedDepositDispatch(ctx context.Context, id int, input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error)
	DeleteFixedDepositDispatch(ctx context.Context, id int) error
	GetFixedDepositDispatch(id int) (*dto.FixedDepositDispatchResponseDTO, error)
	GetFixedDepositDispatchList(filter dto.FixedDepositDispatchFilterDTO) ([]dto.FixedDepositDispatchResponseDTO, *uint64, error)
}

type FixedDepositJudgeService interface {
	CreateFixedDepositJudge(input dto.FixedDepositJudgeDTO) (*dto.FixedDepositJudgeResponseDTO, error)
	UpdateFixedDepositJudge(id int, input dto.FixedDepositJudgeDTO) (*dto.FixedDepositJudgeResponseDTO, error)
	DeleteFixedDepositJudge(id int) error
	GetFixedDepositJudge(id int) (*dto.FixedDepositJudgeResponseDTO, error)
	GetFixedDepositJudgeList(filter dto.FixedDepositJudgeFilterDTO) ([]dto.FixedDepositJudgeResponseDTO, *uint64, error)
}

type FixedDepositWillService interface {
	CreateFixedDepositWill(ctx context.Context, input dto.FixedDepositWillDTO) (*dto.FixedDepositWillResponseDTO, error)
	UpdateFixedDepositWill(ctx context.Context, id int, input dto.FixedDepositWillDTO) (*dto.FixedDepositWillResponseDTO, error)
	DeleteFixedDepositWill(ctx context.Context, id int) error
	GetFixedDepositWill(id int) (*dto.FixedDepositWillResponseDTO, error)
	GetFixedDepositWillList(filter dto.FixedDepositWillFilterDTO) ([]dto.FixedDepositWillResponseDTO, *uint64, error)
}

type FixedDepositWillDispatchService interface {
	CreateFixedDepositWillDispatch(ctx context.Context, input dto.FixedDepositWillDispatchDTO) (*dto.FixedDepositWillDispatchResponseDTO, error)
	UpdateFixedDepositWillDispatch(ctx context.Context, id int, input dto.FixedDepositWillDispatchDTO) (*dto.FixedDepositWillDispatchResponseDTO, error)
	DeleteFixedDepositWillDispatch(ctx context.Context, id int) error
	GetFixedDepositWillDispatch(id int) (*dto.FixedDepositWillDispatchResponseDTO, error)
	GetFixedDepositWillDispatchList(filter dto.FixedDepositWillDispatchFilterDTO) ([]dto.FixedDepositWillDispatchResponseDTO, *uint64, error)
}

type DepositPaymentService interface {
	CreateDepositPayment(ctx context.Context, input dto.DepositPaymentDTO) (*dto.DepositPaymentResponseDTO, error)
	UpdateDepositPayment(ctx context.Context, id int, input dto.DepositPaymentDTO) (*dto.DepositPaymentResponseDTO, error)
	DeleteDepositPayment(ctx context.Context, id int) error
	GetDepositPayment(id int) (*dto.DepositPaymentResponseDTO, error)
	GetDepositPaymentList(filter dto.DepositPaymentFilterDTO) ([]dto.DepositPaymentResponseDTO, *uint64, error)
	GetInitialState(filter dto.DepositInitialStateFilter) ([]dto.DepositPaymentResponseDTO, error)
	GetDepositPaymentByCaseNumber(caseNumber *string, sourceBankAccount *string) (*dto.DepositPaymentResponseDTO, error)
	GetCaseNumber(orgUnitID *int, sourceBankAccount *string) ([]dto.DepositPaymentResponseDTO, error)
}

type DepositPaymentOrderService interface {
	CreateDepositPaymentOrder(ctx context.Context, input dto.DepositPaymentOrderDTO) (*dto.DepositPaymentOrderResponseDTO, error)
	UpdateDepositPaymentOrder(ctx context.Context, id int, input dto.DepositPaymentOrderDTO) (*dto.DepositPaymentOrderResponseDTO, error)
	PayDepositPaymentOrder(ctx context.Context, id int, input dto.DepositPaymentOrderDTO) error
	DeleteDepositPaymentOrder(ctx context.Context, id int) error
	GetDepositPaymentOrder(id int) (*dto.DepositPaymentOrderResponseDTO, error)
	GetDepositPaymentOrderList(filter dto.DepositPaymentOrderFilterDTO) ([]dto.DepositPaymentOrderResponseDTO, *uint64, error)
}

type DepositAdditionalExpenseService interface {
	CreateDepositAdditionalExpense(input dto.DepositAdditionalExpenseDTO) (*dto.DepositAdditionalExpenseResponseDTO, error)
	UpdateDepositAdditionalExpense(id int, input dto.DepositAdditionalExpenseDTO) (*dto.DepositAdditionalExpenseResponseDTO, error)
	DeleteDepositAdditionalExpense(id int) error
	GetDepositAdditionalExpense(id int) (*dto.DepositAdditionalExpenseResponseDTO, error)
	GetDepositAdditionalExpenseList(filter dto.DepositAdditionalExpenseFilterDTO) ([]dto.DepositAdditionalExpenseResponseDTO, *uint64, error)
}

type PaymentOrderService interface {
	CreatePaymentOrder(ctx context.Context, input dto.PaymentOrderDTO) (*dto.PaymentOrderResponseDTO, error)
	UpdatePaymentOrder(ctx context.Context, id int, input dto.PaymentOrderDTO) (*dto.PaymentOrderResponseDTO, error)
	DeletePaymentOrder(ctx context.Context, id int) error
	GetPaymentOrder(id int) (*dto.PaymentOrderResponseDTO, error)
	GetPaymentOrderList(filter dto.PaymentOrderFilterDTO) ([]dto.PaymentOrderResponseDTO, *uint64, error)
	GetAllObligations(filter dto.GetObligationsFilterDTO) ([]dto.ObligationResponse, *uint64, error)
	PayPaymentOrder(ctx context.Context, id int, input dto.PaymentOrderDTO) error
	CancelPaymentOrder(ctx context.Context, id int) error
}

type PaymentOrderItemService interface {
	CreatePaymentOrderItem(input dto.PaymentOrderItemDTO) (*dto.PaymentOrderItemResponseDTO, error)
	UpdatePaymentOrderItem(id int, input dto.PaymentOrderItemDTO) (*dto.PaymentOrderItemResponseDTO, error)
	DeletePaymentOrderItem(id int) error
	GetPaymentOrderItem(id int) (*dto.PaymentOrderItemResponseDTO, error)
	GetPaymentOrderItemList(filter dto.PaymentOrderItemFilterDTO) ([]dto.PaymentOrderItemResponseDTO, *uint64, error)
}

type EnforcedPaymentService interface {
	CreateEnforcedPayment(ctx context.Context, input dto.EnforcedPaymentDTO) (*dto.EnforcedPaymentResponseDTO, error)
	UpdateEnforcedPayment(ctx context.Context, id int, input dto.EnforcedPaymentDTO) (*dto.EnforcedPaymentResponseDTO, error)
	ReturnEnforcedPayment(ctx context.Context, id int, input dto.EnforcedPaymentDTO) error
	DeleteEnforcedPayment(ctx context.Context, id int) error
	GetEnforcedPayment(id int) (*dto.EnforcedPaymentResponseDTO, error)
	GetEnforcedPaymentList(filter dto.EnforcedPaymentFilterDTO) ([]dto.EnforcedPaymentResponseDTO, *uint64, error)
}

type EnforcedPaymentItemService interface {
	CreateEnforcedPaymentItem(input dto.EnforcedPaymentItemDTO) (*dto.EnforcedPaymentItemResponseDTO, error)
	UpdateEnforcedPaymentItem(id int, input dto.EnforcedPaymentItemDTO) (*dto.EnforcedPaymentItemResponseDTO, error)
	DeleteEnforcedPaymentItem(id int) error
	GetEnforcedPaymentItem(id int) (*dto.EnforcedPaymentItemResponseDTO, error)
	GetEnforcedPaymentItemList(filter dto.EnforcedPaymentItemFilterDTO) ([]dto.EnforcedPaymentItemResponseDTO, *uint64, error)
}

type AccountingEntryService interface {
	CreateAccountingEntry(ctx context.Context, input dto.AccountingEntryDTO) (*dto.AccountingEntryResponseDTO, error)
	UpdateAccountingEntry(ctx context.Context, id int, input dto.AccountingEntryDTO) (*dto.AccountingEntryResponseDTO, error)
	DeleteAccountingEntry(ctx context.Context, id int) error
	GetAccountingEntry(id int) (*dto.AccountingEntryResponseDTO, error)
	GetAccountingEntryList(filter dto.AccountingEntryFilterDTO) ([]dto.AccountingEntryResponseDTO, *uint64, error)

	GetObligationsForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.ObligationForAccounting, *uint64, error)
	GetPaymentOrdersForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.PaymentOrdersForAccounting, *uint64, error)
	GetEnforcedPaymentsForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.PaymentOrdersForAccounting, *uint64, error)
	GetReturnedEnforcedPaymentsForAccounting(filter dto.GetObligationsFilterDTO) ([]dto.PaymentOrdersForAccounting, *uint64, error)
	BuildAccountingOrderForObligations(data dto.AccountingOrderForObligationsData) (*dto.AccountingOrderForObligations, error)

	GetAnalyticalCard(filter data.AnalyticalCardFilter) ([]data.AnalyticalCard, error)
}

type ModelsOfAccountingService interface {
	CreateModelsOfAccounting(ctx context.Context, input dto.ModelsOfAccountingDTO) (*dto.ModelsOfAccountingResponseDTO, error)
	UpdateModelsOfAccounting(ctx context.Context, id int, input dto.ModelsOfAccountingDTO) (*dto.ModelsOfAccountingResponseDTO, error)
	GetModelsOfAccounting(id int) (*dto.ModelsOfAccountingResponseDTO, error)
	GetModelsOfAccountingList(filter dto.ModelsOfAccountingFilterDTO) ([]dto.ModelsOfAccountingResponseDTO, *uint64, error)
}

type ModelOfAccountingItemService interface {
	CreateModelOfAccountingItem(input dto.ModelOfAccountingItemDTO) (*dto.ModelOfAccountingItemResponseDTO, error)
	UpdateModelOfAccountingItem(id int, input dto.ModelOfAccountingItemDTO) (*dto.ModelOfAccountingItemResponseDTO, error)
	GetModelOfAccountingItem(id int) (*dto.ModelOfAccountingItemResponseDTO, error)
	GetModelOfAccountingItemList(filter dto.ModelOfAccountingItemFilterDTO) ([]dto.ModelOfAccountingItemResponseDTO, *uint64, error)
}

type AccountingEntryItemService interface {
	CreateAccountingEntryItem(input dto.AccountingEntryItemDTO) (*dto.AccountingEntryItemResponseDTO, error)
	UpdateAccountingEntryItem(id int, input dto.AccountingEntryItemDTO) (*dto.AccountingEntryItemResponseDTO, error)
	DeleteAccountingEntryItem(id int) error
	GetAccountingEntryItem(id int) (*dto.AccountingEntryItemResponseDTO, error)
	GetAccountingEntryItemList(filter dto.AccountingEntryItemFilterDTO) ([]dto.AccountingEntryItemResponseDTO, *uint64, error)
}

type SpendingDynamicService interface {
	CreateSpendingDynamic(ctx context.Context, budgetID, unitID int, input []dto.SpendingDynamicDTO) error
	GetSpendingDynamic(currentBudgetID, budgetID, unitID *int, version *int) ([]dto.SpendingDynamicWithEntryResponseDTO, error)
	GetActual(budgetID, unitID, accountID int) (decimal.Decimal, error)
	GetSpendingDynamicHistory(budgetID, unitID int) ([]dto.SpendingDynamicHistoryResponseDTO, error)
	CreateInititalSpendingDynamicFromCurrentBudget(ctx context.Context, currentBudget *data.CurrentBudget) error
}

type CurrentBudgetService interface {
	CreateCurrentBudget(ctx context.Context, input dto.CurrentBudgetDTO) (*dto.CurrentBudgetResponseDTO, error)
	UpdateActual(ctx context.Context, budgetID, accountID, unitID int, actual decimal.Decimal) (*dto.CurrentBudgetResponseDTO, error)
	UpdateBalance(ctx context.Context, tx up.Session, id int, balance decimal.Decimal) error
	GetCurrentBudget(id int) (*dto.CurrentBudgetResponseDTO, error)
	GetCurrentBudgetList(filter dto.CurrentBudgetFilterDTO) ([]dto.CurrentBudgetResponseDTO, *uint64, error)
	GetAcctualCurrentBudget(organizationUnitID int) ([]dto.CurrentBudgetResponseDTO, error)
}

type SpendingReleaseService interface {
	CreateSpendingRelease(ctx context.Context, budgetID, unitID int, input []dto.SpendingReleaseDTO) ([]dto.SpendingReleaseResponseDTO, error)
	DeleteSpendingRelease(ctx context.Context, input *dto.DeleteSpendingReleaseInput) error
	GetSpendingRelease(id int) (*dto.SpendingReleaseResponseDTO, error)
	GetSpendingReleaseList(filter data.SpendingReleaseFilterDTO) ([]dto.SpendingReleaseResponseDTO, error)
	GetSpendingReleaseOverview(filter dto.SpendingReleaseOverviewFilterDTO) ([]dto.SpendingReleaseOverview, error)
}

type InternalReallocationService interface {
	CreateInternalReallocation(ctx context.Context, input dto.InternalReallocationDTO) (*dto.InternalReallocationResponseDTO, error)
	DeleteInternalReallocation(ctx context.Context, id int) error
	GetInternalReallocation(id int) (*dto.InternalReallocationResponseDTO, error)
	GetInternalReallocationList(filter dto.InternalReallocationFilterDTO) ([]dto.InternalReallocationResponseDTO, *uint64, error)
}

type InternalReallocationItemService interface {
	CreateInternalReallocationItem(input dto.InternalReallocationItemDTO) (*dto.InternalReallocationItemResponseDTO, error)
	DeleteInternalReallocationItem(id int) error
	GetInternalReallocationItem(id int) (*dto.InternalReallocationItemResponseDTO, error)
	GetInternalReallocationItemList(filter dto.InternalReallocationItemFilterDTO) ([]dto.InternalReallocationItemResponseDTO, *uint64, error)
}

type ExternalReallocationService interface {
	CreateExternalReallocation(ctx context.Context, input dto.ExternalReallocationDTO) (*dto.ExternalReallocationResponseDTO, error)
	DeleteExternalReallocation(ctx context.Context, id int) error

	AcceptOUExternalReallocation(ctx context.Context, input dto.ExternalReallocationDTO) (*dto.ExternalReallocationResponseDTO, error)
	RejectOUExternalReallocation(ctx context.Context, id int) error

	AcceptSSSExternalReallocation(ctx context.Context, id int) error
	RejectSSSExternalReallocation(ctx context.Context, id int) error

	GetExternalReallocation(id int) (*dto.ExternalReallocationResponseDTO, error)
	GetExternalReallocationList(filter dto.ExternalReallocationFilterDTO) ([]dto.ExternalReallocationResponseDTO, *uint64, error)
}

type ExternalReallocationItemService interface {
	CreateExternalReallocationItem(input dto.ExternalReallocationItemDTO) (*dto.ExternalReallocationItemResponseDTO, error)
	DeleteExternalReallocationItem(id int) error
	GetExternalReallocationItem(id int) (*dto.ExternalReallocationItemResponseDTO, error)
	GetExternalReallocationItemList(filter dto.ExternalReallocationItemFilterDTO) ([]dto.ExternalReallocationItemResponseDTO, *uint64, error)
}

type LogService interface {
	CreateLog(input dto.LogDTO) (*dto.LogResponseDTO, error)
	UpdateLog(id int, input dto.LogDTO) (*dto.LogResponseDTO, error)
	DeleteLog(id int) error
	GetLog(id int) (*dto.LogResponseDTO, error)
	GetLogList(filter dto.LogFilterDTO) ([]dto.LogResponseDTO, *uint64, error)
}
