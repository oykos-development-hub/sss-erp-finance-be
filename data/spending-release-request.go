package data

import (
	"context"
	"time"

	up "github.com/upper/db/v4"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type SpendingReleaseStatus string

var (
	SprendingReleaseStatusCreated SpendingReleaseStatus = "Kreiran"
	SpendingReleaseStatusAccept   SpendingReleaseStatus = "Prihvaćen"
	SpendingReleaseStatusFilled   SpendingReleaseStatus = "Popunjen"
)

// SpendingReleaseRequest struct
type SpendingReleaseRequest struct {
	ID                     int                   `db:"id,omitempty"`
	Year                   int                   `db:"year"`
	Month                  int                   `db:"month"`
	OrganizationUnitID     int                   `db:"organization_unit_id"`
	OrganizationUnitFileID int                   `db:"organization_unit_file_id"`
	SSSFileID              int                   `db:"sss_file_id"`
	Status                 SpendingReleaseStatus `db:"status"`
	CreatedAt              time.Time             `db:"created_at,omitempty"`
}

// Table returns the table name
func (t *SpendingReleaseRequest) Table() string {
	return "spending_release_requests"
}

// GetAll gets all records from the database, using upper
func (t *SpendingReleaseRequest) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*SpendingReleaseRequest, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*SpendingReleaseRequest
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
func (t *SpendingReleaseRequest) Get(id int) (*SpendingReleaseRequest, error) {
	var one SpendingReleaseRequest
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *SpendingReleaseRequest) Update(tx up.Session, m SpendingReleaseRequest) error {
	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *SpendingReleaseRequest) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *SpendingReleaseRequest) Insert(tx up.Session, m SpendingReleaseRequest) (int, error) {
	m.CreatedAt = time.Now()
	m.Status = SprendingReleaseStatusCreated
	collection := tx.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}

func (t *SpendingReleaseRequest) AcceptSSSRequest(ctx context.Context, id int, fileID int) error {

	/*userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}*/

	err := Upper.Tx(func(sess up.Session) error {
		// Set the user_id variable
		/*query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return newErrors.Wrap(err, "upper exec")
		}*/

		query := `update spending_release_requests
			  set status = $1, sss_file_id = $2 where id = $3`

		_, err := sess.SQL().Query(query, SpendingReleaseStatusAccept, fileID, id)

		if err != nil {
			return newErrors.Wrap(err, "upper exec")
		}

		return nil
	})

	return err
}
