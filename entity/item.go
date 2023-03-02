package entity

type AddItemReq struct {
	Name          string       `json:"item_name"`
	Description   *string      `json:"item_description"`
	PurchasePrice int32        `json:"item_purchase_price"`
	SellPrice     int32        `json:"item_sell_price"`
	SupplierId    int32        `json:"supplier_id"`
	Unit          string       `json:"item_unit"`
	WholeSalers   []WholeSaler `json:"item_wholesalers"`
}

type WholeSaler struct {
	Qty   int32 `json:"wholesaler_qty"`
	Price int32 `json:"wholesaler_price"`
}
