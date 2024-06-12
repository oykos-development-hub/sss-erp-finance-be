package services

import (
	"context"

	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"
)

type ProcedureCostPaymentServiceImpl struct {
	App                             *celeritas.Celeritas
	repo                            data.ProcedureCostPayment
	procedurecostSharedLogicService ProcedureCostSharedLogicService
}

// NewProcedureCostPaymentServiceImpl is a factory function that returns a new instance of ProcedureCostPaymentServiceImpl
func NewProcedureCostPaymentServiceImpl(app *celeritas.Celeritas, repo data.ProcedureCostPayment, procedurecostSharedLogicService ProcedureCostSharedLogicService) ProcedureCostPaymentService {
	return &ProcedureCostPaymentServiceImpl{
		App:                             app,
		repo:                            repo,
		procedurecostSharedLogicService: procedurecostSharedLogicService,
	}
}

// CreateProcedureCostPayment creates a new procedurecost payment
func (h *ProcedureCostPaymentServiceImpl) CreateProcedureCostPayment(ctx context.Context, input dto.ProcedureCostPaymentDTO) (*dto.ProcedureCostPaymentResponseDTO, error) {
	procedurecostPayment := input.ToProcedureCostPayment()
	procedurecostPayment.Status = data.PaidProcedureCostPeymentStatus

	id, err := h.repo.Insert(ctx, *procedurecostPayment)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	procedurecostPayment, err = procedurecostPayment.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToProcedureCostPaymentResponseDTO(*procedurecostPayment)

	_, _, err = h.procedurecostSharedLogicService.CalculateProcedureCostDetailsAndUpdateStatus(ctx, procedurecostPayment.ProcedureCostID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, err
	}

	return &res, nil
}

// GetProcedureCostPayment returns a procedurecost payment by its id
func (h *ProcedureCostPaymentServiceImpl) DeleteProcedureCostPayment(ctx context.Context, id int) error {
	procedurecostPayment, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrNotFound
	}

	err = h.repo.Delete(ctx, id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	_, _, err = h.procedurecostSharedLogicService.CalculateProcedureCostDetailsAndUpdateStatus(ctx, procedurecostPayment.ProcedureCostID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

// UpdateProcedureCostPayment updates a procedurecost payment by its id
func (h *ProcedureCostPaymentServiceImpl) UpdateProcedureCostPayment(ctx context.Context, id int, input dto.ProcedureCostPaymentDTO) (*dto.ProcedureCostPaymentResponseDTO, error) {
	data := input.ToProcedureCostPayment()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToProcedureCostPaymentResponseDTO(*data)

	_, _, err = h.procedurecostSharedLogicService.CalculateProcedureCostDetailsAndUpdateStatus(ctx, data.ProcedureCostID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return &response, nil
}

// GetProcedureCostPaymentList returns a list of procedurecost payments by procedurecost id
func (h *ProcedureCostPaymentServiceImpl) GetProcedureCostPaymentList(input dto.ProcedureCostPaymentFilterDTO) ([]dto.ProcedureCostPaymentResponseDTO, *uint64, error) {

	procedurecostPayments, total, err := h.getProcedureCostPaymentsByProcedureCostID(input.ProcedureCostID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	if len(procedurecostPayments) == 0 {
		return nil, nil, errors.ErrNotFound
	}
	response := dto.ToProcedureCostPaymentListResponseDTO(procedurecostPayments)

	return response, total, nil
}

func (h *ProcedureCostPaymentServiceImpl) getProcedureCostPaymentsByProcedureCostID(procedurecostID int) ([]*data.ProcedureCostPayment, *uint64, error) {
	cond := db.Cond{"procedure_cost_id": procedurecostID}

	procedurecostPayments, total, err := h.repo.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}

	return procedurecostPayments, total, nil
}

// GetProcedureCostPayment returns a procedurecost payment by its id
func (h *ProcedureCostPaymentServiceImpl) GetProcedureCostPayment(id int) (*dto.ProcedureCostPaymentResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToProcedureCostPaymentResponseDTO(*data)

	return &response, nil
}
