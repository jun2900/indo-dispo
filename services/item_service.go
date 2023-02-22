package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type ItemService interface {
	CreateItem(record *model.Item) (results *model.Item, RowsAffected int64, err error)
	GetItem(itemId int32) (result *model.Item, err error)
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
