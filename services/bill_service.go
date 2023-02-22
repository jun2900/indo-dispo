package services

import (
	"anara/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BillService interface {
	GetBill(billId int32) (result *model.Bill, RowsAffected int64, err error)
	GetAllBill(page, pagesize int, order string, dueDate time.Time, status string, vendor string) (results []model.VSupplierBill, totalRows int64, err error)
	CreateBill(record *model.Bill) (result *model.Bill, RowsAffected int64, err error)
	GetAllDraftBillTotal() (billTotal int32)
	GetAllOverdueBillTotal() (billTotal int32)
	GetAllOpenBillTotal() (billTotal int32)
}

func NewBillService(mysqlConnection *gorm.DB) BillService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) GetBill(billId int32) (result *model.Bill, RowsAffected int64, err error) {
	db := r.mysql.First(&result, billId)
	if err = db.Error; err != nil {
		return nil, -1, err
	}
	return result, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetAllBill(page, pagesize int, order string, dueDate time.Time, status string, vendor string) (results []model.VSupplierBill, totalRows int64, err error) {
	resultOrm := r.mysql.Model(&model.VSupplierBill{})
	if len(status) > 0 {
		resultOrm = resultOrm.Where("bill_status = ?", status)
	}
	if len(vendor) > 0 {
		resultOrm = resultOrm.Where("supplier_name LIKE ?", fmt.Sprint("%", vendor, "%"))
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
		return nil, -1, err
	}

	return results, totalRows, nil
}

func (r *mysqlDBRepository) CreateBill(record *model.Bill) (result *model.Bill, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetAllDraftBillTotal() (billTotal int32) {
	r.mysql.Model(&model.Bill{}).Where("bill_status = ?", "DRAFT").Select("sum(bill_total)").Row().Scan(&billTotal)
	return billTotal
}

func (r *mysqlDBRepository) GetAllOverdueBillTotal() (billTotal int32) {
	now := time.Now()
	r.mysql.Model(&model.Bill{}).Where("bill_due_date <= ?", now).Not("bill_status IN ?", []string{"PAID", "DRAFT", "CANCELLED"}).Select("sum(bill_total)").Row().Scan(&billTotal)
	return billTotal
}

func (r *mysqlDBRepository) GetAllOpenBillTotal() (billTotal int32) {
	now := time.Now()
	r.mysql.Model(&model.Bill{}).Where("bill_due_date >= ?", now).Not("bill_status IN ?", []string{"PAID", "DRAFT", "CANCELLED"}).Select("sum(bill_total)").Row().Scan(&billTotal)
	return billTotal
}
