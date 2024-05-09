package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// Salary struct

type SalaryAdditionalExpenseType string

var (
	ContributionsSalaryExpenseType = "contributions"
	TaxesSalaryExpenseType         = "taxes"
	SuspensionsSalaryExpenseType   = "suspensions"
	SubTaxesSalaryExpenseType      = "subtaxes"
	BanksSalaryExpenseType         = "banks"
)

type Salary struct {
	ID                 int       `db:"id,omitempty"`
	ActivityID         int       `db:"activity_id"`
	Month              string    `db:"month"`
	DateOfCalculation  time.Time `db:"date_of_calculation"`
	Description        string    `db:"description"`
	Status             string    `db:"status"`
	Registred          *bool     `db:"registred,omitempty"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	NumberOfEmployees  int       `db:"number_of_employees"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *Salary) Table() string {
	return "salaries"
}

// GetAll gets all records from the database, using upper
func (t *Salary) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*Salary, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*Salary
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
func (t *Salary) Get(id int) (*Salary, error) {
	var one Salary
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Salary) Update(tx up.Session, m Salary) error {
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
func (t *Salary) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *Salary) Insert(tx up.Session, m Salary) (int, error) {
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
