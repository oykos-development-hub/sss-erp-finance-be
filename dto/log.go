package dto

import (
	"encoding/json"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type LogDTO struct {
	ChangedAt time.Time         `json:"changed_at"`
	UserID    int               `json:"user_id"`
	ItemID    int               `json:"item_id"`
	Operation data.LogOperation `json:"operation"`
	Entity    data.LogEntity    `json:"entity"`
	OldState  json.RawMessage   `json:"old_state"`
	NewState  json.RawMessage   `json:"new_state"`
}

type LogResponseDTO struct {
	ID        int               `json:"id,omitempty"`
	ChangedAt time.Time         `json:"changed_at"`
	UserID    int               `json:"user_id"`
	ItemID    int               `json:"item_id"`
	Operation data.LogOperation `json:"operation"`
	Entity    data.LogEntity    `json:"entity"`
	OldState  json.RawMessage   `json:"old_state"`
	NewState  json.RawMessage   `json:"new_state"`
}

type LogFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
	Entity      *string `json:"entity"`
	UserID      *int    `json:"user_id"`
	ItemID      *int    `json:"item_id"`
	Operation   *string `json:"operation"`
}

func (dto LogDTO) ToLog() *data.Log {
	return &data.Log{
		ChangedAt: dto.ChangedAt,
		UserID:    dto.UserID,
		Operation: dto.Operation,
		Entity:    dto.Entity,
		OldState:  dto.OldState,
		NewState:  dto.NewState,
	}
}

func ToLogResponseDTO(data data.Log) LogResponseDTO {
	return LogResponseDTO{
		ID:        data.ID,
		ChangedAt: data.ChangedAt,
		UserID:    data.UserID,
		Operation: data.Operation,
		Entity:    data.Entity,
		OldState:  data.OldState,
		NewState:  data.NewState,
	}
}

func ToLogListResponseDTO(logs []*data.Log) []LogResponseDTO {
	dtoList := make([]LogResponseDTO, len(logs))
	for i, x := range logs {
		dtoList[i] = ToLogResponseDTO(*x)
	}
	return dtoList
}
