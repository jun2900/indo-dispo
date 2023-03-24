package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type RecurringBillService interface {
	CreateRecurringBill(record *model.RecurringBill) (results *model.RecurringBill, RowsAffected int64, err error)
}

func NewRecurringBillService(mysqlConnection *gorm.DB) RecurringBillService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateRecurringBill(record *model.RecurringBill) (results *model.RecurringBill, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}
