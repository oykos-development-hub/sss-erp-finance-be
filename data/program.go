package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// Program struct
type Program struct {
	ID          int       `db:"id,omitempty"`
	Title       string    `db:"title"`
	Code        string    `db:"code"`
	Description string    `db:"description"`
	ParentID    int       `db:"parent_id"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *Program) Table() string {
	return "programs"
}

// GetAll gets all records from the database, using upper
func (t *Program) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*Program, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*Program
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
func (t *Program) Get(id int) (*Program, error) {
	var one Program
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Program) Update(m Program) error {
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
func (t *Program) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *Program) Insert(m Program) (int, error) {
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