package services

import (
	"fmt"
	"time"

	"github.com/jun2900/indo-dispo/model"
	"gorm.io/gorm"
)

type InvoiceService interface {
	CreateInvoice(record *model.Invoice) (results *model.Invoice, RowsAffected int64, err error)
	GetAllInvoices(
		page, pagesize int,
		order, status, customer string,
		dateFrom, dateTo time.Time) (results []model.VSupplierInvoice, totalRows int64, err error)
	GetInvoice(invoiceId int32) (result *model.Invoice, RowsAffected int64, err error)
	WithTrx(*gorm.DB) *mysqlDBRepository
}

func NewInvoiceService(mysqlConnection *gorm.DB) InvoiceService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateInvoice(record *model.Invoice) (result *model.Invoice, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetAllInvoices(
	page, pagesize int,
	order, status, customer string,
	dateFrom, dateTo time.Time) (results []model.VSupplierInvoice, totalRows int64, err error) {

	resultOrm := r.mysql.Model(&model.VSupplierInvoice{})
	if len(status) > 0 {
		resultOrm = resultOrm.Where("invoice_status = ?", status)
	}
	if len(customer) > 0 {
		resultOrm = resultOrm.Where("supplier_name LIKE ?", fmt.Sprint("%", customer, "%"))
	}

	if !dateFrom.IsZero() && !dateTo.IsZero() {
		resultOrm = resultOrm.Where("invoice_start_date >= (? - INTERVAL 1 DAY) AND invoice_start_date < (? + INTERVAL 1 DAY)", dateFrom, dateTo)
	} else if !dateFrom.IsZero() && dateTo.IsZero() {
		resultOrm = resultOrm.Where("invoice_start_date >= (? - INTERVAL 1 DAY)", dateFrom)
	} else if dateFrom.IsZero() && !dateTo.IsZero() {
		resultOrm = resultOrm.Where("invoice_start_date < (? + INTERVAL 1 DAY)", dateTo)
	}

	resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Find(&results).Error; err != nil {
		return nil, -1, err
	}

	return results, totalRows, nil
}

func (r *mysqlDBRepository) GetInvoice(invoiceId int32) (result *model.Invoice, RowsAffected int64, err error) {
	db := r.mysql.First(&result, invoiceId)
	if err = db.Error; err != nil {
		return nil, -1, err
	}
	return result, db.RowsAffected, nil
}
