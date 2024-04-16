package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// DepositPayment struct
type DepositPayment struct {
	ID                        int        `db:"id,omitempty"`
	OrganizationUnitID        int        `db:"organization_unit_id"`
	Payer                     string     `db:"payer"`
	CaseNumber                string     `db:"case_number"`
	PartyName                 string     `db:"party_name"`
	NumberOfBankStatement     string     `db:"number_of_bank_statement"`
	DateOfBankStatement       string     `db:"date_of_bank_statement"`
	AccountID                 int        `db:"account_id"`
	Amount                    float64    `db:"amount"`
	MainBankAccount           bool       `db:"main_bank_account"`
	CurrentBankAccount        string     `db:"current_bank_account"`
	DateOfTransferMainAccount *time.Time `db:"date_of_transfer_main_account"`
	FileID                    *int       `db:"file_id"`
	CreatedAt                 time.Time  `db:"created_at,omitempty"`
	UpdatedAt                 time.Time  `db:"updated_at"`
}

// Table returns the table name
func (t *DepositPayment) Table() string {
	return "deposit_payments"
}

// GetAll gets all records from the database, using upper
func (t *DepositPayment) GetAll(page *int, size *int, condition *up.AndExpr, orders []interface{}) ([]*DepositPayment, *uint64, error) {
	collection := Upper.Collection(t.Table())
	var all []*DepositPayment
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
func (t *DepositPayment) Get(id int) (*DepositPayment, error) {
	var one DepositPayment
	collection := Upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *DepositPayment) Update(tx up.Session, m DepositPayment) error {
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
func (t *DepositPayment) Delete(id int) error {
	collection := Upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *DepositPayment) Insert(tx up.Session, m DepositPayment) (int, error) {
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

func (t *DepositPayment) GetDepositPaymentByCaseNumber(caseNumber string) (DepositPayment, error) {
	var response DepositPayment

	query1 := `select sum(amount) from deposit_payments where case_number = $1`
	query2 := `select sum(price) from deposit_additional_expenses a 
			   left join deposit_payment_orders d on a.payment_order_id = d.id
			   where d.case_number = $1`

	rows1, err := Upper.SQL().Query(query1, caseNumber)
	if err != nil {
		return response, err
	}
	defer rows1.Close()

	var amountPayments float64
	for rows1.Next() {
		err = rows1.Scan(&amountPayments)

		if err != nil {
			return response, err
		}
	}

	rows2, err := Upper.SQL().Query(query2, caseNumber)
	if err != nil {
		return response, err
	}
	defer rows2.Close()

	var amountSpending float64
	for rows2.Next() {
		err = rows2.Scan(&amountSpending)

		if err != nil {
			return response, err
		}
	}

	response.Amount = amountPayments - amountSpending

	return response, nil
}

func (t *DepositPayment) GetCaseNumber(orgUnitID int) ([]*DepositPayment, error) {
	var response []*DepositPayment

	query1 := ` select case_number, sum(amount) from deposit_payments where organization_unit_id = $1 group by case_number`
	query2 := `select sum(price) from deposit_additional_expenses a
			   left join deposit_payment_orders d on a.payment_order_id = d.id
			   where d.case_number = $1`

	rows1, err := Upper.SQL().Query(query1, orgUnitID)
	if err != nil {
		return response, err
	}
	defer rows1.Close()

	for rows1.Next() {
		var item DepositPayment
		var caseNumber string
		var amount float64
		err = rows1.Scan(&caseNumber, &amount)

		item.CaseNumber = caseNumber

		item.Amount = amount

		if err != nil {
			return response, err
		}

		rows2, err := Upper.SQL().Query(query2, &item.CaseNumber)
		if err != nil {
			return response, err
		}
		defer rows2.Close()

		var amountSpending float64
		for rows2.Next() {
			err = rows2.Scan(&amountSpending)

			if err != nil {
				return response, err
			}
		}

		item.Amount -= amountSpending

		if item.Amount > 0 {
			response = append(response, &item)
		}
	}

	return response, nil
}
