package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FixedDepositJudgeDTO struct {
	JudgeID     int        `json:"judge_id"`
	DepositID   int        `json:"deposit_id"`
	WillID      int        `json:"will_id"`
	DateOfStart time.Time  `json:"date_of_start"`
	DateOfEnd   *time.Time `json:"date_of_end"`
}

type FixedDepositJudgeResponseDTO struct {
	ID          int        `json:"id"`
	JudgeID     int        `json:"judge_id"`
	DepositID   int        `json:"deposit_id"`
	WillID      int        `json:"will_id"`
	DateOfStart time.Time  `json:"date_of_start"`
	DateOfEnd   *time.Time `json:"date_of_end"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type FixedDepositJudgeFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	DepositID   *int    `json:"deposit_id"`
	WillID      *int    `json:"will_id"`
}

func (dto FixedDepositJudgeDTO) ToFixedDepositJudge() *data.FixedDepositJudge {
	return &data.FixedDepositJudge{
		JudgeID:     dto.JudgeID,
		DepositID:   dto.DepositID,
		WillID:      dto.WillID,
		DateOfStart: dto.DateOfStart,
		DateOfEnd:   dto.DateOfEnd,
	}
}

func ToFixedDepositJudgeResponseDTO(data data.FixedDepositJudge) FixedDepositJudgeResponseDTO {
	return FixedDepositJudgeResponseDTO{
		ID:          data.ID,
		JudgeID:     data.JudgeID,
		DepositID:   data.DepositID,
		WillID:      data.WillID,
		DateOfStart: data.DateOfStart,
		DateOfEnd:   data.DateOfEnd,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ToFixedDepositJudgeListResponseDTO(fixed_deposit_judges []*data.FixedDepositJudge) []FixedDepositJudgeResponseDTO {
	dtoList := make([]FixedDepositJudgeResponseDTO, len(fixed_deposit_judges))
	for i, x := range fixed_deposit_judges {
		dtoList[i] = ToFixedDepositJudgeResponseDTO(*x)
	}
	return dtoList
}
