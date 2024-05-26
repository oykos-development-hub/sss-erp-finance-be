package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type AccountingOrderItemsTitle string

var (
	MainBillTitle                           AccountingOrderItemsTitle = "Korektivni račun"
	SupplierTitle                           AccountingOrderItemsTitle = "Dobavljač"
	TaxTitle                                AccountingOrderItemsTitle = "Porez"
	SubTaxTitle                             AccountingOrderItemsTitle = "Prirez"
	BankTitle                               AccountingOrderItemsTitle = "Banka"
	SuspensionsTitle                        AccountingOrderItemsTitle = "Obustave"
	PIOContributionsTitle                   AccountingOrderItemsTitle = "Doprinos za PIO"
	UnemployementContributionsTitle         AccountingOrderItemsTitle = "Doprinos za nezaposlenost"
	PIOEmployeeContributionsTitle           AccountingOrderItemsTitle = "Doprinos za PIO (zaposleni)"
	UnemployementEmployeeContributionsTitle AccountingOrderItemsTitle = "Doprinos za nezaposlenost (zaposleni)"
	PIOEmployerContributionsTitle           AccountingOrderItemsTitle = "Doprinos za PIO (poslodavac)"
	UnemployementEmployerContributionsTitle AccountingOrderItemsTitle = "Doprinos za nezaposlenost (poslodavac)"
	LaborContributionsTitle                 AccountingOrderItemsTitle = "Doprinos za Fond rada"
	CostTitle                               AccountingOrderItemsTitle = "Izdatak"
	AllocatedAmountTitle                    AccountingOrderItemsTitle = "Rezervisana sredstva"
	ProcessCostTitle                        AccountingOrderItemsTitle = "Trošak izvršenja"
	LawyerCostTitle                         AccountingOrderItemsTitle = "Trošak advokata"
	EnforcedPaymentTitle                    AccountingOrderItemsTitle = "Prinudna naplata"
)

type ObligationTitles string

var (
	NetTitle                                 ObligationTitles = "Neto"
	ObligationTaxTitle                       ObligationTitles = "Porez"
	ObligationSubTaxTitle                    ObligationTitles = "Prirez"
	LaborFundTitle                           ObligationTitles = "Fond rada"
	ContributionForPIOTitle                  ObligationTitles = "PIO"
	ContributionForUnemploymentTitle         ObligationTitles = "Nezaposlenost"
	ContributionForPIOEmployeeTitle          ObligationTitles = "PIO na teret zaposlenog"
	ContributionForPIOEmployerTitle          ObligationTitles = "PIO na teret poslodavca"
	ContributionForUnemploymentEmployeeTitle ObligationTitles = "Nezaposlenost na teret zaposlenog"
	ContributionForUnemploymentEmployerTitle ObligationTitles = "Nezaposlenost na teret poslodavca"
)

var InitialStateTitle = "Početno stanje"

