package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type InvoicesFilter struct {
	Search             *string                 `json:"search"`
	Page               *int                    `json:"page"`
	Size               *int                    `json:"size"`
	Year               *int                    `json:"year"`
	Status             *data.InvoiceStatus     `json:"status"`
	SupplierID         *int                    `json:"supplier_id"`
	OrganizationUnitID *int                    `json:"organization_unit_id"`
	ActivityID         *int                    `json:"activity_id"`
	OrderID            *int                    `json:"order_id"`
	Type               *data.TypesOfObligation `json:"type"`
	PassedToInventory  *bool                   `json:"passed_to_inventory"`
	Registred          *bool                   `json:"registred"`
}

type InvoiceDTO struct {
	InvoiceNumber          string                 `json:"invoice_number"`
	PassedToInventory      bool                   `json:"passed_to_inventory"`
	PassedToAccounting     bool                   `json:"passed_to_accounting"`
	IsInvoice              bool                   `json:"is_invoice"`
	Issuer                 string                 `json:"issuer"`
	Status                 data.InvoiceStatus     `json:"status"`
	Type                   data.TypesOfObligation `json:"type"`
	TaxAuthorityCodebookID int                    `json:"tax_authority_codebook_id"`
	TypeOfSubject          int                    `json:"type_of_subject"`
	TypeOfContract         int                    `json:"type_of_contract"`
	SourceOfFunding        string                 `json:"source_of_funding"`
	Supplier               string                 `json:"supplier"`
	GrossPrice             float64                `json:"gross_price"`
	MunicipalityID         int                    `json:"municipality_id"`
	VATPrice               float64                `json:"vat_price"`
	Registred              *bool                  `json:"registred"`
	SupplierID             int                    `json:"supplier_id"`
	TypeOfDecision         int                    `json:"type_of_decision"`
	OrderID                int                    `json:"order_id"`
	OrganizationUnitID     int                    `json:"organization_unit_id"`
	ActivityID             int                    `json:"activity_id"`
	ProFormaInvoiceNumber  string                 `json:"pro_forma_invoice_number"`
	ProFormaInvoiceDate    *time.Time             `json:"pro_forma_invoice_date"`
	DateOfInvoice          time.Time              `json:"date_of_invoice"`
	ReceiptDate            time.Time              `json:"receipt_date"`
	DateOfPayment          time.Time              `json:"date_of_payment"`
	DateOfStart            time.Time              `json:"date_of_start"`
	DateOfEnd              time.Time              `json:"date_of_end"`
	SSSInvoiceReceiptDate  *time.Time             `json:"sss_invoice_receipt_date"`
	AdditionalExpenses     []AdditionalExpenseDTO `json:"additional_expenses"`
	FileID                 int                    `json:"file_id"`
	ProFormaInvoiceFileID  int                    `json:"pro_forma_invoice_file_id"`
	BankAccount            string                 `json:"bank_account"`
	Description            string                 `json:"description"`
}

type InvoiceResponseDTO struct {
	ID                     int                            `json:"id"`
	InvoiceNumber          string                         `json:"invoice_number"`
	PassedToInventory      bool                           `json:"passed_to_inventory"`
	PassedToAccounting     bool                           `json:"passed_to_accounting"`
	IsInvoice              bool                           `json:"is_invoice"`
	Type                   data.TypesOfObligation         `json:"type"`
	Registred              *bool                          `json:"registred"`
	TaxAuthorityCodebookID int                            `json:"tax_authority_codebook_id"`
	MunicipalityID         int                            `json:"municipality_id"`
	TypeOfSubject          int                            `json:"type_of_subject"`
	TypeOfContract         int                            `json:"type_of_contract"`
	SourceOfFunding        string                         `json:"source_of_funding"`
	Supplier               string                         `json:"supplier"`
	Issuer                 string                         `json:"issuer"`
	Status                 data.InvoiceStatus             `json:"status"`
	GrossPrice             float64                        `json:"gross_price"`
	TypeOfDecision         int                            `json:"type_of_decision"`
	VATPrice               float64                        `json:"vat_price"`
	NetPrice               float64                        `json:"net_price"`
	SupplierID             int                            `json:"supplier_id"`
	OrderID                int                            `json:"order_id"`
	OrganizationUnitID     int                            `json:"organization_unit_id"`
	ActivityID             int                            `json:"activity_id"`
	ProFormaInvoiceNumber  string                         `json:"pro_forma_invoice_number"`
	ProFormaInvoiceDate    *time.Time                     `json:"pro_forma_invoice_date"`
	DateOfInvoice          time.Time                      `json:"date_of_invoice"`
	ReceiptDate            time.Time                      `json:"receipt_date"`
	DateOfPayment          time.Time                      `json:"date_of_payment"`
	SSSInvoiceReceiptDate  *time.Time                     `json:"sss_invoice_receipt_date"`
	DateOfStart            time.Time                      `json:"date_of_start"`
	DateOfEnd              time.Time                      `json:"date_of_end"`
	FileID                 int                            `json:"file_id"`
	ProFormaInvoiceFileID  int                            `json:"pro_forma_invoice_file_id"`
	BankAccount            string                         `json:"bank_account"`
	Description            string                         `json:"description"`
	Articles               []ArticleResponseDTO           `json:"articles"`
	AdditionalExpenses     []AdditionalExpenseResponseDTO `json:"additional_expenses"`
	CreatedAt              time.Time                      `json:"created_at"`
	UpdatedAt              time.Time                      `json:"updated_at"`
}

