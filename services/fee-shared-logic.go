package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	"github.com/shopspring/decimal"
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

func (h *FeeSharedLogicServiceImpl) CalculateFeeDetailsAndUpdateStatus(feeId int) (*dto.FeeDetailsDTO, data.FeeStatus, error) {
	fee, err := h.repoFee.Get(feeId)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}
	payments, err := h.getFeePaymentsByFeeID(fee.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, 0, errors.ErrNotFound
	}

	details := &dto.FeeDetailsDTO{}

	// count all payments
	for _, payment := range payments {
		if data.FeePaymentStatus(payment.Status) == data.PaidFeePeymentStatus {
			details.FeeAllPaymentAmount = details.FeeAllPaymentAmount.Add(payment.Amount)
		}
	}

	// calculate the rest of the fees
	details.FeeLeftToPayAmount = fee.Amount.Sub(details.FeeAllPaymentAmount)

	var newStatus data.FeeStatus
	tolerance := decimal.NewFromFloat(0.00001)

	feeLeftToPayAmount := decimal.Max(decimal.NewFromInt(0), details.FeeLeftToPayAmount)

	absFeeLeftToPayAmount := feeLeftToPayAmount.Abs()

	// Provera da li je apsolutna vrednost feeLeftToPayAmount manja od tolerancije
	if absFeeLeftToPayAmount.Cmp(tolerance) < 0 {
		newStatus = data.PaidFeeStatus
	} else if feeLeftToPayAmount.Cmp(decimal.NewFromInt(0)) > 0 && details.FeeAllPaymentAmount.Cmp(decimal.NewFromInt(0)) > 0 {
		newStatus = data.PartFeeStatus
	} else {
		newStatus = data.UnpaidFeeStatus
	}

	if newStatus != fee.Status {
		fee.Status = newStatus
		err = h.repoFee.Update(*fee)
		if err != nil {
			return nil, 0, errors.ErrInternalServer
		}
	}

	return details, newStatus, nil
}

func (h *FeeSharedLogicServiceImpl) getFeePaymentsByFeeID(feeID int) ([]*data.FeePayment, error) {
	cond := db.Cond{"fee_id": feeID}

	feePayments, _, err := h.repoFeePayment.GetAll(&cond)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrInternalServer
	}

	return feePayments, nil
}
