package data

import (
	"time"

	up "github.com/upper/db/v4"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// EnforcedPaymentItem struct
type EnforcedPaymentItem struct {
	ID             int       `db:"id,omitempty"`
	PaymentOrderID int       `db:"payment_order_id"`
	Title          string    `db:"title"`
	Amount         float32   `db:"amount"`
	InvoiceID      *int      `db:"invoice_id"`
	AccountID      int       `db:"account_id"`
	CreatedAt      time.Time `db:"created_at,omitempty"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *EnforcedPaymentItem) Table() string {
	return "enforced_payment_items"
}

// GetAll gets all records from the database, using upper
func (t *EnforcedPaymentItem) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*EnforcedPaymentItem, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*EnforcedPaymentItem
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

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper all")
	}

	return all, &total, err
}

// Get gets one record from the database, by id, using upper
func (t *EnforcedPaymentItem) Get(id int) (*EnforcedPaymentItem, error) {
	var one EnforcedPaymentItem
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}

	return &one, nil
}

// Update updates a record in the database, using upper
func (t *EnforcedPaymentItem) Update(tx up.Session, m EnforcedPaymentItem) error {
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return newErrors.Wrap(err, "upper update")
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *EnforcedPaymentItem) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return newErrors.Wrap(err, "upper delete")
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *EnforcedPaymentItem) Insert(tx up.Session, m EnforcedPaymentItem) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, newErrors.Wrap(err, "upper insert")
	}

	id := getInsertId(res.ID())

	return id, nil
}
