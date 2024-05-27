package data

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

// SpendingDynamic struct
type SpendingDynamicEntry struct {
	ID              int             `db:"id,omitempty"`
	CurrentBudgetID int             `db:"current_budget_id"`
	Username        string          `db:"username"`
	January         decimal.Decimal `db:"january"`
	February        decimal.Decimal `db:"february"`
	March           decimal.Decimal `db:"march"`
	April           decimal.Decimal `db:"april"`
	May             decimal.Decimal `db:"may"`
	June            decimal.Decimal `db:"june"`
	July            decimal.Decimal `db:"july"`
	August          decimal.Decimal `db:"august"`
	September       decimal.Decimal `db:"september"`
	October         decimal.Decimal `db:"october"`
	November        decimal.Decimal `db:"november"`
	December        decimal.Decimal `db:"december"`
	Version         int             `db:"version"`
	CreatedAt       time.Time       `db:"created_at,omitempty"`
}

type SpendingDynamicEntryWithCurrentBudget struct {
	SpendingDynamicEntry
	BudgetID  int             `db:"budget_id"`
	UnitID    int             `db:"unit_id"`
	AccountID int             `db:"account_id"`
	Actual    decimal.Decimal `db:"actual"`
}

// Table returns the table name
func (t *SpendingDynamicEntry) Table() string {
	return "spending_dynamic_entries"
}

func (t *SpendingDynamicEntry) SumOfMonths() decimal.Decimal {
	return decimal.Sum(t.January, t.February, t.March, t.April, t.May, t.June, t.July, t.August, t.September, t.October, t.November, t.December)
}

func (t *SpendingDynamicEntry) FindLatestVersion() (int, error) {
	var version int

	row, err := Upper.SQL().QueryRow("SELECT MAX(version) AS version FROM spending_dynamic_entries")
	if err != nil {
		return 0, errors.Wrap(err, "FindLatestVersion")
	}

	err = row.Scan(&version)
	if err != nil {
		return 0, errors.Wrap(err, "FindLatestVersion")
	}

	return version, nil

}

// GetAll gets all records from the database, using upper
func (t *SpendingDynamicEntry) FindAll(currentBudgetID, version, budgetID, unitID *int) ([]SpendingDynamicEntryWithCurrentBudget, error) {
	query := Upper.SQL().Select(
		"sd.id",
		"sd.current_budget_id",
		"sd.username",
		"sd.january",
		"sd.february",
		"sd.march",
		"sd.april",
		"sd.may",
		"sd.june",
		"sd.july",
		"sd.august",
		"sd.september",
		"sd.october",
		"sd.november",
		"sd.december",
		"sd.version",
		"sd.created_at",
		"cb.budget_id",
		"cb.unit_id",
		"cb.account_id",
		"cb.actual",
		"cb.created_at",
	).
		From("spending_dynamic_entries AS sd").
		Join("current_budgets AS cb").On("cb.id = sd.current_budget_id")

	var all []SpendingDynamicEntryWithCurrentBudget

	if currentBudgetID != nil {
		query = query.Where("sd.current_budget_id = ?", *currentBudgetID)
	}
	if budgetID != nil {
		query = query.Where("cb.budget_id = ?", *budgetID)
	}
	if unitID != nil {
		query = query.Where("cb.unit_id = ?", *unitID)
	}
	if version != nil {
		query = query.Where("sd.version = ?", *version)
	}

	fmt.Println(query)

	err := query.OrderBy("-sd.created_at").All(&all)
	if err != nil {
		return nil, err
	}

	return all, err
}

type SpendingDynamicHistory struct {
	CreatedAt time.Time
	Version   int
	Username  string
}

func (t *SpendingDynamicEntry) FindHistoryChanges(budgetID, unitID int) ([]SpendingDynamicHistory, error) {
	query := `SELECT MIN(sd.created_at), MIN(username), version FROM spending_dynamic_entries sd 
	JOIN current_budgets AS cb ON cb.id = sd.current_budget_id
	WHERE cb.budget_id = $1 AND cb.unit_id = $2
	group by version`

	rows, err := Upper.SQL().Query(query, budgetID, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]SpendingDynamicHistory, 0)
	for rows.Next() {
		var history SpendingDynamicHistory
		err = rows.Scan(&history.CreatedAt, &history.Username, &history.Version)
		if err != nil {
			return nil, err
		}

		items = append(items, history)
	}

	return items, nil
}

// ValidateNewEntry validates the new entry against the old entry up to the end of the previous month.
func (t *SpendingDynamicEntry) ValidateNewEntry(oldEntry *SpendingDynamicEntryWithCurrentBudget) bool {
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
