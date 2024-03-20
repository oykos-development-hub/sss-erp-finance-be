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
	CreateAdditionalExpense(input dto.AdditionalExpenseDTO) (*dto.AdditionalExpenseResponseDTO, error)
	UpdateAdditionalExpense(id int, input dto.AdditionalExpenseDTO) (*dto.AdditionalExpenseResponseDTO, error)
	DeleteAdditionalExpense(id int) error
	GetAdditionalExpense(id int) (*dto.AdditionalExpenseResponseDTO, error)
	GetAdditionalExpenseList(filter dto.AdditionalExpenseFilterDTO) ([]dto.AdditionalExpenseResponseDTO, *uint64, error)
}
