package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// Activity struct
type Activity struct {
	ID                 int       `db:"id,omitempty"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	SubProgramID       int       `db:"sub_program_id"`
	Title              string    `db:"title"`
	Code               string    `db:"code"`
	Description        string    `db:"description"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *Activity) Table() string {
	return "activities"
}

// GetAll gets all records from the database, using upper
func (t *Activity) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*Activity, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*Activity
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
func (t *Activity) Get(id int) (*Activity, error) {
	var one Activity
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Activity) Update(ctx context.Context, m Activity) error {
	m.UpdatedAt = time.Now()
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
		res := collection.Find(m.ID)
		if err := res.Update(&m); err != nil {
			return newErrors.Wrap(err, "upper update")
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *Activity) Delete(ctx context.Context, id int) error {
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
func (t *Activity) Insert(ctx context.Context, m Activity) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return 0, newErrors.Wrap(err, "contextuitl get user id from context")
	}

	var id int

	err := Upper.Tx(func(sess up.Session) error {

		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return newErrors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())

		var res up.InsertResult
		var err error

		if res, err = collection.Insert(m); err != nil {
			return newErrors.Wrap(err, "upper insert")
		}

		id = getInsertId(res.ID())

		return nil
	})

	if err != nil {
		return 0, newErrors.Wrap(err, "upper tx")
	}

	return id, nil
}
