package services

import (
	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
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
func (h *FeePaymentServiceImpl) CreateFeePayment(input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error) {
	feePayment := input.ToFeePayment()

	id, err := h.repo.Insert(*feePayment)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	feePayment, err = feePayment.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFeePaymentResponseDTO(*feePayment)

	_, _, err = h.feeSharedLogicService.CalculateFeeDetailsAndUpdateStatus(feePayment.FeeID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	return &res, nil
}

// GetFeePayment returns a fee payment by its id
func (h *FeePaymentServiceImpl) DeleteFeePayment(id int) error {
	feePayment, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrNotFound
	}

	err = h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	_, _, err = h.feeSharedLogicService.CalculateFeeDetailsAndUpdateStatus(feePayment.FeeID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

// UpdateFeePayment updates a fee payment by its id
func (h *FeePaymentServiceImpl) UpdateFeePayment(id int, input dto.FeePaymentDTO) (*dto.FeePaymentResponseDTO, error) {
	data := input.ToFeePayment()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFeePaymentResponseDTO(*data)

	_, _, err = h.feeSharedLogicService.CalculateFeeDetailsAndUpdateStatus(data.FeeID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return &response, nil
}

// GetFeePaymentList returns a list of fee payments by fee id
func (h *FeePaymentServiceImpl) GetFeePaymentList(input dto.FeePaymentFilterDTO) ([]dto.FeePaymentResponseDTO, *uint64, error) {

	feePayments, total, err := h.getFeePaymentsByFeeID(input.FeeID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	if len(feePayments) == 0 {
		return nil, nil, errors.ErrNotFound
	}
	response := dto.ToFeePaymentListResponseDTO(feePayments)

	return response, total, nil
}

func (h *FeePaymentServiceImpl) getFeePaymentsByFeeID(feeID int) ([]*data.FeePayment, *uint64, error) {
	cond := db.Cond{"fee_id": feeID}

	feePayments, total, err := h.repo.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	return feePayments, total, nil
}

// GetFeePayment returns a fee payment by its id
func (h *FeePaymentServiceImpl) GetFeePayment(id int) (*dto.FeePaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFeePaymentResponseDTO(*data)

	return &response, nil
}
