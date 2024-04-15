package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// PaymentOrder struct
type PaymentOrder struct {
	ID                 int        `db:"id,omitempty"`
	OrganizationUnitID int        `db:"organization_unit_id"`
	SupplierID         int        `db:"supplier_id"`
	BankAccount        string     `db:"bank_account"`
	DateOfPayment      time.Time  `db:"date_of_payment"`
	IDOfStatement      *string    `db:"id_of_statement"`
	SAPID              *string    `db:"sap_id"`
	DateOfSAP          *time.Time `db:"date_of_sap"`
	FileID             *int       `db:"file_id"`
	Amount             float64    `db:"amount"`
	CreatedAt          time.Time  `db:"created_at,omitempty"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *PaymentOrder) Table() string {
	return "payment_orders"
}

// GetAll gets all records from the database, using upper
func (t *PaymentOrder) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*PaymentOrder, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*PaymentOrder
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
func (t *PaymentOrder) Get(id int) (*PaymentOrder, error) {
	var one PaymentOrder
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *PaymentOrder) Update(tx up.Session, m PaymentOrder) error {
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *PaymentOrder) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *PaymentOrder) Insert(tx up.Session, m PaymentOrder) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
