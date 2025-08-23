package service

import (
	"fmt"

	"neema.co.za/rest/utils/domains"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/managers"
	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/payloads"

	CustomErrors "neema.co.za/rest/utils/errors"
)

func (s *Service) GetImputationsService(idInvoice int) (any, error) {
	return s.Repository.GetByInvoiceId(idInvoice)
}

func (s *Service) ApplyImputationsService(idInvoice int, payload []*payloads.ImputationPayload) (int, int, int, error) {

	insertedImputationCount := 0
	updateImputationCount := 0
	deletedImputationCount := 0

	invoice, err := s.Imports.GetInvoice(idInvoice)

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting invoice: %v", err))
		return 0, 0, 0, err
	}

	invoiceDomain := domains.NewInvoiceDomain(invoice)

	if err = invoiceDomain.Validate(); err != nil {
		return 0, 0, 0, err
	}

	paymentIds := []int{}
	for _, imputationPayload := range payload {
		paymentIds = append(paymentIds, imputationPayload.IdPayment)
	}

	logger.Info(fmt.Sprintf("CustomerId: %v", invoice.IdCustomer))

	if err = s.Imports.CheckPaymentsOwnership(invoice.IdCustomer, paymentIds); err != nil {
		logger.Error(fmt.Sprintf("Error checking payment ownership: %v", err))
		return 0, 0, 0, err
	}

	TransactionManager := managers.NewTransactionManager(s.Engine)
	err = TransactionManager.Begin()

	if err != nil {
		return 0, 0, 0, CustomErrors.UnknownError(err)
	}

	for _, imputationPayload := range payload {

		imputedDifference := float64(0)
		exists, data, err := s.Repository.GetByPaymentIdAndInvoiceId(imputationPayload.IdPayment, idInvoice)

		if err != nil {
			logger.Error(fmt.Sprintf("Error getting imputation: %v", err))
			return 0, 0, 0, err
		}

		if exists {
			if data.Imputation.AmountApplied != imputationPayload.AmountApplied && (imputationPayload.AmountApplied > 0 || imputationPayload.AmountApplied == 0) {
				imputedDifference = imputationPayload.AmountApplied - data.Imputation.AmountApplied

				paymentDomain := domains.NewPaymentDomain(&data.Payment)

				if err = paymentDomain.Validate(); err != nil {
					return 0, 0, 0, err
				}

				if err = paymentDomain.AllocateAmount(data.Imputation.AmountApplied, imputationPayload.AmountApplied); err != nil {
					return 0, 0, 0, err
				}

				if err = s.Imports.UpdatePayment(TransactionManager.GetTransaction(), paymentDomain.GetPayment()); err != nil {
					logger.Error(fmt.Sprintf("Error saving payment allocation: %v", err))
					TransactionManager.Rollback()
					return 0, 0, 0, err
				}

				data.Imputation.AmountApplied = imputationPayload.AmountApplied

				if imputationPayload.AmountApplied > 0 {
					if err = s.Repository.Update(TransactionManager.GetTransaction(), &data.Imputation); err != nil {
						logger.Error(fmt.Sprintf("Error updating imputation: %v", err))
						TransactionManager.Rollback()
						return 0, 0, 0, err
					}
					updateImputationCount++
				} else {
					if err = s.Repository.DeleteById(TransactionManager.GetTransaction(), data.Imputation.Id); err != nil {
						logger.Error(fmt.Sprintf("Error deleting imputation record: %v", err))
						TransactionManager.Rollback()
						return 0, 0, 0, err
					}
					deletedImputationCount++
				}

				if err = invoiceDomain.ApplyImputation(imputedDifference); err != nil {
					return 0, 0, 0, err
				}

			}

		} else {

			if imputationPayload.AmountApplied == 0 {
				continue
			}

			paymentRecord, err := s.Imports.GetPayment(imputationPayload.IdPayment)

			if err != nil {
				logger.Error(fmt.Sprintf("Error getting payment: %v", err))
				TransactionManager.Rollback()
				return 0, 0, 0, err
			}

			paymentDomain := domains.NewPaymentDomain(paymentRecord)

			if err = paymentDomain.Validate(); err != nil {
				return 0, 0, 0, err
			}

			if err = paymentDomain.AllocateAmount(0, imputationPayload.AmountApplied); err != nil {
				return 0, 0, 0, err
			}

			if err = s.Imports.UpdatePayment(TransactionManager.GetTransaction(), paymentDomain.GetPayment()); err != nil {
				logger.Error(fmt.Sprintf("Error saving payment allocation: %v", err))
				TransactionManager.Rollback()
				return 0, 0, 0, err
			}

			imputation := models.Imputation{}

			imputation.IdInvoice = idInvoice
			imputation.IdPaymentReceived = imputationPayload.IdPayment
			imputation.AmountApplied = imputationPayload.AmountApplied
			imputation.InvoiceAmount = invoiceDomain.GetInvoice().Amount
			imputation.PaymentAmount = paymentDomain.GetPayment().Amount
			imputation.Tag = "3"

			if err = s.Repository.Save(TransactionManager.GetTransaction(), &imputation); err != nil {
				logger.Error(fmt.Sprintf("Error saving imputation: %v", err))
				TransactionManager.Rollback()
				return 0, 0, 0, err
			}

			insertedImputationCount++

			if err = invoiceDomain.ApplyImputation(imputationPayload.AmountApplied); err != nil {
				return 0, 0, 0, err
			}
		}

		if insertedImputationCount > 0 || updateImputationCount > 0 || deletedImputationCount > 0 {
			if err = s.Imports.UpdateInvoice(TransactionManager.GetTransaction(), invoiceDomain.GetInvoice()); err != nil {
				logger.Error(fmt.Sprintf("Error saving invoice: %v", err))
				TransactionManager.Rollback()
				return 0, 0, 0, err
			}
		}

	}

	TransactionManager.Commit()

	logger.Info(fmt.Sprintf("Applying imputations to invoice: %v", *invoice))
	return insertedImputationCount, updateImputationCount, deletedImputationCount, nil
}
