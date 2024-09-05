package dto

type InviteMerchantUserDto struct {
	Email      string `json:"email"`
	RolesId    int    `json:"rolesId"`
	MerchantId string `json:"merchantId"`
	Username   string
}

type EmailDataHtmlDto struct {
	Username string
	Password string `json:"password"`
	Pin      string `json:"pin"`
}

type UpdatePassOrPinDto struct {
	Username string
	Password *string `json:"password"`
	Pin      *string `json:"pin"`
}
