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
	            account_number as accountNumber,
	            alias,
	            ab_key AS abKey,
	            state,
	            tmc_client_number AS tmcClientNumber
	        FROM
	            customer
	        WHERE
	            id = invoice.id_customer) AS customer) AS customer`
const embedTravelItemsSqlQuery = ` 
	(
	SELECT
		jsonb_agg(travelItems)
	FROM (
		SELECT
			id,
			total_price::NUMERIC AS totalPrice,
			itinerary,
			traveler_name AS travelerName,
			ticket_number AS ticketNumber
		FROM
			air_booking
		WHERE
			id_invoice = invoice.id) AS travelItems) AS travelItems`

const invoiceSql = `
	SELECT
	    id,
	    invoice_number,
	    to_char(creation_date, 'yyyy-mm-dd') as creation_date,
	    to_char(due_date, 'yyyy-mm-dd') as due_date,
	    amount::NUMERIC,
	    balance::NUMERIC,
	    credit_apply::NUMERIC,
	    status
	   `

func (r *Repository) Count() (int64, error) {
	logger.Info("Counting invoices")

	totalRowCount, err := r.Where("tag = ?", tag).Count(new(models.Invoice))

	if err != nil {
		logger.Error(fmt.Sprintf("Error counting invoices: %v", err))
		return 0, CustomErrors.RepositoryError(fmt.Errorf("error counting invoice records: %v", err))
	}

	logger.Info(fmt.Sprintf("Total invoice count: %v", totalRowCount))

	return totalRowCount, nil
}

func (r *Repository) GetAll(queryParams *types.GetQueryParams) (*types.GetAllDTO[any], error) {

	embedCustomer := false
	invoiceQuery := invoiceSql + "," + embedTravelItemsSqlQuery
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
			invoiceQuery = invoiceQuery + "," + embedCustomerSqlQuery
		} else {
			invoiceQuery = invoiceQuery + ",id_customer"
		}

	}

	invoiceQuery = invoiceQuery + " FROM invoice WHERE tag = ? ORDER BY invoice_number ASC LIMIT ? OFFSET ? "

	var result any
	var invoices = make([]*struct {
		models.Invoice `xorm:"extends"`
		TravelItems    []models.TravelItem `xorm:"jsonb 'travelItems'" json:"travelItems"`
	}, 0)

	var invoicesWithCustomer = make([]*struct {
		models.Invoice `xorm:"extends"`
		Customer       models.Customer     `xorm:"jsonb 'customer'" json:"customer"`
		TravelItems    []models.TravelItem `xorm:"jsonb 'travelItems'" json:"travelItems"`
	}, 0)

	if embedCustomer {
		err = r.SQL(invoiceQuery, tag, pageSize, pageNumber*pageSize).Find(&invoicesWithCustomer)
		result = invoicesWithCustomer
	} else {
		err = r.SQL(invoiceQuery, tag, pageSize, pageNumber*pageSize).Find(&invoices)
		result = invoices
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting invoices: %v", err))
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error getting invoice records: %v", err))
	}

	return &types.GetAllDTO[any]{
		Data:          result,
		TotalRowCount: int(totalRowCount),
		PageSize:      pageSize,
		PageNumber:    pageNumber,
	}, nil

}

func (r *Repository) GetById(id int, queryParams *types.GetQueryParams, embedTravelItems bool) (any, error) {

	embedCustomer := false
	invoiceQuery := invoiceSql

	if embedTravelItems {
		invoiceQuery = invoiceQuery + "," + embedTravelItemsSqlQuery
	}

	if queryParams != nil {
		if queryParams.Embed != nil && *queryParams.Embed == "customer" {
			embedCustomer = true
			invoiceQuery = invoiceQuery + "," + embedCustomerSqlQuery
		} else {
			invoiceQuery = invoiceQuery + ",id_customer"
		}

	} else {
		invoiceQuery = invoiceQuery + ",id_customer"
	}

	invoiceQuery = invoiceQuery + " FROM invoice WHERE id = ?"

	var result any
	var err error

	if embedTravelItems && !embedCustomer {
		var invoices = make([]*struct {
			models.Invoice `xorm:"extends"`
			TravelItems    []models.TravelItem `xorm:"jsonb 'travelItems'" json:"travelItems,omitempty"`
		}, 0)

		err = r.SQL(invoiceQuery, id).Find(&invoices)
		result = invoices

	} else if !embedTravelItems && embedCustomer {
		var invoices = make([]*struct {
			models.Invoice `xorm:"extends"`
			Customer       models.Customer `xorm:"jsonb 'customer'" json:"customer"`
		}, 0)

		err = r.SQL(invoiceQuery, id).Find(&invoices)
		result = invoices

	} else if embedTravelItems && embedCustomer {
		var invoices = make([]*struct {
			models.Invoice `xorm:"extends"`
			Customer       models.Customer     `xorm:"jsonb 'customer'" json:"customer"`
			TravelItems    []models.TravelItem `xorm:"jsonb 'travelItems'" json:"travelItems,omitempty"`
		}, 0)

		err = r.SQL(invoiceQuery, id).Find(&invoices)
		result = invoices

	} else {
		var invoices = []*models.Invoice{}
		err = r.SQL(invoiceQuery, id).Find(&invoices)
		result = invoices
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting invoice: %v", err))
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error getting invoice records: %v", err))
	}

	if reflect.ValueOf(result).Len() == 0 {
		logger.Error(fmt.Sprintf("Error getting invoice: %v", err))
		return nil, CustomErrors.NotFoundError(fmt.Errorf("invoice record not found"))
	}

	return reflect.ValueOf(result).Index(0).Interface(), nil

}

func (r *Repository) GetByCustomerId(idCustomer int, queryParams *types.GetQueryParams, paid bool) (*types.GetAllDTO[any], error) {

	WhereCondition := "WHERE tag = ? AND id_customer = ?"
	invoiceQuery := invoiceSql

	if paid {
		WhereCondition = WhereCondition + " AND status = 'paid'"
	} else {
		WhereCondition = WhereCondition + " AND status = 'paid'"
	}

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
	}

	invoiceQuery = invoiceQuery + " FROM invoice " + WhereCondition + " ORDER BY invoice_number ASC LIMIT ? OFFSET ? "

	var invoices = make([]*struct {
		models.Invoice `xorm:"extends"`
		TravelItems    []models.TravelItem `xorm:"jsonb 'travelItems'" json:"travelItems"`
	}, 0)

	err = r.SQL(invoiceQuery, tag, idCustomer, pageSize, pageNumber*pageSize).Find(&invoices)

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting invoices: %v", err))
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error getting customer's invoice records: %v", err))
	}

	logger.Info(fmt.Sprintf("Total customer's invoice count: %v", len(invoices)))

	return &types.GetAllDTO[any]{
		Data:          invoices,
		TotalRowCount: int(totalRowCount),
		PageSize:      pageSize,
		PageNumber:    pageNumber,
	}, nil

}

func (r *Repository) Save(transaction *xorm.Session, invoice *models.Invoice) (*models.Invoice, error) {

	insertSqlCommand := `
	INSERT INTO invoice (invoice_number, creation_date, due_date, amount, balance, credit_apply, net_amount, base_amount, status, id_customer, tag)
		VALUES (CONCAT('INV-', NEXTVAL('invoice_sequence')), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING
		id, invoice_number, to_char(creation_date, 'yyyy-mm-dd') AS creation_date, to_char(due_date, 'yyyy-mm-dd') AS due_date, amount::NUMERIC, balance::NUMERIC, credit_apply::NUMERIC, status, id_customer`

	err := transaction.DB().QueryRow(insertSqlCommand, invoice.CreationDate, invoice.DueDate, invoice.Amount, invoice.Balance, invoice.CreditApply, invoice.NetAmount, invoice.BaseAmount, invoice.Status, invoice.IdCustomer, invoice.Tag).ScanStructByName(invoice)

	if err != nil {
		logger.Error(fmt.Sprintf("Error saving invoice: %v", err))
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error saving invoice: %v", err))
	}

	return invoice, nil
}

func (r *Repository) Update(transaction *xorm.Session, invoice *models.Invoice) error {

	updateCount, err := transaction.ID(invoice.Id).Update(invoice)

	if err != nil {
		logger.Error(fmt.Sprintf("Error updating invoice: %v", err))
		return CustomErrors.RepositoryError(fmt.Errorf("error updating invoice: %v", err))
	}

	if updateCount == 0 {
		logger.Error(fmt.Sprintf("Error updating invoice: %v", err))
		return CustomErrors.NotFoundError(fmt.Errorf("invoice record not found"))
	}

	return nil

}
