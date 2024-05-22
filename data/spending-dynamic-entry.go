package data

import (
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
)

// SpendingDynamic struct
type SpendingDynamicEntry struct {
	ID                int             `db:"id,omitempty"`
	SpendingDynamicID int             `db:"spending_dynamic_id"`
	January           decimal.Decimal `db:"january"`
	February          decimal.Decimal `db:"february"`
	March             decimal.Decimal `db:"march"`
	April             decimal.Decimal `db:"april"`
	May               decimal.Decimal `db:"may"`
	June              decimal.Decimal `db:"june"`
	July              decimal.Decimal `db:"july"`
	August            decimal.Decimal `db:"august"`
	September         decimal.Decimal `db:"september"`
	October           decimal.Decimal `db:"october"`
	November          decimal.Decimal `db:"november"`
	December          decimal.Decimal `db:"december"`
	CreatedAt         time.Time       `db:"created_at,omitempty"`
}

// Table returns the table name
func (t *SpendingDynamicEntry) Table() string {
	return "spending_dynamic_entries"
}

func (t *SpendingDynamicEntry) SumOfMonths() decimal.Decimal {
	return decimal.Sum(t.January, t.February, t.March, t.April, t.May, t.June, t.July, t.August, t.September, t.October, t.November, t.December)
}

// GetAll gets all records from the database, using upper
func (t *SpendingDynamicEntry) FindAll(condition *up.Cond) ([]SpendingDynamicEntry, error) {
	collection := Upper.Collection(t.Table())
	var all []SpendingDynamicEntry
	var res up.Result

	if condition != nil {
		res = collection.Find(*condition)
	} else {
		res = collection.Find()
	}

	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, err
}

func (t *SpendingDynamicEntry) FindBy(condition *up.Cond) (*SpendingDynamicEntry, error) {
	collection := Upper.Collection(t.Table())
	var one SpendingDynamicEntry
	var res up.Result

	if condition != nil {
		res = collection.Find(*condition)
	} else {
		res = collection.Find()
	}

	err := res.Limit(1).OrderBy("-created_at").One(&one)
	if err != nil {
		return nil, err
	}

	return &one, err
}

// Get gets one record from the database, by id, using upper
func (t *SpendingDynamicEntry) Get(id int) (*SpendingDynamicEntry, error) {
	var one SpendingDynamicEntry

	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Insert inserts a model into the database, using upper
func (t *SpendingDynamicEntry) Insert(m SpendingDynamicEntry) (int, error) {
	m.CreatedAt = time.Now()

	collection := Upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
