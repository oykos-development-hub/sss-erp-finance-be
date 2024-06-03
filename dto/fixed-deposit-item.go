package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type FixedDepositItemDTO struct {
	DepositID          int             `json:"deposit_id"`
	CategoryID         int             `json:"category_id"`
	TypeID             int             `json:"type_id"`
	Unit               string          `json:"unit"`
	Currency           string          `json:"currency"`
	Amount             decimal.Decimal `json:"amount"`
	SerialNumber       string          `json:"serial_number"`
	DateOfConfiscation *time.Time      `json:"date_of_confiscation"`
	CaseNumber         string          `json:"case_number"`
	JudgeID            int             `json:"judge_id"`
	FileID             int             `json:"file_id"`
}

type FixedDepositItemResponseDTO struct {
	ID                 int             `json:"id"`
	DepositID          int             `json:"deposit_id"`
	CategoryID         int             `json:"category_id"`
	TypeID             int             `json:"type_id"`
	Unit               string          `json:"unit"`
	Currency           string          `json:"currency"`
	Amount             decimal.Decimal `json:"amount"`
	SerialNumber       string          `json:"serial_number"`
	DateOfConfiscation *time.Time      `json:"date_of_confiscation"`
	CaseNumber         string          `json:"case_number"`
	JudgeID            int             `json:"judge_id"`
	FileID             int             `json:"file_id"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type FixedDepositItemFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	DepositID   *int    `json:"deposit_id"`
}

func (dto FixedDepositItemDTO) ToFixedDepositItem() *data.FixedDepositItem {
	return &data.FixedDepositItem{
		DepositID:          dto.DepositID,
		CategoryID:         dto.CategoryID,
		TypeID:             dto.TypeID,
		Currency:           dto.Currency,
		SerialNumber:       dto.SerialNumber,
		DateOfConfiscation: dto.DateOfConfiscation,
		CaseNumber:         dto.CaseNumber,
		Unit:               dto.Unit,
		Amount:             dto.Amount,
		JudgeID:            dto.JudgeID,
		FileID:             dto.FileID,
	}
}

func ToFixedDepositItemResponseDTO(data data.FixedDepositItem) FixedDepositItemResponseDTO {
	return FixedDepositItemResponseDTO{
		ID:                 data.ID,
		DepositID:          data.DepositID,
		CategoryID:         data.CategoryID,
		Currency:           data.Currency,
		TypeID:             data.TypeID,
		Unit:               data.Unit,
		Amount:             data.Amount,
		SerialNumber:       data.SerialNumber,
		DateOfConfiscation: data.DateOfConfiscation,
		CaseNumber:         data.CaseNumber,
		JudgeID:            data.JudgeID,
		FileID:             data.FileID,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
	}
}

func ToFixedDepositItemListResponseDTO(fixed_deposit_items []*data.FixedDepositItem) []FixedDepositItemResponseDTO {
	dtoList := make([]FixedDepositItemResponseDTO, len(fixed_deposit_items))
	for i, x := range fixed_deposit_items {
		dtoList[i] = ToFixedDepositItemResponseDTO(*x)
	}
	return dtoList
}
