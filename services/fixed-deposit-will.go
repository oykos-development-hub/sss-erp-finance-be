package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FixedDepositWillServiceImpl struct {
	App        *celeritas.Celeritas
	repo       data.FixedDepositWill
	judges     FixedDepositJudgeService
	dispatches FixedDepositWillDispatchService
}

func NewFixedDepositWillServiceImpl(app *celeritas.Celeritas, repo data.FixedDepositWill, judges FixedDepositJudgeService, dispatches FixedDepositWillDispatchService) FixedDepositWillService {
	return &FixedDepositWillServiceImpl{
		App:        app,
		repo:       repo,
		dispatches: dispatches,
		judges:     judges,
	}
}

func (h *FixedDepositWillServiceImpl) CreateFixedDepositWill(ctx context.Context, input dto.FixedDepositWillDTO) (*dto.FixedDepositWillResponseDTO, error) {
	dataToInsert := input.ToFixedDepositWill()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(ctx, tx, *dataToInsert)
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

	res := dto.ToFixedDepositWillResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *FixedDepositWillServiceImpl) UpdateFixedDepositWill(ctx context.Context, id int, input dto.FixedDepositWillDTO) (*dto.FixedDepositWillResponseDTO, error) {
	dataToInsert := input.ToFixedDepositWill()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(ctx, tx, *dataToInsert)
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

	response := dto.ToFixedDepositWillResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *FixedDepositWillServiceImpl) DeleteFixedDepositWill(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *FixedDepositWillServiceImpl) GetFixedDepositWill(id int) (*dto.FixedDepositWillResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFixedDepositWillResponseDTO(*data)

	dispatches, _, err := h.dispatches.GetFixedDepositWillDispatchList(
		dto.FixedDepositWillDispatchFilterDTO{
			WillID: &id,
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

func (h *FixedDepositWillServiceImpl) GetFixedDepositWillList(filter dto.FixedDepositWillFilterDTO) ([]dto.FixedDepositWillResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Status != nil && *filter.Status != "" {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"subject ILIKE": likeCondition},
			up.Cond{"case_number_si ILIKE": likeCondition},
			up.Cond{"case_number_rs ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	//	if filter.SortByTitle != nil {
	//		if *filter.SortByTitle == "asc" {
	//			orders = append(orders, "-title")
	//		} else {
	//			orders = append(orders, "title")
	//		}
	//	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToFixedDepositWillListResponseDTO(data)

	return response, total, nil
}
