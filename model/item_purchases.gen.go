// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameItemPurchase = "item_purchases"

// ItemPurchase mapped from table <item_purchases>
type ItemPurchase struct {
	ItemPurchaseID       int32     `gorm:"column:item_purchase_id;primaryKey;autoIncrement:true" json:"item_purchase_id"`
	ItemID               int32     `gorm:"column:item_id;not null" json:"item_id"`
	BillID               int32     `gorm:"column:bill_id;not null" json:"bill_id"`
	ItemPurchaseQty      int32     `gorm:"column:item_purchase_qty;not null" json:"item_purchase_qty"`
	ItemPurchaseTime     time.Time `gorm:"column:item_purchase_time;not null" json:"item_purchase_time"`
	ItemPurchaseDiscount *int32    `gorm:"column:item_purchase_discount" json:"item_purchase_discount"`
}

// TableName ItemPurchase's table name
func (*ItemPurchase) TableName() string {
	return TableNameItemPurchase
}
