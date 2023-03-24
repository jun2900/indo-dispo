package entity

type SupplierAddReq struct {
	Name        string  `json:"supplier_name"`
	Email       *string `json:"supplier_email"`
	Telephone   *string `json:"supplier_telephone"`
	Web         *string `json:"supplier_web"`
	Npwp        *string `json:"supplier_npwp"`
	Address     *string `json:"supplier_address"`
	Type        string  `json:"supplier_type"`
	Whatsapp    *string `json:"supplier_whatsapp"`
	Description *string `json:"supplier_description"`
	City        *string `json:"supplier_city"`
	State       *string `json:"supplier_state"`
	ZipCode     *string `json:"supplier_zip_code"`
	Country     *string `json:"supplier_country"`
}

type ListSupplierByItemReq struct {
	ItemName string `json:"item_name"`
}

type ListSupplierByItemDetails struct {
	SupplierName  string  `json:"supplier_name"`
	Description   *string `json:"item_description"`
	PurchasePrice float64 `json:"item_purchase_price"`
	SellPrice     float64 `json:"item_sell_price"`
}
type ListSupplierByItemResp struct {
	ItemName string                      `json:"item_name"`
	Data     []ListSupplierByItemDetails `json:"details"`
}
