package data

import (
	"time"

	"github.com/lib/pq"
	up "github.com/upper/db/v4"
)

type ProcedureCostType int

const (
	DissisionProcedureCostType ProcedureCostType = 1
	VerdictProcedureCostType   ProcedureCostType = 2
)

type ProcedureCostStatus int

const (
	UnpaidProcedureCostStatus ProcedureCostStatus = 1
	PaidProcedureCostStatus   ProcedureCostStatus = 2
	PartProcedureCostStatus   ProcedureCostStatus = 3
)

const (
	ProcedureCostGracePeriod = 7
)

// ProcedureCost struct
type ProcedureCost struct {
	ID                     int                 `db:"id,omitempty"`
	ProcedureCostType      ProcedureCostType   `db:"procedure_cost_type"`
	DecisionNumber         string              `db:"decision_number"`
	DecisionDate           time.Time           `db:"decision_date"`
	Subject                string              `db:"subject"`
	JMBG                   string              `db:"jmbg"`
	Residence              string              `db:"residence"`
	Amount                 float64             `db:"amount"`
	PaymentReferenceNumber string              `db:"payment_reference_number"`
	DebitReferenceNumber   string              `db:"debit_reference_number"`
	AccountID              int                 `db:"account_id"`
	ExecutionDate          time.Time           `db:"execution_date"`
	PaymentDeadlineDate    time.Time           `db:"payment_deadline_date"`
	Description            string              `db:"description"`
	Status                 ProcedureCostStatus `db:"status"`
	CourtCosts             *float64            `db:"court_costs"`
	CourtAccountID         *int                `db:"court_account_id"`
	File                   pq.Int64Array       `db:"file"`
	CreatedAt              time.Time           `db:"created_at,omitempty"`
	UpdatedAt              time.Time           `db:"updated_at"`
}

// Table returns the table name
func (t *ProcedureCost) Table() string {
	return "procedure_costs"
}

// Get gets one record from the database, by id, using upper
func (t *ProcedureCost) Get(id int) (*ProcedureCost, error) {
	var one ProcedureCost
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Insert inserts a model into the database, using upper
func (t *ProcedureCost) Insert(m ProcedureCost) (int, error) {
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
func (t *ProcedureCost) GetAll(page *int, size *int, condition *up.AndExpr) ([]*ProcedureCost, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*ProcedureCost
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
func (t *ProcedureCost) Update(m ProcedureCost) error {
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
func (t *ProcedureCost) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}
