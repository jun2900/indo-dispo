package services

import (
	"github.com/jun2900/indo-dispo/model"

	"gorm.io/gorm"
)

type WholesalerService interface {
	CreateWholesaler(record []model.Wholesaler) (result []model.Wholesaler, RowsAffected int64, err error)
	DeleteWholesalerByItemId(itemId int32) (result []model.Wholesaler, RowsAffected int64, err error)
	GetWholesalerByItemId(itemId int32) (result []model.Wholesaler, RowsAffected int64, err error)
}

func NewWholesalerService(mysqlConnection *gorm.DB) WholesalerService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateWholesaler(record []model.Wholesaler) (result []model.Wholesaler, RowsAffected int64, err error) {
	db := r.mysql.Save(&record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) DeleteWholesalerByItemId(itemId int32) (result []model.Wholesaler, RowsAffected int64, err error) {
	db := r.mysql.Model(&model.Wholesaler{}).Where("item_id = ?", itemId).First(&result)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	if err := db.Where("item_id = ?", itemId).Delete(&result).Error; err != nil {
		return nil, -1, ErrDeleteFailed
	}

	return result, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetWholesalerByItemId(itemId int32) (result []model.Wholesaler, RowsAffected int64, err error) {
	if err := r.mysql.Model(&model.Wholesaler{}).Where("item_id = ?", itemId).Count(&RowsAffected).Find(&result).Error; err != nil {
		return nil, -1, err
	}
	return result, RowsAffected, nil
}
