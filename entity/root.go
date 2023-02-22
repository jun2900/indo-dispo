package entity

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrRespController struct {
	SourceFunction string `json:"source_function"`
	ErrMessage     string `json:"err_message"`
}

type PagedResults struct {
	Page         int         `json:"page"`
	PageSize     int         `json:"page_size"`
	Data         interface{} `json:"data"`
	TotalRecords int         `json:"total_records"`
}
