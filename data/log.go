package data

import (
	"encoding/json"
	"time"

	up "github.com/upper/db/v4"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type LogOperation string
type LogEntity string

var (
	OperationInsert LogOperation = "INSERT"
	OperationUpdate LogOperation = "UPDATE"
	OperationDelete LogOperation = "DELETE"
)

var (
	EntityBudget                                LogEntity = "budgets"
	EntityAccountingEntries                     LogEntity = "accounting_entries"
	EntityActivities                            LogEntity = "activities"
	EntityBudgetRequests                        LogEntity = "budget_requests"
	EntityCurrentBudget                         LogEntity = "current_budgets"
	EntityDepositPaymentOrders                  LogEntity = "deposit_payment_orders"
	EntityDepositPayments                       LogEntity = "deposit_payments"
	EntityEnforcedPayments                      LogEntity = "enforced_payments"
	EntityExternalReallocations                 LogEntity = "external_reallocations"
	EntityFeePayments                           LogEntity = "fee_payments"
	EntityFees                                  LogEntity = "fees"
	EntityFilledFinancialBudgets                LogEntity = "filled_financial_budgets"
	EntityFinancialBudgets                      LogEntity = "financial_budgets"
	EntityFinePayments                          LogEntity = "fine_payments"
	EntityFines                                 LogEntity = "fines"
	EntityFixedDepositDispatches                LogEntity = "fixed_deposit_dispatches"
	EntityFixedDepositItems                     LogEntity = "fixed_deposit_items"
	EntityFixedDepositWillDispatches            LogEntity = "fixed_deposit_will_dispatches"
	EntityFixedDepositWills                     LogEntity = "fixed_deposit_wills"
	EntityFixedDeposits                         LogEntity = "fixed_deposits"
	EntityFlatRatePayments                      LogEntity = "flat_rate_payments"
	EntityFlatRates                             LogEntity = "flat_rates"
	EntityGoalIndicators                        LogEntity = "goal_indicators"
	EntityInternalReallocations                 LogEntity = "internal_reallocations"
	EntityInvoices                              LogEntity = "invoices"
	EntityModelsOfAccounting                    LogEntity = "models_of_accounting"
	EntityNonFinancialBudgetGoals               LogEntity = "non_financial_budget_goals"
	EntityNonFinancialBudgets                   LogEntity = "non_financial_budgets"
	EntityPaymentOrders                         LogEntity = "payment_orders"
	EntityProcedureCostPayments                 LogEntity = "procedure_cost_payments"
	EntityProcedureCosts                        LogEntity = "procedure_costs"
	EntityPrograms                              LogEntity = "programs"
	EntityPropertyBenefitsConfiscations         LogEntity = "property_benefits_confiscations"
	EntityPropertyBenefitsConfiscationsPayments LogEntity = "property_benefits_confiscations_payments"
	EntitySalaries                              LogEntity = "salaries"
	EntitySpendingDynamicEntries                LogEntity = "spending_dynamic_entries"
	EntitySpendingReleases                      LogEntity = "spending_releases"
	EntityTaxAuthorityCodebooks                 LogEntity = "tax_authority_codebooks"
)

// Log struct
type Log struct {
	ID        int              `db:"id,omitempty"`
	ChangedAt time.Time        `db:"changed_at"`
	UserID    int              `db:"user_id"`
	ItemID    int              `db:"item_id"`
	Operation LogOperation     `db:"operation"`
	Entity    LogEntity        `db:"entity"`
	OldState  *json.RawMessage `db:"old_state"`
	NewState  *json.RawMessage `db:"new_state"`
}

// Table returns the table name
func (t *Log) Table() string {
	return "logs"
}

// GetAll gets all records from the database, using Upper
func (t *Log) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*Log, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*Log
	var res up.Result

	if condition != nil {
		res = collection.Find(condition)
	} else {
		res = collection.Find()
	}
	total, err := res.Count()
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper count")
	}

	if page != nil && size != nil {
		res = paginateResult(res, *page, *size)
	}

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper all")
	}

	return all, &total, err
}

// Get gets one record from the database, by id, using Upper
func (t *Log) Get(id int) (*Log, error) {
	var one Log
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Update updates a record in the database, using Upper
func (t *Log) Update(m Log) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return newErrors.Wrap(err, "upper update")
	}
	return nil
}

// Delete deletes a record from the database by id, using Upper
func (t *Log) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return newErrors.Wrap(err, "upper delete")
	}
	return nil
}

// Insert inserts a model into the database, using Upper
func (t *Log) Insert(m Log) (int, error) {
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, newErrors.Wrap(err, "upper insert")
	}

	id := getInsertId(res.ID())

	return id, nil
}
