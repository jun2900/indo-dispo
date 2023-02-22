package entity

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	AccessToken string `json:"access_token"`
}
