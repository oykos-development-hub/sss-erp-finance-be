package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type PaymentOrder struct {
	ID                 int        `db:"id,omitempty"`
	OrganizationUnitID int        `db:"organization_unit_id"`
	SupplierID         int        `db:"supplier_id"`
	BankAccount        string     `db:"bank_account"`
	DateOfPayment      time.Time  `db:"date_of_payment"`
	DateOfOrder        *time.Time `db:"date_of_order"`
	IDOfStatement      *int       `db:"id_of_statement,omitempty"`
	SAPID              *string    `db:"sap_id"`
	Registred          *bool      `db:"registred,omitempty"`
	DateOfSAP          *time.Time `db:"date_of_sap"`
	SourceOfFunding    string     `db:"source_of_funding"`
	FileID             *int       `db:"file_id"`
	Amount             float64    `db:"amount"`
	Description        string     `db:"description"`
	Status             *string    `db:"status,omitempty"`
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
	DateOfStart        *time.Time         `json:"date_of_start"`
	DateOfEnd          *time.Time         `json:"date_of_end"`
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
	AccountID                 int               `json:"account_id"`
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
		return nil, nil, newErrors.Wrap(err, "upper count")
	}

	if page != nil && size != nil {
		res = paginateResult(res, *page, *size)
	}

	err = res.OrderBy(orders...).All(&all)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "upper all")
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
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

func (t *PaymentOrder) GetByIdOfStatement(id int) (*PaymentOrder, error) {
	var one PaymentOrder
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id_of_statement": id})
	err := res.One(&one)
	if err != nil {
		return nil, newErrors.Wrap(err, "upper one")
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *PaymentOrder) Update(ctx context.Context, tx up.Session, m PaymentOrder) error {
	m.UpdatedAt = time.Now()
	collection := tx.Collection(t.Table())

	order, err := t.Get(m.ID)

	if err != nil {
		return err
	}

	m.IDOfStatement = order.IDOfStatement

	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return newErrors.Wrap(err, "upper exec")
	}

	res := collection.Find(m.ID)
	if err := res.Update(&m); err != nil {
		return newErrors.Wrap(err, "upper update")
	}

	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *PaymentOrder) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return newErrors.Wrap(err, "upper exec")
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(id)
		if err := res.Delete(); err != nil {
			return newErrors.Wrap(err, "upper delete")
		}

		return nil
	})

	if err != nil {
		return newErrors.Wrap(err, "upper tx")
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *PaymentOrder) Insert(ctx context.Context, tx up.Session, m PaymentOrder) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return 0, newErrors.Wrap(err, "contextuitl get user id from context")
	}

	var id int

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return 0, newErrors.Wrap(err, "upper exec")
	}

	collection := tx.Collection(t.Table())

	var res up.InsertResult
	var err error

	status := "Kreiran"

	m.Status = &status

	if res, err = collection.Insert(m); err != nil {
		return 0, newErrors.Wrap(err, "upper insert")
	}

	id = getInsertId(res.ID())

	return id, nil
}

