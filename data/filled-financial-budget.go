package data

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/errors"
)

// FilledFinancialBudget struct
type FilledFinancialBudget struct {
	ID              int                 `db:"id,omitempty"`
	BudgetRequestID int                 `db:"budget_request_id"`
	AccountID       int                 `db:"account_id"`
	CurrentYear     decimal.Decimal     `db:"current_year"`
	NextYear        decimal.Decimal     `db:"next_year"`
	YearAfterNext   decimal.Decimal     `db:"year_after_next"`
	Actual          decimal.NullDecimal `db:"actual"`
	Description     string              `db:"description"`
	CreatedAt       time.Time           `db:"created_at,omitempty"`
	UpdatedAt       time.Time           `db:"updated_at"`
}

// Table returns the table name
func (t *FilledFinancialBudget) Table() string {
	return "filled_financial_budgets"
}

// GetAll gets all records from the database, using upper
func (t *FilledFinancialBudget) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]FilledFinancialBudget, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []FilledFinancialBudget
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
func (t *FilledFinancialBudget) GetSummaryFilledFinancialRequests(budgetID int, requestType RequestType) ([]FilledFinancialBudget, error) {
	var res []FilledFinancialBudget

	const query = `SELECT f.account_id, SUM(f.current_year) AS current_year, SUM(f.next_year) AS next_year, SUM(f.year_after_next) AS year_after_next, SUM(f.actual) AS actual
						FROM filled_financial_budgets f
						JOIN budget_requests r ON f.budget_request_id = r.id
						WHERE r.budget_id = $1 AND r.request_type = $2 AND r.status = 7
						GROUP BY f.account_id;`

	rows, err := Upper.SQL().Query(query, budgetID, requestType)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item FilledFinancialBudget
		err = rows.Scan(
			&item.AccountID,
			&item.CurrentYear,
			&item.NextYear,
			&item.YearAfterNext,
			&item.Actual,
		)

		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil
}

// Get gets one record from the database, by id, using upper
func (t *FilledFinancialBudget) Get(id int) (*FilledFinancialBudget, error) {
	var one FilledFinancialBudget
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FilledFinancialBudget) Update(ctx context.Context, m FilledFinancialBudget) error {
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.ErrUnauthorized
	}

	err := Upper.Tx(func(sess up.Session) error {

		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return err
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(m.ID)
		if err := res.Update(&m); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// Update updates a record in the database, using upper
func (t *FilledFinancialBudget) UpdateActual(ctx context.Context, id int, actual decimal.Decimal) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.ErrUnauthorized
	}

	err := Upper.Tx(func(sess up.Session) error {
		// Set the user_id variable
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return err
		}

		updateQuery := fmt.Sprintf("UPDATE %s SET actual = $1, updated_at = $2 WHERE id = $3", t.Table())

		res, err := sess.SQL().Exec(updateQuery, actual, time.Now(), id)
		if err != nil {
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected != 1 {
			return errors.ErrNotFound
		}

		return err
	})

	return err
}

// Delete deletes a record from the database by id, using upper
func (t *FilledFinancialBudget) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.ErrUnauthorized
	}

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return err
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(id)
		if err := res.Delete(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FilledFinancialBudget) Insert(ctx context.Context, m FilledFinancialBudget) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return 0, errors.ErrUnauthorized
	}

	var id int

	err := Upper.Tx(func(sess up.Session) error {

		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return err
		}

		collection := sess.Collection(t.Table())

		var res up.InsertResult
		var err error

		if res, err = collection.Insert(m); err != nil {
			return err
		}

		id = getInsertId(res.ID())

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
