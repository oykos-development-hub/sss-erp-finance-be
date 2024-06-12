package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type ProcedureCostPaymentMethod int

const (
	PaymentProcedureCostPeymentMethod    ProcedureCostPaymentMethod = 1
	ForcedProcedureCostPeymentMethod     ProcedureCostPaymentMethod = 2
	CourtCostsProcedureCostPeymentMethod ProcedureCostPaymentMethod = 3
)

type ProcedureCostPaymentStatus int

const (
	PaidProcedureCostPeymentStatus      ProcedureCostPaymentStatus = 1
	CancelledProcedureCostPeymentStatus ProcedureCostPaymentStatus = 2
	RetunedProcedureCostPeymentStatus   ProcedureCostPaymentStatus = 3
)

// ProcedureCostPayment struct
type ProcedureCostPayment struct {
	ID                     int                        `db:"id,omitempty"`
	ProcedureCostID        int                        `db:"procedure_cost_id"`
	PaymentMethod          ProcedureCostPaymentMethod `db:"payment_method"`
	Amount                 float64                    `db:"amount"`
	PaymentDate            time.Time                  `db:"payment_date"`
	PaymentDueDate         time.Time                  `db:"payment_due_date"`
	ReceiptNumber          string                     `db:"receipt_number"`
	PaymentReferenceNumber string                     `db:"payment_reference_number"`
	DebitReferenceNumber   string                     `db:"debit_reference_number"`
	Status                 ProcedureCostPaymentStatus `db:"status"`
	CreatedAt              time.Time                  `db:"created_at,omitempty"`
	UpdatedAt              time.Time                  `db:"updated_at"`
}

// Table returns the table name
func (t *ProcedureCostPayment) Table() string {
	return "procedure_cost_payments"
}

// Insert inserts a model into the database, using upper
func (t *ProcedureCostPayment) Insert(ctx context.Context, m ProcedureCostPayment) (int, error) {
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

// Get gets one record from the database, by id, using upper
func (t *ProcedureCostPayment) Get(id int) (*ProcedureCostPayment, error) {
	var one ProcedureCostPayment
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// GetAll gets all records from the database, using upper
func (t *ProcedureCostPayment) GetAll(condition *up.Cond) ([]*ProcedureCostPayment, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*ProcedureCostPayment
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

	err = res.OrderBy("created_at desc").All(&all)
	if err != nil {
		return nil, nil, err
	}

	return all, &total, err
}

// Delete deletes a record from the database by id, using upper
func (t *ProcedureCostPayment) Delete(ctx context.Context, id int) error {
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

// Update updates a record in the database, using upper
func (t *ProcedureCostPayment) Update(ctx context.Context, m ProcedureCostPayment) error {
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
