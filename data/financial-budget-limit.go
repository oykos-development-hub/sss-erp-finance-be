package data

import (
	"time"

	up "github.com/upper/db/v4"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// FinancialBudgetLimit struct
type FinancialBudgetLimit struct {
	ID                 int       `db:"id,omitempty"`
	BudgetID           int       `db:"budget_id"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	Limit              int       `db:"limit_value"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *FinancialBudgetLimit) Table() string {
	return "financial_budget_limits"
}

// GetAll gets all records from the database, using upper
func (t *FinancialBudgetLimit) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*FinancialBudgetLimit, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*FinancialBudgetLimit
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

// Get gets one record from the database, by id, using upper
func (t *FinancialBudgetLimit) Get(id int) (*FinancialBudgetLimit, error) {
	var one FinancialBudgetLimit
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FinancialBudgetLimit) Update(m FinancialBudgetLimit) error {
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return newErrors.Wrap(err, "upper update")
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *FinancialBudgetLimit) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return newErrors.Wrap(err, "upper delete")
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FinancialBudgetLimit) Insert(m FinancialBudgetLimit) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, newErrors.Wrap(err, "upper insert")
	}

	id := getInsertId(res.ID())

	return id, nil
}