func (dto InvoiceDTO) ToInvoice() *data.Invoice {
	return &data.Invoice{
		InvoiceNumber:          dto.InvoiceNumber,
		PassedToInventory:      dto.PassedToInventory,
		PassedToAccounting:     dto.PassedToAccounting,
		Registred:              dto.Registred,
		IsInvoice:              dto.IsInvoice,
		Status:                 dto.Status,
		GrossPrice:             dto.GrossPrice,
		VATPrice:               dto.VATPrice,
		TaxAuthorityCodebookID: dto.TaxAuthorityCodebookID,
		ActivityID:             dto.ActivityID,
		SupplierID:             dto.SupplierID,
		MunicipalityID:         dto.MunicipalityID,
		OrderID:                dto.OrderID,
		Issuer:                 dto.Issuer,
		OrganizationUnitID:     dto.OrganizationUnitID,
		ProFormaInvoiceNumber:  dto.ProFormaInvoiceNumber,
		ProFormaInvoiceDate:    dto.ProFormaInvoiceDate,
		DateOfInvoice:          dto.DateOfInvoice,
		ReceiptDate:            dto.ReceiptDate,
		DateOfPayment:          dto.DateOfPayment,
		SSSInvoiceReceiptDate:  dto.SSSInvoiceReceiptDate,
		FileID:                 dto.FileID,
		ProFormaInvoiceFileID:  dto.ProFormaInvoiceFileID,
		BankAccount:            dto.BankAccount,
		Description:            dto.Description,
		TypeOfDecision:         dto.TypeOfDecision,
		Type:                   dto.Type,
		TypeOfSubject:          dto.TypeOfSubject,
		SourceOfFunding:        dto.SourceOfFunding,
		Supplier:               dto.Supplier,
		TypeOfContract:         dto.TypeOfContract,
		DateOfStart:            dto.DateOfStart,
		DateOfEnd:              dto.DateOfEnd,
	}
}

func ToInvoiceResponseDTO(data data.Invoice) InvoiceResponseDTO {
	return InvoiceResponseDTO{
		ID:                     data.ID,
		PassedToInventory:      data.PassedToInventory,
		PassedToAccounting:     data.PassedToAccounting,
		Registred:              data.Registred,
		IsInvoice:              data.IsInvoice,
		InvoiceNumber:          data.InvoiceNumber,
		Status:                 data.Status,
		Issuer:                 data.Issuer,
		GrossPrice:             data.GrossPrice,
		VATPrice:               data.VATPrice,
		TaxAuthorityCodebookID: data.TaxAuthorityCodebookID,
		ActivityID:             data.ActivityID,
		SupplierID:             data.SupplierID,
		MunicipalityID:         data.MunicipalityID,
		OrderID:                data.OrderID,
		OrganizationUnitID:     data.OrganizationUnitID,
		DateOfInvoice:          data.DateOfInvoice,
		ReceiptDate:            data.ReceiptDate,
		DateOfPayment:          data.DateOfPayment,
		SSSInvoiceReceiptDate:  data.SSSInvoiceReceiptDate,
		FileID:                 data.FileID,
		ProFormaInvoiceFileID:  data.ProFormaInvoiceFileID,
		BankAccount:            data.BankAccount,
		TypeOfDecision:         data.TypeOfDecision,
		Description:            data.Description,
		Type:                   data.Type,
		TypeOfSubject:          data.TypeOfSubject,
		SourceOfFunding:        data.SourceOfFunding,
		Supplier:               data.Supplier,
		TypeOfContract:         data.TypeOfContract,
		DateOfStart:            data.DateOfStart,
		DateOfEnd:              data.DateOfEnd,
		ProFormaInvoiceNumber:  data.ProFormaInvoiceNumber,
		ProFormaInvoiceDate:    data.ProFormaInvoiceDate,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
	}
}

func ToInvoiceListResponseDTO(invoices []*data.Invoice) []InvoiceResponseDTO {
	dtoList := make([]InvoiceResponseDTO, len(invoices))
	for i, x := range invoices {
		dtoList[i] = ToInvoiceResponseDTO(*x)
	}
	return dtoList
}
