package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ErrorLogDTO struct {
	Error  string                `json:"error"`
	Code   int                   `json:"code"`
	Entity data.HandlersEntities `json:"entity"`
}

type ErrorLogResponseDTO struct {
	ID        int                   `json:"id"`
	Error     string                `json:"error"`
	Code      int                   `json:"code"`
	Entity    data.HandlersEntities `json:"entity"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type ErrorLogFilterDTO struct {
	Page        *int       `json:"page"`
	Size        *int       `json:"size"`
	DateOfStart *time.Time `json:"date_of_start"`
	DateOfEnd   *time.Time `json:"date_of_end"`
	Entity      *string    `json:"entity"`
}

func (dto ErrorLogDTO) ToErrorLog() *data.ErrorLog {
	return &data.ErrorLog{
		Error:  dto.Error,
		Code:   dto.Code,
		Entity: dto.Entity,
	}
}

func ToErrorLogResponseDTO(data data.ErrorLog) ErrorLogResponseDTO {
	return ErrorLogResponseDTO{
		ID:        data.ID,
		Error:     data.Error,
		Code:      data.Code,
		Entity:    data.Entity,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func ToErrorLogListResponseDTO(error_logs []*data.ErrorLog) []ErrorLogResponseDTO {
	dtoList := make([]ErrorLogResponseDTO, len(error_logs))
	for i, x := range error_logs {
		dtoList[i] = ToErrorLogResponseDTO(*x)
	}
	return dtoList
}
