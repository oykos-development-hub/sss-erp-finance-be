package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
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
	AmountForBank      float64               `db:"amount_for_bank"`
	AgentID            int                   `db:"agent_id"`
	ExecutionNumber    string                `db:"execution_number"`
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
func (t *EnforcedPayment) Update(ctx context.Context, tx up.Session, m EnforcedPayment) error {
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

	order, err := t.Get(m.ID)

	if err != nil {
		return err
	}

	m.IDOfStatement = order.IDOfStatement
	if err := res.Update(&m); err != nil {
		return err
	}

	return nil
}

func (t *EnforcedPayment) ReturnEnforcedPayment(ctx context.Context, tx up.Session, m EnforcedPayment) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	// Set the user_id variable
	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return err
	}

	m.UpdatedAt = time.Now()
	m.Status = EnforcedPaymentStatusStatusReturn

	query = `update enforced_payments set status = $2, return_date = $3, return_file_id = $4, return_amount = $5 where id = $1`

	_, err := tx.SQL().Query(query, m.ID, m.Status, m.ReturnDate, m.ReturnFileID, m.ReturnAmount)

	return err

}

// Delete deletes a record from the database by id, using upper
func (t *EnforcedPayment) Delete(ctx context.Context, id int) error {
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
func (t *EnforcedPayment) Insert(ctx context.Context, tx up.Session, m EnforcedPayment) (int, error) {
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

	return id, err
}
