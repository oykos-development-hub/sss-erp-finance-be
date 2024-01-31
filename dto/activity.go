package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ActivityDTO struct {
	Title       string `json:"title" validate:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ActivityResponseDTO struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ActivityFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
}

func (dto ActivityDTO) ToActivity() *data.Activity {
	return &data.Activity{
		Title:       dto.Title,
		Code:        dto.Code,
		Description: dto.Description,
	}
}

func ToActivityResponseDTO(data data.Activity) ActivityResponseDTO {
	return ActivityResponseDTO{
		ID:          data.ID,
		Title:       data.Title,
		Code:        data.Code,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ToActivityListResponseDTO(activities []*data.Activity) []ActivityResponseDTO {
	dtoList := make([]ActivityResponseDTO, len(activities))
	for i, x := range activities {
		dtoList[i] = ToActivityResponseDTO(*x)
	}
	return dtoList
}
