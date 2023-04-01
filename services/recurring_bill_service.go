package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type RecurringBillService interface {
	GetAllRecurringBill(page, pagesize int, order string) (results []model.RecurringBill, TotalRows int64, err error)
	CreateRecurringBill(record *model.RecurringBill) (results *model.RecurringBill, RowsAffected int64, err error)
	UpdateRecurringBill(recurringBillId int32, updated *model.RecurringBill) (result *model.RecurringBill, RowsAffected int64, err error)
}

func NewRecurringBillService(mysqlConnection *gorm.DB) RecurringBillService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) GetAllRecurringBill(page, pagesize int, order string) (results []model.RecurringBill, totalRows int64, err error) {
	resultOrm := r.mysql.Model(&model.RecurringBill{})

	resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Find(&results).Error; err != nil {
		err = ErrNotFound
		return nil, -1, err
	}

	return results, totalRows, nil
}
func (r *mysqlDBRepository) CreateRecurringBill(record *model.RecurringBill) (results *model.RecurringBill, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) UpdateRecurringBill(recurringBillId int32, updated *model.RecurringBill) (result *model.RecurringBill, RowsAffected int64, err error) {
	result = &model.RecurringBill{}
	db := r.mysql.First(result, recurringBillId)
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
