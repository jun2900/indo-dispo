package services

import (
	"github.com/jun2900/indo-dispo/model"

	"gorm.io/gorm"
)

type BalanceService interface {
	GetBalance(balanceId int32) (result *model.Balance, RowsAffected int64, err error)
	UpdateBalanceAmount(balanceId int32, BalanceAmount float64) (result *model.Balance, RowsAffected int64, err error)
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

func (r *mysqlDBRepository) UpdateBalanceAmount(balanceId int32, BalanceAmount float64) (result *model.Balance, RowsAffected int64, err error) {
	result = &model.Balance{}
	db := r.mysql.First(&result, balanceId)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	result.BalanceAmount = BalanceAmount

	if err := db.Save(&result).Error; err != nil {
		return nil, -1, err
	}

	return result, db.RowsAffected, nil
}
