package services

import (
	"context"

	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
)

type FlatRatePaymentServiceImpl struct {
	App                        *celeritas.Celeritas
	repo                       data.FlatRatePayment
	FlatRateSharedLogicService FlatRateSharedLogicService
}

// NewFlatRatePaymentServiceImpl is a factory function that returns a new instance of FlatRatePaymentServiceImpl
func NewFlatRatePaymentServiceImpl(app *celeritas.Celeritas, repo data.FlatRatePayment, FlatRateSharedLogicService FlatRateSharedLogicService) FlatRatePaymentService {
	return &FlatRatePaymentServiceImpl{
		App:                        app,
		repo:                       repo,
		FlatRateSharedLogicService: FlatRateSharedLogicService,
	}
}

// CreateFlatRatePayment creates a new FlatRate payment
func (h *FlatRatePaymentServiceImpl) CreateFlatRatePayment(ctx context.Context, input dto.FlatRatePaymentDTO) (*dto.FlatRatePaymentResponseDTO, error) {
	flatRatePayment := input.ToFlatRatePayment()
	flatRatePayment.Status = data.PaidFlatRatePeymentStatus

	id, err := h.repo.Insert(ctx, *flatRatePayment)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	flatRatePayment, err = flatRatePayment.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFlatRatePaymentResponseDTO(*flatRatePayment)

	_, _, err = h.FlatRateSharedLogicService.CalculateFlatRateDetailsAndUpdateStatus(ctx, flatRatePayment.FlatRateID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	return &res, nil
}

// GetFlatRatePayment returns a FlatRate payment by its id
func (h *FlatRatePaymentServiceImpl) DeleteFlatRatePayment(ctx context.Context, id int) error {
	FlatRatePayment, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrNotFound
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	_, _, err = h.FlatRateSharedLogicService.CalculateFlatRateDetailsAndUpdateStatus(ctx, FlatRatePayment.FlatRateID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

// UpdateFlatRatePayment updates a FlatRate payment by its id
func (h *FlatRatePaymentServiceImpl) UpdateFlatRatePayment(ctx context.Context, id int, input dto.FlatRatePaymentDTO) (*dto.FlatRatePaymentResponseDTO, error) {
	data := input.ToFlatRatePayment()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFlatRatePaymentResponseDTO(*data)

	_, _, err = h.FlatRateSharedLogicService.CalculateFlatRateDetailsAndUpdateStatus(ctx, data.FlatRateID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return &response, nil
}

// GetFlatRatePaymentList returns a list of FlatRate payments by FlatRate id
func (h *FlatRatePaymentServiceImpl) GetFlatRatePaymentList(input dto.FlatRatePaymentFilterDTO) ([]dto.FlatRatePaymentResponseDTO, *uint64, error) {

	FlatRatePayments, total, err := h.getFlatRatePaymentsByFlatRateID(input.FlatRateID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	if len(FlatRatePayments) == 0 {
		return nil, nil, errors.ErrNotFound
	}
	response := dto.ToFlatRatePaymentListResponseDTO(FlatRatePayments)

	return response, total, nil
}

func (h *FlatRatePaymentServiceImpl) getFlatRatePaymentsByFlatRateID(FlatRateID int) ([]*data.FlatRatePayment, *uint64, error) {
	cond := db.Cond{"flat_rate_id": FlatRateID}

	FlatRatePayments, total, err := h.repo.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	return FlatRatePayments, total, nil
}

// GetFlatRatePayment returns a FlatRate payment by its id
func (h *FlatRatePaymentServiceImpl) GetFlatRatePayment(id int) (*dto.FlatRatePaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFlatRatePaymentResponseDTO(*data)

	return &response, nil
}
