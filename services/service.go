package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
)

type BaseService interface {
	RandomString(n int) string
	Encrypt(text string) (string, error)
	Decrypt(crypto string) (string, error)
}

type InvoiceService interface {
	CreateInvoice(input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error)
	UpdateInvoice(id int, input dto.InvoiceDTO) (*dto.InvoiceResponseDTO, error)
	DeleteInvoice(id int) error
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
	CreateBudget(input dto.BudgetDTO) (*dto.BudgetResponseDTO, error)
	UpdateBudget(id int, input dto.BudgetDTO) (*dto.BudgetResponseDTO, error)
	DeleteBudget(id int) error
	GetBudget(id int) (*dto.BudgetResponseDTO, error)
	GetBudgetList(input dto.GetBudgetListInput) ([]dto.BudgetResponseDTO, error)
}

type FinancialBudgetService interface {
	CreateFinancialBudget(input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error)
	UpdateFinancialBudget(id int, input dto.FinancialBudgetDTO) (*dto.FinancialBudgetResponseDTO, error)
	DeleteFinancialBudget(id int) error
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
	CreateNonFinancialBudget(input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error)
	UpdateNonFinancialBudget(id int, input dto.NonFinancialBudgetDTO) (*dto.NonFinancialBudgetResponseDTO, error)
	DeleteNonFinancialBudget(id int) error
	GetNonFinancialBudget(id int) (*dto.NonFinancialBudgetResponseDTO, error)
	GetNonFinancialBudgetList(filter dto.NonFinancialBudgetFilterDTO) ([]dto.NonFinancialBudgetResponseDTO, *uint64, error)
}

