package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// Article struct
type Article struct {
	ID            int       `db:"id,omitempty"`
	Title         string    `db:"title"`
	NetPrice      float64   `db:"net_price"`
	VatPrice      float64   `db:"vat_price"`
	Description   string    `db:"description"`
	InvoiceID     int       `db:"invoice_id"`
	AccountID     int       `db:"account_id"`
	CostAccountID int       `db:"cost_account_id"`
	CreatedAt     time.Time `db:"created_at,omitempty"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *Article) Table() string {
	return "articles"
}

// GetAll gets all records from the database, using upper
func (t *Article) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*Article, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*Article
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
func (t *Article) Get(id int) (*Article, error) {
	var one Article
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Article) Update(m Article) error {
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
func (t *Article) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *Article) Insert(m Article) (int, error) {
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
