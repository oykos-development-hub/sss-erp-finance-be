package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// NonFinancialBudgetGoal struct
type NonFinancialBudgetGoal struct {
	ID                   int       `db:"id,omitempty"`
	NonFinancialBudgetID int       `db:"non_financial_budget_id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	CreatedAt            time.Time `db:"created_at,omitempty"`
	UpdatedAt            time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *NonFinancialBudgetGoal) Table() string {
	return "non_financial_budget_goals"
}

// GetAll gets all records from the database, using upper
func (t *NonFinancialBudgetGoal) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*NonFinancialBudgetGoal, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*NonFinancialBudgetGoal
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
func (t *NonFinancialBudgetGoal) Get(id int) (*NonFinancialBudgetGoal, error) {
	var one NonFinancialBudgetGoal
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *NonFinancialBudgetGoal) Update(m NonFinancialBudgetGoal) error {
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
func (t *NonFinancialBudgetGoal) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *NonFinancialBudgetGoal) Insert(m NonFinancialBudgetGoal) (int, error) {
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
