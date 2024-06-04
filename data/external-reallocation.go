package data

import (
	"time"

	up "github.com/upper/db/v4"
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
func (t *ExternalReallocation) Update(tx up.Session, m ExternalReallocation) error {
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
func (t *ExternalReallocation) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *ExternalReallocation) Insert(tx up.Session, m ExternalReallocation) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.Status = ReallocationStatusCreated
	collection := tx.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
