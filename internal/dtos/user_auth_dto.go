package dtos

type LoginForm struct {
	Account  string `json:"account" valid:"Required;MaxSize(100)"`
	Password string `json:"password" valid:"Required;MinSize(8);MaxSize(100)"`
}

type LineLoginForm struct {
	Code string `json:"code" valid:"Required;MaxSize(100)"`
}
