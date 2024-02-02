package entity

type AddItemReq struct {
	Name          string       `json:"item_name"`
	Code          string       `json:"item_code"`
	Description   *string      `json:"item_description"`
	PurchasePrice float64      `json:"item_purchase_price"`
	SellPrice     float64      `json:"item_sell_price"`
	SupplierId    int32        `json:"supplier_id"`
	Unit          string       `json:"item_unit"`
	WholeSalers   []WholeSaler `json:"item_wholesalers"`
}

type WholeSaler struct {
	Qty   int32   `json:"wholesaler_qty"`
	Price float64 `json:"wholesaler_price"`
}

type ListItemBySupplierResp struct {
	Id            int32            `json:"item_id"`
	Code          string           `json:"item_code"`
	PurchasePrice float64          `json:"item_purchase_price"`
	SellPrice     float64          `json:"item_sell_price"`
	Unit          string           `json:"item_unit"`
	Name          string           `json:"item_name"`
	Description   *string          `json:"item_description"`
	WholeSaler    []ListWholeSaler `json:"wholesalers"`
}

type ListWholeSaler struct {
	Id    int32   `json:"wholesaler_id"`
	Qty   int32   `json:"wholesaler_qty"`
	Price float64 `json:"wholesaler_price"`
}

type UpdateItemReq struct {
	Name          string       `json:"item_name"`
	Description   string       `json:"item_description"`
	PurchasePrice float64      `json:"item_purchase_price"`
	SellPrice     float64      `json:"item_sell_price"`
	WholeSalers   []WholeSaler `json:"item_wholesalers"`
	Unit          string       `json:"item_unit"`
}

type WholesaleForUpdateReq struct {
	Id    int32 `json:"wholesaler_id"`
	Qty   int32 `json:"wholesaler_qty"`
	Price int32 `json:"wholesaler_price"`
}
