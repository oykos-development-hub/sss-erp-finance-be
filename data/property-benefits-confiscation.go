package data

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

type PropBenConfType int

const (
	DissisionPropBenConfType PropBenConfType = 1
	VerdictPropBenConfType   PropBenConfType = 2
)

type PropBenConfStatus int

const (
	UnpaidPropBenConfStatus PropBenConfStatus = 1
	PaidPropBenConfStatus   PropBenConfStatus = 2
	PartPropBenConfStatus   PropBenConfStatus = 3
)

const (
	PropBenConfGracePeriod = 7
)

// PropBenConf struct
type PropBenConf struct {
	ID                     int               `db:"id,omitempty"`
	PropBenConfType        PropBenConfType   `db:"property_benefits_confiscation_type"`
	DecisionNumber         string            `db:"decision_number"`
	DecisionDate           time.Time         `db:"decision_date"`
	Subject                string            `db:"subject"`
	JMBG                   string            `db:"jmbg"`
	Residence              string            `db:"residence"`
	Amount                 decimal.Decimal   `db:"amount"`
	PaymentReferenceNumber string            `db:"payment_reference_number"`
	DebitReferenceNumber   string            `db:"debit_reference_number"`
	AccountID              int               `db:"account_id"`
	ExecutionDate          time.Time         `db:"execution_date"`
	PaymentDeadlineDate    time.Time         `db:"payment_deadline_date"`
	Description            string            `db:"description"`
	Status                 PropBenConfStatus `db:"status"`
	CourtCosts             *decimal.Decimal  `db:"court_costs"`
	CourtAccountID         *int              `db:"court_account_id"`
	File                   pq.Int64Array     `db:"file"`
	CreatedAt              time.Time         `db:"created_at,omitempty"`
	UpdatedAt              time.Time         `db:"updated_at"`
}

// Table returns the table name
func (t *PropBenConf) Table() string {
	return "property_benefits_confiscations"
}

// Get gets one record from the database, by id, using upper
func (t *PropBenConf) Get(id int) (*PropBenConf, error) {
	var one PropBenConf
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Insert inserts a model into the database, using upper
func (t *PropBenConf) Insert(m PropBenConf) (int, error) {
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
func (t *PropBenConf) GetAll(page *int, size *int, condition *up.AndExpr) ([]*PropBenConf, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*PropBenConf
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
func (t *PropBenConf) Update(m PropBenConf) error {
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
func (t *PropBenConf) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}
