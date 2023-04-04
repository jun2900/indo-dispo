package services

import (
	"time"

	"github.com/jun2900/indo-dispo/model"

	"gorm.io/gorm"
)

type BalanceLogService interface {
	CreateBalanceLog(record *model.BalanceLog) (results *model.BalanceLog, RowsAffected int64, err error)
	GetAllBalanceLog(page, pagesize int, order string, logTimeStart, logTimeEnd time.Time) (results []model.BalanceLog, totalRows int64, err error)
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

func (r *mysqlDBRepository) GetAllBalanceLog(page, pagesize int, order string, logTimeStart, logTimeEnd time.Time) (results []model.BalanceLog, totalRows int64, err error) {
	resultOrm := r.mysql.Model(&model.BalanceLog{})

	if !logTimeStart.IsZero() && !logTimeEnd.IsZero() {
		resultOrm = resultOrm.Where("balance_log_time_added >= ? AND balance_log_time_added < (? + INTERVAL 1 DAY)", logTimeStart, logTimeEnd)
	} else if !logTimeStart.IsZero() && logTimeEnd.IsZero() {
		resultOrm = resultOrm.Where("balance_log_time_added >= ?", logTimeStart)
	} else if logTimeStart.IsZero() && !logTimeEnd.IsZero() {
		resultOrm = resultOrm.Where("balance_log_time_added < (? + INTERVAL 1 DAY)", logTimeEnd)
	}

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
