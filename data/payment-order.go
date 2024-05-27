package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type PaymentOrder struct {
	ID                 int        `db:"id,omitempty"`
	OrganizationUnitID int        `db:"organization_unit_id"`
	SupplierID         int        `db:"supplier_id"`
	BankAccount        string     `db:"bank_account"`
	DateOfPayment      time.Time  `db:"date_of_payment"`
	DateOfOrder        *time.Time `db:"date_of_order"`
	IDOfStatement      *string    `db:"id_of_statement"`
	SAPID              *string    `db:"sap_id"`
	Registred          *bool      `db:"registred,omitempty"`
	DateOfSAP          *time.Time `db:"date_of_sap"`
	SourceOfFunding    string     `db:"source_of_funding"`
	FileID             *int       `db:"file_id"`
	Amount             float64    `db:"amount"`
	Description        string     `db:"description"`
	CreatedAt          time.Time  `db:"created_at,omitempty"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

type ObligationsFilter struct {
	Page               *int               `json:"page"`
	Size               *int               `json:"size"`
	OrganizationUnitID int                `json:"organization_unit_id"`
	SupplierID         int                `json:"supplier_id"`
	Type               *TypesOfObligation `json:"type"`
	Search             *string            `json:"search"`
}

type Obligation struct {
	InvoiceID                 *int              `json:"invoice_id"`
	AdditionalExpenseID       *int              `json:"additional_expense_id"`
	SalaryAdditionalExpenseID *int              `json:"salary_additional_expense_id"`
	Type                      TypesOfObligation `json:"type"`
	Title                     string            `json:"title"`
	TotalPrice                float64           `json:"total_price"`
	RemainPrice               float64           `json:"remain_price"`
	Status                    string            `json:"status"`
	CreatedAt                 time.Time         `json:"created_at"`
}

// Table returns the table name
func (t *PaymentOrder) Table() string {
	return "payment_orders"
}

// GetAll gets all records from the database, using upper
func (t *PaymentOrder) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*PaymentOrder, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*PaymentOrder
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
func (t *PaymentOrder) Get(id int) (*PaymentOrder, error) {
	var one PaymentOrder
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *PaymentOrder) Update(tx up.Session, m PaymentOrder) error {
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
func (t *PaymentOrder) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *PaymentOrder) Insert(tx up.Session, m PaymentOrder) (int, error) {
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

func (t *PaymentOrder) GetAllObligations(filter ObligationsFilter) ([]Obligation, *uint64, error) {
	var items []Obligation

	queryForInvoices := `select i.id, sum((a.net_price +a.net_price*a.vat_percentage/100)*a.amount) as sum, 
						i.invoice_number, i.status, i.created_at
						from invoices i
						left join articles a on a.invoice_id = i.id
						where i.supplier_id = $1 and
						i.organization_unit_id = $2 and i.type = $4 and i.status <> $3 
						group by i.id;`

	queryForPaidInvoices := `select sum(p.amount) from payment_order_items pi 
							left join payment_orders p on p.id = pi.payment_order_id
							where pi.invoice_id = $1`

	queryForAdditionalExpenses := `select a.id, a.price, a.title, i.type, i.invoice_number, a.status, a.created_at
	                               from additional_expenses a
	                               left join invoices i on a.invoice_id = i.id
	                               where a.invoice_id = i.id and a.subject_id = $1 and
	                               i.organization_unit_id = $2 and a.status <> $3
	                               group by a.id, a.title, i.type, i.invoice_number order by a.id;`

	queryForPaidAdditionalExpenses := `select sum(p.amount) from payment_order_items pi 
								   left join payment_orders p on p.id = pi.payment_order_id
								   where pi.additional_expense_id = $1`

	queryForSalaryAdditionalExpenses := `select a.id, a.amount, a.title, a.status, a.created_at, s.month
	                                     from salary_additional_expenses a
	                                     left join salaries s on s.id = a.salary_id
	                                     where  a.subject_id = $1 and
	                                     s.organization_unit_id = $2 and a.status <> $3
	                                     group by a.id, a.title, s.month order by a.id;`

	queryForPaidSalaryAdditionalExpenses := `select sum(p.amount) from payment_order_items pi 
										 left join payment_orders p on p.id = pi.payment_order_id
										 where pi.salary_additional_expense_id = $1`

	if filter.Type == nil || *filter.Type == TypeInvoice {
		rows, err := Upper.SQL().Query(queryForInvoices, filter.SupplierID, filter.OrganizationUnitID, InvoiceStatusFull, TypeInvoice)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var obligation Obligation
			var paid *float64
			err = rows.Scan(&obligation.InvoiceID, &obligation.TotalPrice, &obligation.Title, &obligation.Status, &obligation.CreatedAt)

			if err != nil {
				return nil, nil, err
			}

			rows1, err := Upper.SQL().Query(queryForPaidInvoices, obligation.InvoiceID)

			if err != nil {
				return nil, nil, err
			}

			for rows1.Next() {
				err = rows1.Scan(&paid)

				if err != nil {
					return nil, nil, err
				}

				if paid != nil {
					obligation.RemainPrice = obligation.TotalPrice - *paid
				} else {
					obligation.RemainPrice = obligation.TotalPrice
				}
			}
			obligation.Type = TypeInvoice
			obligation.Title = "Račun broj " + obligation.Title + " Neto"
			items = append(items, obligation)
		}
	}

	if filter.Type == nil || (*filter.Type == TypeDecision || *filter.Type == TypeContract) {
		rows, err := Upper.SQL().Query(queryForAdditionalExpenses, filter.SupplierID, filter.OrganizationUnitID, InvoiceStatusFull)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var obligation Obligation
			var paid *float64
			var title string
			err = rows.Scan(&obligation.AdditionalExpenseID, &obligation.TotalPrice, &obligation.Title, &obligation.Type, &title, &obligation.Status, &obligation.CreatedAt)

			if err != nil {
				return nil, nil, err
			}

			rows1, err := Upper.SQL().Query(queryForPaidAdditionalExpenses, obligation.AdditionalExpenseID)

			if err != nil {
				return nil, nil, err
			}

			for rows1.Next() {
				err = rows1.Scan(&paid)

				if err != nil {
					return nil, nil, err
				}

				if paid != nil {
					obligation.RemainPrice = obligation.TotalPrice - *paid
				} else {
					obligation.RemainPrice = obligation.TotalPrice
				}
			}

			if obligation.Type == TypeDecision {
				obligation.Title = "Rješenje broj " + title + " " + obligation.Title
			} else {
				obligation.Title = "Ugovor broj " + title + " " + obligation.Title
			}

			if filter.Type == nil || *filter.Type == obligation.Type {
				items = append(items, obligation)
			}
		}
	}

	if filter.Type == nil || *filter.Type == TypeSalary {
		rows, err := Upper.SQL().Query(queryForSalaryAdditionalExpenses, filter.SupplierID, filter.OrganizationUnitID, InvoiceStatusFull)
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var obligation Obligation
			var paid *float64
			var title string
			err = rows.Scan(&obligation.SalaryAdditionalExpenseID, &obligation.TotalPrice, &title, &obligation.Status, &obligation.CreatedAt, &obligation.Title)

			if err != nil {
				return nil, nil, err
			}

			rows1, err := Upper.SQL().Query(queryForPaidSalaryAdditionalExpenses, obligation.AdditionalExpenseID)

			if err != nil {
				return nil, nil, err
			}

			for rows1.Next() {
				err = rows1.Scan(&paid)

				if err != nil {
					return nil, nil, err
				}

				if paid != nil {
					obligation.RemainPrice = obligation.TotalPrice - *paid
				} else {
					obligation.RemainPrice = obligation.TotalPrice
				}

				obligation.Title = "Zarada " + obligation.Title + " " + title
				obligation.Type = TypeSalary
			}

			items = append(items, obligation)

		}
	}

	total := uint64(len(items))

	return items, &total, nil
}

func (t *PaymentOrder) PayPaymentOrder(tx up.Session, id int, SAPID string, DateOfSAP time.Time) error {
	query := `update payment_orders set sap_id = $1, date_of_sap = $2 where id = $3`

	rows, err := tx.SQL().Query(query, SAPID, DateOfSAP, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
