package entity

type AddRecurringBillReq struct {
	SupplierId    int32          `json:"supplier_id"`
	Frequency     string         `json:"recurring_bill_frequency"`
	Items         []ItemPurchase `json:"recurring_bill_items"`
	BankName      string         `json:"recurring_bill_bank_name"`
	AccountNumber string         `json:"recurring_bill_account_number"`
	ShippingCost  float64        `json:"recurring_bill_shipping_cost"`
	Notes         *string        `json:"recurring_bill_notes"`
	StartDate     string         `json:"recurring_bill_start_date"`
	EndDate       *string        `json:"recurring_bill_end_date"`
	Attachments   []Attachment   `json:"bill_attachments"`
	PaymentDue    int            `json:"recurring_bill_payment_due"`
}
