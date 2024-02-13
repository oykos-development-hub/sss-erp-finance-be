package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type BudgetRequestStatus int

const (
	BudgetRequestSentStatus     BudgetRequestStatus = 1
	BudgetRequestFinishedStatus BudgetRequestStatus = 2
)

type RequestType int

const (
	DonationFinancialRequestType RequestType = 1
	CurrentFinancialRequestType  RequestType = 2
	NonFinancialRequestType      RequestType = 3
)

// BudgetRequest struct
type BudgetRequest struct {
	ID                 int                 `db:"id,omitempty"`
	OrganizationUnitID int                 `db:"organization_unit_id"`
	BudgetID           int                 `db:"budget_id"`
	RequestType        RequestType         `db:"request_type"`
	Status             BudgetRequestStatus `db:"status"`
	CreatedAt          time.Time           `db:"created_at,omitempty"`
	UpdatedAt          time.Time           `db:"updated_at"`
}

// Table returns the table name
func (t *BudgetRequest) Table() string {
	return "budget_requests"
}

// GetAll gets all records from the database, using upper
func (t *BudgetRequest) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*BudgetRequest, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*BudgetRequest
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
func (t *BudgetRequest) Get(id int) (*BudgetRequest, error) {
	var one BudgetRequest
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *BudgetRequest) Update(m BudgetRequest) error {
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *BudgetRequest) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *BudgetRequest) Insert(m BudgetRequest) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
