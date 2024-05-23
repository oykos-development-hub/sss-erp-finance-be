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

	err := res.OrderBy("-created_at").One(&one)
	if err != nil {
		return nil, err
	}

	return &one, err
}

// ValidateNewEntry validates the new entry against the old entry up to the end of the previous month.
func (t *SpendingDynamicEntry) ValidateNewEntry(oldEntry *SpendingDynamicEntry) bool {
	now := time.Now()
	currentMonth := now.Month()

	// If it's January, no previous month to validate
	if currentMonth == time.January {
		return true
	}

	// Get the month values for the current and old entries.
	newValues := t.monthValues()
	oldValues := oldEntry.monthValues()

	// Validate up to the end of the previous month
	for month := time.January; month < currentMonth; month++ {
		if !newValues[month].Equal(oldValues[month]) {
			return false
		}
	}

	return true
}

func (s *SpendingDynamicEntry) monthValues() map[time.Month]decimal.Decimal {
	return map[time.Month]decimal.Decimal{
		time.January:   s.January,
		time.February:  s.February,
		time.March:     s.March,
		time.April:     s.April,
		time.May:       s.May,
		time.June:      s.June,
		time.July:      s.July,
		time.August:    s.August,
		time.September: s.September,
		time.October:   s.October,
		time.November:  s.November,
		time.December:  s.December,
	}
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