func (t *PaymentOrder) GetAllObligations(filter ObligationsFilter) ([]Obligation, *uint64, error) {
	var items []Obligation

	queryForInvoices := `select i.id, COALESCE(SUM((a.net_price + a.net_price * a.vat_percentage / 100) * a.amount), 0) as sum, 
						i.invoice_number, i.pro_forma_invoice_number, i.status, i.created_at
						from invoices i
						left join articles a on a.invoice_id = i.id
						where i.supplier_id = $1 and
						i.organization_unit_id = $2 and i.type = $4 and i.status <> $3 and i.status <> $5
						group by i.id;`

	queryForPaidInvoices := `select COALESCE(sum(p.amount),0) as sum from payment_order_items pi 
							left join payment_orders p on p.id = pi.payment_order_id
							where pi.invoice_id = $1 and (p.status is null or p.status <> 'Storniran');`

	queryForAdditionalExpenses := `select a.id, a.price, a.title, i.type, i.invoice_number, a.status, a.created_at, a.account_id
	                               from additional_expenses a
	                               left join invoices i on a.invoice_id = i.id
	                               where a.invoice_id = i.id and a.subject_id = $1 and
	                               i.organization_unit_id = $2 and a.status <> $3
	                               group by a.id, a.title, i.type, i.invoice_number order by a.id;`

	queryForPaidAdditionalExpenses := `select COALESCE(sum(p.amount),0) as sum from payment_order_items pi 
								   left join payment_orders p on p.id = pi.payment_order_id
								   where pi.additional_expense_id = $1  and (p.status is null or p.status <> 'Storniran');`

	queryForSalaryAdditionalExpenses := `select a.id, a.amount, a.title, a.status, a.created_at, s.month, a.account_id
	                                     from salary_additional_expenses a
	                                     left join salaries s on s.id = a.salary_id
	                                     where  a.subject_id = $1 and
	                                     s.organization_unit_id = $2 and a.status <> $3
	                                     group by a.id, a.title, s.month order by a.id;`

	queryForPaidSalaryAdditionalExpenses := `select COALESCE(sum(p.amount),0) as sum from payment_order_items pi 
										 left join payment_orders p on p.id = pi.payment_order_id
										 where pi.salary_additional_expense_id = $1  and (p.status is null or p.status <> 'Storniran');`

	if filter.Type == nil || *filter.Type == TypeInvoice {
		rows, err := Upper.SQL().Query(queryForInvoices, filter.SupplierID, filter.OrganizationUnitID, InvoiceStatusFull, TypeInvoice, InvoiceStatusIncomplete)
		if err != nil {
			return nil, nil, newErrors.Wrap(err, "upper exec")
		}
		defer rows.Close()

		for rows.Next() {
			var obligation Obligation
			var paid *float64
			var invoiceNumber *string
			err = rows.Scan(&obligation.InvoiceID, &obligation.TotalPrice, &obligation.Title, &invoiceNumber, &obligation.Status, &obligation.CreatedAt)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "upper scan")
			}

			rows1, err := Upper.SQL().Query(queryForPaidInvoices, obligation.InvoiceID)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "upper exec")
			}

			for rows1.Next() {
				err = rows1.Scan(&paid)

				if err != nil {
					return nil, nil, newErrors.Wrap(err, "upper scan")
				}

				if paid != nil {
					obligation.RemainPrice = obligation.TotalPrice - *paid
				} else {
					obligation.RemainPrice = obligation.TotalPrice
				}
			}
			obligation.Type = TypeInvoice

			if invoiceNumber != nil && *invoiceNumber != "" && obligation.Title == "" {
				obligation.Title = *invoiceNumber
			}

			items = append(items, obligation)
		}
	}

	if filter.Type == nil || (*filter.Type == TypeDecision || *filter.Type == TypeContract) {
		rows, err := Upper.SQL().Query(queryForAdditionalExpenses, filter.SupplierID, filter.OrganizationUnitID, InvoiceStatusFull)
		if err != nil {
			return nil, nil, newErrors.Wrap(err, "upper exec")
		}
		defer rows.Close()

		for rows.Next() {
			var obligation Obligation
			var paid *float64
			var title string
			err = rows.Scan(&obligation.AdditionalExpenseID, &obligation.TotalPrice, &obligation.Title, &obligation.Type, &title, &obligation.Status, &obligation.CreatedAt, &obligation.AccountID)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "upper scan")
			}

			rows1, err := Upper.SQL().Query(queryForPaidAdditionalExpenses, obligation.AdditionalExpenseID)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "upper exec")
			}

			for rows1.Next() {
				err = rows1.Scan(&paid)

				if err != nil {
					return nil, nil, newErrors.Wrap(err, "upper scan")
				}

				if paid != nil {
					obligation.RemainPrice = obligation.TotalPrice - *paid
				} else {
					obligation.RemainPrice = obligation.TotalPrice
				}
			}

			if filter.Type == nil || *filter.Type == obligation.Type {
				items = append(items, obligation)
			}
		}
	}

	if filter.Type == nil || *filter.Type == TypeSalary {
		rows, err := Upper.SQL().Query(queryForSalaryAdditionalExpenses, filter.SupplierID, filter.OrganizationUnitID, InvoiceStatusFull)
		if err != nil {
			return nil, nil, newErrors.Wrap(err, "upper exec")
		}
		defer rows.Close()

		for rows.Next() {
			var obligation Obligation
			var paid *float64
			var title string
			err = rows.Scan(&obligation.SalaryAdditionalExpenseID, &obligation.TotalPrice, &title, &obligation.Status, &obligation.CreatedAt, &obligation.Title, &obligation.AccountID)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "upper scan")
			}

			rows1, err := Upper.SQL().Query(queryForPaidSalaryAdditionalExpenses, obligation.AdditionalExpenseID)

			if err != nil {
				return nil, nil, newErrors.Wrap(err, "upper exec")
			}

			for rows1.Next() {
				err = rows1.Scan(&paid)

				if err != nil {
					return nil, nil, newErrors.Wrap(err, "upper scan")
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

func (t *PaymentOrder) PayPaymentOrder(ctx context.Context, tx up.Session, id int, SAPID string, DateOfSAP time.Time) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return newErrors.Wrap(err, "upper exec")
	}

	query = `update payment_orders set sap_id = $1, date_of_sap = $2, status = 'Plaćen' where id = $3`

	_, err := tx.SQL().Query(query, SAPID, DateOfSAP, id)

	if err != nil {
		return newErrors.Wrap(err, "upper exec")
	}

	return nil
}

func (t *PaymentOrder) CancelPaymentOrder(ctx context.Context, tx up.Session, id int) error {

	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		err := newErrors.New("user ID not found in context")
		return newErrors.Wrap(err, "contextuitl get user id from context")
	}

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return newErrors.Wrap(err, "upper exec")
	}

	query = `update payment_orders set status = 'Storniran' where id = $1`

	_, err := tx.SQL().Query(query, id)

	if err != nil {
		return newErrors.Wrap(err, "upper exec")
	}

	return nil

}
