package services

import (
	"github.com/jun2900/indo-dispo/model"

	"gorm.io/gorm"
)

type ItemPurchaseService interface {
	GetAllItemPurchaseByBillId(billId int32) (results []model.ItemPurchase, totalRows int64, err error)
	CreateItemPurchase(record []model.ItemPurchase) (results []model.ItemPurchase, RowsAffected int64, err error)
	DeleteItemPurchasesByBillId(billId int32) (err error)
}

func NewItemPurchaseService(mysqlConnection *gorm.DB) ItemPurchaseService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) GetAllItemPurchaseByBillId(billId int32) (results []model.ItemPurchase, totalRows int64, err error) {
	if err = r.mysql.Model(&model.ItemPurchase{}).Where("bill_id = ?", billId).Count(&totalRows).Find(&results).Error; err != nil {
		return nil, -1, err
	}

	return results, totalRows, nil
}

func (r *mysqlDBRepository) CreateItemPurchase(record []model.ItemPurchase) (results []model.ItemPurchase, RowsAffected int64, err error) {
	db := r.mysql.Save(&record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) DeleteItemPurchasesByBillId(billId int32) (err error) {
	if err := r.mysql.Where("bill_id = ?", billId).Delete(&model.ItemPurchase{}).Error; err != nil {
		return err
	}
	return nil
}
