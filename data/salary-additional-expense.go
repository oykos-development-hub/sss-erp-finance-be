package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// SalaryAdditionalExpense struct
type SalaryAdditionalExpense struct {
	ID                  int                         `db:"id,omitempty"`
	SalaryID            int                         `db:"salary_id"`
	AccountID           int                         `db:"account_id"`
	Amount              float64                     `db:"amount"`
	SubjectID           int                         `db:"subject_id"`
	BankAccount         string                      `db:"bank_account"`
	Status              InvoiceStatus               `db:"status"`
	OrganizationUnitID  int                         `db:"organization_unit_id"`
	DebtorID            int                         `db:"debtor_id"`
	Type                SalaryAdditionalExpenseType `db:"type"`
	Title               AccountingOrderItemsTitle   `db:"title"`
	IdentificatorNumber string                      `db:"identificator_number"`
	CreatedAt           time.Time                   `db:"created_at,omitempty"`
	UpdatedAt           time.Time                   `db:"updated_at"`
}

// Table returns the table name
func (t *SalaryAdditionalExpense) Table() string {
	return "salary_additional_expenses"
}

// GetAll gets all records from the database, using upper
func (t *SalaryAdditionalExpense) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*SalaryAdditionalExpense, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*SalaryAdditionalExpense
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
func (t *SalaryAdditionalExpense) Get(id int) (*SalaryAdditionalExpense, error) {
	var one SalaryAdditionalExpense
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *SalaryAdditionalExpense) Update(tx up.Session, m SalaryAdditionalExpense) error {
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
func (t *SalaryAdditionalExpense) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *SalaryAdditionalExpense) Insert(tx up.Session, m SalaryAdditionalExpense) (int, error) {
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
