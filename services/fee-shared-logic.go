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

type FeeSharedLogicServiceImpl struct {
	App            *celeritas.Celeritas
	repoFee        data.Fee
	repoFeePayment data.FeePayment
}

// NewFeeSharedLogicServiceImpl  creates a new instance of FeeSharedLogicService
func NewFeeSharedLogicServiceImpl(app *celeritas.Celeritas, repoFee data.Fee, repoFeePayment data.FeePayment) FeeSharedLogicService {
	return &FeeSharedLogicServiceImpl{
		App:            app,
		repoFee:        repoFee,
		repoFeePayment: repoFeePayment,
	}
}

func (h *FeeSharedLogicServiceImpl) CalculateFeeDetailsAndUpdateStatus(ctx context.Context, feeId int) (*dto.FeeDetailsDTO, data.FeeStatus, error) {
	fee, err := h.repoFee.Get(feeId)
	if err != nil {
		return nil, 0, newErrors.Wrap(err, "repo fine get")
	}
	payments, err := h.getFeePaymentsByFeeID(fee.ID)
	if err != nil {
		return nil, 0, newErrors.Wrap(err, "get fine payments by fine id")
	}

	details := &dto.FeeDetailsDTO{}

	var paidDuringGracePeriod float64

	details.FeeAmountGracePeriodDueDate = fee.DecisionDate.AddDate(0, 0, data.FineGracePeriod)
	details.FeeAmountGracePeriod = math.Ceil(float64(fee.Amount) * 2 / 3)

	// count all payments and court costs payments
	for _, payment := range payments {
		if data.FeePaymentStatus(payment.Status) == data.PaidFeePeymentStatus {

			details.FeeAllPaymentAmount += payment.Amount
			if payment.PaymentDate.Before(details.FeeAmountGracePeriodDueDate) {
				paidDuringGracePeriod += payment.Amount
			}
		}
	}

	// calculate the rest of the fees
	details.FeeLeftToPayAmount = fee.Amount - details.FeeAllPaymentAmount

	var newStatus data.FeeStatus
	const tolerance = 0.00001

	if time.Until(details.FeeAmountGracePeriodDueDate) > 0 || (paidDuringGracePeriod+tolerance > details.FeeAmountGracePeriod) {
		details.FeeAmountGracePeriodAvailable = true
		details.FeeLeftToPayAmount = details.FeeAmountGracePeriod - details.FeeAllPaymentAmount
	}

	feeLeftToPayAmount := math.Max(0, details.FeeLeftToPayAmount)

	if math.Abs(feeLeftToPayAmount-0) < tolerance {
		newStatus = data.PaidFeeStatus
	} else if feeLeftToPayAmount > 0 && (details.FeeAllPaymentAmount > 0) {
		newStatus = data.PartFeeStatus
	} else {
		newStatus = data.UnpaidFeeStatus
	}

	if newStatus != fee.Status {
		fee.Status = newStatus
		err = h.repoFee.Update(ctx, *fee)
		if err != nil {
			return nil, 0, newErrors.Wrap(err, "repo fine update")
		}
	}

	details.FeeAmountGracePeriodDueDate = details.FeeAmountGracePeriodDueDate.Add(-24 * time.Hour)

	return details, newStatus, nil
}

func (h *FeeSharedLogicServiceImpl) getFeePaymentsByFeeID(feeID int) ([]*data.FeePayment, error) {
	cond := db.Cond{"fee_id": feeID}

	feePayments, _, err := h.repoFeePayment.GetAll(&cond)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo fee payment get all")
	}

	return feePayments, nil
}
