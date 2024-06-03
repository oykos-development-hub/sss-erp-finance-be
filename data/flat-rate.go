package data

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type FlatRateType int

const (
	DissisionType FlatRateType = 1
	VerdictType   FlatRateType = 2
)

type FlatRateStatus int

const (
	UnpaidFlatRateStatus FlatRateStatus = 1
	PaidFlatRateStatus   FlatRateStatus = 2
	PartFlatRateStatus   FlatRateStatus = 3
)

const (
	FlatRateGracePeriod = 7
)

// FlatRate struct
type FlatRate struct {
	ID                     int              `db:"id,omitempty"`
	FlatRateType           FlatRateType     `db:"flat_rate_type"`
	DecisionNumber         string           `db:"decision_number"`
	DecisionDate           time.Time        `db:"decision_date"`
	Subject                string           `db:"subject"`
	JMBG                   string           `db:"jmbg"`
	Residence              string           `db:"residence"`
	Amount                 decimal.Decimal  `db:"amount"`
	PaymentReferenceNumber string           `db:"payment_reference_number"`
	DebitReferenceNumber   string           `db:"debit_reference_number"`
	AccountID              int              `db:"account_id"`
	ExecutionDate          time.Time        `db:"execution_date"`
	PaymentDeadlineDate    time.Time        `db:"payment_deadline_date"`
	Description            string           `db:"description"`
	Status                 FlatRateStatus   `db:"status"`
	CourtCosts             *decimal.Decimal `db:"court_costs"`
	CourtAccountID         *int             `db:"court_account_id"`
	File                   pq.Int64Array    `db:"file"`
	CreatedAt              time.Time        `db:"created_at,omitempty"`
	UpdatedAt              time.Time        `db:"updated_at"`
}

// Table returns the table name
func (t *FlatRate) Table() string {
	return "flat_rates"
}

// Get gets one record from the database, by id, using upper
func (t *FlatRate) Get(id int) (*FlatRate, error) {
	var one FlatRate
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Insert inserts a model into the database, using upper
func (t *FlatRate) Insert(m FlatRate) (int, error) {
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

// GetAll gets all records from the database, using upper
func (t *FlatRate) GetAll(page *int, size *int, condition *up.AndExpr) ([]*FlatRate, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*FlatRate
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

	err = res.OrderBy("created_at desc").All(&all)
	if err != nil {
		return nil, nil, err
	}

	return all, &total, err
}

// Update updates a record in the database, using upper
func (t *FlatRate) Update(m FlatRate) error {
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
func (t *FlatRate) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}
