package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type TaxAuthorityCodebookDTO struct {
	Title                                string  `json:"title"`
	Code                                 string  `json:"code"`
	Percentage                           float64 `json:"percentage"`
	PreviousIncomePercentageLessThan700  float64 `json:"previous_income_percentage_less_than_700"`
	PreviousIncomePercentageLessThan1000 float64 `json:"previous_income_percentage_less_than_1000"`
	PreviousIncomePercentageMoreThan1000 float64 `json:"previous_income_percentage_more_than_1000"`
}

type TaxAuthorityCodebookResponseDTO struct {
	ID                                   int       `json:"id"`
	Title                                string    `json:"title"`
	Code                                 string    `json:"code"`
	Percentage                           float64   `json:"percentage"`
	PreviousIncomePercentageLessThan700  float64   `json:"previous_income_percentage_less_than_700"`
	PreviousIncomePercentageLessThan1000 float64   `json:"previous_income_percentage_less_than_1000"`
	PreviousIncomePercentageMoreThan1000 float64   `json:"previous_income_percentage_more_than_1000"`
	CreatedAt                            time.Time `json:"created_at"`
	UpdatedAt                            time.Time `json:"updated_at"`
}

type TaxAuthorityCodebookFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	Search      *string `json:"search"`
}

func (dto TaxAuthorityCodebookDTO) ToTaxAuthorityCodebook() *data.TaxAuthorityCodebook {
	return &data.TaxAuthorityCodebook{
		Title:                                dto.Title,
		Code:                                 dto.Code,
		Percentage:                           dto.Percentage,
		PreviousIncomePercentageLessThan700:  dto.PreviousIncomePercentageLessThan700,
		PreviousIncomePercentageLessThan1000: dto.PreviousIncomePercentageLessThan1000,
		PreviousIncomePercentageMoreThan1000: dto.PreviousIncomePercentageMoreThan1000,
	}
}

func ToTaxAuthorityCodebookResponseDTO(data data.TaxAuthorityCodebook) TaxAuthorityCodebookResponseDTO {
	return TaxAuthorityCodebookResponseDTO{
		ID:                                   data.ID,
		Title:                                data.Title,
		Code:                                 data.Code,
		Percentage:                           data.Percentage,
		PreviousIncomePercentageLessThan700:  data.PreviousIncomePercentageLessThan700,
		PreviousIncomePercentageLessThan1000: data.PreviousIncomePercentageLessThan1000,
		PreviousIncomePercentageMoreThan1000: data.PreviousIncomePercentageMoreThan1000,
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
