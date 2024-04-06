package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type FixedDepositWillDispatchDTO struct {
	WillID         int       `json:"will_id"`
	DispatchTypeID int       `json:"dispatch_type_id"`
	JudgeID        int       `json:"judge_id"`
	CaseNumber     string    `json:"case_number"`
	DateOfDispatch time.Time `json:"date_of_dispatch"`
	FileID         int       `json:"file_id"`
}

type FixedDepositWillDispatchResponseDTO struct {
	ID             int       `json:"id"`
	WillID         int       `json:"will_id"`
	DispatchTypeID int       `json:"dispatch_type_id"`
	JudgeID        int       `json:"judge_id"`
	CaseNumber     string    `json:"case_number"`
	DateOfDispatch time.Time `json:"date_of_dispatch"`
	FileID         int       `json:"file_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FixedDepositWillDispatchFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	WillID      *int    `json:"will_id"`
}

func (dto FixedDepositWillDispatchDTO) ToFixedDepositWillDispatch() *data.FixedDepositWillDispatch {
	return &data.FixedDepositWillDispatch{
		WillID:         dto.WillID,
		DispatchTypeID: dto.DispatchTypeID,
		JudgeID:        dto.JudgeID,
		CaseNumber:     dto.CaseNumber,
		DateOfDispatch: dto.DateOfDispatch,
		FileID:         dto.FileID,
	}
}

func ToFixedDepositWillDispatchResponseDTO(data data.FixedDepositWillDispatch) FixedDepositWillDispatchResponseDTO {
	return FixedDepositWillDispatchResponseDTO{
		ID:             data.ID,
		WillID:         data.WillID,
		DispatchTypeID: data.DispatchTypeID,
		JudgeID:        data.JudgeID,
		CaseNumber:     data.CaseNumber,
		DateOfDispatch: data.DateOfDispatch,
		FileID:         data.FileID,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
}

func ToFixedDepositWillDispatchListResponseDTO(fixed_deposit_will_dispatches []*data.FixedDepositWillDispatch) []FixedDepositWillDispatchResponseDTO {
	dtoList := make([]FixedDepositWillDispatchResponseDTO, len(fixed_deposit_will_dispatches))
	for i, x := range fixed_deposit_will_dispatches {
		dtoList[i] = ToFixedDepositWillDispatchResponseDTO(*x)
	}
	return dtoList
}
