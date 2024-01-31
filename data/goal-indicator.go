package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// GoalIndicator struct
type GoalIndicator struct {
	ID                       int       `db:"id,omitempty"`
	GoalID                   int       `db:"goal_id"`
	PerformanceIndicatorCode string    `db:"performance_indicator_code"`
	IndicatorSource          string    `db:"indicator_source"`
	BaseYear                 string    `db:"base_year"`
	GenderEquality           string    `db:"gender_equality"`
	BaseValue                string    `db:"base_value"`
	SourceOfInformation      string    `db:"source_of_information"`
	UnitOfMeasure            string    `db:"unit_of_measure"`
	IndicatorDescription     string    `db:"indicator_description"`
	PlannedValue1            string    `db:"planned_value_1"`
	RevisedValue1            string    `db:"revised_value_1"`
	AchievedValue1           string    `db:"achieved_value_1"`
	PlannedValue2            string    `db:"planned_value_2"`
	RevisedValue2            string    `db:"revised_value_2"`
	AchievedValue2           string    `db:"achieved_value_2"`
	PlannedValue3            string    `db:"planned_value_3"`
	RevisedValue3            string    `db:"revised_value_3"`
	AchievedValue3           string    `db:"achieved_value_3"`
	CreatedAt                time.Time `db:"created_at,omitempty"`
	UpdatedAt                time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *GoalIndicator) Table() string {
	return "goal_indicators"
}

// GetAll gets all records from the database, using upper
func (t *GoalIndicator) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*GoalIndicator, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*GoalIndicator
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
func (t *GoalIndicator) Get(id int) (*GoalIndicator, error) {
	var one GoalIndicator
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *GoalIndicator) Update(m GoalIndicator) error {
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
func (t *GoalIndicator) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *GoalIndicator) Insert(m GoalIndicator) (int, error) {
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
