package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type ItemService interface {
	CreateItem(record *model.Item) (results *model.Item, RowsAffected int64, err error)
	GetItem(itemId int32) (result *model.Item, err error)
	GetItemWithItemIdAndSupplierId(itemId, supplierId int32) (result *model.Item, err error)
	GetItemByItemName(itemName string) (result *model.Item, err error)
	GetItemBySupplierId(supplierId int32) (result []model.Item, totalRows int64, err error)
}

func NewItemService(mysqlConnection *gorm.DB) ItemService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateItem(record *model.Item) (results *model.Item, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetItem(itemId int32) (result *model.Item, err error) {
	if err = r.mysql.First(&result, itemId).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mysqlDBRepository) GetItemByItemName(itemName string) (result *model.Item, err error) {
	if err = r.mysql.Where("item_name = ?", itemName).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mysqlDBRepository) GetItemWithItemIdAndSupplierId(itemId, supplierId int32) (result *model.Item, err error) {
	if err = r.mysql.Where("item_id = ? AND supplier_id = ?", itemId, supplierId).First(&result, itemId).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mysqlDBRepository) GetItemBySupplierId(supplierId int32) (result []model.Item, totalRows int64, err error) {
	if err = r.mysql.Model(&model.Item{}).Where("supplier_id = ?", supplierId).Count(&totalRows).Find(&result).Error; err != nil {
		return nil, -1, err
	}
	return result, totalRows, nil
}
