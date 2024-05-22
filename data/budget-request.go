package data

import (
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
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
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

func (t *BudgetRequest) GetActual(budgetID, unitID int) (decimal.NullDecimal, error) {
	var actual decimal.NullDecimal

	query := `SELECT FFR.actual
		FROM budget_requests BR
		JOIN filled_financial_budgets FFR ON BR.id = FFR.budget_request_id
		WHERE BR.budget_id = $1 AND BR.organization_unit_id = $2 AND BR.request_type = $3;`

	row, err := Upper.SQL().QueryRow(query, budgetID, unitID, RequestTypeCurrentFinancial)
	if err != nil {
		return actual, err
	}

	err = row.Scan(&actual)
	if err != nil {
		return decimal.NullDecimal{}, err
	}

	return actual, nil
}

// Update updates a record in the database, using upper
func (t *BudgetRequest) Update(m BudgetRequest) error {
	m.UpdatedAt = time.Now()
	collection := Upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *BudgetRequest) Delete(id int) error {
	collection := Upper.Collection(t.Table())
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
	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
