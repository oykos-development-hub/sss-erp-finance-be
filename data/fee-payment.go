package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type FeePaymentMethod int

const (
	PaymentFeePeymentMethod    FeePaymentMethod = 1
	ForcedFeePeymentMethod     FeePaymentMethod = 2
	CourtCostsFeePeymentMethod FeePaymentMethod = 3
)

type FeePaymentStatus int

const (
	PaidFeePeymentStatus      FeePaymentStatus = 1
	CancelledFeePeymentStatus FeePaymentStatus = 2
	RetunedFeePeymentStatus   FeePaymentStatus = 3
)

// FeePayment struct
type FeePayment struct {
	ID                     int              `db:"id,omitempty"`
	FeeID                  int              `db:"fee_id"`
	PaymentMethod          FeePaymentMethod `db:"payment_method"`
	Amount                 float64          `db:"amount"`
	PaymentDate            time.Time        `db:"payment_date"`
	PaymentDueDate         time.Time        `db:"payment_due_date"`
	ReceiptNumber          string           `db:"receipt_number"`
	PaymentReferenceNumber string           `db:"payment_reference_number"`
	DebitReferenceNumber   string           `db:"debit_reference_number"`
	Status                 FeePaymentStatus `db:"status"`
	CreatedAt              time.Time        `db:"created_at,omitempty"`
	UpdatedAt              time.Time        `db:"updated_at"`
}

// Table returns the table name
func (t *FeePayment) Table() string {
	return "fee_payments"
}

// Insert inserts a model into the database, using upper
func (t *FeePayment) Insert(m FeePayment) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}

// Get gets one record from the database, by id, using upper
func (t *FeePayment) Get(id int) (*FeePayment, error) {
	var one FeePayment
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// GetAll gets all records from the database, using upper
func (t *FeePayment) GetAll(condition *up.Cond) ([]*FeePayment, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*FeePayment
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
func (t *FeePayment) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Update updates a record in the database, using upper
func (t *FeePayment) Update(m FeePayment) error {
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}
