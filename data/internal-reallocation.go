package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// InternalReallocation struct

type InternalReallocation struct {
	ID                 int       `db:"id,omitempty"`
	Title              string    `db:"title"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	DateOfRequest      time.Time `db:"date_of_request"`
	RequestedBy        int       `db:"requested_by"`
	FileID             int       `db:"file_id"`
	BudgetID           int       `db:"budget_id"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *InternalReallocation) Table() string {
	return "internal_reallocations"
}

// GetAll gets all records from the database, using upper
func (t *InternalReallocation) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*InternalReallocation, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*InternalReallocation
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
func (t *InternalReallocation) Get(id int) (*InternalReallocation, error) {
	var one InternalReallocation
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *InternalReallocation) Update(tx up.Session, m InternalReallocation) error {
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
func (t *InternalReallocation) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *InternalReallocation) Insert(tx up.Session, m InternalReallocation) (int, error) {
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
