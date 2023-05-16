package services

import (
	"github.com/jun2900/indo-dispo/model"

	"gorm.io/gorm"
)

type ItemService interface {
	CreateItem(record *model.Item) (results *model.Item, RowsAffected int64, err error)
	GetItem(itemId int32) (result *model.Item, err error)
	GetItemWithItemIdAndSupplierId(itemId, supplierId int32) (result *model.Item, err error)
	GetItemByItemName(itemName string) (result *model.Item, err error)
	GetItemBySupplierId(supplierId int32) (result []model.Item, totalRows int64, err error)
	UpdateItem(itemId int32, updated *model.Item) (result *model.Item, RowsAffected int64, err error)
	DeleteItem(itemId int32) (result *model.Item, RowsAffected int64, err error)
	GetAllItemByItemName(itemName string) (result []model.Item, totalRows int64, err error)
	WithTrx(*gorm.DB) *mysqlDBRepository
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

func (r *mysqlDBRepository) UpdateItem(itemId int32, updated *model.Item) (result *model.Item, RowsAffected int64, err error) {
	result = &model.Item{}
	db := r.mysql.First(result, itemId)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	if err = Copy(result, updated); err != nil {
		return nil, -1, ErrUpdateFailed
	}

	db = db.Save(result)
	if err = db.Error; err != nil {
		return nil, -1, ErrUpdateFailed
	}

	return result, db.RowsAffected, nil
}

func (r *mysqlDBRepository) DeleteItem(itemId int32) (result *model.Item, RowsAffected int64, err error) {
	db := r.mysql.Model(&model.Item{}).First(&result, itemId)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	if err := db.Where("item_id = ?", itemId).Delete(&result).Error; err != nil {
		return nil, -1, ErrDeleteFailed
	}

	return result, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetAllItemByItemName(itemName string) (result []model.Item, totalRows int64, err error) {
	if err = r.mysql.Model(&model.Item{}).Where("item_name = ?", itemName).Count(&totalRows).Find(&result).Error; err != nil {
		return nil, -1, err
	}
	return result, totalRows, nil
}
