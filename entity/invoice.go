package entity

type AddInvoiceReq struct {
	CustomerId    int32          `json:"customer_id"`
	StartDate     string         `json:"invoice_start_date"`
	DueDate       string         `json:"invoice_due_date"`
	Attachments   []Attachment   `json:"invoice_attachments"`
	Items         []ItemPurchase `json:"invoice_items"`
	BankName      *string        `json:"invoice_bank_name"`
	AccountNumber *string        `json:"invoice_account_number"`
	ShippingCost  float64        `json:"invoice_shipping_cost"`
	Notes         *string        `json:"invoice_notes"`
	Logo          *[]byte        `json:"invoice_logo"`
	Title         *string        `json:"invoice_title"`
	Subheading    *string        `json:"invoice_subheading"`
}

type InvoiceDetailResp struct {
	StartDate          string       `json:"invoice_start_date"`
	DueDate            string       `json:"invoice_due_date"`
	InvoiceNumber      string       `json:"invoice_number"`
	InvoiceOrderNumber *string      `json:"invoice_order_number"`
	Attachments        []Attachment `json:"invoice_attachments"`
	Items              []ItemBill   `json:"invoice_items"`
	InvoiceStatus      string       `json:"invoice_status"`
	InvoiceSubTotal    int64        `json:"invoice_subtotal"`
	InvoiceTotal       int64        `json:"invoice_total"`
	Logo               *[]byte      `json:"invoice_logo"`
	Title              *string      `json:"invoice_title"`
	Subheading         *string      `json:"invoice_subheading"`
	ShippingCost       float64      `json:"invoice_shipping_cost"`
	BankName           *string      `json:"invoice_bank_name"`
	Notes              *string      `json:"invoice_notes"`
}

type InvoiceHeaderResp struct {
	Overdue float64 `json:"invoice_overdue"`
	Open    float64 `json:"invoice_open"`
}

type InvoiceUpdateStatusReq struct {
	Status string `json:"invoice_status"`
}

type UpdateInvoiceReq struct {
	SupplierId    int32          `json:"customer_id"`
	StartDate     string         `json:"invoice_start_date"`
	DueDate       string         `json:"invoice_due_date"`
	Attachments   []Attachment   `json:"invoice_attachments"`
	Items         []ItemPurchase `json:"invoice_items"`
	BankName      *string        `json:"invoice_bank_name"`
	AccountNumber *string        `json:"invoice_account_number"`
	ShippingCost  float64        `json:"invoice_shipping_cost"`
	Logo          *[]byte        `json:"invoice_logo"`
	Title         *string        `json:"invoice_title"`
	Subheading    *string        `json:"invoice_subheading"`
}
