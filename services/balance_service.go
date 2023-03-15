package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type BalanceService interface {
	GetBalance(balanceId int32) (result *model.Balance, RowsAffected int64, err error)
}

func NewBalanceService(mysqlConnection *gorm.DB) BalanceService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) GetBalance(balanceId int32) (result *model.Balance, RowsAffected int64, err error) {
	db := r.mysql.First(&result, balanceId)
	if err = db.Error; err != nil {
		return nil, -1, err
	}
	return result, db.RowsAffected, nil
}
