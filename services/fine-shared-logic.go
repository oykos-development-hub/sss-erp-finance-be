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

type FineSharedLogicServiceImpl struct {
	App             *celeritas.Celeritas
	repoFine        data.Fine
	repoFinePayment data.FinePayment
}

// FineSharedLogicServiceImpl  creates a new instance of FineService
func NewFineSharedLogicServiceImpl(app *celeritas.Celeritas, repoFine data.Fine, repoFinePayment data.FinePayment) FineSharedLogicService {
	return &FineSharedLogicServiceImpl{
		App:             app,
		repoFine:        repoFine,
		repoFinePayment: repoFinePayment,
	}
}

func (h *FineSharedLogicServiceImpl) CalculateFineDetailsAndUpdateStatus(fineId int) (*dto.FineFeeDetailsDTO, data.FineStatus, error) {
	fine, err := h.repoFine.Get(fineId)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}
	payments, err := h.getFinePaymentsByFineID(fine.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}

	details := &dto.FineFeeDetailsDTO{}

	// count all payments and court costs payments
	for _, payment := range payments {
		if data.FinePaymentStatus(payment.Status) == data.PaidFinePeymentStatus {
			if data.FinePaymentMethod(payment.PaymentMethod) == data.CourtCostsFinePeymentMethod {
				details.FeeCourtCostsPaid = details.FeeCourtCostsPaid.Add(payment.Amount)
			} else {
				details.FeeAllPaymentAmount = details.FeeAllPaymentAmount.Add(payment.Amount)
			}
		}
	}

	// calculate the rest of the fees
	details.FeeLeftToPayAmount = fine.Amount.Sub(details.FeeAllPaymentAmount)
	if fine.CourtCosts != nil {
		details.FeeCourtCostsLeftToPayAmount = fine.CourtCosts.Sub(details.FeeCourtCostsPaid)
	}

	details.FeeAmountGracePeriodDueDate = fine.DecisionDate.AddDate(0, 0, data.FineGracePeriod)
	twoThirds := decimal.NewFromFloat(2.0).Div(decimal.NewFromFloat(3.0))
	details.FeeAmountGracePeriod = fine.Amount.Mul(twoThirds).Ceil()

	if time.Until(details.FeeAmountGracePeriodDueDate) > 0 {
		details.FeeAmountGracePeriodAvailable = true
		details.FeeLeftToPayAmount = details.FeeAmountGracePeriod.Sub(details.FeeAllPaymentAmount)
	}

	var newStatus data.FineStatus
	tolerance := decimal.NewFromFloat(0.00001)

	zero := decimal.NewFromInt(0)
	feeLeftToPayAmount := decimal.Max(zero, details.FeeLeftToPayAmount)
	feeCourtCostsLeftToPayAmount := decimal.Max(zero, details.FeeCourtCostsLeftToPayAmount)

	absFeeLeftToPayAmount := feeLeftToPayAmount.Abs()
	absFeeCourtCostsLeftToPayAmount := feeCourtCostsLeftToPayAmount.Abs()

	if absFeeLeftToPayAmount.Cmp(tolerance) < 0 && absFeeCourtCostsLeftToPayAmount.Cmp(tolerance) < 0 {
		newStatus = data.PaidFineStatus
	} else if (feeLeftToPayAmount.Cmp(decimal.NewFromInt(0)) > 0 || feeCourtCostsLeftToPayAmount.Cmp(decimal.NewFromInt(0)) > 0) &&
		(details.FeeAllPaymentAmount.Cmp(decimal.NewFromInt(0)) > 0 || details.FeeCourtCostsPaid.Cmp(decimal.NewFromInt(0)) > 0) {
		newStatus = data.PartFineStatus
	} else {
		newStatus = data.UnpaidFineStatus
	}

	if newStatus != fine.Status {
		fine.Status = newStatus
		err = h.repoFine.Update(*fine)
		if err != nil {
			return nil, 0, errors.ErrInternalServer
		}
	}

	return details, newStatus, nil
}

func (h *FineSharedLogicServiceImpl) getFinePaymentsByFineID(fineID int) ([]*data.FinePayment, error) {
	cond := db.Cond{"fine_id": fineID}

	finePayments, _, err := h.repoFinePayment.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return finePayments, nil
}
