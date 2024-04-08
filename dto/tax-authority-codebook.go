package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type TaxAuthorityCodebookDTO struct {
	Title                                string  `json:"title"`
	Code                                 string  `json:"code"`
	Active                               bool    `json:"active"`
	TaxPercentage                        float64 `json:"tax_percentage"`
	TaxSupplierID                        int     `json:"tax_supplier_id"`
	ReleasePercentage                    float64 `json:"release_percentage"`
	PioPercentage                        float64 `json:"pio_percentage"`
	PioSupplierID                        int     `json:"pio_supplier_id"`
	PioPercentageEmployerPercentage      float64 `json:"pio_percentage_employer_percentage"`
	PioEmployerSupplierID                int     `json:"pio_employer_supplier_id"`
	PioPercentageEmployeePercentage      float64 `json:"pio_percentage_employee_percentage"`
	PioEmployeeSupplierID                int     `json:"pio_employee_supplier_id"`
	UnemploymentPercentage               float64 `json:"unemployment_percentage"`
	UnemploymentSupplierID               int     `json:"unemployment_supplier_id"`
	UnemploymentEmployerPercentage       float64 `json:"unemployment_employer_percentage"`
	UnemploymentEmployerSupplierID       int     `json:"unemployment_employer_supplier_id"`
	UnemploymentEmployeePercentage       float64 `json:"unemployment_employee_percentage"`
	UnemploymentEmployeeSupplierID       int     `json:"unemployment_employee_supplier_id"`
	LaborFund                            float64 `json:"labor_fund"`
	LaborFundSupplierID                  int     `json:"labor_fund_supplier_id"`
	PreviousIncomePercentageLessThan700  float64 `json:"previous_income_percentage_less_than_700"`
	PreviousIncomePercentageLessThan1000 float64 `json:"previous_income_percentage_less_than_1000"`
	PreviousIncomePercentageMoreThan1000 float64 `json:"previous_income_percentage_more_than_1000"`
	Coefficient                          float64 `json:"coefficient"`
	CoefficientLess700                   float64 `json:"coefficient_less_700"`
	CoefficientLess1000                  float64 `json:"coefficient_less_1000"`
	CoefficientMore1000                  float64 `json:"coefficient_more_1000"`
}

type TaxAuthorityCodebookResponseDTO struct {
	ID                                   int       `json:"id"`
	Title                                string    `json:"title"`
	Code                                 string    `json:"code"`
	Active                               bool      `json:"active"`
	TaxPercentage                        float64   `json:"tax_percentage"`
	TaxSupplierID                        int       `json:"tax_supplier_id"`
	ReleasePercentage                    float64   `json:"release_percentage"`
	PioPercentage                        float64   `json:"pio_percentage"`
	PioSupplierID                        int       `json:"pio_supplier_id"`
	PioPercentageEmployerPercentage      float64   `json:"pio_percentage_employer_percentage"`
	PioEmployerSupplierID                int       `json:"pio_employer_supplier_id"`
	PioPercentageEmployeePercentage      float64   `json:"pio_percentage_employee_percentage"`
	PioEmployeeSupplierID                int       `json:"pio_employee_supplier_id"`
	UnemploymentPercentage               float64   `json:"unemployment_percentage"`
	UnemploymentSupplierID               int       `json:"unemployment_supplier_id"`
	UnemploymentEmployerPercentage       float64   `json:"unemployment_employer_percentage"`
	UnemploymentEmployerSupplierID       int       `json:"unemployment_employer_supplier_id"`
	UnemploymentEmployeePercentage       float64   `json:"unemployment_employee_percentage"`
	UnemploymentEmployeeSupplierID       int       `json:"unemployment_employee_supplier_id"`
	LaborFund                            float64   `json:"labor_fund"`
	LaborFundSupplierID                  int       `json:"labor_fund_supplier_id"`
	PreviousIncomePercentageLessThan700  float64   `json:"previous_income_percentage_less_than_700"`
	PreviousIncomePercentageLessThan1000 float64   `json:"previous_income_percentage_less_than_1000"`
	PreviousIncomePercentageMoreThan1000 float64   `json:"previous_income_percentage_more_than_1000"`
	Coefficient                          float64   `json:"coefficient"`
	CoefficientLess700                   float64   `json:"coefficient_less_700"`
	CoefficientLess1000                  float64   `json:"coefficient_less_1000"`
	CoefficientMore1000                  float64   `json:"coefficient_more_1000"`
	CreatedAt                            time.Time `json:"created_at"`
	UpdatedAt                            time.Time `json:"updated_at"`
}

type TaxAuthorityCodebookFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	Search      *string `json:"search"`
	Active      *bool   `json:"active"`
}

