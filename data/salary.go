package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// Salary struct

type SalaryAdditionalExpenseType string

var (
	ContributionsSalaryExpenseType = "contributions"
	TaxesSalaryExpenseType         = "taxes"
	SuspensionsSalaryExpenseType   = "suspensions"
	SubTaxesSalaryExpenseType      = "subtaxes"
	BanksSalaryExpenseType         = "banks"
)

type Salary struct {
	ID                 int       `db:"id,omitempty"`
	ActivityID         int       `db:"activity_id"`
	Month              string    `db:"month"`
	DateOfCalculation  time.Time `db:"date_of_calculation"`
	Description        string    `db:"description"`
	Status             string    `db:"status"`
	Registred          *bool     `db:"registred,omitempty"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	NumberOfEmployees  int       `db:"number_of_employees"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *Salary) Table() string {
	return "salaries"
}

// GetAll gets all records from the database, using upper
func (t *Salary) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*Salary, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*Salary
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
func (t *Salary) Get(id int) (*Salary, error) {
	var one Salary
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Salary) Update(ctx context.Context, tx up.Session, m Salary) error {
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
func (t *Salary) Delete(ctx context.Context, id int) error {
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
func (t *Salary) Insert(ctx context.Context, tx up.Session, m Salary) (int, error) {
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
