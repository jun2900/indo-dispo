// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameInvoice = "invoices"

// Invoice mapped from table <invoices>
type Invoice struct {
	InvoicesID           int32     `gorm:"column:invoices_id;primaryKey;autoIncrement:true" json:"invoices_id"`
	SupplierID           int32     `gorm:"column:supplier_id;not null" json:"supplier_id"`
	InvoiceStartDate     time.Time `gorm:"column:invoice_start_date;not null" json:"invoice_start_date"`
	InvoiceDueDate       time.Time `gorm:"column:invoice_due_date;not null" json:"invoice_due_date"`
	InvoiceNumber        string    `gorm:"column:invoice_number;not null" json:"invoice_number"`
	InvoiceOrderNumber   *string   `gorm:"column:invoice_order_number" json:"invoice_order_number"`
	InvoiceTitle         *string   `gorm:"column:invoice_title" json:"invoice_title"`
	InvoiceSubheading    *string   `gorm:"column:invoice_subheading" json:"invoice_subheading"`
	InvoiceLogo          *[]byte   `gorm:"column:invoice_logo" json:"invoice_logo"`
	InvoiceShippingCost  float64   `gorm:"column:invoice_shipping_cost;not null" json:"invoice_shipping_cost"`
	InvoiceAccountNumber *string   `gorm:"column:invoice_account_number" json:"invoice_account_number"`
	InvoiceBankName      *string   `gorm:"column:invoice_bank_name" json:"invoice_bank_name"`
	InvoiceTotal         float64   `gorm:"column:invoice_total;not null" json:"invoice_total"`
	InvoiceStatus        string    `gorm:"column:invoice_status;not null" json:"invoice_status"`
}

// TableName Invoice's table name
func (*Invoice) TableName() string {
	return TableNameInvoice
}
