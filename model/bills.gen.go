// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameBill = "bills"

// Bill mapped from table <bills>
type Bill struct {
	BillID          int32     `gorm:"column:bill_id;primaryKey;autoIncrement:true" json:"bill_id"`
	SupplierID      int32     `gorm:"column:supplier_id;not null" json:"supplier_id"`
	BillStartDate   time.Time `gorm:"column:bill_start_date;not null" json:"bill_start_date"`
	BillDueDate     time.Time `gorm:"column:bill_due_date;not null" json:"bill_due_date"`
	BillNumber      string    `gorm:"column:bill_number;not null" json:"bill_number"`
	BillOrderNumber *string   `gorm:"column:bill_order_number" json:"bill_order_number"`
	BillDiscount    *int32    `gorm:"column:bill_discount" json:"bill_discount"`
	BillTotal       int32     `gorm:"column:bill_total;not null" json:"bill_total"`
	BillStatus      string    `gorm:"column:bill_status;not null" json:"bill_status"`
	BillType        string    `gorm:"column:bill_type;not null" json:"bill_type"`
}

// TableName Bill's table name
func (*Bill) TableName() string {
	return TableNameBill
}
