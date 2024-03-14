package data

import (
	"time"

	"github.com/lib/pq"
	up "github.com/upper/db/v4"
)

type FineActType int

const (
	DissisionActType FineActType = 1
	VerdictActType   FineActType = 2
)

type FineStatus int

const (
	UnpaidFineStatus FineStatus = 1
	PaidFineStatus   FineStatus = 2
	PartFineStatus   FineStatus = 3
)

const (
	FineGracePeriod = 7
)

// Fine struct
type Fine struct {
	ID                     int           `db:"id,omitempty"`
	ActType                FineActType   `db:"act_type"`
	DecisionNumber         string        `db:"decision_number"`
	DecisionDate           time.Time     `db:"decision_date"`
	Subject                string        `db:"subject"`
	JMBG                   string        `db:"jmbg"`
	Residence              string        `db:"residence"`
	Amount                 float64       `db:"amount"`
	PaymentReferenceNumber string        `db:"payment_reference_number"`
	DebitReferenceNumber   string        `db:"debit_reference_number"`
	AccountID              int           `db:"account_id"`
	ExecutionDate          time.Time     `db:"execution_date"`
	PaymentDeadlineDate    time.Time     `db:"payment_deadline_date"`
	Description            string        `db:"description"`
	Status                 FineStatus    `db:"status"`
	CourtCosts             *float64      `db:"court_costs"`
	CourtAccountID         *int          `db:"court_account_id"`
	File                   pq.Int64Array `db:"file"`
	CreatedAt              time.Time     `db:"created_at,omitempty"`
	UpdatedAt              time.Time     `db:"updated_at"`
}

// Table returns the table name
func (t *Fine) Table() string {
	return "fines"
}

// Get gets one record from the database, by id, using upper
func (t *Fine) Get(id int) (*Fine, error) {
	var one Fine
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Insert inserts a model into the database, using upper
func (t *Fine) Insert(m Fine) (int, error) {
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

// GetAll gets all records from the database, using upper
func (t *Fine) GetAll(page *int, size *int, condition *up.AndExpr) ([]*Fine, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*Fine
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
func (t *Fine) Update(m Fine) error {
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
func (t *Fine) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}
