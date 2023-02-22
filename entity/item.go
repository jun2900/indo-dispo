package entity

type AddItemReq struct {
	Name          string  `json:"item_name"`
	Description   *string `json:"item_description"`
	PurchasePrice *int32  `json:"item_purchase_price"`
	SellPrice     *int32  `json:"item_sell_price"`
}
