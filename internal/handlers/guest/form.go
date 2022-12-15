package guest

type LoginUsersForm struct {
	Username *string `json:"username" form:"username" validate:"required"`
	Password *string `json:"password" form:"password" validate:"required"`
}
