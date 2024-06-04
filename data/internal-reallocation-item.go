package data

import (
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

// InternalReallocationItem struct
type InternalReallocationItem struct {
	ID                   int             `db:"id,omitempty"`
	ReallocationID       int             `db:"reallocation_id"`
	SourceAccountID      int             `db:"source_account_id"`
	DestinationAccountID int             `db:"destination_account_id"`
	Amount               decimal.Decimal `db:"amount"`
	CreatedAt            time.Time       `db:"created_at,omitempty"`
	UpdatedAt            time.Time       `db:"updated_at"`
}

// Table returns the table name
func (t *InternalReallocationItem) Table() string {
	return "internal_reallocation_items"
}

// GetAll gets all records from the database, using upper
func (t *InternalReallocationItem) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*InternalReallocationItem, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*InternalReallocationItem
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
func (t *InternalReallocationItem) Get(id int) (*InternalReallocationItem, error) {
	var one InternalReallocationItem
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *InternalReallocationItem) Update(tx up.Session, m InternalReallocationItem) error {
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
func (t *InternalReallocationItem) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *InternalReallocationItem) Insert(tx up.Session, m InternalReallocationItem) (int, error) {
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
