package data

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type BudgetRequestStatus int

const (
	BudgetRequestSentStatus            BudgetRequestStatus = 1
	BudgetRequestFilledStatus          BudgetRequestStatus = 2
	BudgetRequestSentOnReviewStatus    BudgetRequestStatus = 3
	BudgetRequestAcceptedStatus        BudgetRequestStatus = 4
	BudgetRequestRejectedStatus        BudgetRequestStatus = 5
	BudgetRequestWaitingForActual      BudgetRequestStatus = 6
	BudgetRequestCompletedActualStatus BudgetRequestStatus = 7
)

type RequestType int

const (
	RequestTypeGeneral           RequestType = 1
	RequestTypeNonFinancial      RequestType = 2
	RequestTypeFinancial         RequestType = 3
	RequestTypeCurrentFinancial  RequestType = 4
	RequestTypeDonationFinancial RequestType = 5
)

// BudgetRequest struct
type BudgetRequest struct {
	ID                 int                 `db:"id,omitempty"`
	ParentID           *int                `db:"parent_id"`
	OrganizationUnitID int                 `db:"organization_unit_id"`
	BudgetID           int                 `db:"budget_id"`
	RequestType        RequestType         `db:"request_type"`
	Status             BudgetRequestStatus `db:"status"`
	Comment            string              `db:"comment"`
	CreatedAt          time.Time           `db:"created_at,omitempty"`
	UpdatedAt          time.Time           `db:"updated_at"`
}

// Table returns the table name
func (t *BudgetRequest) Table() string {
	return "budget_requests"
}

// GetAll gets all records from the database, using upper
func (t *BudgetRequest) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*BudgetRequest, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*BudgetRequest
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
func (t *BudgetRequest) Get(id int) (*BudgetRequest, error) {
	var one BudgetRequest
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

func (t *BudgetRequest) GetActual(budgetID, unitID, accountID int) (decimal.NullDecimal, error) {
	var actual decimal.NullDecimal

	query := `SELECT FFb.actual
		FROM budget_requests BR
		JOIN filled_financial_budgets FFB ON BR.id = FFB.budget_request_id 
		WHERE BR.budget_id = $1 AND BR.organization_unit_id = $2 AND FFB.account_id = $3 AND BR.request_type = $4;`

	row, err := Upper.SQL().QueryRow(query, budgetID, unitID, accountID, RequestTypeCurrentFinancial)
	if err != nil {
		return actual, err
	}

	err = row.Scan(&actual)
	if err != nil {
		return actual, err
	}

	return actual, nil
}

// Update updates a record in the database, using upper
func (t *BudgetRequest) Update(ctx context.Context, m BudgetRequest) error {
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
func (t *BudgetRequest) Delete(ctx context.Context, id int) error {
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
func (t *BudgetRequest) Insert(ctx context.Context, m BudgetRequest) (int, error) {
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
