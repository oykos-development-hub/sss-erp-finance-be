package dto

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.sudovi.me/erp/finance-api/data"
)

type FixedDepositDispatchDTO struct {
	DepositID    int             `json:"deposit_id"`
	CategoryID   int             `json:"category_id"`
	TypeID       int             `json:"type_id"`
	Unit         string          `json:"unit"`
	Currency     string          `json:"currency"`
	Amount       decimal.Decimal `json:"amount"`
	SerialNumber string          `json:"serial_number"`
	DateOfAction *time.Time      `json:"date_of_action"`
	Subject      string          `json:"subject"`
	Action       string          `json:"action"`
	CaseNumber   string          `json:"case_number"`
	JudgeID      int             `json:"judge_id"`
	FileID       int             `json:"file_id"`
}

type FixedDepositDispatchResponseDTO struct {
	ID           int             `json:"id"`
	DepositID    int             `json:"deposit_id"`
	CategoryID   int             `json:"category_id"`
	TypeID       int             `json:"type_id"`
	Unit         string          `json:"unit"`
	Currency     string          `json:"currency"`
	Amount       decimal.Decimal `json:"amount"`
	SerialNumber string          `json:"serial_number"`
	DateOfAction *time.Time      `json:"date_of_action"`
	Subject      string          `json:"subject"`
	Action       string          `json:"action"`
	CaseNumber   string          `json:"case_number"`
	JudgeID      int             `json:"judge_id"`
	FileID       int             `json:"file_id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type FixedDepositDispatchFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	DepositID   *int    `json:"deposit_id"`
}

func (dto FixedDepositDispatchDTO) ToFixedDepositDispatch() *data.FixedDepositDispatch {
	return &data.FixedDepositDispatch{
		DepositID:    dto.DepositID,
		CategoryID:   dto.CategoryID,
		TypeID:       dto.TypeID,
		Unit:         dto.Unit,
		Currency:     dto.Currency,
		Amount:       dto.Amount,
		SerialNumber: dto.SerialNumber,
		DateOfAction: dto.DateOfAction,
		Subject:      dto.Subject,
		Action:       dto.Action,
		CaseNumber:   dto.CaseNumber,
		FileID:       dto.FileID,
		JudgeID:      dto.JudgeID,
	}
}

func ToFixedDepositDispatchResponseDTO(data data.FixedDepositDispatch) FixedDepositDispatchResponseDTO {
	return FixedDepositDispatchResponseDTO{
		ID:           data.ID,
		DepositID:    data.DepositID,
		CategoryID:   data.CategoryID,
		TypeID:       data.TypeID,
		Unit:         data.Unit,
		Currency:     data.Currency,
		Amount:       data.Amount,
		SerialNumber: data.SerialNumber,
		DateOfAction: data.DateOfAction,
		Subject:      data.Subject,
		Action:       data.Action,
		CaseNumber:   data.CaseNumber,
		JudgeID:      data.JudgeID,
		FileID:       data.FileID,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
}

func ToFixedDepositDispatchListResponseDTO(fixed_deposit_dispatches []*data.FixedDepositDispatch) []FixedDepositDispatchResponseDTO {
	dtoList := make([]FixedDepositDispatchResponseDTO, len(fixed_deposit_dispatches))
	for i, x := range fixed_deposit_dispatches {
		dtoList[i] = ToFixedDepositDispatchResponseDTO(*x)
	}
	return dtoList
}
