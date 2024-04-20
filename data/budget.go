package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type BudgetType int

const (
	CapitalBudgetType BudgetType = 1
	CurrentBudgetType BudgetType = 2
)

type BudgetStatus int

const (
	BudgetCreatedStatus  BudgetStatus = 1
	BudgetSentStatus     BudgetStatus = 2
	BudgetSentOnReview   BudgetStatus = 3
	BudgetRejectedStatus BudgetStatus = 4
	BudgetAcceptedStatus BudgetStatus = 5
)

// Budget struct
type Budget struct {
	ID           int          `db:"id,omitempty"`
	Year         int          `db:"year"`
	BudgetType   BudgetType   `db:"budget_type"`
	BudgetStatus BudgetStatus `db:"budget_status"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}

// Table returns the table name
func (t *Budget) Table() string {
	return "budgets"
}

// GetAll gets all records from the database, using upper
func (t *Budget) GetAll(condition *up.Cond, orders []any) ([]*Budget, error) {
	collection := Upper.Collection(t.Table())
	var all []*Budget
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

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, err
	}

	return all, err
}

// Get gets one record from the database, by id, using upper
func (t *Budget) Get(id int) (*Budget, error) {
	var one Budget
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Budget) Update(m Budget) error {
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *Budget) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *Budget) Insert(m Budget) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
