package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// NonFinancialBudget struct
type NonFinancialBudget struct {
	ID         int `db:"id,omitempty"`
	BudetID    int `db:"budget_id"`
	ActivityID int `db:"activity_id"`

	ImplContactFullName     string `db:"impl_contact_fullname"`
	ImplContactWorkingPlace string `db:"impl_contact_working_place"`
	ImplContactPhone        string `db:"impl_contact_phone"`
	ImplContactEmail        string `db:"impl_contact_email"`

	ContactFullName     string `db:"contact_fullname"`
	ContactWorkingPlace string `db:"contact_working_place"`
	ContactPhone        string `db:"contact_phone"`
	ContactEmail        string `db:"contact_email"`

	CreatedAt time.Time `db:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *NonFinancialBudget) Table() string {
	return "non_financial_budgets"
}

// GetAll gets all records from the database, using upper
func (t *NonFinancialBudget) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*NonFinancialBudget, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*NonFinancialBudget
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
func (t *NonFinancialBudget) Get(id int) (*NonFinancialBudget, error) {
	var one NonFinancialBudget
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *NonFinancialBudget) Update(m NonFinancialBudget) error {
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
func (t *NonFinancialBudget) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *NonFinancialBudget) Insert(m NonFinancialBudget) (int, error) {
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
