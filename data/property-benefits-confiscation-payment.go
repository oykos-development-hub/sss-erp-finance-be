package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type PropBenConfPaymentMethod int

const (
	PaymentPropBenConfPeymentMethod    PropBenConfPaymentMethod = 1
	ForcedPropBenConfPeymentMethod     PropBenConfPaymentMethod = 2
	CourtCostsPropBenConfPeymentMethod PropBenConfPaymentMethod = 3
)

type PropBenConfPaymentStatus int

const (
	PaidPropBenConfPeymentStatus      PropBenConfPaymentStatus = 1
	CancelledPropBenConfPeymentStatus PropBenConfPaymentStatus = 2
	RetunedPropBenConfPeymentStatus   PropBenConfPaymentStatus = 3
)

// PropBenConfPayment struct
type PropBenConfPayment struct {
	ID                     int                      `db:"id,omitempty"`
	PropBenConfID          int                      `db:"property_benefits_confiscation_id"`
	PaymentMethod          PropBenConfPaymentMethod `db:"payment_method"`
	Amount                 float64                  `db:"amount"`
	PaymentDate            time.Time                `db:"payment_date"`
	PaymentDueDate         time.Time                `db:"payment_due_date"`
	ReceiptNumber          string                   `db:"receipt_number"`
	PaymentReferenceNumber string                   `db:"payment_reference_number"`
	DebitReferenceNumber   string                   `db:"debit_reference_number"`
	Status                 PropBenConfPaymentStatus `db:"status"`
	CreatedAt              time.Time                `db:"created_at,omitempty"`
	UpdatedAt              time.Time                `db:"updated_at"`
}

// Table returns the table name
func (t *PropBenConfPayment) Table() string {
	return "property_benefits_confiscations_payments"
}

// Insert inserts a model into the database, using upper
func (t *PropBenConfPayment) Insert(m PropBenConfPayment) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}

// Get gets one record from the database, by id, using upper
func (t *PropBenConfPayment) Get(id int) (*PropBenConfPayment, error) {
	var one PropBenConfPayment
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// GetAll gets all records from the database, using upper
func (t *PropBenConfPayment) GetAll(condition *up.Cond) ([]*PropBenConfPayment, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*PropBenConfPayment
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
func (t *PropBenConfPayment) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Update updates a record in the database, using upper
func (t *PropBenConfPayment) Update(m PropBenConfPayment) error {
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}
