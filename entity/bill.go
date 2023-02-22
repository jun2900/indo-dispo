package entity

type AddBillReq struct {
	SupplierId      int32          `json:"supplier_id"`
	StartDate       string         `json:"bill_start_date"`
	DueDate         string         `json:"bill_due_date"`
	BillNumber      string         `json:"bill_number"`
	BillOrderNumber *string        `json:"bill_order_number"`
	BillType        string         `json:"bill_type"`
	Attachments     []Attachment   `json:"bill_attachments"`
	Items           []ItemPurchase `json:"bill_items"`
	Discount        *int32         `json:"bill_discount"`
}

type ItemPurchase struct {
	ItemId  int32 `json:"item_id"`
	ItemQty int32 `json:"item_qty"`
}

type BillDetailsResp struct {
	StartDate       string       `json:"bill_start_date"`
	DueDate         string       `json:"bill_due_date"`
	BillNumber      string       `json:"bill_number"`
	BillOrderNumber *string      `json:"bill_order_number"`
	BillType        string       `json:"bill_type"`
	Attachments     []Attachment `json:"bill_attachments"`
	Items           []ItemBill   `json:"bill_items"`
	BillStatus      string       `json:"bill_status"`
	BillSubTotal    int64        `json:"bill_subtotal"`
	BillTotal       int64        `json:"bill_total"`
	BillDiscount    *int32       `json:"bill_discount"`
}

type ItemBill struct {
	Name        string  `json:"item_name"`
	Description *string `json:"item_description"`
	Qty         int32   `json:"item_qty"`
	Price       int32   `json:"item_price"`
	Amount      int32   `json:"item_amount"`
}

type BillHeaderResp struct {
	Overdue   int32 `json:"bill_overdue"`
	Open      int32 `json:"bill_open"`
	BillDraft int32 `json:"bill_draft"`
}
