package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type FinePaymentMethod int

const (
	PaymentFinePeymentMethod    FinePaymentMethod = 1
	ForcedFinePeymentMethod     FinePaymentMethod = 2
	CourtCostsFinePeymentMethod FinePaymentMethod = 3
)

type FinePaymentStatus int

const (
	PaidFinePeymentStatus      FinePaymentStatus = 1
	CancelledFinePeymentStatus FinePaymentStatus = 2
	RetunedFinePeymentStatus   FinePaymentStatus = 3
)

// FinePayment struct
type FinePayment struct {
	ID                     int               `db:"id,omitempty"`
	FineID                 int               `db:"fine_id"`
	PaymentMethod          FinePaymentMethod `db:"payment_method"`
	Amount                 float64           `db:"amount"`
	PaymentDate            time.Time         `db:"payment_date"`
	PaymentDueDate         time.Time         `db:"payment_due_date"`
	ReceiptNumber          string            `db:"receipt_number"`
	PaymentReferenceNumber string            `db:"payment_reference_number"`
	DebitReferenceNumber   string            `db:"debit_reference_number"`
	Status                 FinePaymentStatus `db:"status"`
	CreatedAt              time.Time         `db:"created_at,omitempty"`
	UpdatedAt              time.Time         `db:"updated_at"`
}

// Table returns the table name
func (t *FinePayment) Table() string {
	return "fine_payments"
}

// Insert inserts a model into the database, using upper
func (t *FinePayment) Insert(ctx context.Context, m FinePayment) (int, error) {
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

// Get gets one record from the database, by id, using upper
func (t *FinePayment) Get(id int) (*FinePayment, error) {
	var one FinePayment
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

// GetAll gets all records from the database, using upper
func (t *FinePayment) GetAll(condition *up.Cond) ([]*FinePayment, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*FinePayment
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

	err = res.OrderBy("created_at desc").All(&all)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper all")
	}

	return all, &total, err
}

// Delete deletes a record from the database by id, using upper
func (t *FinePayment) Delete(ctx context.Context, id int) error {
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

// Update updates a record in the database, using upper
func (t *FinePayment) Update(ctx context.Context, m FinePayment) error {
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
