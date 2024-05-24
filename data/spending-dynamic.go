package data

import (
	goerrors "errors"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/errors"
)

// SpendingDynamic struct
type SpendingDynamic struct {
	ID              int             `db:"id,omitempty"`
	CurrentBudgetID int             `db:"current_budget_id"`
	PlannedTotal    decimal.Decimal `db:"actual"`
}

// Table returns the table name
func (t *SpendingDynamic) Table() string {
	return "spending_dynamics"
}

func (t *SpendingDynamic) List(cond *up.AndExpr, orders []any) ([]SpendingDynamic, error) {
	if cond == nil {
		return nil, goerrors.New("cond param is required")
	}
	collection := Upper.Collection(t.Table())

	var all []SpendingDynamic
	err := collection.
		Find(cond).
		OrderBy(orders...).
		All(&all)
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return all, nil
}

func (t *SpendingDynamic) GetBy(cond *up.AndExpr, orders []any) (*SpendingDynamic, error) {
	if cond == nil {
		return nil, goerrors.New("cond param is required")
	}
	collection := Upper.Collection(t.Table())

	var one SpendingDynamic
	err := collection.
		Find(cond).
		OrderBy(orders...).
		One(&one)
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &one, nil
}

func (t *SpendingDynamic) Get(id int) (*SpendingDynamic, error) {
	var one SpendingDynamic
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Insert inserts a model into the database, using upper
func (t *SpendingDynamic) Insert(m SpendingDynamic) (int, error) {
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
