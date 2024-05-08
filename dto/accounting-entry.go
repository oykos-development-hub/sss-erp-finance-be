package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type AccountingEntryDTO struct {
	Title              string                   `json:"title"`
	OrganizationUnitID int                      `json:"organization_unit_id"`
	DateOfBooking      time.Time                `json:"date_of_booking"`
	Type               data.TypesOfObligation   `json:"type"`
	IDOfEntry          int                      `json:"id_of_entry"`
	Items              []AccountingEntryItemDTO `json:"items"`
}

type AccountingEntryResponseDTO struct {
	ID                 int                              `json:"id"`
	Title              string                           `json:"title"`
	IDOfEntry          int                              `json:"id_of_entry"`
	OrganizationUnitID int                              `json:"organization_unit_id"`
	DateOfBooking      time.Time                        `json:"date_of_booking"`
	Type               data.TypesOfObligation           `json:"type"`
	CreditAmount       float64                          `json:"credit_amount"`
	DebitAmount        float64                          `json:"debit_amount"`
	Items              []AccountingEntryItemResponseDTO `json:"items"`
	CreatedAt          time.Time                        `json:"created_at"`
	UpdatedAt          time.Time                        `json:"updated_at"`
}

type AccountingEntryFilterDTO struct {
	Page               *int                    `json:"page"`
	Size               *int                    `json:"size"`
	SortByTitle        *string                 `json:"sort_by_title"`
	OrganizationUnitID *int                    `json:"organization_unit_id"`
	Type               *data.TypesOfObligation `json:"type"`
}

type ObligationForAccounting struct {
	InvoiceID  *int                   `json:"invoice_id"`
	SalaryID   *int                   `json:"salary_id"`
	SupplierID *int                   `json:"supplier_id"`
	Date       time.Time              `json:"date"`
	Type       data.TypesOfObligation `json:"type"`
	Title      string                 `json:"title"`
	Price      float64                `json:"price"`
	Status     string                 `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
}

type PaymentOrdersForAccounting struct {
	PaymentOrderID int       `json:"payment_order_id"`
	SupplierID     *int      `json:"supplier_id"`
	Date           time.Time `json:"date"`
	Title          string    `json:"title"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
}

type AccountingOrderForObligationsData struct {
	InvoiceID               []int     `json:"invoice_id"`
	SalaryID                []int     `json:"salary_id"`
	PaymentOrderID          []int     `json:"payment_order_id"`
	EnforcedPaymentID       []int     `json:"enforced_payment_id"`
	ReturnEnforcedPaymentID []int     `json:"return_enforced_payment_id"`
	DateOfBooking           time.Time `json:"date_of_booking"`
	OrganizationUnitID      int       `json:"organization_unit_id"`
}

type AccountingOrderForObligations struct {
	OrganizationUnitID int                                  `json:"organization_unit_id"`
	DateOfBooking      time.Time                            `json:"date_of_booking"`
	CreditAmount       float32                              `json:"credit_amount"`
	DebitAmount        float32                              `json:"debit_amount"`
	Items              []AccountingOrderItemsForObligations `json:"items"`
}

type AccountingOrderItemsForObligations struct {
	AccountID             int                    `json:"account_id"`
	Title                 string                 `json:"title"`
	CreditAmount          float32                `json:"credit_amount"`
	DebitAmount           float32                `json:"debit_amount"`
	Type                  data.TypesOfObligation `json:"type"`
	SupplierID            int                    `json:"supplier_id"`
	Invoice               DropdownSimple         `json:"invoice"`
	Salary                DropdownSimple         `json:"salary"`
	PaymentOrder          DropdownSimple         `json:"payment_order"`
	EnforcedPayment       DropdownSimple         `json:"enforced_payment"`
	ReturnEnforcedPayment DropdownSimple         `json:"return_enforced_payment"`
}

type DropdownSimple struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (dto AccountingEntryDTO) ToAccountingEntry() *data.AccountingEntry {
	return &data.AccountingEntry{
		Title:              dto.Title,
		OrganizationUnitID: dto.OrganizationUnitID,
		DateOfBooking:      dto.DateOfBooking,
		IDOfEntry:          dto.IDOfEntry,
		Type:               dto.Type,
	}
}

func ToAccountingEntryResponseDTO(data data.AccountingEntry) AccountingEntryResponseDTO {
	return AccountingEntryResponseDTO{
		ID:                 data.ID,
		IDOfEntry:          data.IDOfEntry,
		Title:              data.Title,
		Type:               data.Type,
		OrganizationUnitID: data.OrganizationUnitID,
		DateOfBooking:      data.DateOfBooking,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToAccountingEntryListResponseDTO(accounting_entries []*data.AccountingEntry) []AccountingEntryResponseDTO {
	dtoList := make([]AccountingEntryResponseDTO, len(accounting_entries))
	for i, x := range accounting_entries {
		dtoList[i] = ToAccountingEntryResponseDTO(*x)
	}
	return dtoList
}
