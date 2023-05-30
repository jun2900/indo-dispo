package services

import (
	"github.com/jun2900/indo-dispo/model"
	"gorm.io/gorm"
)

type ItemSellService interface {
	CreateItemSell(record []model.ItemSell) (results []model.ItemSell, RowsAffected int64, err error)
	GetAllItemSellsByInvoiceId(invoiceId int32) (results []model.ItemSell, totalRows int64, err error)
	DeleteItemSellsByInvoiceId(invoiceId int32) (err error)
	WithTrx(*gorm.DB) *mysqlDBRepository
}

func NewItemSellService(mysqlConnection *gorm.DB) ItemSellService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateItemSell(record []model.ItemSell) (results []model.ItemSell, RowsAffected int64, err error) {
	db := r.mysql.Save(&record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}
	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetAllItemSellsByInvoiceId(invoiceId int32) (results []model.ItemSell, totalRows int64, err error) {
	if err = r.mysql.Model(&model.ItemSell{}).Where("invoice_id = ?", invoiceId).Count(&totalRows).Find(&results).Error; err != nil {
		return nil, -1, err
	}

	return results, totalRows, nil
}

func (r *mysqlDBRepository) DeleteItemSellsByInvoiceId(invoiceId int32) (err error) {
	if err := r.mysql.Where("invoice_id = ?", invoiceId).Delete(&model.ItemSell{}).Error; err != nil {
		return err
	}
	return nil
}
