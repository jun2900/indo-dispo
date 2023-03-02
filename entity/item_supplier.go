package entity

type ListItemBySupplierResp struct {
	PurchasePrice int32   `json:"item_purchase_price"`
	SellPrice     int32   `json:"item_sell_price"`
	Unit          string  `json:"item_unit"`
	Name          string  `json:"item_name"`
	Description   *string `json:"item_description"`
}
