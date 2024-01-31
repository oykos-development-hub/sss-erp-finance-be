package dto

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
)

type GoalIndicatorDTO struct {
	GoalID                   int    `json:"goal_id"`
	PerformanceIndicatorCode string `json:"performance_indicator_code"`
	IndicatorSource          string `json:"indicator_source"`
	BaseYear                 string `json:"base_year"`
	GenderEquality           string `json:"gender_equality"`
	BaseValue                string `json:"base_value"`
	SourceOfInformation      string `json:"source_of_information"`
	UnitOfMeasure            string `json:"unit_of_measure"`
	IndicatorDescription     string `json:"indicator_description"`
	PlannedValue1            string `json:"planned_value_1"`
	RevisedValue1            string `json:"revised_value_1"`
	AchievedValue1           string `json:"achieved_value_1"`
	PlannedValue2            string `json:"planned_value_2"`
	RevisedValue2            string `json:"revised_value_2"`
	AchievedValue2           string `json:"achieved_value_2"`
	PlannedValue3            string `json:"planned_value_3"`
	RevisedValue3            string `json:"revised_value_3"`
	AchievedValue3           string `json:"achieved_value_3"`
}

type GoalIndicatorResponseDTO struct {
	ID                       int       `json:"id"`
	GoalID                   int       `json:"goal_id"`
	PerformanceIndicatorCode string    `json:"performance_indicator_code"`
	IndicatorSource          string    `json:"indicator_source"`
	BaseYear                 string    `json:"base_year"`
	GenderEquality           string    `json:"gender_equality"`
	BaseValue                string    `json:"base_value"`
	SourceOfInformation      string    `json:"source_of_information"`
	UnitOfMeasure            string    `json:"unit_of_measure"`
	IndicatorDescription     string    `json:"indicator_description"`
	PlannedValue1            string    `json:"planned_value_1"`
	RevisedValue1            string    `json:"revised_value_1"`
	AchievedValue1           string    `json:"achieved_value_1"`
	PlannedValue2            string    `json:"planned_value_2"`
	RevisedValue2            string    `json:"revised_value_2"`
	AchievedValue2           string    `json:"achieved_value_2"`
	PlannedValue3            string    `json:"planned_value_3"`
	RevisedValue3            string    `json:"revised_value_3"`
	AchievedValue3           string    `json:"achieved_value_3"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type GoalIndicatorFilterDTO struct {
	Page        *int    `json:"page"`
	Size        *int    `json:"size"`
	SortByTitle *string `json:"sort_by_title"`
}

func (dto GoalIndicatorDTO) ToGoalIndicator() *data.GoalIndicator {
	return &data.GoalIndicator{
		GoalID:                   dto.GoalID,
		PerformanceIndicatorCode: dto.PerformanceIndicatorCode,
		IndicatorSource:          dto.IndicatorSource,
		BaseYear:                 dto.BaseYear,
		GenderEquality:           dto.GenderEquality,
		BaseValue:                dto.BaseValue,
		SourceOfInformation:      dto.SourceOfInformation,
		UnitOfMeasure:            dto.UnitOfMeasure,
		IndicatorDescription:     dto.IndicatorDescription,
		PlannedValue1:            dto.PlannedValue1,
		AchievedValue1:           dto.AchievedValue1,
		RevisedValue1:            dto.RevisedValue1,
		PlannedValue2:            dto.PlannedValue2,
		AchievedValue2:           dto.AchievedValue2,
		RevisedValue2:            dto.RevisedValue2,
		PlannedValue3:            dto.PlannedValue3,
		AchievedValue3:           dto.AchievedValue3,
		RevisedValue3:            dto.RevisedValue3,
	}
}

func ToGoalIndicatorResponseDTO(data data.GoalIndicator) GoalIndicatorResponseDTO {
	return GoalIndicatorResponseDTO{
		ID:                       data.ID,
		GoalID:                   data.GoalID,
		PerformanceIndicatorCode: data.PerformanceIndicatorCode,
		IndicatorSource:          data.IndicatorSource,
		BaseYear:                 data.BaseYear,
		GenderEquality:           data.GenderEquality,
		BaseValue:                data.BaseValue,
		SourceOfInformation:      data.SourceOfInformation,
		UnitOfMeasure:            data.UnitOfMeasure,
		IndicatorDescription:     data.IndicatorDescription,
		PlannedValue1:            data.PlannedValue1,
		AchievedValue1:           data.AchievedValue1,
		RevisedValue1:            data.RevisedValue1,
		PlannedValue2:            data.PlannedValue2,
		AchievedValue2:           data.AchievedValue2,
		RevisedValue2:            data.RevisedValue2,
		PlannedValue3:            data.PlannedValue3,
		AchievedValue3:           data.AchievedValue3,
		RevisedValue3:            data.RevisedValue3,
		CreatedAt:                data.CreatedAt,
		UpdatedAt:                data.UpdatedAt,
	}
}

func ToGoalIndicatorListResponseDTO(goalindicators []*data.GoalIndicator) []GoalIndicatorResponseDTO {
	dtoList := make([]GoalIndicatorResponseDTO, len(goalindicators))
	for i, x := range goalindicators {
		dtoList[i] = ToGoalIndicatorResponseDTO(*x)
	}
	return dtoList
}
