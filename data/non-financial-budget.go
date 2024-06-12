package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// NonFinancialBudget struct
type NonFinancialBudget struct {
	ID        int `db:"id,omitempty"`
	RequestID int `db:"request_id"`

	ImplContactFullName     string `db:"impl_contact_fullname"`
	ImplContactWorkingPlace string `db:"impl_contact_working_place"`
	ImplContactPhone        string `db:"impl_contact_phone"`
	ImplContactEmail        string `db:"impl_contact_email"`

	ContactFullName     string `db:"contact_fullname"`
	ContactWorkingPlace string `db:"contact_working_place"`
	ContactPhone        string `db:"contact_phone"`
	ContactEmail        string `db:"contact_email"`

	Statement string `db:"statement"`

	CreatedAt time.Time `db:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *NonFinancialBudget) Table() string {
	return "non_financial_budgets"
}

// GetAll gets all records from the database, using upper
func (t *NonFinancialBudget) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*NonFinancialBudget, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*NonFinancialBudget
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
func (t *NonFinancialBudget) Get(id int) (*NonFinancialBudget, error) {
	var one NonFinancialBudget
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *NonFinancialBudget) Update(ctx context.Context, m NonFinancialBudget) error {
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
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

// Delete deletes a record from the database by id, using upper
func (t *NonFinancialBudget) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
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
func (t *NonFinancialBudget) Insert(ctx context.Context, m NonFinancialBudget) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return 0, errors.New("user ID not found in context")
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
