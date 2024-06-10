package data

import (
	goerrors "errors"
	"time"

	"github.com/shopspring/decimal"
	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type SpendingReleaseFilterDTO struct {
	CurrentBudgetID *int `json:"current_budget_id"`
	BudgetID        *int `json:"budget_id"`
	UnitID          *int `json:"unit_id"`
	AccountID       *int `json:"account_id"`
	Year            *int `json:"year"`
	Month           *int `json:"month"`
}

type SpendingReleaseOverview struct {
	Month     int             `db:"month"`
	Year      int             `db:"year"`
	CreatedAt time.Time       `db:"created_at"`
	Value     decimal.Decimal `db:"value"`
}

// SpendingRelease struct
type SpendingRelease struct {
	ID              int             `db:"id,omitempty"`
	CurrentBudgetID int             `db:"current_budget_id"`
	Year            int             `db:"year"`
	Month           int             `db:"month"`
	Value           decimal.Decimal `db:"value"`
	CreatedAt       time.Time       `db:"created_at,omitempty"`
	Username        string          `db:"username"`
}

type SpendingReleaseWithCurrentBudget struct {
	SpendingRelease
	BudgetID  int `db:"budget_id"`
	UnitID    int `db:"unit_id"`
	AccountID int `db:"account_id"`
	Actual    int `db:"actual"`
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
func (t *SpendingRelease) GetAll(filter SpendingReleaseFilterDTO) ([]SpendingReleaseWithCurrentBudget, error) {
	query := Upper.SQL().Select(
		"sr.id",
		"sr.current_budget_id",
		"sr.year",
		"sr.month",
		"sr.value",
		"sr.created_at",
		"sr.username",
		"cb.budget_id",
		"cb.unit_id",
		"cb.account_id",
	).
		From("spending_releases AS sr").
		Join("current_budgets AS cb").On("cb.id = sr.current_budget_id")

	var all []SpendingReleaseWithCurrentBudget

	if filter.CurrentBudgetID != nil {
		query = query.Where("sr.current_budget_id = ?", *filter.CurrentBudgetID)
	}
	if filter.BudgetID != nil {
		query = query.Where("cb.budget_id = ?", *filter.BudgetID)
	}
	if filter.UnitID != nil {
		query = query.Where("cb.unit_id = ?", *filter.UnitID)
	}
	if filter.Year != nil {
		query = query.Where("sr.year = ?", *filter.Year)
	}
	if filter.Month != nil {
		query = query.Where("sr.month = ?", *filter.Month)
	}

	err := query.OrderBy("-sr.created_at").All(&all)
	if err != nil {
		return nil, err
	}

	return all, err
}

// GetAll gets all records from the database, using upper
func (t *SpendingRelease) GetAllSum(month, year, budgetID, unitID int) ([]SpendingReleaseOverview, error) {
	var all []SpendingReleaseOverview

	query := Upper.SQL().Select(
		up.Raw("SUM(sr.value) AS value"),
		"sr.month",
		"sr.year",
		up.Raw("MAX(sr.created_at) AS created_at"),
	).
		From("spending_releases AS sr").
		Join("current_budgets AS cb").On("cb.id = sr.current_budget_id").
		Where("cb.budget_id", budgetID).And("cb.unit_id", unitID).
		GroupBy("cb.budget_id", "cb.unit_id", "sr.month", "sr.year")

	if month != 0 {
		query = query.Where("sr.month = ?", month)
	}

	if year != 0 {
		query = query.Where("sr.year = ?", year)
	}

	err := query.All(&all)
	if err != nil {
		return nil, errors.Wrap(err, "get release overview")
	}

	return all, err
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
func (t *SpendingRelease) Get(id int) (*SpendingReleaseWithCurrentBudget, error) {
	query := Upper.SQL().Select(
		"sr.id",
		"sr.current_budget_id",
		"sr.year",
		"sr.month",
		"sr.value",
		"sr.created_at",
		"sr.username",
		"cb.budget_id",
		"cb.unit_id",
		"cb.account_id",
	).
		From("spending_releases AS sr").
		Join("current_budgets AS cb").On("cb.id = sr.current_budget_id").
		Where("sr.id = ?", id)

	var one SpendingReleaseWithCurrentBudget

	err := query.OrderBy("-sr.created_at").One(&one)
	if err != nil {
		return nil, err
	}

	return &one, err
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
