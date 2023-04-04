package entity

type GetBalanceAmountResp struct {
	Amount float64 `json:"net_amount"`
}

type UpdateBalanceAmountReq struct {
	Amount     float64 `json:"balance_amount"`
	Notes      *string `json:"balance_notes"`
	Attachment *string `json:"balance_attachment"`
	DateAdded  string  `json:"balance_date_added"`
}
