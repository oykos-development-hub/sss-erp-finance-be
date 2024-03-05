package services

import (
	"math"
	"time"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
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
				details.FeeCourtCostsPaid += payment.Amount
			} else {
				details.FeeAllPaymentAmount += payment.Amount
			}
		}
	}

	// calculate the rest of the fees
	details.FeeLeftToPayAmount = fine.Amount - details.FeeAllPaymentAmount
	if fine.CourtCosts != nil {
		details.FeeCourtCostsLeftToPayAmount = *fine.CourtCosts - details.FeeCourtCostsPaid
	}

	details.FeeAmountGracePeriodDueDate = fine.DecisionDate.AddDate(0, 0, data.FineGracePeriod)
	details.FeeAmountGracePeriod = fine.Amount * 2 / 3

	if time.Until(details.FeeAmountGracePeriodDueDate) > 0 {
		details.FeeAmountGracePeriodAvailable = true
		details.FeeLeftToPayAmount = details.FeeAmountGracePeriod - details.FeeAllPaymentAmount
	}

	var newStatus data.FineStatus
	const tolerance = 0.00001

	feeLeftToPayAmount := math.Max(0, details.FeeLeftToPayAmount)
	feeCourtCostsLeftToPayAmount := math.Max(0, details.FeeCourtCostsLeftToPayAmount)

	if math.Abs(feeLeftToPayAmount-0) < tolerance && math.Abs(feeCourtCostsLeftToPayAmount-0) < tolerance {
		newStatus = data.PaidFineStatus
	} else if feeLeftToPayAmount > 0 || feeCourtCostsLeftToPayAmount > 0 {
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
