package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ArticleDTO struct {
	Title         string  `json:"title"`
	NetPrice      float64 `json:"net_price"`
	VatPrice      float64 `json:"vat_price"`
	Description   string  `json:"description"`
	InvoiceID     int     `json:"invoice_id"`
	AccountID     int     `json:"account_id"`
	CostAccountID int     `json:"cost_account_id"`
}

type ArticleResponseDTO struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	NetPrice      float64   `json:"net_price"`
	VatPrice      float64   `json:"vat_price"`
	Description   string    `json:"description"`
	InvoiceID     int       `json:"invoice_id"`
	AccountID     int       `json:"account_id"`
	CostAccountID int       `json:"cost_account_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ArticleFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	InvoiceID   *int    `json:"invoice_id"`
	SortByTitle *string `json:"sort_by_title"`
}

func (dto ArticleDTO) ToArticle() *data.Article {
	return &data.Article{
		Title:         dto.Title,
		NetPrice:      dto.NetPrice,
		VatPrice:      dto.VatPrice,
		Description:   dto.Description,
		InvoiceID:     dto.InvoiceID,
		AccountID:     dto.AccountID,
		CostAccountID: dto.CostAccountID,
	}
}

func ToArticleResponseDTO(data data.Article) ArticleResponseDTO {
	return ArticleResponseDTO{
		ID:            data.ID,
		Title:         data.Title,
		NetPrice:      data.NetPrice,
		VatPrice:      data.VatPrice,
		Description:   data.Description,
		InvoiceID:     data.InvoiceID,
		AccountID:     data.AccountID,
		CostAccountID: data.CostAccountID,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}
}

func ToArticleListResponseDTO(articles []*data.Article) []ArticleResponseDTO {
	dtoList := make([]ArticleResponseDTO, len(articles))
	for i, x := range articles {
		dtoList[i] = ToArticleResponseDTO(*x)
	}
	return dtoList
}
