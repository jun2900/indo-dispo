package entity

type AddBillReq struct {
	SupplierId    int32          `json:"supplier_id"`
	StartDate     string         `json:"bill_start_date"`
	DueDate       string         `json:"bill_due_date"`
	BillType      string         `json:"bill_type"`
	Attachments   []Attachment   `json:"bill_attachments"`
	Items         []ItemPurchase `json:"bill_items"`
	BankName      *string        `json:"bill_bank_name"`
	AccountNumber *string        `json:"bill_account_number"`
	ShippingCost  float64        `json:"bill_shipping_cost"`
	BillNote      *string        `json:"bill_notes"`
}

type UpdateBillReq struct {
	SupplierId    int32          `json:"supplier_id"`
	StartDate     string         `json:"bill_start_date"`
	DueDate       string         `json:"bill_due_date"`
	Attachments   []Attachment   `json:"bill_attachments"`
	Items         []ItemPurchase `json:"bill_items"`
	BankName      *string        `json:"bill_bank_name"`
	AccountNumber *string        `json:"bill_account_number"`
	ShippingCost  float64        `json:"bill_shipping_cost"`
	BillNote      *string        `json:"bill_notes"`
}

type ItemPurchase struct {
	ItemId       int32    `json:"item_id"`
	ItemQty      int32    `json:"item_qty"`
	ItemDiscount *float64 `json:"item_discount"`
	ItemPpn      bool     `json:"item_ppn"`
}

type BillDetailsResp struct {
	StartDate        string       `json:"bill_start_date"`
	DueDate          string       `json:"bill_due_date"`
	BillNumber       string       `json:"bill_number"`
	BillOrderNumber  *string      `json:"bill_order_number"`
	BillType         string       `json:"bill_type"`
	Attachments      []Attachment `json:"bill_attachments"`
	Items            []ItemBill   `json:"bill_items"`
	BillStatus       string       `json:"bill_status"`
	BillSubTotal     int64        `json:"bill_subtotal"`
	BillTotal        int64        `json:"bill_total"`
	BillShippingCost float64      `json:"bill_shipping_cost"`
}

type ItemBill struct {
	Name        string  `json:"item_name"`
	Description *string `json:"item_description"`
	Qty         int32   `json:"item_qty"`
	Price       float64 `json:"item_price"`
	Amount      float64 `json:"item_amount"`
	ItemPpn     bool    `json:"item_ppn"`
}

type BillHeaderResp struct {
	Overdue   float64 `json:"bill_overdue"`
	Open      float64 `json:"bill_open"`
	BillDraft float64 `json:"bill_draft"`
}

type BillUpdateStatusReq struct {
	Status string `json:"bill_status"`
}
