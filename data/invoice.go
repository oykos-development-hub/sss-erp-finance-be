package data

import (
	"context"
	"fmt"
	"time"

	up "github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/contextutil"
	"gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type InvoiceStatus string

var (
	InvoiceStatusIncomplete InvoiceStatus = "Nepotpun"
	InvoiceStatusCreated    InvoiceStatus = "Kreiran"
	InvoiceStatusPart       InvoiceStatus = "Djelimično na nalogu"
	InvoiceStatusFull       InvoiceStatus = "Na nalogu"
)

// Invoice struct
type Invoice struct {
	ID                            int               `db:"id,omitempty"`
	InvoiceNumber                 string            `db:"invoice_number"`
	PassedToInventory             bool              `db:"passed_to_inventory"`
	PassedToAccounting            bool              `db:"passed_to_accounting"`
	IsInvoice                     bool              `db:"is_invoice"`
	ProFormaInvoiceNumber         string            `db:"pro_forma_invoice_number"`
	ProFormaInvoiceDate           *time.Time        `db:"pro_forma_invoice_date"`
	Type                          TypesOfObligation `db:"type"`
	Issuer                        string            `db:"issuer"`
	TaxAuthorityCodebookID        int               `db:"tax_authority_codebook_id"`
	TypeOfSubject                 int               `db:"type_of_subject"`
	TypeOfContract                int               `db:"type_of_contract"`
	SourceOfFunding               string            `db:"source_of_funding"`
	Supplier                      string            `db:"supplier"`
	Status                        InvoiceStatus     `db:"status,omitempty"`
	Registred                     *bool             `db:"registred,omitempty"`
	GrossPrice                    float64           `db:"gross_price"`
	VATPrice                      float64           `db:"vat_price"`
	SupplierID                    int               `db:"supplier_id"`
	MunicipalityID                int               `db:"municipality_id"`
	TypeOfDecision                int               `db:"type_of_decision"`
	ActivityID                    int               `db:"activity_id"`
	OrderID                       *int              `db:"order_id,omitempty"`
	OrganizationUnitID            int               `db:"organization_unit_id"`
	DateOfInvoice                 time.Time         `db:"date_of_invoice"`
	ReceiptDate                   time.Time         `db:"receipt_date"`
	DateOfPayment                 time.Time         `db:"date_of_payment"`
	SSSInvoiceReceiptDate         *time.Time        `db:"sss_invoice_receipt_date"`
	SSSProFormaInvoiceReceiptDate *time.Time        `db:"sss_pro_forma_invoice_receipt_date"`
	DateOfStart                   time.Time         `db:"date_of_start"`
	DateOfEnd                     time.Time         `db:"date_of_end"`
	ProFormaInvoiceFileID         int               `db:"pro_forma_invoice_file_id"`
	FileID                        int               `db:"file_id"`
	BankAccount                   string            `db:"bank_account"`
	Description                   string            `db:"description"`
	CreatedAt                     time.Time         `db:"created_at,omitempty"`
	UpdatedAt                     time.Time         `db:"updated_at"`
}

// Table returns the table name
func (t *Invoice) Table() string {
	return "invoices"
}

// GetAll gets all records from the database, using upper
func (t *Invoice) GetAll(page *int, size *int, condition *up.AndExpr) ([]*Invoice, *uint64, error) {
	collection := Upper.Collection(t.Table())
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
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *Invoice) Update(ctx context.Context, tx up.Session, m Invoice) error {
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return err
	}

	collection := tx.Collection(t.Table())
	res := collection.Find(m.ID)
	if err := res.Update(&m); err != nil {
		return err
	}

	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *Invoice) Delete(ctx context.Context, id int) error {
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return errors.New("user ID not found in context")
	}

	err := Upper.Tx(func(sess up.Session) error {
		query := fmt.Sprintf("SET myapp.user_id = %d", userID)
		if _, err := sess.SQL().Exec(query); err != nil {
			return err
		}

		collection := sess.Collection(t.Table())
		res := collection.Find(id)
		if err := res.Delete(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *Invoice) Insert(ctx context.Context, tx up.Session, m Invoice) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	userID, ok := contextutil.GetUserIDFromContext(ctx)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}

	var id int

	query := fmt.Sprintf("SET myapp.user_id = %d", userID)
	if _, err := tx.SQL().Exec(query); err != nil {
		return 0, err
	}

	collection := tx.Collection(t.Table())

	var res up.InsertResult
	var err error

	if res, err = collection.Insert(m); err != nil {
		return 0, err
	}

	id = getInsertId(res.ID())

	return id, nil
}
