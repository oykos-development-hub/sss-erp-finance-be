package data

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type FeeType int

const (
	LawsuitFeeType  FeeType = 1
	JudgmentFeeType FeeType = 2
)

type FeeStatus int

const (
	UnpaidFeeStatus FeeStatus = 1
	PaidFeeStatus   FeeStatus = 2
	PartFeeStatus   FeeStatus = 3
)

type FeeCategory int

const (
	CopyingFeeCategory FeeCategory = 1
)

// Fee struct
type Fee struct {
	ID                     int           `db:"id,omitempty"`
	FeeTypeID              FeeType       `db:"fee_type_id"`
	FeeSubcategoryID       FeeCategory   `db:"fee_subcategory_id"`
	DecisionNumber         string        `db:"decision_number"`
	DecisionDate           time.Time     `db:"decision_date"`
	Subject                string        `db:"subject"`
	JMBG                   string        `db:"jmbg"`
	Amount                 float64       `db:"amount"`
	PaymentReferenceNumber string        `db:"payment_reference_number"`
	DebitReferenceNumber   string        `db:"debit_reference_number"`
	ExecutionDate          time.Time     `db:"execution_date"`
	PaymentDeadlineDate    time.Time     `db:"payment_deadline_date"`
	Description            string        `db:"description"`
	Status                 FeeStatus     `db:"status"`
	CourtAccountID         *int          `db:"court_account_id"`
	File                   pq.Int64Array `db:"file"`
	CreatedAt              time.Time     `db:"created_at,omitempty"`
	UpdatedAt              time.Time     `db:"updated_at"`
}

// Table returns the table name
func (t *Fee) Table() string {
	return "fees"
}

// Get gets one record from the database, by id, using upper
func (t *Fee) Get(id int) (*Fee, error) {
	var one Fee
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Insert inserts a model into the database, using upper
func (t *Fee) Insert(ctx context.Context, m Fee) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return 0, newErrors.Wrap(err, "contextuitl get user id from context")
	}

	var id int

	err := Upper.Tx(func(sess up.Session) error {

		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return newErrors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())

		var res up.InsertResult
		var err error

		if res, err = collection.Insert(m); err != nil {
			return newErrors.Wrap(err, "upper insert")
		}

		id = getInsertId(res.ID())

		return nil
	})

	if err != nil {
		return 0, newErrors.Wrap(err, "upper tx")
	}

	return id, nil
}

// GetAll gets all records from the database, using upper
func (t *Fee) GetAll(page *int, size *int, condition *up.AndExpr) ([]*Fee, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*Fee
	var res up.Result

	if condition != nil {
		res = collection.Find(condition)
	} else {
		res = collection.Find()
	}
	total, err := res.Count()
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper count")
	}

	if page != nil && size != nil {
		res = paginateResult(res, *page, *size)
	}

	err = res.OrderBy("created_at desc").All(&all)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper all")
	}

	return all, &total, err
}

// Update updates a record in the database, using upper
func (t *Fee) Update(ctx context.Context, m Fee) error {
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	err := Upper.Tx(func(sess up.Session) error {

		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return newErrors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(m.ID)
		if err := res.Update(&m); err != nil {
			return newErrors.Wrap(err, "upper update")
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *Fee) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return newErrors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(id)
		if err := res.Delete(); err != nil {
			return newErrors.Wrap(err, "upper delete")
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}
	return nil
}
