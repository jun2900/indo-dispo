package services

import (
	"github.com/jun2900/indo-dispo/model"
	"gorm.io/gorm"
)

type ItemSellService interface {
	CreateItemSell(record []model.ItemSell) (results []model.ItemSell, RowsAffected int64, err error)
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
