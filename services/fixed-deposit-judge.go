package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type FixedDepositJudgeServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.FixedDepositJudge
}

func NewFixedDepositJudgeServiceImpl(app *celeritas.Celeritas, repo data.FixedDepositJudge) FixedDepositJudgeService {
	return &FixedDepositJudgeServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FixedDepositJudgeServiceImpl) CreateFixedDepositJudge(input dto.FixedDepositJudgeDTO) (*dto.FixedDepositJudgeResponseDTO, error) {
	dataToInsert := input.ToFixedDepositJudge()

	var id int
	err := data.Upper.Tx(func(tx up.Session) error {
		var err error
		id, err = h.repo.Insert(tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo fixed deposit judge insert")
		}

		return nil
	})

	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit judge get")
	}

	res := dto.ToFixedDepositJudgeResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *FixedDepositJudgeServiceImpl) UpdateFixedDepositJudge(id int, input dto.FixedDepositJudgeDTO) (*dto.FixedDepositJudgeResponseDTO, error) {
	dataToInsert := input.ToFixedDepositJudge()
	dataToInsert.ID = id

	err := data.Upper.Tx(func(tx up.Session) error {
		err := h.repo.Update(tx, *dataToInsert)
		if err != nil {
			return newErrors.Wrap(err, "repo fixed deposit judge update")
		}
		return nil
	})
	if err != nil {
		return nil, newErrors.Wrap(err, "upper tx")
	}

	dataToInsert, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit judge get")
	}

	response := dto.ToFixedDepositJudgeResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *FixedDepositJudgeServiceImpl) DeleteFixedDepositJudge(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		return newErrors.Wrap(err, "repo fixed deposit judge delete")
	}

	return nil
}

func (h *FixedDepositJudgeServiceImpl) GetFixedDepositJudge(id int) (*dto.FixedDepositJudgeResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fixed deposit judge get")
	}

	response := dto.ToFixedDepositJudgeResponseDTO(*data)

	return &response, nil
}

func (h *FixedDepositJudgeServiceImpl) GetFixedDepositJudgeList(filter dto.FixedDepositJudgeFilterDTO) ([]dto.FixedDepositJudgeResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.DepositID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"deposit_id": *filter.DepositID})
	}

	if filter.WillID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"will_id": *filter.WillID})
	}

	if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo fixed deposit judge get all")
	}
	response := dto.ToFixedDepositJudgeListResponseDTO(data)

	return response, total, nil
}
