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

type FeePaymentServiceImpl struct {
	App                   *celeritas.Celeritas
	repo                  data.FeePayment
	feeSharedLogicService FeeSharedLogicService
}

// NewFeePaymentServiceImpl is a factory function that returns a new instance of FeePaymentServiceImpl
func NewFeePaymentServiceImpl(app *celeritas.Celeritas, repo data.FeePayment, feeSharedLogicService FeeSharedLogicService) FeePaymentService {
	return &FeePaymentServiceImpl{
		App:                   app,
		repo:                  repo,
		feeSharedLogicService: feeSharedLogicService,
	}
}

// CreateFeePayment creates a new fee payment
func (h *FeePaymentServiceImpl) CreateFeePayment(ctx context.Context, input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error) {
	feePayment := input.ToFeePayment()
	feePayment.Status = data.PaidFeePeymentStatus

	id, err := h.repo.Insert(ctx, *feePayment)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fee payment insert")
	}

	feePayment, err = feePayment.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fee payment get")
	}

	res := dto.ToFeePaymentResponseDTO(*feePayment)

	_, _, err = h.feeSharedLogicService.CalculateFeeDetailsAndUpdateStatus(ctx, feePayment.FeeID)
	if err != nil {
		return nil, newErrors.Wrap(err, "fee shared logic service calculate fee details and update status")
	}

	return &res, nil
}

// GetFeePayment returns a fee payment by its id
func (h *FeePaymentServiceImpl) DeleteFeePayment(ctx context.Context, id int) error {
	feePayment, err := h.repo.Get(id)
	if err != nil {
		return newErrors.Wrap(err, "repo fee payment get")
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo fee payment delete")
	}

	_, _, err = h.feeSharedLogicService.CalculateFeeDetailsAndUpdateStatus(ctx, feePayment.FeeID)
	if err != nil {
		return newErrors.Wrap(err, "fee shared logic service calculate fee details and update status")
	}

	return nil
}

// UpdateFeePayment updates a fee payment by its id
func (h *FeePaymentServiceImpl) UpdateFeePayment(ctx context.Context, id int, input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error) {
	data := input.ToFeePayment()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fee payment update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fee payment get")
	}

	response := dto.ToFeePaymentResponseDTO(*data)

	_, _, err = h.feeSharedLogicService.CalculateFeeDetailsAndUpdateStatus(ctx, data.FeeID)
	if err != nil {
		return nil, newErrors.Wrap(err, "fee shared logic service calculate fee details and update status")
	}

	return &response, nil
}

// GetFeePaymentList returns a list of fee payments by fee id
func (h *FeePaymentServiceImpl) GetFeePaymentList(input dto.FeePaymentFilterDTO) ([]dto.FeePaymentResponseDTO, *uint64, error) {

	feePayments, total, err := h.getFeePaymentsByFeeID(input.FeeID)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "get fee payments by fee id")
	}

	if len(feePayments) == 0 {
		return nil, nil, newErrors.Wrap(errors.ErrNotFound, "get fee payments by fee id")
	}
	response := dto.ToFeePaymentListResponseDTO(feePayments)

	return response, total, nil
}

func (h *FeePaymentServiceImpl) getFeePaymentsByFeeID(feeID int) ([]*data.FeePayment, *uint64, error) {
	cond := db.Cond{"fee_id": feeID}

	feePayments, total, err := h.repo.GetAll(&cond)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo fee payments get all")
	}

	return feePayments, total, nil
}

// GetFeePayment returns a fee payment by its id
func (h *FeePaymentServiceImpl) GetFeePayment(id int) (*dto.FeePaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fee payments get")
	}

	response := dto.ToFeePaymentResponseDTO(*data)

	return &response, nil
}
