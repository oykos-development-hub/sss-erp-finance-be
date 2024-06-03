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

type FlatRateSharedLogicServiceImpl struct {
	App                 *celeritas.Celeritas
	repoFlatRate        data.FlatRate
	repoFlatRatePayment data.FlatRatePayment
}

// FlatRateSharedLogicServiceImpl  creates a new instance of FlatRateService
func NewFlatRateSharedLogicServiceImpl(app *celeritas.Celeritas, repoFlatRate data.FlatRate, repoFlatRatePayment data.FlatRatePayment) FlatRateSharedLogicService {
	return &FlatRateSharedLogicServiceImpl{
		App:                 app,
		repoFlatRate:        repoFlatRate,
		repoFlatRatePayment: repoFlatRatePayment,
	}
}

func (h *FlatRateSharedLogicServiceImpl) CalculateFlatRateDetailsAndUpdateStatus(flatrateId int) (*dto.FlatRateDetailsDTO, data.FlatRateStatus, error) {
	flatrate, err := h.repoFlatRate.Get(flatrateId)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}
	payments, err := h.getFlatRatePaymentsByFlatRateID(flatrate.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}

	details := &dto.FlatRateDetailsDTO{}

	// count all payments and court costs payments
	for _, payment := range payments {
		if data.FlatRatePaymentStatus(payment.Status) == data.PaidFlatRatePeymentStatus {
			if data.FlatRatePaymentMethod(payment.PaymentMethod) == data.CourtCostsFlatRatePeymentMethod {
				details.CourtCostsPaid = details.CourtCostsPaid.Add(payment.Amount)
			} else {
				details.AllPaymentAmount = details.AllPaymentAmount.Add(payment.Amount)
			}
		}
	}

	// calculate the rest of the flatrates
	details.LeftToPayAmount = flatrate.Amount.Sub(details.AllPaymentAmount)
	if flatrate.CourtCosts != nil {
		details.CourtCostsLeftToPayAmount = flatrate.CourtCosts.Sub(details.CourtCostsPaid)
	}

	details.AmountGracePeriodDueDate = flatrate.DecisionDate.AddDate(0, 0, data.FlatRateGracePeriod)
	twoThirds := decimal.NewFromFloat(2.0).Div(decimal.NewFromFloat(3.0))
	details.AmountGracePeriod = flatrate.Amount.Mul(twoThirds).Ceil()

	if time.Until(details.AmountGracePeriodDueDate) > 0 {
		details.AmountGracePeriodAvailable = true
		details.LeftToPayAmount = details.AmountGracePeriod.Sub(details.AllPaymentAmount)
	}

	var newStatus data.FlatRateStatus
	tolerance := decimal.NewFromFloat(0.00001)

	zero := decimal.NewFromInt(0)
	flatrateLeftToPayAmount := decimal.Max(zero, details.LeftToPayAmount)
	flatrateCourtCostsLeftToPayAmount := decimal.Max(zero, details.CourtCostsLeftToPayAmount)

	if flatrateLeftToPayAmount.Abs().Cmp(tolerance) < 0 && flatrateCourtCostsLeftToPayAmount.Abs().Cmp(tolerance) < 0 {
		newStatus = data.PaidFlatRateStatus
	} else if (flatrateLeftToPayAmount.Cmp(zero) > 0 || flatrateCourtCostsLeftToPayAmount.Cmp(zero) > 0) &&
		(details.AllPaymentAmount.Cmp(zero) > 0 || details.CourtCostsPaid.Cmp(zero) > 0) {
		newStatus = data.PartFlatRateStatus
	} else {
		newStatus = data.UnpaidFlatRateStatus
	}

	if newStatus != flatrate.Status {
		flatrate.Status = newStatus
		err = h.repoFlatRate.Update(*flatrate)
		if err != nil {
			return nil, 0, errors.ErrInternalServer
		}
	}

	return details, newStatus, nil
}

func (h *FlatRateSharedLogicServiceImpl) getFlatRatePaymentsByFlatRateID(flatrateID int) ([]*data.FlatRatePayment, error) {
	cond := db.Cond{"flat_rate_id": flatrateID}

	flatratePayments, _, err := h.repoFlatRatePayment.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return flatratePayments, nil
}
