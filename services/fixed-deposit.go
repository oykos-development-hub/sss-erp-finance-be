package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FixedDepositServiceImpl struct {
	App        *celeritas.Celeritas
	repo       data.FixedDeposit
	items      FixedDepositItemService
	dispatches FixedDepositDispatchService
	judges     FixedDepositJudgeService
}

func NewFixedDepositServiceImpl(app *celeritas.Celeritas, repo data.FixedDeposit, items FixedDepositItemService, dispatches FixedDepositDispatchService, judges FixedDepositJudgeService) FixedDepositService {
	return &FixedDepositServiceImpl{
		App:        app,
		repo:       repo,
		items:      items,
		dispatches: dispatches,
		judges:     judges,
	}
}

func (h *FixedDepositServiceImpl) CreateFixedDeposit(input dto.FixedDepositDTO) (*dto.FixedDepositResponseDTO, error) {
	dataToInsert := input.ToFixedDeposit()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFixedDepositResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *FixedDepositServiceImpl) UpdateFixedDeposit(id int, input dto.FixedDepositDTO) (*dto.FixedDepositResponseDTO, error) {
	dataToInsert := input.ToFixedDeposit()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(tx, *dataToInsert)
		if err != nil {
			return errors.ErrInternalServer
		}
		return nil
	})
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFixedDepositResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *FixedDepositServiceImpl) DeleteFixedDeposit(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *FixedDepositServiceImpl) GetFixedDeposit(id int) (*dto.FixedDepositResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFixedDepositResponseDTO(*data)

	items, _, err := h.items.GetFixedDepositItemList(dto.FixedDepositItemFilterDTO{
		DepositID: &id,
	})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	response.Items = items

	dispatches, _, err := h.dispatches.GetFixedDepositDispatchList(
		dto.FixedDepositDispatchFilterDTO{
			DepositID: &id,
		})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	response.Dispatches = dispatches

	judges, _, err := h.judges.GetFixedDepositJudgeList(dto.FixedDepositJudgeFilterDTO{DepositID: &id})

	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	response.Judges = judges

	return &response, nil
}

func (h *FixedDepositServiceImpl) GetFixedDepositList(filter dto.FixedDepositFilterDTO) ([]dto.FixedDepositResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.JudgeID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"judge_id": *filter.JudgeID})
	}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Status != nil {
		switch *filter.Status {
		case "U radu":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_end is": nil})
		case "Zaključen":
			conditionAndExp = up.And(conditionAndExp, &up.Cond{"date_of_end is not": nil})
		}
	}

	if filter.Subject != nil && *filter.Subject != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"subject": *filter.Subject})
	}

	if filter.Type != nil && *filter.Type != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"type": *filter.Type})
	}

	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"case_number ILIKE": likeCondition},
			up.Cond{"subject ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	/*	if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}*/

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToFixedDepositListResponseDTO(data)

	return response, total, nil
}
