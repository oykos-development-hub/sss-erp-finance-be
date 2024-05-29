package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type EnforcedPaymentStatus string

var (
	EnforcedPaymentStatusCreated      EnforcedPaymentStatus = "Kreiran"
	EnforcedPaymentStatusStatusReturn EnforcedPaymentStatus = "Povraćaj"
)

// EnforcedPayment struct
type EnforcedPayment struct {
	ID                 int                   `db:"id,omitempty"`
	OrganizationUnitID int                   `db:"organization_unit_id"`
	SupplierID         int                   `db:"supplier_id"`
	BankAccount        string                `db:"bank_account"`
	DateOfPayment      time.Time             `db:"date_of_payment"`
	DateOfOrder        *time.Time            `db:"date_of_order"`
	IDOfStatement      *int                  `db:"id_of_statement,omitempty"`
	SAPID              *string               `db:"sap_id"`
	Status             EnforcedPaymentStatus `db:"status"`
	Registred          *bool                 `db:"registred,omitempty"`
	RegistredReturn    *bool                 `db:"registred_return,omitempty"`
	ReturnAmount       *float64              `db:"return_amount"`
	DateOfSAP          *time.Time            `db:"date_of_sap"`
	FileID             *int                  `db:"file_id"`
	ReturnFileID       *int                  `db:"return_file_id"`
	ReturnDate         *time.Time            `db:"return_date"`
	Amount             float64               `db:"amount"`
	AmountForLawyer    float64               `db:"amount_for_lawyer"`
	AmountForAgent     float64               `db:"amount_for_agent"`
	Description        string                `db:"description"`
	CreatedAt          time.Time             `db:"created_at,omitempty"`
	UpdatedAt          time.Time             `db:"updated_at"`
}

// Table returns the table name
func (t *EnforcedPayment) Table() string {
	return "enforced_payments"
}

// GetAll gets all records from the database, using upper
func (t *EnforcedPayment) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*EnforcedPayment, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*EnforcedPayment
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
func (t *EnforcedPayment) Get(id int) (*EnforcedPayment, error) {
	var one EnforcedPayment
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *EnforcedPayment) Update(tx up.Session, m EnforcedPayment) error {
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())
	order, err := t.Get(m.ID)

	if err != nil {
		return err
	}

	m.IDOfStatement = order.IDOfStatement

	res := collection.Find(m.ID)
	err = res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

func (t *EnforcedPayment) ReturnEnforcedPayment(tx up.Session, m EnforcedPayment) error {
	m.UpdatedAt = time.Now()
	m.Status = EnforcedPaymentStatusStatusReturn

	query := `update enforced_payments set status = $2, return_date = $3, return_file_id = $4, return_amount = $5 where id = $1`

	_, err := Upper.SQL().Query(query, m.ID, m.Status, m.ReturnDate, m.ReturnFileID, m.ReturnAmount)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *EnforcedPayment) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *EnforcedPayment) Insert(tx up.Session, m EnforcedPayment) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.Status = EnforcedPaymentStatusCreated
	collection := tx.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
