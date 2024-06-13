package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// DepositPaymentOrder struct
type DepositPaymentOrder struct {
	ID                     int        `db:"id,omitempty"`
	OrganizationUnitID     int        `db:"organization_unit_id"`
	CaseNumber             string     `db:"case_number"`
	SupplierID             int        `db:"supplier_id"`
	SubjectTypeID          int        `db:"subject_type_id"`
	NetAmount              float64    `db:"net_amount"`
	BankAccount            string     `db:"bank_account"`
	SourceBankAccount      string     `db:"source_bank_account"`
	DateOfPayment          time.Time  `db:"date_of_payment"`
	DateOfStatement        *time.Time `db:"date_of_statement"`
	IDOfStatement          *string    `db:"id_of_statement"`
	MunicipalityID         *int       `db:"municipality_id"`
	TaxAuthorityCodebookID *int       `db:"tax_authority_codebook_id"`
	FileID                 *int       `db:"file_id"`
	CreatedAt              time.Time  `db:"created_at,omitempty"`
	UpdatedAt              time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *DepositPaymentOrder) Table() string {
	return "deposit_payment_orders"
}

// GetAll gets all records from the database, using upper
func (t *DepositPaymentOrder) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*DepositPaymentOrder, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*DepositPaymentOrder
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
func (t *DepositPaymentOrder) Get(id int) (*DepositPaymentOrder, error) {
	var one DepositPaymentOrder
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *DepositPaymentOrder) Update(ctx context.Context, tx up.Session, m DepositPaymentOrder) error {
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return err
	}

	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	if err := res.Update(&m); err != nil {
		return err
	}

	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *DepositPaymentOrder) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return err
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(id)
		if err := res.Delete(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *DepositPaymentOrder) Insert(ctx context.Context, tx up.Session, m DepositPaymentOrder) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}

	var id int

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return 0, err
	}

	collection := tx.Collection(t.Table())

	var res up.InsertResult
	var err error

	if res, err = collection.Insert(m); err != nil {
		return 0, err
	}

	id = getInsertId(res.ID())

	return id, nil
}

func (t *DepositPaymentOrder) PayDepositPaymentOrder(ctx context.Context, tx up.Session, id int, IDOfStatement string, DateOfStatement time.Time) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	// Set the user_id variable
	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return err
	}

	query = `update deposit_payment_orders set id_of_statement = $1, date_of_statement = $2 where id = $3`

	_, err := tx.SQL().Query(query, IDOfStatement, DateOfStatement, id)
	if err != nil {
		return err
	}

	return err
}
