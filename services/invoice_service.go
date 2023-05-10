package services

import (
	"github.com/jun2900/indo-dispo/model"
	"gorm.io/gorm"
)

type InvoiceService interface {
	CreateInvoice(record *model.Invoice) (results *model.Invoice, RowsAffected int64, err error)
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
