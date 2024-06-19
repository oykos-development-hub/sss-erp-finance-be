package services

import (
	"context"

	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"
)

type PropBenConfPaymentServiceImpl struct {
	App                           *celeritas.Celeritas
	repo                          data.PropBenConfPayment
	propbenconfSharedLogicService PropBenConfSharedLogicService
}

// NewPropBenConfPaymentServiceImpl is a factory function that returns a new instance of PropBenConfPaymentServiceImpl
func NewPropBenConfPaymentServiceImpl(app *celeritas.Celeritas, repo data.PropBenConfPayment, propbenconfSharedLogicService PropBenConfSharedLogicService) PropBenConfPaymentService {
	return &PropBenConfPaymentServiceImpl{
		App:                           app,
		repo:                          repo,
		propbenconfSharedLogicService: propbenconfSharedLogicService,
	}
}

// CreatePropBenConfPayment creates a new propbenconf payment
func (h *PropBenConfPaymentServiceImpl) CreatePropBenConfPayment(ctx context.Context, input dto.PropBenConfPaymentDTO) (*dto.PropBenConfPaymentResponseDTO, error) {
	propbenconfPayment := input.ToPropBenConfPayment()
	propbenconfPayment.Status = data.PaidPropBenConfPeymentStatus

	id, err := h.repo.Insert(ctx, *propbenconfPayment)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf payment insert")
	}

	propbenconfPayment, err = propbenconfPayment.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf payment get")
	}

	res := dto.ToPropBenConfPaymentResponseDTO(*propbenconfPayment)

	_, _, err = h.propbenconfSharedLogicService.CalculatePropBenConfDetailsAndUpdateStatus(ctx, propbenconfPayment.PropBenConfID)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf shared logic service calcualte prop ben conf details and update status")
	}

	return &res, nil
}

// GetPropBenConfPayment returns a propbenconf payment by its id
func (h *PropBenConfPaymentServiceImpl) DeletePropBenConfPayment(ctx context.Context, id int) error {
	propbenconfPayment, err := h.repo.Get(id)
	if err != nil {
		return newErrors.Wrap(err, "repo prop ben conf payment get")
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo prop ben conf payment delete")
	}

	_, _, err = h.propbenconfSharedLogicService.CalculatePropBenConfDetailsAndUpdateStatus(ctx, propbenconfPayment.PropBenConfID)
	if err != nil {
		return newErrors.Wrap(err, "repo prop ben conf shared logic service calculate prop ben conf details and update status")
	}

	return nil
}

// UpdatePropBenConfPayment updates a propbenconf payment by its id
func (h *PropBenConfPaymentServiceImpl) UpdatePropBenConfPayment(ctx context.Context, id int, input dto.PropBenConfPaymentDTO) (*dto.PropBenConfPaymentResponseDTO, error) {
	data := input.ToPropBenConfPayment()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf payment update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf payment get")
	}

	response := dto.ToPropBenConfPaymentResponseDTO(*data)

	_, _, err = h.propbenconfSharedLogicService.CalculatePropBenConfDetailsAndUpdateStatus(ctx, data.PropBenConfID)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf shared logic service calculate prop ben conf details and update status")
	}

	return &response, nil
}

// GetPropBenConfPaymentList returns a list of propbenconf payments by propbenconf id
func (h *PropBenConfPaymentServiceImpl) GetPropBenConfPaymentList(input dto.PropBenConfPaymentFilterDTO) ([]dto.PropBenConfPaymentResponseDTO, *uint64, error) {

	propbenconfPayments, total, err := h.getPropBenConfPaymentsByPropBenConfID(input.PropBenConfID)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "get prop ben conf payments by pro ben conf id")
	}

	if len(propbenconfPayments) == 0 {
		return nil, nil, newErrors.Wrap(errors.ErrNotFound, "get prop ben conf payments by pro ben conf id")
	}
	response := dto.ToPropBenConfPaymentListResponseDTO(propbenconfPayments)

	return response, total, nil
}

func (h *PropBenConfPaymentServiceImpl) getPropBenConfPaymentsByPropBenConfID(propbenconfID int) ([]*data.PropBenConfPayment, *uint64, error) {
	cond := db.Cond{"property_benefits_confiscation_id": propbenconfID}

	propbenconfPayments, total, err := h.repo.GetAll(&cond)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "rep prop ben conf payments get all")
	}

	return propbenconfPayments, total, nil
}

// GetPropBenConfPayment returns a propbenconf payment by its id
func (h *PropBenConfPaymentServiceImpl) GetPropBenConfPayment(id int) (*dto.PropBenConfPaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf payments get")
	}

	response := dto.ToPropBenConfPaymentResponseDTO(*data)

	return &response, nil
}
