package services

import (
	"context"
	"math"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
)

type ProcedureCostSharedLogicServiceImpl struct {
	App                      *celeritas.Celeritas
	repoProcedureCost        data.ProcedureCost
	repoProcedureCostPayment data.ProcedureCostPayment
}

// ProcedureCostSharedLogicServiceImpl  creates a new instance of ProcedureCostService
func NewProcedureCostSharedLogicServiceImpl(app *celeritas.Celeritas, repoProcedureCost data.ProcedureCost, repoProcedureCostPayment data.ProcedureCostPayment) ProcedureCostSharedLogicService {
	return &ProcedureCostSharedLogicServiceImpl{
		App:                      app,
		repoProcedureCost:        repoProcedureCost,
		repoProcedureCostPayment: repoProcedureCostPayment,
	}
}

func (h *ProcedureCostSharedLogicServiceImpl) CalculateProcedureCostDetailsAndUpdateStatus(ctx context.Context, procedurecostId int) (*dto.ProcedureCostDetailsDTO, data.ProcedureCostStatus, error) {
	procedurecost, err := h.repoProcedureCost.Get(procedurecostId)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}
	payments, err := h.getProcedureCostPaymentsByProcedureCostID(procedurecost.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}

	details := &dto.ProcedureCostDetailsDTO{}

	// count all payments and court costs payments
	for _, payment := range payments {
		if data.ProcedureCostPaymentStatus(payment.Status) == data.PaidProcedureCostPeymentStatus {
			if data.ProcedureCostPaymentMethod(payment.PaymentMethod) == data.CourtCostsProcedureCostPeymentMethod {
				details.CourtCostsPaid += payment.Amount
			} else {
				details.AllPaymentAmount += payment.Amount
			}
		}
	}

	details.LeftToPayAmount = procedurecost.Amount - details.AllPaymentAmount
	if procedurecost.CourtCosts != nil {
		details.CourtCostsLeftToPayAmount = *procedurecost.CourtCosts - details.CourtCostsPaid
	}

	details.AmountGracePeriodDueDate = procedurecost.DecisionDate.AddDate(0, 0, data.ProcedureCostGracePeriod)
	details.AmountGracePeriod = math.Ceil(float64(procedurecost.Amount) * 2 / 3)

	if time.Until(details.AmountGracePeriodDueDate) > 0 {
		details.AmountGracePeriodAvailable = true
		details.LeftToPayAmount = details.AmountGracePeriod - details.AllPaymentAmount
	}

	var newStatus data.ProcedureCostStatus
	const tolerance = 0.00001

	leftToPayAmount := math.Max(0, details.LeftToPayAmount)
	courtCostsLeftToPayAmount := math.Max(0, details.CourtCostsLeftToPayAmount)

	if math.Abs(leftToPayAmount-0) < tolerance && math.Abs(courtCostsLeftToPayAmount-0) < tolerance {
		newStatus = data.PaidProcedureCostStatus
	} else if (leftToPayAmount > 0 || courtCostsLeftToPayAmount > 0) && (details.CourtCostsPaid > 0 || details.AllPaymentAmount > 0) {
		newStatus = data.PartProcedureCostStatus
	} else {
		newStatus = data.UnpaidProcedureCostStatus
	}

	if newStatus != procedurecost.Status {
		procedurecost.Status = newStatus
		err = h.repoProcedureCost.Update(ctx, *procedurecost)
		if err != nil {
			return nil, 0, errors.ErrInternalServer
		}
	}

	return details, newStatus, nil
}

func (h *ProcedureCostSharedLogicServiceImpl) getProcedureCostPaymentsByProcedureCostID(procedurecostID int) ([]*data.ProcedureCostPayment, error) {
	cond := db.Cond{"procedure_cost_id": procedurecostID}

	procedurecostPayments, _, err := h.repoProcedureCostPayment.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return procedurecostPayments, nil
}
