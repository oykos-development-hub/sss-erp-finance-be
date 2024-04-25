package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type TypesOfModels string

var (
	TypeInvoice  TypesOfModels = "invoices"
	TypeContract TypesOfModels = "contracts"
	TypeDecision TypesOfModels = "decisions"
	TypeSalary   TypesOfModels = "salaries"
)

// ModelsOfAccounting struct
type ModelsOfAccounting struct {
	ID        int           `db:"id,omitempty"`
	Title     string        `db:"title"`
	Type      TypesOfModels `db:"type"`
	CreatedAt time.Time     `db:"created_at,omitempty"`
	UpdatedAt time.Time     `db:"updated_at"`
}

// Table returns the table name
func (t *ModelsOfAccounting) Table() string {
	return "models_of_accounting"
}

// GetAll gets all records from the database, using upper
func (t *ModelsOfAccounting) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*ModelsOfAccounting, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*ModelsOfAccounting
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
func (t *ModelsOfAccounting) Get(id int) (*ModelsOfAccounting, error) {
	var one ModelsOfAccounting
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *ModelsOfAccounting) Update(tx up.Session, m ModelsOfAccounting) error {
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
func (t *ModelsOfAccounting) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *ModelsOfAccounting) Insert(tx up.Session, m ModelsOfAccounting) (int, error) {
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
