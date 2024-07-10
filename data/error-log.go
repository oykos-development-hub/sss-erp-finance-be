package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type HandlersEntities string

var (
	AccountingEntryItemEntity HandlersEntities = "AccountingEntryItem"
)

// ErrorLog struct
type ErrorLog struct {
	ID        int              `db:"id,omitempty"`
	Error     string           `db:"error"`
	Code      int              `db:"code"`
	Entity    HandlersEntities `db:"entity"`
	CreatedAt time.Time        `db:"created_at,omitempty"`
	UpdatedAt time.Time        `db:"updated_at"`
}

// Table returns the table name
func (t *ErrorLog) Table() string {
	return "error_logs"
}

// GetAll gets all records from the database, using upper
func (t *ErrorLog) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*ErrorLog, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*ErrorLog
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
func (t *ErrorLog) Get(id int) (*ErrorLog, error) {
	var one ErrorLog
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *ErrorLog) Update(tx up.Session, m ErrorLog) error {
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
func (t *ErrorLog) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *ErrorLog) Insert(tx up.Session, m ErrorLog) (int, error) {
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
