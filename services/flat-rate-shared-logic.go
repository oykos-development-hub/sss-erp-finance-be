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

func (h *FlatRateSharedLogicServiceImpl) CalculateFlatRateDetailsAndUpdateStatus(ctx context.Context, flatrateId int) (*dto.FlatRateDetailsDTO, data.FlatRateStatus, error) {
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
				details.CourtCostsPaid += payment.Amount
			} else {
				details.AllPaymentAmount += payment.Amount
			}
		}
	}

	// calculate the rest of the flatrates
	details.LeftToPayAmount = flatrate.Amount - details.AllPaymentAmount
	if flatrate.CourtCosts != nil {
		details.CourtCostsLeftToPayAmount = *flatrate.CourtCosts - details.CourtCostsPaid
	}

	details.AmountGracePeriodDueDate = flatrate.DecisionDate.AddDate(0, 0, data.FlatRateGracePeriod)
	details.AmountGracePeriod = math.Ceil(float64(flatrate.Amount) * 2 / 3)

	if time.Until(details.AmountGracePeriodDueDate) > 0 {
		details.AmountGracePeriodAvailable = true
		details.LeftToPayAmount = details.AmountGracePeriod - details.AllPaymentAmount
	}

	var newStatus data.FlatRateStatus
	const tolerance = 0.00001

	flatrateLeftToPayAmount := math.Max(0, details.LeftToPayAmount)
	flatrateCourtCostsLeftToPayAmount := math.Max(0, details.CourtCostsLeftToPayAmount)

	if math.Abs(flatrateLeftToPayAmount-0) < tolerance && math.Abs(flatrateCourtCostsLeftToPayAmount-0) < tolerance {
		newStatus = data.PaidFlatRateStatus
	} else if (flatrateLeftToPayAmount > 0 || flatrateCourtCostsLeftToPayAmount > 0) && (details.AllPaymentAmount > 0 || details.CourtCostsPaid > 0) {
		newStatus = data.PartFlatRateStatus
	} else {
		newStatus = data.UnpaidFlatRateStatus
	}

	if newStatus != flatrate.Status {
		flatrate.Status = newStatus
		err = h.repoFlatRate.Update(ctx, *flatrate)
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
