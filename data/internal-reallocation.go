package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// InternalReallocation struct

type InternalReallocation struct {
	ID                 int       `db:"id,omitempty"`
	Title              string    `db:"title"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	DateOfRequest      time.Time `db:"date_of_request"`
	RequestedBy        int       `db:"requested_by"`
	FileID             int       `db:"file_id"`
	BudgetID           int       `db:"budget_id"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *InternalReallocation) Table() string {
	return "internal_reallocations"
}

// GetAll gets all records from the database, using upper
func (t *InternalReallocation) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*InternalReallocation, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*InternalReallocation
	var res up.Result

	if condition != nil {
		res = collection.Find(condition)
	} else {
		res = collection.Find()
	}
	total, err := res.Count()
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper count")
	}

	if page != nil && size != nil {
		res = paginateResult(res, *page, *size)
	}

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper all")
	}

	return all, &total, err
}

// Get gets one record from the database, by id, using upper
func (t *InternalReallocation) Get(id int) (*InternalReallocation, error) {
	var one InternalReallocation
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *InternalReallocation) Update(ctx context.Context, tx up.Session, m InternalReallocation) error {
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return newErrors.Wrap(err, "upper exec")
	}

	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	if err := res.Update(&m); err != nil {
		return newErrors.Wrap(err, "upper update")
	}

	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *InternalReallocation) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return newErrors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(id)
		if err := res.Delete(); err != nil {
			return newErrors.Wrap(err, "upper delete")
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *InternalReallocation) Insert(ctx context.Context, tx up.Session, m InternalReallocation) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.DateOfRequest = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return 0, newErrors.Wrap(err, "contextuitl get user id from context")
	}

	var id int

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return 0, newErrors.Wrap(err, "upper exec")
	}

	collection := tx.Collection(t.Table())

	var res up.InsertResult
	var err error

	if res, err = collection.Insert(m); err != nil {
		return 0, newErrors.Wrap(err, "upper insert")
	}

	id = getInsertId(res.ID())

	return id, nil
}
