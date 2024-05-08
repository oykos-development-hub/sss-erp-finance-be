package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// ModelOfAccountingItem struct
type ModelOfAccountingItem struct {
	ID              int                       `db:"id,omitempty"`
	Title           AccountingOrderItemsTitle `db:"title"`
	ModelID         int                       `db:"model_id"`
	DebitAccountID  int                       `db:"debit_account_id"`
	CreditAccountID int                       `db:"credit_account_id"`
	CreatedAt       time.Time                 `db:"created_at,omitempty"`
	UpdatedAt       time.Time                 `db:"updated_at"`
}

// Table returns the table name
func (t *ModelOfAccountingItem) Table() string {
	return "model_of_accounting_items"
}

// GetAll gets all records from the database, using upper
func (t *ModelOfAccountingItem) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*ModelOfAccountingItem, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*ModelOfAccountingItem
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

	orders = append(orders, "id asc")
	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, err
	}

	return all, &total, err
}

// Get gets one record from the database, by id, using upper
func (t *ModelOfAccountingItem) Get(id int) (*ModelOfAccountingItem, error) {
	var one ModelOfAccountingItem
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *ModelOfAccountingItem) Update(tx up.Session, m ModelOfAccountingItem) error {
	m.UpdatedAt = time.Now()
	query1 := `update model_of_accounting_items set debit_account_id = $1, updated_at = $2 where id = $3`
	query2 := `update model_of_accounting_items set credit_account_id = $1, updated_at = $2 where id = $3`

	if m.DebitAccountID != 0 {
		rows, err := Upper.SQL().Query(query1, m.DebitAccountID, m.UpdatedAt, m.ID)
		if err != nil {
			return err
		}
		defer rows.Close()
	} else {
		rows, err := Upper.SQL().Query(query2, m.CreditAccountID, m.UpdatedAt, m.ID)
		if err != nil {
			return err
		}
		defer rows.Close()
	}

	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *ModelOfAccountingItem) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *ModelOfAccountingItem) Insert(tx up.Session, m ModelOfAccountingItem) (int, error) {
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
