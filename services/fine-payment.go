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

type FinePaymentServiceImpl struct {
	App                    *celeritas.Celeritas
	repo                   data.FinePayment
	fineSharedLogicService FineSharedLogicService
}

// NewFinePaymentServiceImpl is a factory function that returns a new instance of FinePaymentServiceImpl
func NewFinePaymentServiceImpl(app *celeritas.Celeritas, repo data.FinePayment, fineSharedLogicService FineSharedLogicService) FinePaymentService {
	return &FinePaymentServiceImpl{
		App:                    app,
		repo:                   repo,
		fineSharedLogicService: fineSharedLogicService,
	}
}

// CreateFinePayment creates a new fine payment
func (h *FinePaymentServiceImpl) CreateFinePayment(ctx context.Context, input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error) {
	finePayment := input.ToFinePayment()
	finePayment.Status = data.PaidFinePeymentStatus

	id, err := h.repo.Insert(ctx, *finePayment)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine payment insert")
	}

	finePayment, err = finePayment.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine payment get")
	}

	res := dto.ToFinePaymentResponseDTO(*finePayment)

	_, _, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(ctx, finePayment.FineID)
	if err != nil {
		return nil, newErrors.Wrap(err, "fine shared logic calculate fine details and update status")
	}

	return &res, nil
}

// GetFinePayment returns a fine payment by its id
func (h *FinePaymentServiceImpl) DeleteFinePayment(ctx context.Context, id int) error {
	finePayment, err := h.repo.Get(id)
	if err != nil {
		return newErrors.Wrap(err, "repo fine payment get")
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo fine payment delete")
	}

	_, _, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(ctx, finePayment.FineID)
	if err != nil {
		return newErrors.Wrap(err, "fine shared logic calculate fine details and update status")
	}

	return nil
}

// UpdateFinePayment updates a fine payment by its id
func (h *FinePaymentServiceImpl) UpdateFinePayment(ctx context.Context, id int, input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error) {
	data := input.ToFinePayment()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine payment update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine payment get")
	}

	response := dto.ToFinePaymentResponseDTO(*data)

	_, _, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(ctx, data.FineID)
	if err != nil {
		return nil, newErrors.Wrap(err, "fine shared logic calculate fine details and update status")
	}

	return &response, nil
}

// GetFinePaymentList returns a list of fine payments by fine id
func (h *FinePaymentServiceImpl) GetFinePaymentList(input dto.FinePaymentFilterDTO) ([]dto.FinePaymentResponseDTO, *uint64, error) {

	finePayments, total, err := h.getFinePaymentsByFineID(input.FineID)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "get fine payment by fine id")
	}

	if len(finePayments) == 0 {
		return nil, nil, newErrors.Wrap(errors.ErrNotFound, "get fine payment by fine id")
	}
	response := dto.ToFinePaymentListResponseDTO(finePayments)

	return response, total, nil
}

func (h *FinePaymentServiceImpl) getFinePaymentsByFineID(fineID int) ([]*data.FinePayment, *uint64, error) {
	cond := db.Cond{"fine_id": fineID}

	finePayments, total, err := h.repo.GetAll(&cond)
	if err != nil {
		return nil, nil, newErrors.Wrap(err, "repo fine payment get all")
	}

	return finePayments, total, nil
}

// GetFinePayment returns a fine payment by its id
func (h *FinePaymentServiceImpl) GetFinePayment(id int) (*dto.FinePaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fine payment get")
	}

	response := dto.ToFinePaymentResponseDTO(*data)

	return &response, nil
}
