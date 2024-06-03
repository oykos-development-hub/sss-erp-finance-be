package services

import (
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
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

func (h *ProcedureCostSharedLogicServiceImpl) CalculateProcedureCostDetailsAndUpdateStatus(procedurecostId int) (*dto.ProcedureCostDetailsDTO, data.ProcedureCostStatus, error) {
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
				details.CourtCostsPaid = details.CourtCostsPaid.Add(payment.Amount)
			} else {
				details.AllPaymentAmount = details.AllPaymentAmount.Add(payment.Amount)
			}
		}
	}

	details.LeftToPayAmount = procedurecost.Amount.Sub(details.AllPaymentAmount)
	if procedurecost.CourtCosts != nil {
		details.CourtCostsLeftToPayAmount = procedurecost.CourtCosts.Sub(details.CourtCostsPaid)
	}

	details.AmountGracePeriodDueDate = procedurecost.DecisionDate.AddDate(0, 0, data.ProcedureCostGracePeriod)
	twoThirds := decimal.NewFromFloat(2.0).Div(decimal.NewFromFloat(3.0))
	details.AmountGracePeriod = procedurecost.Amount.Mul(twoThirds).Ceil()

	if time.Until(details.AmountGracePeriodDueDate) > 0 {
		details.AmountGracePeriodAvailable = true
		details.LeftToPayAmount = details.AmountGracePeriod.Sub(details.AllPaymentAmount)
	}

	var newStatus data.ProcedureCostStatus

	tolerance := decimal.NewFromFloat(0.00001)

	zero := decimal.NewFromInt(0)
	leftToPayAmount := decimal.Max(zero, details.LeftToPayAmount)
	courtCostsLeftToPayAmount := decimal.Max(zero, details.CourtCostsLeftToPayAmount)

	// Provera uslova i postavljanje statusa
	if leftToPayAmount.Abs().Cmp(tolerance) < 0 && courtCostsLeftToPayAmount.Abs().Cmp(tolerance) < 0 {
		newStatus = data.PaidProcedureCostStatus
	} else if (leftToPayAmount.Cmp(zero) > 0 || courtCostsLeftToPayAmount.Cmp(zero) > 0) &&
		(details.CourtCostsPaid.Cmp(zero) > 0 || details.AllPaymentAmount.Cmp(zero) > 0) {
		newStatus = data.PartProcedureCostStatus
	} else {
		newStatus = data.UnpaidProcedureCostStatus
	}

	if newStatus != procedurecost.Status {
		procedurecost.Status = newStatus
		err = h.repoProcedureCost.Update(*procedurecost)
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
