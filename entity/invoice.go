package entity

type AddInvoiceReq struct {
	CustomerId    int32          `json:"customer_id"`
	StartDate     string         `json:"invoice_start_date"`
	DueDate       string         `json:"invoice_due_date"`
	BillType      string         `json:"invoice_type"`
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
