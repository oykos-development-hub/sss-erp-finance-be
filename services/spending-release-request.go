package services

import (
	"context"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type SpendingReleaseRequestServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.SpendingReleaseRequest
}

func NewSpendingReleaseRequestServiceImpl(app *celeritas.Celeritas, repo data.SpendingReleaseRequest) SpendingReleaseRequestService {
	return &SpendingReleaseRequestServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *SpendingReleaseRequestServiceImpl) CreateSpendingReleaseRequest(input dto.SpendingReleaseRequestDTO) (*dto.SpendingReleaseRequestResponseDTO, error) {
	dataToInsert := input.ToSpendingReleaseRequest()

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

	res := dto.ToSpendingReleaseRequestResponseDTO(*dataToInsert)

	return &res, nil
}

func (h *SpendingReleaseRequestServiceImpl) UpdateSpendingReleaseRequest(id int, input dto.SpendingReleaseRequestDTO) (*dto.SpendingReleaseRequestResponseDTO, error) {
	dataToInsert := input.ToSpendingReleaseRequest()
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

	response := dto.ToSpendingReleaseRequestResponseDTO(*dataToInsert)

	return &response, nil
}

func (h *SpendingReleaseRequestServiceImpl) DeleteSpendingReleaseRequest(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *SpendingReleaseRequestServiceImpl) GetSpendingReleaseRequest(id int) (*dto.SpendingReleaseRequestResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToSpendingReleaseRequestResponseDTO(*data)

	return &response, nil
}

func (h *SpendingReleaseRequestServiceImpl) GetSpendingReleaseRequestList(filter dto.SpendingReleaseRequestFilterDTO) ([]dto.SpendingReleaseRequestResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}

	if filter.Status != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"status": *filter.Status})
	}

	if filter.Month != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"month": *filter.Month})
	}

	if filter.Year != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"year": *filter.Year})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToSpendingReleaseRequestListResponseDTO(data)

	return response, total, nil
}

func (h *SpendingReleaseRequestServiceImpl) AcceptSSSRequest(ctx context.Context, id int, fileID int) error {

	err := h.repo.AcceptSSSRequest(ctx, id, fileID)
	if err != nil {
		return newErrors.Wrap(err, "repo spending release request accept sss request")
	}

	return nil
}