type AccountingEntry struct {
	ID                 int               `db:"id,omitempty"`
	Title              string            `db:"title"`
	Type               TypesOfObligation `db:"type"`
	IDOfEntry          int               `db:"id_of_entry"`
	OrganizationUnitID int               `db:"organization_unit_id"`
	DateOfBooking      time.Time         `db:"date_of_booking"`
	CreatedAt          time.Time         `db:"created_at,omitempty"`
	UpdatedAt          time.Time         `db:"updated_at"`
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

type PaymentOrdersForAccounting struct {
	PaymentOrderID int       `json:"payment_order_id"`
	SupplierID     *int      `json:"supplier_id"`
	Date           time.Time `json:"date"`
	Title          string    `json:"title"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
}

type AnalyticalCardFilter struct {
	SupplierID         *int       `json:"supplier_id"`
	OrganizationUnitID int        `json:"organization_unit_id"`
	DateOfStart        *time.Time `json:"date_of_start"`
	DateOfEnd          *time.Time `json:"date_of_end"`
	DateOfStartBooking *time.Time `json:"date_of_start_booking"`
	DateOfEndBooking   *time.Time `json:"date_of_end_booking"`
}

type AnalyticalCard struct {
	InitialState            float64               `json:"initial_state"`
	SumCreditAmount         float64               `json:"sum_credit_amount"`
	SumDebitAmount          float64               `json:"sum_debit_amount"`
	SumCreditAmountInPeriod float64               `json:"sum_credit_amount_in_period"`
	SumDebitAmountInPeriod  float64               `json:"sum_debit_amount_in_period"`
	SupplierID              int                   `json:"supplier_id"`
	Items                   []AnalyticalCardItems `json:"items"`
}

type AnalyticalCardItems struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	CreditAmount   float64   `json:"credit_amount"`
	DebitAmount    float64   `json:"debit_amount"`
	Balance        float64   `json:"balance"`
	DateOfBooking  time.Time `json:"date_of_booking"`
	IDOfEntry      int       `json:"id_of_entry"`
	Date           time.Time `json:"date"`
	Type           string    `json:"type"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
						and i.invoice_number is not null and i.invoice_number <> ''
						and (COALESCE($3, '') = '' OR i.invoice_number LIKE '%' || $3 || '%')
						group by i.id;`

	queryForAdditionalExpenses := `select i.id, sum(a.price), i.type, i.invoice_number, i.supplier_id, i.date_of_invoice
	                               from additional_expenses a
	                               left join invoices i on a.invoice_id = i.id
	                               where a.invoice_id = i.id and
	                               i.organization_unit_id = $1 and i.registred = false 
								   and (COALESCE($2, '') = '' OR i.type = $2)
								   and (COALESCE($3, '') = '' OR i.invoice_number LIKE '%' || $3 || '%')
	                               group by i.id, i.invoice_number order by i.id;`

	queryForSalaryAdditionalExpenses := `select s.id, sum(a.amount), s.month, s.date_of_calculation
	                                     from salary_additional_expenses a
	                                     left join salaries s on s.id = a.salary_id
	                                     where s.organization_unit_id = $1 and s.registred = false
										 and (COALESCE($2, '') = '' OR s.month LIKE '%' || $2 || '%')
	                                     group by s.id, s.month order by s.id;`

	if filter.Type == nil || *filter.Type == TypeInvoice {
		rows, err := Upper.SQL().Query(queryForInvoices, filter.OrganizationUnitID, TypeInvoice, filter.Search)
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
	}

	rows, err := Upper.SQL().Query(queryForAdditionalExpenses, filter.OrganizationUnitID, filter.Type, filter.Search)
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

	if filter.Type == nil || *filter.Type == TypeSalary {
		rows, err = Upper.SQL().Query(queryForSalaryAdditionalExpenses, filter.OrganizationUnitID, filter.Search)
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
	}

	total := uint64(len(items))

	return items, &total, nil
}

func (t *AccountingEntry) GetPaymentOrdersForAccounting(filter ObligationsFilter) ([]PaymentOrdersForAccounting, *uint64, error) {
	var items []PaymentOrdersForAccounting

	query := `select id, supplier_id, sap_id, date_of_sap, amount
			  from payment_orders 
			  where registred = false and sap_id is not null and date_of_sap is not null and sap_id <> '' and date_of_sap <> '0001-01-01'
			  and organization_unit_id = $1 and (COALESCE($2, '') = '' OR sap_id LIKE '%' || $2 || '%');`

	rows, err := Upper.SQL().Query(query, filter.OrganizationUnitID, filter.Search)

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paymentOrder PaymentOrdersForAccounting
		err = rows.Scan(&paymentOrder.PaymentOrderID, &paymentOrder.SupplierID, &paymentOrder.Title, &paymentOrder.Date, &paymentOrder.Price)

		if err != nil {
			return nil, nil, err
		}

		items = append(items, paymentOrder)
	}

	total := uint64(len(items))

	return items, &total, nil
}

func (t *AccountingEntry) GetEnforcedPaymentsForAccounting(filter ObligationsFilter) ([]PaymentOrdersForAccounting, *uint64, error) {
	var items []PaymentOrdersForAccounting

	query := `select id, supplier_id, sap_id, date_of_sap, amount + amount_for_lawyer + amount_for_agent
			  from enforced_payments 
			  where registred = false 
			  and organization_unit_id = $1 and (COALESCE($2, '') = '' OR sap_id LIKE '%' || $2 || '%');`

	rows, err := Upper.SQL().Query(query, filter.OrganizationUnitID, filter.Search)

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paymentOrder PaymentOrdersForAccounting
		err = rows.Scan(&paymentOrder.PaymentOrderID, &paymentOrder.SupplierID, &paymentOrder.Title, &paymentOrder.Date, &paymentOrder.Price)

		if err != nil {
			return nil, nil, err
		}

		items = append(items, paymentOrder)
	}

	total := uint64(len(items))

	return items, &total, nil
}

func (t *AccountingEntry) GetReturnedEnforcedPaymentsForAccounting(filter ObligationsFilter) ([]PaymentOrdersForAccounting, *uint64, error) {
	var items []PaymentOrdersForAccounting

	query := `select id, supplier_id, sap_id, date_of_sap, return_amount
			  from enforced_payments 
			  where registred_return = false 
			  and return_date is not null and return_date <> '0001-01-01' 
			  and return_file_id is not null and return_file_id <> 0 
			  and organization_unit_id = $1 and (COALESCE($2, '') = '' OR sap_id LIKE '%' || $2 || '%');`

	rows, err := Upper.SQL().Query(query, filter.OrganizationUnitID, filter.Search)

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paymentOrder PaymentOrdersForAccounting
		err = rows.Scan(&paymentOrder.PaymentOrderID, &paymentOrder.SupplierID, &paymentOrder.Title, &paymentOrder.Date, &paymentOrder.Price)

		if err != nil {
			return nil, nil, err
		}

		items = append(items, paymentOrder)
	}

	total := uint64(len(items))

	return items, &total, nil
}

func (t *AccountingEntry) GetAnalyticalCard(filter AnalyticalCardFilter) (*AnalyticalCard, error) {
	var item AnalyticalCard
	var sumDebitAmount float64
	var sumCreditAmount float64

	queryForInitialState := `SELECT SUM(a.credit_amount) - SUM(a.debit_amount) AS saldo 
    						 FROM accounting_entry_items a
    						 LEFT JOIN accounting_entries ae ON ae.id = a.entry_id
    						 WHERE a.supplier_id = $1 AND ae.organization_unit_id = $2 
    						 AND ((cast($3 AS timestamp) IS NOT NULL AND a.date <  cast($3 AS timestamp)) OR 
    						      (cast($4 AS timestamp) IS NOT NULL AND ae.date_of_booking < cast($4 AS timestamp)));`

	queryForItems := `select a.date, a.title, ae.date_of_booking, a.debit_amount, a.credit_amount, 
						COALESCE(i.invoice_number, s.month, p.sap_id, e.sap_id, ep.sap_id) as document_number, a.type, ae.id_of_entry
						from accounting_entry_items a
						left join accounting_entries ae on ae.id = a.entry_id
						left join invoices i on i.id = a.invoice_id
						left join salaries s on s.id = a.salary_id
						left join payment_orders p on p.id = a.payment_order_id
						left join enforced_payments e on e.id = a.enforced_payment_id
						left join enforced_payments ep on ep.id = a.return_enforced_payment_id
						where a.supplier_id = $1 and ae.organization_unit_id = $2 and
						((cast($3 AS timestamp) is not null and a.date >= cast($3 AS timestamp) and cast($4 AS timestamp) is not null and a.date <= cast($4 AS timestamp)) or
						(cast($5 AS timestamp) is not null and ae.date_of_booking >= cast($5 AS timestamp) and cast($6 AS timestamp) is not null and ae.date_of_booking <= cast($6 AS timestamp)));`

	rows, err := Upper.SQL().Query(queryForInitialState, filter.SupplierID, filter.OrganizationUnitID, filter.DateOfStart, filter.DateOfStartBooking)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	id := 0

	for rows.Next() {
		var balance *float64
		err = rows.Scan(&balance)

		if err != nil {
			return nil, err
		}

		if balance == nil {
			item.InitialState = 0
		} else {
			item.InitialState = *balance
		}
		var dateOfInitialState time.Time

		if filter.DateOfStart != nil {
			dateOfInitialState = *filter.DateOfStart
		} else if filter.DateOfStartBooking != nil {
			dateOfInitialState = *filter.DateOfStartBooking
		}

		item.Items = append(item.Items, AnalyticalCardItems{
			ID:            id,
			Title:         InitialStateTitle,
			CreditAmount:  item.InitialState,
			DateOfBooking: dateOfInitialState,
			Balance:       item.InitialState,
		})

		id++
	}

	rows, err = Upper.SQL().Query(queryForItems, filter.SupplierID, filter.OrganizationUnitID, filter.DateOfStart, filter.DateOfEnd, filter.DateOfStartBooking, filter.DateOfEndBooking)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var analyticItem AnalyticalCardItems
		err = rows.Scan(&analyticItem.Date, &analyticItem.Title, &analyticItem.DateOfBooking, &analyticItem.DebitAmount, &analyticItem.CreditAmount,
			&analyticItem.DocumentNumber, &analyticItem.Type, &analyticItem.IDOfEntry)

		if err != nil {
			return nil, err
		}

		analyticItem.ID = id
		analyticItem.Balance = item.Items[id-1].Balance + analyticItem.CreditAmount - analyticItem.DebitAmount
		sumDebitAmount += analyticItem.DebitAmount
		sumCreditAmount += analyticItem.CreditAmount

		id++

		item.Items = append(item.Items, analyticItem)
	}

	item.SumCreditAmountInPeriod = sumCreditAmount
	item.SumDebitAmountInPeriod = sumDebitAmount

	item.SumDebitAmount = sumDebitAmount
	item.SumCreditAmount = sumCreditAmount + item.InitialState

	return &item, nil
}

func (t *AccountingEntry) GetAllSuppliers(filter AnalyticalCardFilter) ([]int, error) {
	var items []int

	/*queryForInitialState := `SELECT SUM(a.credit_amount) - SUM(a.debit_amount) AS saldo
	FROM accounting_entry_items a
	LEFT JOIN accounting_entries ae ON ae.id = a.entry_id
	WHERE a.supplier_id = $1 AND ae.organization_unit_id = $2
	AND (($3 IS NOT NULL AND a.date < $3) OR
	     ($4 IS NOT NULL AND ae.date_of_booking < $4));`*/

	queryForItems := `select a.supplier_id
						from accounting_entry_items a
						left join accounting_entries ae on ae.id = a.entry_id
						where ae.organization_unit_id = $1 and
						((cast($2 AS timestamp) is not null and a.date <= cast($2 AS timestamp)) or
						(cast($3 AS timestamp) is not null and ae.date_of_booking <= cast($3 AS timestamp)))
						group by a.supplier_id;`

	rows, err := Upper.SQL().Query(queryForItems, filter.OrganizationUnitID, filter.DateOfEnd, filter.DateOfEndBooking)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var supplierID int
		err = rows.Scan(&supplierID)

		if err != nil {
			return nil, err
		}
		if supplierID != 0 {
			items = append(items, supplierID)
		}
	}

	return items, nil
}
