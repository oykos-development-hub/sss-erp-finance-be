package services

import (
	"context"
	"math"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/upper/db/v4"
)

type PropBenConfSharedLogicServiceImpl struct {
	App                    *celeritas.Celeritas
	repoPropBenConf        data.PropBenConf
	repoPropBenConfPayment data.PropBenConfPayment
}

// PropBenConfSharedLogicServiceImpl  creates a new instance of PropBenConfService
func NewPropBenConfSharedLogicServiceImpl(app *celeritas.Celeritas, repoPropBenConf data.PropBenConf, repoPropBenConfPayment data.PropBenConfPayment) PropBenConfSharedLogicService {
	return &PropBenConfSharedLogicServiceImpl{
		App:                    app,
		repoPropBenConf:        repoPropBenConf,
		repoPropBenConfPayment: repoPropBenConfPayment,
	}
}

func (h *PropBenConfSharedLogicServiceImpl) CalculatePropBenConfDetailsAndUpdateStatus(ctx context.Context, propbenconfId int) (*dto.PropBenConfDetailsDTO, data.PropBenConfStatus, error) {
	propbenconf, err := h.repoPropBenConf.Get(propbenconfId)
	if err != nil {
		return nil, 0, newErrors.Wrap(err, "repo prop ben conf get")
	}
	payments, err := h.getPropBenConfPaymentsByPropBenConfID(propbenconf.ID)
	if err != nil {
		return nil, 0, newErrors.Wrap(err, "get prop ben conf payments by prop ben conf id")
	}

	details := &dto.PropBenConfDetailsDTO{}

	// count all payments and court costs payments
	for _, payment := range payments {
		if data.PropBenConfPaymentStatus(payment.Status) == data.PaidPropBenConfPeymentStatus {
			if data.PropBenConfPaymentMethod(payment.PaymentMethod) == data.CourtCostsPropBenConfPeymentMethod {
				details.CourtCostsPaid += payment.Amount
			} else {
				details.AllPaymentAmount += payment.Amount
			}
		}
	}

	// calculate the rest of the fees
	details.LeftToPayAmount = propbenconf.Amount - details.AllPaymentAmount
	if propbenconf.CourtCosts != nil {
		details.CourtCostsLeftToPayAmount = *propbenconf.CourtCosts - details.CourtCostsPaid
	}

	details.AmountGracePeriodDueDate = propbenconf.DecisionDate.AddDate(0, 0, data.PropBenConfGracePeriod)
	details.AmountGracePeriod = math.Ceil(float64(propbenconf.Amount) * 2 / 3)

	if time.Until(details.AmountGracePeriodDueDate) > 0 {
		details.AmountGracePeriodAvailable = true
		details.LeftToPayAmount = details.AmountGracePeriod - details.AllPaymentAmount
	}

	var newStatus data.PropBenConfStatus
	const tolerance = 0.00001

	feeLeftToPayAmount := math.Max(0, details.LeftToPayAmount)
	feeCourtCostsLeftToPayAmount := math.Max(0, details.CourtCostsLeftToPayAmount)

	if math.Abs(feeLeftToPayAmount-0) < tolerance && math.Abs(feeCourtCostsLeftToPayAmount-0) < tolerance {
		newStatus = data.PaidPropBenConfStatus
	} else if (feeLeftToPayAmount > 0 || feeCourtCostsLeftToPayAmount > 0) && (details.AllPaymentAmount > 0 || details.CourtCostsPaid > 0) {
		newStatus = data.PartPropBenConfStatus
	} else {
		newStatus = data.UnpaidPropBenConfStatus
	}

	if newStatus != propbenconf.Status {
		propbenconf.Status = newStatus
		err = h.repoPropBenConf.Update(ctx, *propbenconf)
		if err != nil {
			return nil, 0, newErrors.Wrap(err, "repo prop ben conf update")
		}
	}

	return details, newStatus, nil
}

func (h *PropBenConfSharedLogicServiceImpl) getPropBenConfPaymentsByPropBenConfID(propbenconfID int) ([]*data.PropBenConfPayment, error) {
	cond := db.Cond{"property_benefits_confiscation_id": propbenconfID}

	propbenconfPayments, _, err := h.repoPropBenConfPayment.GetAll(&cond)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo prop ben conf payments get all")
	}

	return propbenconfPayments, nil
}