func (dto TaxAuthorityCodebookDTO) ToTaxAuthorityCodebook() *data.TaxAuthorityCodebook {
	return &data.TaxAuthorityCodebook{
		Title:                                dto.Title,
		Code:                                 dto.Code,
		TaxPercentage:                        dto.TaxPercentage,
		TaxSupplierID:                        dto.TaxSupplierID,
		ReleasePercentage:                    dto.ReleasePercentage,
		PioPercentage:                        dto.PioPercentage,
		PioSupplierID:                        dto.PioSupplierID,
		PioPercentageEmployerPercentage:      dto.PioPercentageEmployerPercentage,
		PioEmployerSupplierID:                dto.PioEmployerSupplierID,
		PioPercentageEmployeePercentage:      dto.PioPercentageEmployeePercentage,
		PioEmployeeSupplierID:                dto.PioEmployeeSupplierID,
		UnemploymentPercentage:               dto.UnemploymentPercentage,
		UnemploymentSupplierID:               dto.UnemploymentSupplierID,
		UnemploymentEmployerPercentage:       dto.UnemploymentEmployerPercentage,
		UnemploymentEmployerSupplierID:       dto.UnemploymentEmployerSupplierID,
		UnemploymentEmployeePercentage:       dto.UnemploymentEmployeePercentage,
		UnemploymentEmployeeSupplierID:       dto.UnemploymentEmployeeSupplierID,
		LaborFund:                            dto.LaborFund,
		LaborFundSupplierID:                  dto.LaborFundSupplierID,
		PreviousIncomePercentageLessThan700:  dto.PreviousIncomePercentageLessThan700,
		PreviousIncomePercentageLessThan1000: dto.PreviousIncomePercentageLessThan1000,
		PreviousIncomePercentageMoreThan1000: dto.PreviousIncomePercentageMoreThan1000,
		Coefficient:                          dto.Coefficient,
		Active:                               dto.Active,
		CoefficientLess700:                   dto.CoefficientLess700,
		CoefficientLess1000:                  dto.CoefficientLess1000,
		CoefficientMore1000:                  dto.CoefficientMore1000,
	}
}

func ToTaxAuthorityCodebookResponseDTO(data data.TaxAuthorityCodebook) TaxAuthorityCodebookResponseDTO {
	return TaxAuthorityCodebookResponseDTO{
		ID:                                   data.ID,
		Title:                                data.Title,
		Code:                                 data.Code,
		TaxPercentage:                        data.TaxPercentage,
		TaxSupplierID:                        data.TaxSupplierID,
		ReleasePercentage:                    data.ReleasePercentage,
		PioPercentage:                        data.PioPercentage,
		PioSupplierID:                        data.PioSupplierID,
		PioPercentageEmployerPercentage:      data.PioPercentageEmployerPercentage,
		PioEmployerSupplierID:                data.PioEmployerSupplierID,
		PioPercentageEmployeePercentage:      data.PioPercentageEmployeePercentage,
		PioEmployeeSupplierID:                data.PioEmployeeSupplierID,
		UnemploymentPercentage:               data.UnemploymentPercentage,
		UnemploymentSupplierID:               data.UnemploymentSupplierID,
		UnemploymentEmployerPercentage:       data.UnemploymentEmployerPercentage,
		UnemploymentEmployerSupplierID:       data.UnemploymentEmployerSupplierID,
		UnemploymentEmployeePercentage:       data.UnemploymentEmployeePercentage,
		UnemploymentEmployeeSupplierID:       data.UnemploymentEmployeeSupplierID,
		LaborFund:                            data.LaborFund,
		LaborFundSupplierID:                  data.LaborFundSupplierID,
		PreviousIncomePercentageLessThan700:  data.PreviousIncomePercentageLessThan700,
		PreviousIncomePercentageLessThan1000: data.PreviousIncomePercentageLessThan1000,
		PreviousIncomePercentageMoreThan1000: data.PreviousIncomePercentageMoreThan1000,
		Coefficient:                          data.Coefficient,
		Active:                               data.Active,
		CoefficientLess700:                   data.CoefficientLess700,
		CoefficientLess1000:                  data.CoefficientLess1000,
		CoefficientMore1000:                  data.CoefficientMore1000,
		CreatedAt:                            data.CreatedAt,
		UpdatedAt:                            data.UpdatedAt,
	}
}

func ToTaxAuthorityCodebookListResponseDTO(tax_authority_codebooks []*data.TaxAuthorityCodebook) []TaxAuthorityCodebookResponseDTO {
	dtoList := make([]TaxAuthorityCodebookResponseDTO, len(tax_authority_codebooks))
	for i, x := range tax_authority_codebooks {
		dtoList[i] = ToTaxAuthorityCodebookResponseDTO(*x)
	}
	return dtoList
}
