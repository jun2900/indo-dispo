package services

import (
	"anara/model"
	"fmt"

	"gorm.io/gorm"
)

type SupplierService interface {
	GetAllSupplierByType(supplierType string, page, pagesize int, order, name, email, address string) (result []*model.Supplier, totalRows int64, err error)
	CreateSupplier(record *model.Supplier) (result *model.Supplier, RowsAffected int64, err error)
	GetSupplier(supplierId int32) (result *model.Supplier, RowsAffected int64, err error)
}

func NewSupplierService(mysqlConnection *gorm.DB) SupplierService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) GetAllSupplierByType(supplierType string, page, pagesize int, order, name, email, address string) (result []*model.Supplier, totalRows int64, err error) {
	resultOrm := r.mysql.Model(&model.Supplier{})

	if len(name) > 0 {
		resultOrm = resultOrm.Where("supplier_name LIKE ?", fmt.Sprint("%", name, "%"))
	}
	if len(email) > 0 {
		resultOrm = resultOrm.Where("supplier_email LIKE ?", fmt.Sprint("%", email, "%"))
	}
	if len(address) > 0 {
		resultOrm = resultOrm.Where("supplier_address LIKE ?", fmt.Sprint("%", address, "%"))
	}

	resultOrm.Where("supplier_type = ?", supplierType).Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if err = resultOrm.Where("supplier_type = ?", supplierType).Find(&result).Error; err != nil {
		return nil, -1, err
	}

	return result, totalRows, nil

}

func (r *mysqlDBRepository) CreateSupplier(record *model.Supplier) (result *model.Supplier, RowsAffected int64, err error) {
	db := r.mysql.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetSupplier(supplierId int32) (result *model.Supplier, RowsAffected int64, err error) {
	db := r.mysql.First(&result, supplierId)
	if err = db.Error; err != nil {
		return nil, -1, err
	}
	return result, db.RowsAffected, nil
}
