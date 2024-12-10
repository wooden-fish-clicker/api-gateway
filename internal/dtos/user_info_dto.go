package dtos

type RegisterForm struct {
	Account  string `json:"account" valid:"Required;MaxSize(100)"`
	Email    string `json:"email" valid:"Required;Email;MaxSize(255)"`
	Password string `json:"password" valid:"Required;MinSize(8);MaxSize(100)"`
}

type UpdateUserForm struct {
	Account string `json:"account" valid:"Required;MaxSize(100)"`
	Email   string `json:"email" valid:"Required;Email;MaxSize(255)"`
	Name    string `json:"name" valid:"Required;MaxSize(100)"`
	Country string `json:"country" valid:"Required;MaxSize(100)"`
}

type UpdateUserPasswordForm struct {
	OldPassword string `json:"old_password" valid:"Required;MinSize(8);MaxSize(100)"`
	NewPassword string `json:"new_password" valid:"Required;MinSize(8);MaxSize(100)"`
}

type GetCurrentUserInfoResponse struct {
	ID       string       `json:"id"`
	Account  string       `json:"account"`
	Email    string       `json:"email"`
	UserInfo UserInfoData `json:"user_info"`
}

type GetUserInfoResponse struct {
	ID       string       `json:"id"`
	UserInfo UserInfoData `json:"user_info"`
}

type UserInfoData struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Points  int64  `json:"points"`
	Hp      int32  `json:"hp"`
}
