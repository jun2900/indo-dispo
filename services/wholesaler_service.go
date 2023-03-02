package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type WholesalerService interface {
	CreateWholesaler(record []model.Wholesaler) (result []model.Wholesaler, RowsAffected int64, err error)
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
