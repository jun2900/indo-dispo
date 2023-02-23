package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type ItemSupplierService interface {
	CreateItemSupplier(record *model.ItemSupplier) (results *model.ItemSupplier, RowsAffected int64, err error)
	GetItemSupplierByItemIdAndSupplierId(itemId, supplierId int32) (result *model.ItemSupplier, RowsAffected int64, err error)
}

func NewItemSupplierService(mysqlConnection *gorm.DB) ItemSupplierService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateItemSupplier(record *model.ItemSupplier) (results *model.ItemSupplier, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetItemSupplierByItemIdAndSupplierId(itemId, supplierId int32) (result *model.ItemSupplier, RowsAffected int64, err error) {
	if err = r.mysql.Where("item_id = ? AND supplier_id = ?", itemId, supplierId).First(&result).Error; err != nil {
		return nil, -1, err
	}
	return result, RowsAffected, nil
}
