package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// FixedDepositItem struct
type FixedDepositItem struct {
	ID                 int        `db:"id,omitempty"`
	DepositID          int        `db:"deposit_id"`
	CategoryID         int        `db:"category_id"`
	TypeID             int        `db:"type_id"`
	UnitID             int        `db:"unit_id"`
	Amount             float32    `db:"amount"`
	CurencyID          int        `db:"curency_id"`
	SerialNumber       string     `db:"serial_number"`
	DateOfConfiscation *time.Time `db:"date_of_confiscation"`
	CaseNumber         string     `db:"case_number"`
	JudgeID            int        `db:"judge_id"`
	CreatedAt          time.Time  `db:"created_at,omitempty"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *FixedDepositItem) Table() string {
	return "fixed_deposit_items"
}

// GetAll gets all records from the database, using upper
func (t *FixedDepositItem) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*FixedDepositItem, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*FixedDepositItem
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
func (t *FixedDepositItem) Get(id int) (*FixedDepositItem, error) {
	var one FixedDepositItem
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FixedDepositItem) Update(tx up.Session, m FixedDepositItem) error {
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
func (t *FixedDepositItem) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FixedDepositItem) Insert(tx up.Session, m FixedDepositItem) (int, error) {
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
