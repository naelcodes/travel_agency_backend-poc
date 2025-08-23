package repository

import (
	"fmt"
	"reflect"

	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/types"
	"xorm.io/xorm"
)

const tag = "3"
const embedCustomerSqlQuery = `
	(
	    SELECT
	        to_jsonb (customer)
	    FROM (
	        SELECT
	            id,
	            customer_name AS customerName,
	            account_number AS accountNumber,
	            alias,
	            ab_key AS abKey,
	            state,
	            tmc_client_number AS tmcClientNumber
	        FROM
	            customer
	        WHERE
	            id = payment_received.id_customer) AS customer) AS customer`

const paymentSql = `SELECT id,number,to_char(date,'yyyy-mm-dd') as date,balance::numeric,amount::numeric,used_amount::numeric,fop,status`

func (r *Repository) Count() (int64, error) {
	logger.Info("Counting payments")

	totalRowCount, err := r.Where("tag = ?", tag).Count(new(models.Payment))

	if err != nil {
		logger.Error(fmt.Sprintf("Error counting payments: %v", err))
		return 0, CustomErrors.RepositoryError(fmt.Errorf("error counting payment records: %v", err))
	}

	logger.Info(fmt.Sprintf("Total payment count: %v", totalRowCount))

	return totalRowCount, nil
}

func (r *Repository) GetAll(queryParams *types.GetQueryParams) (*types.GetAllDTO[any], error) {

	embedCustomer := false
	paymentSqlQuery := paymentSql

	totalRowCount, err := r.Count()

	if err != nil {
		return nil, err
	}

	pageNumber := 0
	pageSize := int(totalRowCount)

	if queryParams != nil {
		if queryParams.PageNumber != nil && queryParams.PageSize != nil {
			pageNumber = *queryParams.PageNumber
			pageSize = *queryParams.PageSize
		}

		if queryParams.Embed != nil && *queryParams.Embed == "customer" {
			embedCustomer = true
			paymentSqlQuery = paymentSqlQuery + "," + embedCustomerSqlQuery
		} else {
			paymentSqlQuery = paymentSqlQuery + ",id_customer"
		}

	}

	paymentSqlQuery = paymentSqlQuery + " FROM payment_received  WHERE tag = ?  ORDER BY number DESC LIMIT ? OFFSET ?"

	var result any
	var payments = make([]*models.Payment, 0)
	var paymentsWithCustomer = make([]*struct {
		models.Payment `xorm:"extends"`
		Customer       models.Customer `xorm:"jsonb 'customer'" json:"customer"`
	}, 0)

	if embedCustomer {
		err = r.SQL(paymentSqlQuery, tag, pageSize, pageNumber*pageSize).Find(&paymentsWithCustomer)
		result = paymentsWithCustomer

	} else {
		err = r.SQL(paymentSqlQuery, tag, pageSize, pageNumber*pageSize).Find(&payments)
		result = payments
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting payments: %v", err))
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error getting payment records: %v", err))
	}

	return &types.GetAllDTO[any]{
		Data:          result,
		TotalRowCount: int(totalRowCount),
		PageSize:      pageSize,
		PageNumber:    pageNumber,
	}, nil

}

func (r *Repository) GetById(id int, queryParams *types.GetQueryParams) (any, error) {

	embedCustomer := false
	paymentSqlQuery := paymentSql

	if queryParams != nil {
		if queryParams.Embed != nil && *queryParams.Embed == "customer" {
			embedCustomer = true
			paymentSqlQuery = paymentSqlQuery + "," + embedCustomerSqlQuery
		} else {
			paymentSqlQuery = paymentSqlQuery + ",id_customer"
		}
	}

	paymentSqlQuery = paymentSqlQuery + " FROM payment_received WHERE tag = ? AND id = ?"

	var result any
	var payments = make([]*models.Payment, 0)
	var paymentsWithCustomer = make([]*struct {
		models.Payment `xorm:"extends"`
		Customer       models.Customer `xorm:"jsonb 'customer'" json:"customer"`
	}, 0)

	var err error

	if embedCustomer {
		err = r.SQL(paymentSqlQuery, tag, id).Find(&paymentsWithCustomer)
		result = paymentsWithCustomer

	} else {
		err = r.SQL(paymentSqlQuery, tag, id).Find(&payments)
		result = payments
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting payment: %v", err))
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error getting payment records: %v", err))
	}

	if reflect.ValueOf(result).Len() == 0 {
		logger.Error(fmt.Sprintf("Payment not found: %v", id))
		return nil, CustomErrors.NotFoundError(fmt.Errorf("payment record not found"))
	}

	return reflect.ValueOf(result).Index(0).Interface(), nil

}

func (r *Repository) Save(transaction *xorm.Session, payment *models.Payment) (*models.Payment, error) {

	err := transaction.DB().QueryRow(`
		INSERT INTO payment_received (number, date, balance, amount, base_amount, used_amount, type, fop, status, id_chart_of_accounts, id_currency, id_customer, tag)
		    VALUES (CONCAT('PR-', NEXTVAL('payment_sequence')), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING
		    id`, payment.PaymentDate, payment.Balance, payment.Amount, payment.BaseAmount, payment.UsedAmount, payment.Type, payment.PaymentMode, payment.Status, payment.IdChartOfAccounts, payment.IdCurrency, payment.IdCustomer, payment.Tag).Scan(&payment.Id)

	if err != nil {
		logger.Error(fmt.Sprintf("Error saving payment: %v", err))
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error saving payment: %v", err))
	}

	paymentRecord, err := r.GetById(payment.Id, nil)

	if err != nil {
		return nil, err
	}
	return paymentRecord.(*models.Payment), nil
}

func (r *Repository) CheckPaymentsOwnership(idCustomer int, paymentIds []int) error {

	paymentCount, err := r.Where("tag = ? AND id_customer = ?", tag, idCustomer).In("id", paymentIds).Count(new(models.Payment))

	if err != nil {
		logger.Error(fmt.Sprintf("Error checking payment ownership: %v", err))
		return CustomErrors.RepositoryError(fmt.Errorf("error checking payment ownership: %v", err))
	}

	if paymentCount != int64(len(paymentIds)) {
		logger.Error(fmt.Sprintf("paymentCounts: %v - paymentIdsCount %v", paymentCount, len(paymentIds)))
		return CustomErrors.NotFoundError(fmt.Errorf("some payments are not owned by the customer"))
	}

	return nil

}

func (r *Repository) Update(transaction *xorm.Session, payment *models.Payment) error {

	updateCount, err := transaction.ID(payment.Id).Update(payment)

	if err != nil {
		logger.Error(fmt.Sprintf("Error updating payment: %v", err))
		return CustomErrors.RepositoryError(fmt.Errorf("error updating payment: %v", err))
	}

	if updateCount == 0 {
		logger.Error(fmt.Sprintf("Error updating payment: %v", err))
		return CustomErrors.NotFoundError(fmt.Errorf("payment record not found"))
	}
	return nil

}
