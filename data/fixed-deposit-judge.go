package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// FixedDepositJudge struct
type FixedDepositJudge struct {
	ID          int        `db:"id,omitempty"`
	JudgeID     int        `db:"judge_id"`
	DepositID   int        `db:"deposit_id"`
	WillID      int        `db:"will_id"`
	DateOfStart time.Time  `db:"date_of_start"`
	DateOfEnd   *time.Time `db:"date_of_end"`
	FileID      int        `db:"file_id"`
	CreatedAt   time.Time  `db:"created_at,omitempty"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *FixedDepositJudge) Table() string {
	return "fixed_deposit_judges"
}

// GetAll gets all records from the database, using upper
func (t *FixedDepositJudge) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*FixedDepositJudge, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*FixedDepositJudge
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
func (t *FixedDepositJudge) Get(id int) (*FixedDepositJudge, error) {
	var one FixedDepositJudge
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FixedDepositJudge) Update(tx up.Session, m FixedDepositJudge) error {
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
func (t *FixedDepositJudge) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FixedDepositJudge) Insert(tx up.Session, m FixedDepositJudge) (int, error) {
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
