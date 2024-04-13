package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// DepositAdditionalExpense struct
type DepositAdditionalExpense struct {
	ID                   int       `db:"id,omitempty"`
	Title                string    `db:"title"`
	AccountID            int       `db:"account_id"`
	SubjectID            int       `db:"subject_id"`
	BankAccount          string    `db:"bank_account"`
	PaymentOrderID       int       `db:"payment_order_id"`
	PayingPaymentOrderID int       `db:"paying_payment_order_id"`
	OrganizationUnitID   int       `db:"organization_unit_id"`
	Price                float32   `db:"price"`
	Status               string    `db:"status"`
	CreatedAt            time.Time `db:"created_at,omitempty"`
	UpdatedAt            time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *DepositAdditionalExpense) Table() string {
	return "deposit_additional_expenses"
}

// GetAll gets all records from the database, using upper
func (t *DepositAdditionalExpense) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*DepositAdditionalExpense, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*DepositAdditionalExpense
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
func (t *DepositAdditionalExpense) Get(id int) (*DepositAdditionalExpense, error) {
	var one DepositAdditionalExpense
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *DepositAdditionalExpense) Update(tx up.Session, m DepositAdditionalExpense) error {
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
func (t *DepositAdditionalExpense) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *DepositAdditionalExpense) Insert(tx up.Session, m DepositAdditionalExpense) (int, error) {
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
