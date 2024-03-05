package services

import (
	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
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
func (h *FinePaymentServiceImpl) CreateFinePayment(input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error) {
	finePayment := input.ToFinePayment()

	id, err := h.repo.Insert(*finePayment)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	finePayment, err = finePayment.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToFinePaymentResponseDTO(*finePayment)

	_, _, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(finePayment.FineID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	return &res, nil
}

// GetFinePayment returns a fine payment by its id
func (h *FinePaymentServiceImpl) DeleteFinePayment(id int) error {
	finePayment, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrNotFound
	}

	err = h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	_, _, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(finePayment.FineID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

// UpdateFinePayment updates a fine payment by its id
func (h *FinePaymentServiceImpl) UpdateFinePayment(id int, input dto.FinePaymentDTO) (*dto.FinePaymentResponseDTO, error) {
	data := input.ToFinePayment()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToFinePaymentResponseDTO(*data)

	_, _, err = h.fineSharedLogicService.CalculateFineDetailsAndUpdateStatus(data.FineID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return &response, nil
}

// GetFinePaymentList returns a list of fine payments by fine id
func (h *FinePaymentServiceImpl) GetFinePaymentList(input dto.FinePaymentFilterDTO) ([]dto.FinePaymentResponseDTO, *uint64, error) {

	finePayments, total, err := h.getFinePaymentsByFineID(input.FineID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	if len(finePayments) == 0 {
		return nil, nil, errors.ErrNotFound
	}
	response := dto.ToFinePaymentListResponseDTO(finePayments)

	return response, total, nil
}

func (h *FinePaymentServiceImpl) getFinePaymentsByFineID(fineID int) ([]*data.FinePayment, *uint64, error) {
	cond := db.Cond{"fine_id": fineID}

	finePayments, total, err := h.repo.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	return finePayments, total, nil
}

// GetFinePayment returns a fine payment by its id
func (h *FinePaymentServiceImpl) GetFinePayment(id int) (*dto.FinePaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToFinePaymentResponseDTO(*data)

	return &response, nil
}
