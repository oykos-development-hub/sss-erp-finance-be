package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type BudgetRequestServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.BudgetRequest
}

func NewBudgetRequestServiceImpl(app *celeritas.Celeritas, repo data.BudgetRequest) BudgetRequestService {
	return &BudgetRequestServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *BudgetRequestServiceImpl) CreateBudgetRequest(input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error) {
	data := input.ToBudgetRequest()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToBudgetRequestResponseDTO(*data)

	return &res, nil
}

func (h *BudgetRequestServiceImpl) UpdateBudgetRequest(id int, input dto.BudgetRequestDTO) (*dto.BudgetRequestResponseDTO, error) {
	data := input.ToBudgetRequest()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToBudgetRequestResponseDTO(*data)

	return &response, nil
}

func (h *BudgetRequestServiceImpl) DeleteBudgetRequest(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *BudgetRequestServiceImpl) GetBudgetRequest(id int) (*dto.BudgetRequestResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToBudgetRequestResponseDTO(*data)

	return &response, nil
}

func (h *BudgetRequestServiceImpl) GetBudgetRequestList(filter dto.BudgetRequestFilterDTO) ([]dto.BudgetRequestResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	conditionAndExp = up.And(conditionAndExp, &up.Cond{"budget_id": filter.BudgetID})
	if filter.OrganizationUnitID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"organization_unit_id": *filter.OrganizationUnitID})
	}
	if filter.RequestType != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"request_type": *filter.RequestType})
	}
	if filter.RequestTypes != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"request_type": up.In(*filter.RequestTypes)})
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToBudgetRequestListResponseDTO(data)

	return response, total, nil
}
