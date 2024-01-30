package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// FinancialBudget struct
type FinancialBudget struct {
	ID             int       `db:"id,omitempty"`
	AccountVersion int       `db:"account_version"`
	BudgetID       int       `db:"budget_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *FinancialBudget) Table() string {
	return "financial_budgets"
}

// GetAll gets all records from the database, using upper
func (t *FinancialBudget) GetAll(condition *up.Cond) ([]*FinancialBudget, error) {
	collection := upper.Collection(t.Table())
	var all []*FinancialBudget
	var res up.Result

	if condition != nil {
		res = collection.Find(*condition)
	} else {
		res = collection.Find()
	}

	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, err
}

// Get gets one record from the database, by id, using upper
func (t *FinancialBudget) Get(id int) (*FinancialBudget, error) {
	var one FinancialBudget
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FinancialBudget) Update(m FinancialBudget) error {
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
func (t *FinancialBudget) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FinancialBudget) Insert(m FinancialBudget) (int, error) {
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
