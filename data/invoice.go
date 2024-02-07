package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// Invoice struct
type Invoice struct {
	ID                    int        `db:"id,omitempty"`
	InvoiceNumber         string     `db:"invoice_number"`
	Status                string     `db:"status,omitempty"`
	GrossPrice            float64    `db:"gross_price"`
	VATPrice              float64    `db:"vat_price"`
	SupplierID            int        `db:"supplier_id"`
	OrderID               int        `db:"order_id"`
	OrganizationUnitID    int        `db:"organization_unit_id"`
	DateOfInvoice         time.Time  `db:"date_of_invoice"`
	ReceiptDate           time.Time  `db:"receipt_date"`
	DateOfPayment         time.Time  `db:"date_of_payment"`
	SSSInvoiceReceiptDate *time.Time `db:"sss_invoice_receipt_date"`
	FileID                int        `db:"file_id"`
	BankAccount           string     `db:"bank_account"`
	Description           string     `db:"description"`
	CreatedAt             time.Time  `db:"created_at,omitempty"`
	UpdatedAt             time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *Invoice) Table() string {
	return "invoices"
}

// GetAll gets all records from the database, using upper
func (t *Invoice) GetAll(page *int, size *int, condition *up.AndExpr) ([]*Invoice, *uint64, error) {
	collection := upper.Collection(t.Table())
	var all []*Invoice
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

// Get gets one record from the database, by id, using upper
func (t *Invoice) Get(id int) (*Invoice, error) {
	var one Invoice
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Invoice) Update(m Invoice) error {
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
func (t *Invoice) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *Invoice) Insert(m Invoice) (int, error) {
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

// Builder is an example of using upper's sql builder
func (t *Invoice) Builder(id int) ([]*Invoice, error) {
	collection := upper.Collection(t.Table())

	var result []*Invoice

	err := collection.Session().
		SQL().
		SelectFrom(t.Table()).
		Where("id > ?", id).
		OrderBy("id").
		All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}