type NonFinancialBudgetGoalService interface {
	CreateNonFinancialBudgetGoal(input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	UpdateNonFinancialBudgetGoal(id int, input dto.NonFinancialBudgetGoalDTO) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	DeleteNonFinancialBudgetGoal(id int) error
	GetNonFinancialBudgetGoal(id int) (*dto.NonFinancialBudgetGoalResponseDTO, error)
	GetNonFinancialBudgetGoalList(filter dto.NonFinancialBudgetGoalFilterDTO) ([]dto.NonFinancialBudgetGoalResponseDTO, *uint64, error)
}

type ProgramService interface {
	CreateProgram(input dto.ProgramDTO) (*dto.ProgramResponseDTO, error)
	UpdateProgram(id int, input dto.ProgramDTO) (*dto.ProgramResponseDTO, error)
	DeleteProgram(id int) error
	GetProgram(id int) (*dto.ProgramResponseDTO, error)
	GetProgramList(filter dto.ProgramFilterDTO) ([]dto.ProgramResponseDTO, *uint64, error)
}

type ActivityService interface {
	CreateActivity(input dto.ActivityDTO) (*dto.ActivityResponseDTO, error)
	UpdateActivity(id int, input dto.ActivityDTO) (*dto.ActivityResponseDTO, error)
	DeleteActivity(id int) error
	GetActivity(id int) (*dto.ActivityResponseDTO, error)
	GetActivityList(filter dto.ActivityFilterDTO) ([]dto.ActivityResponseDTO, *uint64, error)
}

type GoalIndicatorService interface {
	CreateGoalIndicator(input dto.GoalIndicatorDTO) (*dto.GoalIndicatorResponseDTO, error)
	UpdateGoalIndicator(id int, input dto.GoalIndicatorDTO) (*dto.GoalIndicatorResponseDTO, error)
	DeleteGoalIndicator(id int) error
	GetGoalIndicator(id int) (*dto.GoalIndicatorResponseDTO, error)
	GetGoalIndicatorList(filter dto.GoalIndicatorFilterDTO) ([]dto.GoalIndicatorResponseDTO, *uint64, error)
}

type FilledFinancialBudgetService interface {
	CreateFilledFinancialBudget(input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error)
	UpdateFilledFinancialBudget(id int, input dto.FilledFinancialBudgetDTO) (*dto.FilledFinancialBudgetResponseDTO, error)
	DeleteFilledFinancialBudget(id int) error
	GetFilledFinancialBudget(id int) (*dto.FilledFinancialBudgetResponseDTO, error)
	GetFilledFinancialBudgetList(filter dto.FilledFinancialBudgetFilterDTO) ([]dto.FilledFinancialBudgetResponseDTO, *uint64, error)
}

type BudgetRequestService interface {
	CreateBudgetRequest(input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error)
	UpdateBudgetRequest(id int, input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error)
	DeleteBudgetRequest(id int) error
	GetBudgetRequest(id int) (*dto.BudgetRequestResponseDTO, error)
	GetBudgetRequestList(filter dto.BudgetRequestFilterDTO) ([]dto.BudgetRequestResponseDTO, *uint64, error)
}

type FeeSharedLogicService interface {
	CalculateFeeDetailsAndUpdateStatus(fineId int) (*dto.FeeDetailsDTO, data.FeeStatus, error)
}

type FeeService interface {
	CreateFee(input dto.FeeDTO) (*dto.FeeResponseDTO, error)
	GetFee(id int) (*dto.FeeResponseDTO, error)
	GetFeeList(filter dto.FeeFilterDTO) ([]dto.FeeResponseDTO, *uint64, error)
	UpdateFee(id int, input dto.FeeDTO) (*dto.FeeResponseDTO, error)
	DeleteFee(id int) error
}

type FeePaymentService interface {
	CreateFeePayment(input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error)
	DeleteFeePayment(id int) error
	UpdateFeePayment(id int, input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error)
	GetFeePayment(id int) (*dto.FeePaymentResponseDTO, error)
	GetFeePaymentList(filter dto.FeePaymentFilterDTO) ([]dto.FeePaymentResponseDTO, *uint64, error)
}

type FineSharedLogicService interface {
	CalculateFineDetailsAndUpdateStatus(fineId int) (*dto.FineFeeDetailsDTO, data.FineStatus, error)
}

type FineService interface {
	CreateFine(input dto.FineDTO) (*dto.FineResponseDTO, error)
	GetFine(id int) (*dto.FineResponseDTO, error)
	GetFineList(filter dto.FineFilterDTO) ([]dto.FineResponseDTO, *uint64, error)
	UpdateFine(id int, input dto.FineDTO) (*dto.FineResponseDTO, error)
	DeleteFine(id int) error
}

type FinePaymentService interface {
	CreateFinePayment(input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error)
	DeleteFinePayment(id int) error
	UpdateFinePayment(id int, input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error)
	GetFinePayment(id int) (*dto.FinePaymentResponseDTO, error)
	GetFinePaymentList(filter dto.FinePaymentFilterDTO) ([]dto.FinePaymentResponseDTO, *uint64, error)
}

type ProcedureCostSharedLogicService interface {
	CalculateProcedureCostDetailsAndUpdateStatus(procedureCostId int) (*dto.ProcedureCostDetailsDTO, data.ProcedureCostStatus, error)
}

type ProcedureCostService interface {
	CreateProcedureCost(input dto.ProcedureCostDTO) (*dto.ProcedureCostResponseDTO, error)
	GetProcedureCost(id int) (*dto.ProcedureCostResponseDTO, error)
	GetProcedureCostList(filter dto.ProcedureCostFilterDTO) ([]dto.ProcedureCostResponseDTO, *uint64, error)
	UpdateProcedureCost(id int, input dto.ProcedureCostDTO) (*dto.ProcedureCostResponseDTO, error)
	DeleteProcedureCost(id int) error
}

type ProcedureCostPaymentService interface {
	CreateProcedureCostPayment(input dto.ProcedureCostPaymentDTO) (*dto.ProcedureCostPaymentResponseDTO, error)
	DeleteProcedureCostPayment(id int) error
	UpdateProcedureCostPayment(id int, input dto.ProcedureCostPaymentDTO) (*dto.ProcedureCostPaymentResponseDTO, error)
	GetProcedureCostPayment(id int) (*dto.ProcedureCostPaymentResponseDTO, error)
	GetProcedureCostPaymentList(filter dto.ProcedureCostPaymentFilterDTO) ([]dto.ProcedureCostPaymentResponseDTO, *uint64, error)
}

type AdditionalExpenseService interface {
	DeleteAdditionalExpense(id int) error
	GetAdditionalExpense(id int) (*dto.AdditionalExpenseResponseDTO, error)
	GetAdditionalExpenseList(filter dto.AdditionalExpenseFilterDTO) ([]dto.AdditionalExpenseResponseDTO, *uint64, error)
}

type FlatRateSharedLogicService interface {
	CalculateFlatRateDetailsAndUpdateStatus(flatRateId int) (*dto.FlatRateDetailsDTO, data.FlatRateStatus, error)
}

type FlatRateService interface {
	CreateFlatRate(input dto.FlatRateDTO) (*dto.FlatRateResponseDTO, error)
	GetFlatRate(id int) (*dto.FlatRateResponseDTO, error)
	GetFlatRateList(filter dto.FlatRateFilterDTO) ([]dto.FlatRateResponseDTO, *uint64, error)
	UpdateFlatRate(id int, input dto.FlatRateDTO) (*dto.FlatRateResponseDTO, error)
	DeleteFlatRate(id int) error
}

type FlatRatePaymentService interface {
	CreateFlatRatePayment(input dto.FlatRatePaymentDTO) (*dto.FlatRatePaymentResponseDTO, error)
	DeleteFlatRatePayment(id int) error
	UpdateFlatRatePayment(id int, input dto.FlatRatePaymentDTO) (*dto.FlatRatePaymentResponseDTO, error)
	GetFlatRatePayment(id int) (*dto.FlatRatePaymentResponseDTO, error)
	GetFlatRatePaymentList(filter dto.FlatRatePaymentFilterDTO) ([]dto.FlatRatePaymentResponseDTO, *uint64, error)
}

type PropBenConfSharedLogicService interface {
	CalculatePropBenConfDetailsAndUpdateStatus(propbenconfId int) (*dto.PropBenConfDetailsDTO, data.PropBenConfStatus, error)
}

type PropBenConfService interface {
	CreatePropBenConf(input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error)
	GetPropBenConf(id int) (*dto.PropBenConfResponseDTO, error)
	GetPropBenConfList(filter dto.PropBenConfFilterDTO) ([]dto.PropBenConfResponseDTO, *uint64, error)
	UpdatePropBenConf(id int, input dto.PropBenConfDTO) (*dto.PropBenConfResponseDTO, error)
	DeletePropBenConf(id int) error
}

type PropBenConfPaymentService interface {
	CreatePropBenConfPayment(input dto.PropBenConfPaymentDTO) (*dto.PropBenConfPaymentResponseDTO, error)
	DeletePropBenConfPayment(id int) error
	UpdatePropBenConfPayment(id int, input dto.PropBenConfPaymentDTO) (*dto.PropBenConfPaymentResponseDTO, error)
	GetPropBenConfPayment(id int) (*dto.PropBenConfPaymentResponseDTO, error)
	GetPropBenConfPaymentList(filter dto.PropBenConfPaymentFilterDTO) ([]dto.PropBenConfPaymentResponseDTO, *uint64, error)
}

type TaxAuthorityCodebookService interface {
	CreateTaxAuthorityCodebook(input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error)
	UpdateTaxAuthorityCodebook(id int, input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error)
	DeleteTaxAuthorityCodebook(id int) error
	GetTaxAuthorityCodebook(id int) (*dto.TaxAuthorityCodebookResponseDTO, error)
	GetTaxAuthorityCodebookList(filter dto.TaxAuthorityCodebookFilterDTO) ([]dto.TaxAuthorityCodebookResponseDTO, *uint64, error)
}

type SalaryService interface {
	CreateSalary(input dto.SalaryDTO) (*dto.SalaryResponseDTO, error)
	UpdateSalary(id int, input dto.SalaryDTO) (*dto.SalaryResponseDTO, error)
	DeleteSalary(id int) error
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
	CreateFixedDeposit(input dto.FixedDepositDTO) (*dto.FixedDepositResponseDTO, error)
	UpdateFixedDeposit(id int, input dto.FixedDepositDTO) (*dto.FixedDepositResponseDTO, error)
	DeleteFixedDeposit(id int) error
	GetFixedDeposit(id int) (*dto.FixedDepositResponseDTO, error)
	GetFixedDepositList(filter dto.FixedDepositFilterDTO) ([]dto.FixedDepositResponseDTO, *uint64, error)
}

type FixedDepositItemService interface {
	CreateFixedDepositItem(input dto.FixedDepositItemDTO) (*dto.FixedDepositItemResponseDTO, error)
	UpdateFixedDepositItem(id int, input dto.FixedDepositItemDTO) (*dto.FixedDepositItemResponseDTO, error)
	DeleteFixedDepositItem(id int) error
	GetFixedDepositItem(id int) (*dto.FixedDepositItemResponseDTO, error)
	GetFixedDepositItemList(filter dto.FixedDepositItemFilterDTO) ([]dto.FixedDepositItemResponseDTO, *uint64, error)
}

type FixedDepositDispatchService interface {
	CreateFixedDepositDispatch(input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error)
	UpdateFixedDepositDispatch(id int, input dto.FixedDepositDispatchDTO) (*dto.FixedDepositDispatchResponseDTO, error)
	DeleteFixedDepositDispatch(id int) error
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
	CreateFixedDepositWill(input dto.FixedDepositWillDTO) (*dto.FixedDepositWillResponseDTO, error)
	UpdateFixedDepositWill(id int, input dto.FixedDepositWillDTO) (*dto.FixedDepositWillResponseDTO, error)
	DeleteFixedDepositWill(id int) error
	GetFixedDepositWill(id int) (*dto.FixedDepositWillResponseDTO, error)
	GetFixedDepositWillList(filter dto.FixedDepositWillFilterDTO) ([]dto.FixedDepositWillResponseDTO, *uint64, error)
}

type FixedDepositWillDispatchService interface {
	CreateFixedDepositWillDispatch(input dto.FixedDepositWillDispatchDTO) (*dto.FixedDepositWillDispatchResponseDTO, error)
	UpdateFixedDepositWillDispatch(id int, input dto.FixedDepositWillDispatchDTO) (*dto.FixedDepositWillDispatchResponseDTO, error)
	DeleteFixedDepositWillDispatch(id int) error
	GetFixedDepositWillDispatch(id int) (*dto.FixedDepositWillDispatchResponseDTO, error)
	GetFixedDepositWillDispatchList(filter dto.FixedDepositWillDispatchFilterDTO) ([]dto.FixedDepositWillDispatchResponseDTO, *uint64, error)
}
