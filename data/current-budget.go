package data

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/errors"
)

// CurrentBudget struct
type CurrentBudget struct {
	ID            int             `db:"id,omitempty"`
	BudgetID      int             `db:"budget_id"`
	UnitID        int             `db:"unit_id"`
	AccountID     int             `db:"account_id"`
	InitialActual decimal.Decimal `db:"initial_actual"`
	Actual        decimal.Decimal `db:"actual"`
	Balance       decimal.Decimal `db:"balance"`
	CreatedAt     time.Time       `db:"created_at,omitempty"`
}

// Table returns the table name
func (t *CurrentBudget) Table() string {
	return "current_budgets"
}

// GetAll gets all records from the database, using upper
func (t *CurrentBudget) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*CurrentBudget, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*CurrentBudget
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

// GetAll gets all records from the database, using upper
func (t *CurrentBudget) GetBy(condition up.AndExpr) (*CurrentBudget, error) {
	collection := Upper.Collection(t.Table())
	var one CurrentBudget

	res := collection.Find(&condition)

	err := res.One(&one)

	return &one, err
}

// Get gets one record from the database, by id, using upper
func (t *CurrentBudget) Get(id int) (*CurrentBudget, error) {
	var one CurrentBudget
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *CurrentBudget) UpdateActual(currentBudgetID int, actual decimal.Decimal) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET actual = $1 WHERE id = $2", t.Table())

	res, err := Upper.SQL().Exec(updateQuery, actual, currentBudgetID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		return errors.ErrNotFound
	}

	return nil
}

// Update updates a record in the database, using upper
func (t *CurrentBudget) UpdateBalance(currentBudgetID int, balance decimal.Decimal) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", t.Table())

	res, err := Upper.SQL().Exec(updateQuery, balance, currentBudgetID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		return errors.ErrNotFound
	}

	return nil
}

// Insert inserts a model into the database, using upper
func (t *CurrentBudget) Insert(m CurrentBudget) (int, error) {
	m.CreatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
