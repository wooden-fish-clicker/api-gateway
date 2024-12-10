package dtos

type ReadNotificationForm struct {
	IDs []string `json:"ids" valid:"Required"`
}
