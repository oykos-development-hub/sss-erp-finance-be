package data

import (
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
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
	Amount                 decimal.Decimal            `db:"amount"`
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
func (t *ProcedureCostPayment) Insert(m ProcedureCostPayment) (int, error) {
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
func (t *ProcedureCostPayment) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Update updates a record in the database, using upper
func (t *ProcedureCostPayment) Update(m ProcedureCostPayment) error {
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}
