package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// FilledFinancialBudget struct
type FilledFinancialBudget struct {
	ID                 int       `db:"id,omitempty"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	FinanceBudgetID    int       `db:"finance_budget_id"`
	AccountID          int       `db:"account_id"`
	CurrentYear        int       `db:"current_year"`
	NextYear           int       `db:"next_year"`
	YearAfterNext      int       `db:"year_after_next"`
	Description        string    `db:"description"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *FilledFinancialBudget) Table() string {
	return "filled_financial_budgets"
}

// GetAll gets all records from the database, using upper
func (t *FilledFinancialBudget) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*FilledFinancialBudget, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*FilledFinancialBudget
	var res up.Result

	if condition != nil {
		res = collection.Find(condition)
	} else {
		res = collection.Find()
	}
	total, err := res.Count()
	if err != nil {
		return nil, nil, err
	}

	if page != nil && size != nil {
		res = paginateResult(res, *page, *size)
	}

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, err
	}

	return all, &total, err
}

// Get gets one record from the database, by id, using upper
func (t *FilledFinancialBudget) Get(id int) (*FilledFinancialBudget, error) {
	var one FilledFinancialBudget
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FilledFinancialBudget) Update(m FilledFinancialBudget) error {
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *FilledFinancialBudget) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FilledFinancialBudget) Insert(m FilledFinancialBudget) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}