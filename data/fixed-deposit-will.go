package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// FixedDepositWill struct
type FixedDepositWill struct {
	ID                 int        `db:"id,omitempty"`
	OrganizationUnitID int        `db:"organization_unit_id"`
	Subject            string     `db:"subject"`
	FatherName         string     `db:"father_name"`
	DateOfBirth        time.Time  `db:"date_of_birth"`
	JMBG               string     `db:"jmbg"`
	CaseNumberSI       string     `db:"case_number_si"`
	CaseNumberRS       string     `db:"case_number_rs"`
	DateOfReceiptSI    *time.Time `db:"date_of_receipt_si"`
	DateOfReceiptRS    *time.Time `db:"date_of_receipt_rs"`
	DateOfEnd          *time.Time `db:"date_of_end"`
	Status             string     `db:"status"`
	FileID             int        `db:"file_id"`
	CreatedAt          time.Time  `db:"created_at,omitempty"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *FixedDepositWill) Table() string {
	return "fixed_deposit_wills"
}

// GetAll gets all records from the database, using upper
func (t *FixedDepositWill) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*FixedDepositWill, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*FixedDepositWill
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
func (t *FixedDepositWill) Get(id int) (*FixedDepositWill, error) {
	var one FixedDepositWill
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *FixedDepositWill) Update(tx up.Session, m FixedDepositWill) error {
	m.UpdatedAt = time.Now()

	if m.DateOfEnd != nil {
		m.Status = "Zaključen"
	}

	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *FixedDepositWill) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *FixedDepositWill) Insert(tx up.Session, m FixedDepositWill) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	if m.DateOfEnd != nil {
		m.Status = "Zaključen"
	}
	collection := tx.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
