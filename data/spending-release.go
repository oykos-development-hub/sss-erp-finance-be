package data

import (
	goerrors "errors"
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// SpendingRelease struct
type SpendingRelease struct {
	ID              int             `db:"id,omitempty"`
	CurrentBudgetID int             `db:"current_budget_id"`
	Year            int             `db:"year"`
	Month           int             `db:"month"`
	Value           decimal.Decimal `db:"value"`
	CreatedAt       time.Time       `db:"created_at,omitempty"`
}

// Table returns the table name
func (t *SpendingRelease) Table() string {
	return "spending_releases"
}

// ValidateNewEntry validates if not expired
func (t *SpendingRelease) ValidateNewRelease() bool {
	now := time.Now()

	// TODO: add validations for year and first 5 days if needed.
	// day := now.Day()
	// year := now.Year()
	currentMonth := int(now.Month())

	return t.Month == currentMonth
}

// GetAll gets all records from the database, using upper
func (t *SpendingRelease) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*SpendingRelease, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*SpendingRelease
	var res up.Result

	if condition != nil {
		res = collection.Find(condition)
	} else {
		res = collection.Find()
	}
	total, err := res.Count()
	if err != nil {
		return nil, nil, errors.Wrap(err, "repo spending-release get all")
	}

	if page != nil && size != nil {
		res = paginateResult(res, *page, *size)
	}

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, errors.Wrap(err, "repo spending-release getAll")
	}

	return all, &total, err
}

// GetAll gets all records from the database, using upper
func (t *SpendingRelease) GetBy(condition up.AndExpr) (*SpendingRelease, error) {
	collection := Upper.Collection(t.Table())
	var one SpendingRelease

	res := collection.Find(&condition)

	err := res.One(&one)
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.WrapNotFoundError(err, "repo spending-release getBy")
		}
		return nil, errors.Wrap(err, "repo spending-release getBy")
	}

	return &one, err
}

// Get gets one record from the database, by id, using upper
func (t *SpendingRelease) Get(id int) (*SpendingRelease, error) {
	var one SpendingRelease
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.WrapNotFoundError(err, "repo spending-release get")
		}
		return nil, errors.Wrap(err, "repo spending-release get")
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *SpendingRelease) Update(m SpendingRelease) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return errors.Wrap(err, "repo spending-release update")
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *SpendingRelease) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		if goerrors.Is(err, up.ErrNoMoreRows) {
			return errors.WrapNotFoundError(err, "repo spending-release delete")
		}
		return errors.Wrap(err, "repo spending-release delete")
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *SpendingRelease) Insert(m SpendingRelease) (int, error) {
	m.CreatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, errors.Wrap(err, "repo spending-release insert")
	}

	id := getInsertId(res.ID())

	return id, nil
}
