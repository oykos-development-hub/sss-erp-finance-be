package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type ProgramDTO struct {
	Title       string `json:"title" validate:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
	ParentID    int    `json:"parent_id"`
}

type ProgramResponseDTO struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	ParentID    int       `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProgramFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
}

func (dto ProgramDTO) ToProgram() *data.Program {
	return &data.Program{
		Title:       dto.Title,
		Code:        dto.Code,
		Description: dto.Description,
		ParentID:    dto.ParentID,
	}
}

func ToProgramResponseDTO(data data.Program) ProgramResponseDTO {
	return ProgramResponseDTO{
		ID:          data.ID,
		Title:       data.Title,
		Code:        data.Code,
		Description: data.Description,
		ParentID:    data.ParentID,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ToProgramListResponseDTO(programs []*data.Program) []ProgramResponseDTO {
	dtoList := make([]ProgramResponseDTO, len(programs))
	for i, x := range programs {
		dtoList[i] = ToProgramResponseDTO(*x)
	}
	return dtoList
}
