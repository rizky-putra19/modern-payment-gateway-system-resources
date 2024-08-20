package dto

type InviteMerchantUserDto struct {
	Email      string `json:"email"`
	RolesId    int    `json:"rolesId"`
	MerchantId string `json:"merchantId"`
}

type EmailDataHtmlDto struct {
	Username string
	Password string
	Pin      string
}
