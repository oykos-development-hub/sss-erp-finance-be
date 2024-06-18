package data

import (
	"context"
	goerrors "errors"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type BudgetType int

const (
	CapitalBudgetType BudgetType = 1
	CurrentBudgetType BudgetType = 2
)

type BudgetStatus int

const (
	BudgetCreatedStatus   BudgetStatus = 1
	BudgetSentStatus      BudgetStatus = 2
	BudgetAcceptedStatus  BudgetStatus = 3
	BudgetCompletedStatus BudgetStatus = 4
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

	err := res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, errors.Wrap(err, "upper all")
	}

	return all, nil
}

// Get gets one record from the database, by id, using upper
func (t *Budget) Get(id int) (*Budget, error) {
	var one Budget
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.WrapNotFoundError(err, "repo get")
		}
		return nil, errors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Budget) Update(ctx context.Context, m Budget) error {
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user id not found in context")
	}

	err := Upper.Tx(func(sess up.Session) error {

		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return errors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(m.ID)
		if err := res.Update(&m); err != nil {
			return errors.Wrap(err, "upper update")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "upper tx")
	}

	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *Budget) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user id not found in context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return errors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(id)
		if err := res.Delete(); err != nil {
			return errors.Wrap(err, "upper delete")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "upper tx")
	}

	return nil
}

// Insert inserts a model into the database, using upper
func (t *Budget) Insert(ctx context.Context, m Budget) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	var id int

	err := Upper.Tx(func(sess up.Session) error {

		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return errors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())

		var res up.InsertResult
		var err error

		if res, err = collection.Insert(m); err != nil {
			return errors.Wrap(err, "upper insert")
		}

		id = getInsertId(res.ID())

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "upper tx")
	}

	return id, nil
}
