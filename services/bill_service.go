package services

import (
	"anara/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BillService interface {
	GetBill(billId int32) (result *model.Bill, RowsAffected int64, err error)
	GetAllBill(page, pagesize int, order string, dueDate time.Time, status string, vendor, billType string) (results []model.VSupplierBill, totalRows int64, err error)
	CreateBill(record *model.Bill) (result *model.Bill, RowsAffected int64, err error)
	GetAllPaidAndOperasionalBillTotal() (billTotal float64)
	GetAllMenungguPembayaranBillTotalWithBillType(billType string) (billTotal float64)
	GetAllOverdueBillTotalWithBillType(billType string) (billTotal float64)
	GetAllOpenBillTotal() (billTotal float64)
	UpdateBillStatus(billId int32, billStatus string) (result *model.Bill, RowsAffected int64, err error)
	UpdateBill(itemId int32, updated *model.Bill) (result *model.Bill, RowsAffected int64, err error)
	DeleteBill(billId int32) (err error)
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

func (r *mysqlDBRepository) GetAllBill(page, pagesize int, order string, dueDate time.Time, status string, vendor, billType string) (results []model.VSupplierBill, totalRows int64, err error) {
	resultOrm := r.mysql.Model(&model.VSupplierBill{})
	if len(status) > 0 {
		resultOrm = resultOrm.Where("bill_status = ?", status)
	}
	if len(vendor) > 0 {
		resultOrm = resultOrm.Where("supplier_name LIKE ?", fmt.Sprint("%", vendor, "%"))
	}
	if len(billType) > 0 {
		resultOrm = resultOrm.Where("bill_type = ?", billType)
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

func (r *mysqlDBRepository) GetAllMenungguPembayaranBillTotalWithBillType(billType string) (billTotal float64) {
	now := time.Now()
	r.mysql.Model(&model.Bill{}).Where("bill_status = ? AND bill_type = ? AND bill_due_date >= ?", "MENUNGGU PEMBAYARAN", billType, now).Select("sum(bill_total)").Row().Scan(&billTotal)
	return billTotal
}

func (r *mysqlDBRepository) GetAllOverdueBillTotalWithBillType(billType string) (billTotal float64) {
	now := time.Now()
	r.mysql.Model(&model.Bill{}).Where("bill_due_date < ? AND bill_type = ?", now, billType).Not("bill_status IN ?", []string{"SUDAH DIBAYAR", "CANCELLED"}).Select("sum(bill_total)").Row().Scan(&billTotal)
	return billTotal
}

func (r *mysqlDBRepository) GetAllOpenBillTotal() (billTotal float64) {
	now := time.Now()
	r.mysql.Model(&model.Bill{}).Where("bill_due_date >= ?", now).Not("bill_status IN ?", []string{"SUDAH DIBAYAR", "MENUNGGU PEMBAYARAN", "CANCELLED"}).Select("sum(bill_total)").Row().Scan(&billTotal)
	return billTotal
}

func (r *mysqlDBRepository) GetAllPaidAndOperasionalBillTotal() (billTotal float64) {
	r.mysql.Model(&model.Bill{}).Where("bill_status = ? AND bill_type = ?", "SUDAH DIBAYAR", "operasional").Select("sum(bill_total)").Row().Scan(&billTotal)
	return billTotal
}

func (r *mysqlDBRepository) UpdateBillStatus(billId int32, billStatus string) (result *model.Bill, RowsAffected int64, err error) {
	result = &model.Bill{}
	db := r.mysql.First(&result, billId)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	result.BillStatus = billStatus

	if err := db.Save(&result).Error; err != nil {
		return nil, -1, err
	}

	return result, db.RowsAffected, nil
}

func (r *mysqlDBRepository) DeleteBill(billId int32) (err error) {
	if err := r.mysql.Delete(&model.Bill{}, billId).Error; err != nil {
		return err
	}
	return nil
}

func (r *mysqlDBRepository) UpdateBill(itemId int32, updated *model.Bill) (result *model.Bill, RowsAffected int64, err error) {
	result = &model.Bill{}
	db := r.mysql.First(result, itemId)
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
