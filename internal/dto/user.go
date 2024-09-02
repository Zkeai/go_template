package dto

type UserRegisterReq struct {
	Wallet string `json:"wallet"`
}

type UserQueryReq struct {
	Wallet string `form:"wallet"`
}
type UserLoginReq struct {
	*UserRegisterReq
}
