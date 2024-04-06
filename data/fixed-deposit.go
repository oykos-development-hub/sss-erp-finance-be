package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// FixedDeposit struct
type FixedDeposit struct {
	ID                   int        `db:"id,omitempty"`
	OrganizationUnitID   int        `db:"organization_unit_id"`
	JudgeID              int        `db:"judge_id"`
	Subject              string     `db:"subject"`
	CaseNumber           string     `db:"case_number"`
	DateOfRecipiet       *time.Time `db:"date_of_recipiet"`
	DateOfCase           *time.Time `db:"date_of_case"`
	DateOfFinality       *time.Time `db:"date_of_finality"`
	DateOfEnforceability *time.Time `db:"date_of_enforceability"`
	DateOfEnd            *time.Time `db:"date_of_end"`
	AccountID            int        `db:"account_id"`
	FileID               int        `db:"file_id"`
	Status               string     `db:"status"`
	Type                 string     `db:"type"`
	CreatedAt            time.Time  `db:"created_at,omitempty"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *FixedDeposit) Table() string {
	return "fixed_deposits"
}

// GetAll gets all records from the database, using upper
func (t *FixedDeposit) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*FixedDeposit, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*FixedDeposit
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
func (t *FixedDeposit) Get(id int) (*FixedDeposit, error) {
	var one FixedDeposit
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FixedDeposit) Update(tx up.Session, m FixedDeposit) error {
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
func (t *FixedDeposit) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FixedDeposit) Insert(tx up.Session, m FixedDeposit) (int, error) {
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
