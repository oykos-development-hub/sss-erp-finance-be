package data

import (
	"context"
	"fmt"
	"time"

	goerrors "errors"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
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

func (t *CurrentBudget) Vault() decimal.Decimal {
	return t.Actual.Sub(t.Balance)
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
		return nil, nil, errors.Wrap(err, "upper count")
	}

	if page != nil && size != nil {
		res = paginateResult(res, *page, *size)
	}

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, errors.Wrap(err, "upper all")
	}

	return all, &total, nil
}

// GetAll gets all records from the database, using upper
func (t *CurrentBudget) GetBy(condition up.AndExpr) (*CurrentBudget, error) {
	collection := Upper.Collection(t.Table())
	var one CurrentBudget

	res := collection.Find(&condition)

	err := res.One(&one)
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.WrapNotFoundError(err, "upper one")
		}
		return nil, errors.Wrap(err, "upper one")
	}

	return &one, nil
}

// Get gets one record from the database, by id, using upper
func (t *CurrentBudget) Get(id int) (*CurrentBudget, error) {
	var one CurrentBudget
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.WrapNotFoundError(err, "upper one")
		}
		return nil, errors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *CurrentBudget) UpdateActual(ctx context.Context, currentBudgetID int, actual decimal.Decimal) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user id not found in context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		// Set the user_id variable
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return errors.Wrap(err, "upper exec")
		}

		updateQuery := fmt.Sprintf("UPDATE %s SET actual = $1 WHERE id = $2", t.Table())

		res, err := sess.SQL().Exec(updateQuery, actual, currentBudgetID)
		if err != nil {
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected != 1 {
			return errors.NewNotFoundError("upper no rows affected")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "upper tx")
	}

	return nil
}

func (t *CurrentBudget) UpdateBalanceWithTx(ctx context.Context, tx up.Session, currentBudgetID int, balance decimal.Decimal) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user id not found in context")
	}

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return errors.Wrap(err, "upper exec")
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", t.Table())

	res, err := tx.SQL().Exec(updateQuery, balance, currentBudgetID)
	if err != nil {
		return errors.Wrap(err, "upper tx")
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		return errors.NewNotFoundError("upper no rows affected")
	}

	return nil
}

// Update updates a record in the database, using upper
func (t *CurrentBudget) UpdateBalance(currentBudgetID int, balance decimal.Decimal) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", t.Table())

	res, err := Upper.SQL().Exec(updateQuery, balance, currentBudgetID)
	if err != nil {
		return errors.Wrap(err, "upper exec")
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		return errors.NewNotFoundError("upper no rows affected")
	}

	return nil
}

// Insert inserts a model into the database, using upper
func (t *CurrentBudget) Insert(ctx context.Context, m CurrentBudget) (int, error) {
	m.CreatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	var id int

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return errors.New("upper exec")
		}

		collection := sess.Collection(t.Table())

		var res up.InsertResult
		var err error

		if res, err = collection.Insert(m); err != nil {
			return errors.New("upper insert")
		}

		id = getInsertId(res.ID())

		return nil
	})
	if err != nil {
		return 0, errors.New("upper tx")
	}

	return id, nil
}

func (t *CurrentBudget) GetActualCurrentBudget(organizationUnitID int) ([]*CurrentBudget, error) {
	var response []*CurrentBudget

	query := `WITH sorted_data AS (
   			  SELECT *
   			  FROM current_budgets
			  WHERE unit_id = $1
   			  ORDER BY id DESC)
		SELECT budget_id, account_id, actual, balance, initial_actual
		FROM sorted_data
		WHERE budget_id = (SELECT budget_id FROM sorted_data LIMIT 1) and unit_id = $1;`

	rows, err := Upper.SQL().Query(query, organizationUnitID)
	if err != nil {
		return nil, errors.New("upper query")
	}

	defer rows.Close()

	for rows.Next() {
		var item CurrentBudget

		err = rows.Scan(&item.BudgetID, &item.AccountID, &item.Actual, &item.Balance, &item.InitialActual)

		if err != nil {
			return nil, errors.New("sql scan")
		}

		response = append(response, &item)
	}

	return response, nil
}
