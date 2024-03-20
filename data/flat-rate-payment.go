package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type FlatRatePaymentMethod int

const (
	PaymentFlatRatePeymentMethod    FlatRatePaymentMethod = 1
	ForcedFlatRatePeymentMethod     FlatRatePaymentMethod = 2
	CourtCostsFlatRatePeymentMethod FlatRatePaymentMethod = 3
)

type FlatRatePaymentStatus int

const (
	PaidFlatRatePeymentStatus      FlatRatePaymentStatus = 1
	CancelledFlatRatePeymentStatus FlatRatePaymentStatus = 2
	RetunedFlatRatePeymentStatus   FlatRatePaymentStatus = 3
)

// FlatRatePayment struct
type FlatRatePayment struct {
	ID                     int                   `db:"id,omitempty"`
	FlatRateID             int                   `db:"flat_rate_id"`
	PaymentMethod          FlatRatePaymentMethod `db:"payment_method"`
	Amount                 float64               `db:"amount"`
	PaymentDate            time.Time             `db:"payment_date"`
	PaymentDueDate         time.Time             `db:"payment_due_date"`
	ReceiptNumber          string                `db:"receipt_number"`
	PaymentReferenceNumber string                `db:"payment_reference_number"`
	DebitReferenceNumber   string                `db:"debit_reference_number"`
	Status                 FlatRatePaymentStatus `db:"status"`
	CreatedAt              time.Time             `db:"created_at,omitempty"`
	UpdatedAt              time.Time             `db:"updated_at"`
}

// Table returns the table name
func (t *FlatRatePayment) Table() string {
	return "flat_rate_payments"
}

// Insert inserts a model into the database, using upper
func (t *FlatRatePayment) Insert(m FlatRatePayment) (int, error) {
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
func (t *FlatRatePayment) Get(id int) (*FlatRatePayment, error) {
	var one FlatRatePayment
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// GetAll gets all records from the database, using upper
func (t *FlatRatePayment) GetAll(condition *up.Cond) ([]*FlatRatePayment, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*FlatRatePayment
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
func (t *FlatRatePayment) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Update updates a record in the database, using upper
func (t *FlatRatePayment) Update(m FlatRatePayment) error {
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}
