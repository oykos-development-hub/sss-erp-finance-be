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

func (h *PropBenConfSharedLogicServiceImpl) CalculatePropBenConfDetailsAndUpdateStatus(propbenconfId int) (*dto.PropBenConfDetailsDTO, data.PropBenConfStatus, error) {
	propbenconf, err := h.repoPropBenConf.Get(propbenconfId)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}
	payments, err := h.getPropBenConfPaymentsByPropBenConfID(propbenconf.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}

	details := &dto.PropBenConfDetailsDTO{}

	// count all payments and court costs payments
	for _, payment := range payments {
		if data.PropBenConfPaymentStatus(payment.Status) == data.PaidPropBenConfPeymentStatus {
			if data.PropBenConfPaymentMethod(payment.PaymentMethod) == data.CourtCostsPropBenConfPeymentMethod {
				details.CourtCostsPaid = details.CourtCostsPaid.Add(payment.Amount)
			} else {
				details.AllPaymentAmount = details.AllPaymentAmount.Add(payment.Amount)
			}
		}
	}

	// calculate the rest of the fees
	details.LeftToPayAmount = propbenconf.Amount.Sub(details.AllPaymentAmount)
	if propbenconf.CourtCosts != nil {
		details.CourtCostsLeftToPayAmount = propbenconf.CourtCosts.Sub(details.CourtCostsPaid)
	}

	details.AmountGracePeriodDueDate = propbenconf.DecisionDate.AddDate(0, 0, data.PropBenConfGracePeriod)

	twoThirds := decimal.NewFromFloat(2.0).Div(decimal.NewFromFloat(3.0))
	details.AmountGracePeriod = propbenconf.Amount.Mul(twoThirds).Ceil()

	if time.Until(details.AmountGracePeriodDueDate) > 0 {
		details.AmountGracePeriodAvailable = true
		details.LeftToPayAmount = details.AmountGracePeriod.Sub(details.AllPaymentAmount)
	}

	var newStatus data.PropBenConfStatus
	tolerance := decimal.NewFromFloat(0.00001)

	// Korišćenje decimal.Max za poređenje sa 0
	zero := decimal.NewFromInt(0)
	feeLeftToPayAmount := decimal.Max(zero, details.LeftToPayAmount)
	feeCourtCostsLeftToPayAmount := decimal.Max(zero, details.CourtCostsLeftToPayAmount)

	// Provera uslova i postavljanje statusa
	if feeLeftToPayAmount.Abs().Cmp(tolerance) < 0 && feeCourtCostsLeftToPayAmount.Abs().Cmp(tolerance) < 0 {
		newStatus = data.PaidPropBenConfStatus
	} else if (feeLeftToPayAmount.Cmp(zero) > 0 || feeCourtCostsLeftToPayAmount.Cmp(zero) > 0) &&
		(details.AllPaymentAmount.Cmp(zero) > 0 || details.CourtCostsPaid.Cmp(zero) > 0) {
		newStatus = data.PartPropBenConfStatus
	} else {
		newStatus = data.UnpaidPropBenConfStatus
	}

	if newStatus != propbenconf.Status {
		propbenconf.Status = newStatus
		err = h.repoPropBenConf.Update(*propbenconf)
		if err != nil {
			return nil, 0, errors.ErrInternalServer
		}
	}

	return details, newStatus, nil
}

func (h *PropBenConfSharedLogicServiceImpl) getPropBenConfPaymentsByPropBenConfID(propbenconfID int) ([]*data.PropBenConfPayment, error) {
	cond := db.Cond{"property_benefits_confiscation_id": propbenconfID}

	propbenconfPayments, _, err := h.repoPropBenConfPayment.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return propbenconfPayments, nil
}
