package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// TaxAuthorityCodebook struct
type TaxAuthorityCodebook struct {
	ID                                   int       `db:"id,omitempty"`
	Title                                string    `db:"title"`
	Active                               bool      `db:"active"`
	Code                                 string    `db:"code"`
	TaxPercentage                        float64   `db:"tax_percentage"`
	TaxSupplierID                        int       `db:"tax_supplier_id"`
	ReleasePercentage                    float64   `db:"release_percentage"`
	PioPercentage                        float64   `db:"pio_percentage"`
	PioSupplierID                        int       `db:"pio_supplier_id"`
	PioPercentageEmployerPercentage      float64   `db:"pio_percentage_employer_percentage"`
	PioEmployerSupplierID                int       `db:"pio_employer_supplier_id"`
	PioPercentageEmployeePercentage      float64   `db:"pio_percentage_employee_percentage"`
	PioEmployeeSupplierID                int       `db:"pio_employee_supplier_id"`
	UnemploymentPercentage               float64   `db:"unemployment_percentage"`
	UnemploymentSupplierID               int       `db:"unemployment_supplier_id"`
	UnemploymentEmployerPercentage       float64   `db:"unemployment_employer_percentage"`
	UnemploymentEmployerSupplierID       int       `db:"unemployment_employer_supplier_id"`
	UnemploymentEmployeePercentage       float64   `db:"unemployment_employee_percentage"`
	UnemploymentEmployeeSupplierID       int       `db:"unemployment_employee_supplier_id"`
	LaborFund                            float64   `db:"labor_fund"`
	LaborFundSupplierID                  int       `db:"labor_fund_supplier_id"`
	PreviousIncomePercentageLessThan700  float64   `db:"previous_income_percentage_less_than_700"`
	PreviousIncomePercentageLessThan1000 float64   `db:"previous_income_percentage_less_than_1000"`
	PreviousIncomePercentageMoreThan1000 float64   `db:"previous_income_percentage_more_than_1000"`
	Coefficient                          float64   `db:"coefficient"`
	CreatedAt                            time.Time `db:"created_at,omitempty"`
	UpdatedAt                            time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *TaxAuthorityCodebook) Table() string {
	return "tax_authority_codebooks"
}

// GetAll gets all records from the database, using upper
func (t *TaxAuthorityCodebook) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*TaxAuthorityCodebook, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*TaxAuthorityCodebook
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
func (t *TaxAuthorityCodebook) Get(id int) (*TaxAuthorityCodebook, error) {
	var one TaxAuthorityCodebook
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *TaxAuthorityCodebook) Update(m TaxAuthorityCodebook) error {
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
func (t *TaxAuthorityCodebook) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *TaxAuthorityCodebook) Insert(m TaxAuthorityCodebook) (int, error) {
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
