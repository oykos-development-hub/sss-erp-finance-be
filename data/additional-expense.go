package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type AdditionalExpenseStatus int

const (
	AdditionalExpenseStatusCreated        AdditionalExpenseStatus = 1
	AdditionalExpenseStatusWaitingPayment AdditionalExpenseStatus = 2
	AdditionalExpenseStatusPaid           AdditionalExpenseStatus = 3
)

// AdditionalExpense struct
type AdditionalExpense struct {
	ID          int                     `db:"id,omitempty"`
	Title       string                  `db:"title"`
	AccountID   int                     `db:"account_id"`
	SubjectID   int                     `db:"subject_id"`
	BankAccount int                     `db:"bank_account"`
	InvoiceID   int                     `db:"invoice_id"`
	Price       float32                 `db:"price"`
	Status      AdditionalExpenseStatus `db:"status"`
	CreatedAt   time.Time               `db:"created_at,omitempty"`
	UpdatedAt   time.Time               `db:"updated_at"`
}

// Table returns the table name
func (t *AdditionalExpense) Table() string {
	return "additional_expenses"
}

// GetAll gets all records from the database, using upper
func (t *AdditionalExpense) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*AdditionalExpense, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*AdditionalExpense
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
func (t *AdditionalExpense) Get(id int) (*AdditionalExpense, error) {
	var one AdditionalExpense
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *AdditionalExpense) Update(m AdditionalExpense) error {
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *AdditionalExpense) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *AdditionalExpense) Insert(m AdditionalExpense) (int, error) {
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