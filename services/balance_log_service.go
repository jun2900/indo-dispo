package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type BalanceLogService interface {
	CreateBalanceLog(record *model.BalanceLog) (results *model.BalanceLog, RowsAffected int64, err error)
}

func NewBalanceLogService(mysqlConnection *gorm.DB) BalanceLogService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateBalanceLog(record *model.BalanceLog) (results *model.BalanceLog, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}
