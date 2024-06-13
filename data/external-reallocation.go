package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type ReallocationStatus string

var (
	ReallocationStatusCreated    ReallocationStatus = "Kreiran"
	ReallocationStatusOUAccept   ReallocationStatus = "Prihvaćen OJ"
	ReallocationStatusOUDecline  ReallocationStatus = "Odbijen OJ"
	ReallocationStatusSSSAccept  ReallocationStatus = "Prihvaćen SSS"
	ReallocationStatusSSSDecline ReallocationStatus = "Odbijen SSS"
)

// ExternalReallocation struct
type ExternalReallocation struct {
	ID                            int                `db:"id,omitempty"`
	Title                         string             `db:"title"`
	Status                        ReallocationStatus `db:"status"`
	SourceOrganizationUnitID      int                `db:"source_organization_unit_id"`
	DestinationOrganizationUnitID int                `db:"destination_organization_unit_id"`
	DateOfRequest                 time.Time          `db:"date_of_request"`
	DateOfActionDestOrgUnit       time.Time          `db:"date_of_action_dest_org_unit"`
	DateOfActionSSS               time.Time          `db:"date_of_action_sss"`
	RequestedBy                   int                `db:"requested_by"`
	AcceptedBy                    int                `db:"accepted_by"`
	FileID                        int                `db:"file_id"`
	DestinationOrgUnitFileID      int                `db:"destination_org_unit_file_id"`
	SSSFileID                     int                `db:"sss_file_id"`
	BudgetID                      int                `db:"budget_id"`
	CreatedAt                     time.Time          `db:"created_at,omitempty"`
	UpdatedAt                     time.Time          `db:"updated_at"`
}

// Table returns the table name
func (t *ExternalReallocation) Table() string {
	return "external_reallocations"
}

// GetAll gets all records from the database, using upper
func (t *ExternalReallocation) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*ExternalReallocation, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*ExternalReallocation
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
func (t *ExternalReallocation) Get(id int) (*ExternalReallocation, error) {
	var one ExternalReallocation
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *ExternalReallocation) Update(ctx context.Context, tx up.Session, m ExternalReallocation) error {
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
func (t *ExternalReallocation) Delete(ctx context.Context, id int) error {
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
func (t *ExternalReallocation) Insert(ctx context.Context, tx up.Session, m ExternalReallocation) (int, error) {
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

func (t *ExternalReallocation) AcceptOUExternalReallocation(ctx context.Context, tx up.Session, m ExternalReallocation) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	// Set the user_id variable
	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return err
	}

	query = `update external_reallocations
		set status = $1, date_of_action_dest_org_unit =$2, accepted_by = $3, destination_org_unit_file_id = $4
		where id = $5`

	_, err := tx.SQL().Query(query, ReallocationStatusOUAccept, m.DateOfActionDestOrgUnit, m.AcceptedBy, m.DestinationOrgUnitFileID, m.ID)

	return err

}

func (t *ExternalReallocation) RejectOUExternalReallocation(ctx context.Context, id int) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		// Set the user_id variable
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return err
		}

		query = `update external_reallocations
			  set status = $1, date_of_action_dest_org_unit = NOW() where id = $2`

		_, err := sess.SQL().Query(query, ReallocationStatusOUDecline, id)

		return err
	})

	return err
}

func (t *ExternalReallocation) AcceptSSSExternalReallocation(ctx context.Context, tx up.Session, id int) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	// Set the user_id variable
	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return err
	}

	query = `update external_reallocations
			  set status = $1, date_of_action_sss = NOW() where id = $2`

	_, err := tx.SQL().Query(query, ReallocationStatusSSSAccept, id)

	return err
}

func (t *ExternalReallocation) RejectSSSExternalReallocation(ctx context.Context, tx up.Session, id int) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}
	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return err
	}

	query = `update external_reallocations
			  set status = $1, date_of_action_sss = NOW() where id = $2`

	_, err := tx.SQL().Query(query, ReallocationStatusSSSDecline, id)

	return err
}
