package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type AccountingOrderItemsTitle string

var (
	MainBillTitle AccountingOrderItemsTitle = "Korektivni račun"
	SupplierTitle AccountingOrderItemsTitle = "Dobavljač"
	TaxTitle      AccountingOrderItemsTitle = "Porez"
	SubTaxTitle   AccountingOrderItemsTitle = "Prirez"
)

type AccountingEntry struct {
	ID                 int       `db:"id,omitempty"`
	Title              string    `db:"title"`
	IDOfEntry          int       `db:"id_of_entry"`
	OrganizationUnitID int       `db:"organization_unit_id"`
	DateOfBooking      time.Time `db:"date_of_booking"`
	CreatedAt          time.Time `db:"created_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at"`
}

type ObligationForAccounting struct {
	InvoiceID  *int              `json:"invoice_id"`
	SalaryID   *int              `json:"salary_id"`
	SupplierID *int              `json:"supplier_id"`
	Date       time.Time         `json:"date"`
	Type       TypesOfObligation `json:"type"`
	Title      string            `json:"title"`
	Price      float64           `json:"price"`
	Status     string            `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
}

// Table returns the table name
func (t *AccountingEntry) Table() string {
	return "accounting_entries"
}

// GetAll gets all records from the database, using upper
func (t *AccountingEntry) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*AccountingEntry, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*AccountingEntry
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
func (t *AccountingEntry) Get(id int) (*AccountingEntry, error) {
	var one AccountingEntry
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *AccountingEntry) Update(tx up.Session, m AccountingEntry) error {
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *AccountingEntry) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *AccountingEntry) Insert(tx up.Session, m AccountingEntry) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}

func (t *AccountingEntry) GetObligationsForAccounting(filter ObligationsFilter) ([]ObligationForAccounting, *uint64, error) {
	var items []ObligationForAccounting

	queryForInvoices := `select i.id, sum((a.net_price +a.net_price*a.vat_percentage/100)*a.amount) as sum, i.invoice_number, i.supplier_id, i.date_of_invoice
						from invoices i
						left join articles a on a.invoice_id = i.id
						where 
						i.organization_unit_id = $1 and i.type = $2 and i.registred = false 
						group by i.id;`

	queryForAdditionalExpenses := `select i.id, sum(a.price), i.type, i.invoice_number, i.supplier_id, i.date_of_invoice
	                               from additional_expenses a
	                               left join invoices i on a.invoice_id = i.id
	                               where a.invoice_id = i.id and
	                               i.organization_unit_id = $1 and i.registred = false
	                               group by i.id, i.invoice_number order by i.id;`

	queryForSalaryAdditionalExpenses := `select s.id, sum(a.amount), s.month, s.date_of_calculation
	                                     from salary_additional_expenses a
	                                     left join salaries s on s.id = a.salary_id
	                                     where s.organization_unit_id = $1 and s.registred = false
	                                     group by s.id, s.month order by s.id;`

	rows, err := Upper.SQL().Query(queryForInvoices, filter.OrganizationUnitID, TypeInvoice)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var obligation ObligationForAccounting
		err = rows.Scan(&obligation.InvoiceID, &obligation.Price, &obligation.Title, &obligation.SupplierID, &obligation.Date)

		if err != nil {
			return nil, nil, err
		}

		obligation.Type = TypeInvoice
		obligation.Title = "Račun broj " + obligation.Title
		items = append(items, obligation)
	}

	rows, err = Upper.SQL().Query(queryForAdditionalExpenses, filter.OrganizationUnitID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var obligation ObligationForAccounting
		err = rows.Scan(&obligation.InvoiceID, &obligation.Price, &obligation.Type, &obligation.Title, &obligation.SupplierID, &obligation.Date)

		if err != nil {
			return nil, nil, err
		}

		if obligation.Type == TypeDecision {
			obligation.Title = "Rješenje broj " + obligation.Title
		} else {
			obligation.Title = "Ugovor broj " + obligation.Title
		}

		if filter.Type == nil || *filter.Type == obligation.Type {
			items = append(items, obligation)
		}
	}

	rows, err = Upper.SQL().Query(queryForSalaryAdditionalExpenses, filter.OrganizationUnitID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var obligation ObligationForAccounting
		err = rows.Scan(&obligation.SalaryID, &obligation.Price, &obligation.Title, &obligation.Date)

		if err != nil {
			return nil, nil, err
		}

		obligation.Title = "Zarada " + obligation.Title
		obligation.Type = TypeSalary

		items = append(items, obligation)

	}

	total := uint64(len(items))

	return items, &total, nil
}
