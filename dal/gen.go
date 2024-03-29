// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q                = new(Query)
	Admin            *admin
	Attachment       *attachment
	Balance          *balance
	BalanceLog       *balanceLog
	Bill             *bill
	Invoice          *invoice
	Item             *item
	ItemPurchase     *itemPurchase
	ItemSell         *itemSell
	RecurringBill    *recurringBill
	Role             *role
	Supplier         *supplier
	VSupplierBill    *vSupplierBill
	VSupplierInvoice *vSupplierInvoice
	Wholesaler       *wholesaler
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Admin = &Q.Admin
	Attachment = &Q.Attachment
	Balance = &Q.Balance
	BalanceLog = &Q.BalanceLog
	Bill = &Q.Bill
	Invoice = &Q.Invoice
	Item = &Q.Item
	ItemPurchase = &Q.ItemPurchase
	ItemSell = &Q.ItemSell
	RecurringBill = &Q.RecurringBill
	Role = &Q.Role
	Supplier = &Q.Supplier
	VSupplierBill = &Q.VSupplierBill
	VSupplierInvoice = &Q.VSupplierInvoice
	Wholesaler = &Q.Wholesaler
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:               db,
		Admin:            newAdmin(db, opts...),
		Attachment:       newAttachment(db, opts...),
		Balance:          newBalance(db, opts...),
		BalanceLog:       newBalanceLog(db, opts...),
		Bill:             newBill(db, opts...),
		Invoice:          newInvoice(db, opts...),
		Item:             newItem(db, opts...),
		ItemPurchase:     newItemPurchase(db, opts...),
		ItemSell:         newItemSell(db, opts...),
		RecurringBill:    newRecurringBill(db, opts...),
		Role:             newRole(db, opts...),
		Supplier:         newSupplier(db, opts...),
		VSupplierBill:    newVSupplierBill(db, opts...),
		VSupplierInvoice: newVSupplierInvoice(db, opts...),
		Wholesaler:       newWholesaler(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Admin            admin
	Attachment       attachment
	Balance          balance
	BalanceLog       balanceLog
	Bill             bill
	Invoice          invoice
	Item             item
	ItemPurchase     itemPurchase
	ItemSell         itemSell
	RecurringBill    recurringBill
	Role             role
	Supplier         supplier
	VSupplierBill    vSupplierBill
	VSupplierInvoice vSupplierInvoice
	Wholesaler       wholesaler
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:               db,
		Admin:            q.Admin.clone(db),
		Attachment:       q.Attachment.clone(db),
		Balance:          q.Balance.clone(db),
		BalanceLog:       q.BalanceLog.clone(db),
		Bill:             q.Bill.clone(db),
		Invoice:          q.Invoice.clone(db),
		Item:             q.Item.clone(db),
		ItemPurchase:     q.ItemPurchase.clone(db),
		ItemSell:         q.ItemSell.clone(db),
		RecurringBill:    q.RecurringBill.clone(db),
		Role:             q.Role.clone(db),
		Supplier:         q.Supplier.clone(db),
		VSupplierBill:    q.VSupplierBill.clone(db),
		VSupplierInvoice: q.VSupplierInvoice.clone(db),
		Wholesaler:       q.Wholesaler.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:               db,
		Admin:            q.Admin.replaceDB(db),
		Attachment:       q.Attachment.replaceDB(db),
		Balance:          q.Balance.replaceDB(db),
		BalanceLog:       q.BalanceLog.replaceDB(db),
		Bill:             q.Bill.replaceDB(db),
		Invoice:          q.Invoice.replaceDB(db),
		Item:             q.Item.replaceDB(db),
		ItemPurchase:     q.ItemPurchase.replaceDB(db),
		ItemSell:         q.ItemSell.replaceDB(db),
		RecurringBill:    q.RecurringBill.replaceDB(db),
		Role:             q.Role.replaceDB(db),
		Supplier:         q.Supplier.replaceDB(db),
		VSupplierBill:    q.VSupplierBill.replaceDB(db),
		VSupplierInvoice: q.VSupplierInvoice.replaceDB(db),
		Wholesaler:       q.Wholesaler.replaceDB(db),
	}
}

type queryCtx struct {
	Admin            *adminDo
	Attachment       *attachmentDo
	Balance          *balanceDo
	BalanceLog       *balanceLogDo
	Bill             *billDo
	Invoice          *invoiceDo
	Item             *itemDo
	ItemPurchase     *itemPurchaseDo
	ItemSell         *itemSellDo
	RecurringBill    *recurringBillDo
	Role             *roleDo
	Supplier         *supplierDo
	VSupplierBill    *vSupplierBillDo
	VSupplierInvoice *vSupplierInvoiceDo
	Wholesaler       *wholesalerDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Admin:            q.Admin.WithContext(ctx),
		Attachment:       q.Attachment.WithContext(ctx),
		Balance:          q.Balance.WithContext(ctx),
		BalanceLog:       q.BalanceLog.WithContext(ctx),
		Bill:             q.Bill.WithContext(ctx),
		Invoice:          q.Invoice.WithContext(ctx),
		Item:             q.Item.WithContext(ctx),
		ItemPurchase:     q.ItemPurchase.WithContext(ctx),
		ItemSell:         q.ItemSell.WithContext(ctx),
		RecurringBill:    q.RecurringBill.WithContext(ctx),
		Role:             q.Role.WithContext(ctx),
		Supplier:         q.Supplier.WithContext(ctx),
		VSupplierBill:    q.VSupplierBill.WithContext(ctx),
		VSupplierInvoice: q.VSupplierInvoice.WithContext(ctx),
		Wholesaler:       q.Wholesaler.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
