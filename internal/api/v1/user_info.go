package v1

import (
	"api-gateway/grpc/user_info"
	"api-gateway/internal/api"
	"api-gateway/pkg/app"
	"api-gateway/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

var UserInfoService user_info.UserInfoServiceClient

func Register(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form RegisterForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	api.HandleValidError(httpCode, errors, &appG)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserInfoService.UserRegister(ctx, &user_info.CreateUserRequest{
		User: &user_info.User{
			Account:  form.Account,
			Email:    form.Email,
			Password: form.Password,
		},
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "建立成功", nil, map[string]string{"id": response.GetId()})
}

func UpdateUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form UpdateUserForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	api.HandleValidError(httpCode, errors, &appG)

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "更新失敗", map[string]string{"發生錯誤": "c.Get(jwtClaims) 不存在"}, nil)
		return
	}

	claims := jwtClaims.(*utils.Claims)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := UserInfoService.UpdateUser(ctx, &user_info.UpdateUserRequest{
		User: &user_info.User{
			Id:      claims.Subject,
			Account: form.Account,
			Email:   form.Email,
			UserInfo: &user_info.UserInfo{
				Name:    form.Name,
				Country: form.Country,
			},
		},
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "更新成功", nil, nil)

}

func UpdateUserPassword(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form UpdateUserPasswordForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	api.HandleValidError(httpCode, errors, &appG)

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "更新失敗", map[string]string{"發生錯誤": "c.Get(jwtClaims) 不存在"}, nil)
		return
	}

	claims := jwtClaims.(*utils.Claims)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := UserInfoService.UpdateUserPassword(ctx, &user_info.UpdateUserPasswordRequest{
		Id:          claims.Subject,
		NewPassword: form.NewPassword,
		OldPassword: form.OldPassword,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "更新成功", nil, nil)

}

type GetCurrentUserInfoResponse struct {
	ID       string   `json:"id"`
	Account  string   `json:"account"`
	Email    string   `json:"email"`
	UserInfo UserInfo `json:"user_info"`
}

type UserInfo struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Points  int64  `json:"points"`
	Hp      int32  `json:"hp"`
}

func GetCurrentUserInfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "取得資料成功", map[string]string{"發生錯誤": "c.Get(jwtClaims) 不存在"}, nil)
		return
	}

	claims := jwtClaims.(*utils.Claims)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserInfoService.GetUserDeteil(ctx, &user_info.GetUserRequest{
		Id: claims.Subject,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	user := GetCurrentUserInfoResponse{
		ID:      response.GetUser().GetId(),
		Account: response.GetUser().GetAccount(),
		Email:   response.GetUser().GetEmail(),
		UserInfo: UserInfo{
			Name:    response.GetUser().GetUserInfo().GetName(),
			Country: response.GetUser().GetUserInfo().GetCountry(),
			Points:  response.GetUser().GetUserInfo().GetPoints(),
			Hp:      response.GetUser().GetUserInfo().GetHp(),
		},
	}

	appG.Response(http.StatusOK, true, "取得資料成功", nil, user)

}

type GetUserInfoResponse struct {
	ID       string   `json:"id"`
	UserInfo UserInfo `json:"user_info"`
}

func GetUserInfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserInfoService.GetUserDeteil(ctx, &user_info.GetUserRequest{
		Id: c.Param("id"),
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	user := GetUserInfoResponse{
		ID: response.GetUser().GetId(),
		UserInfo: UserInfo{
			Name:    response.GetUser().GetUserInfo().GetName(),
			Country: response.GetUser().GetUserInfo().GetCountry(),
			Points:  response.GetUser().GetUserInfo().GetPoints(),
			Hp:      response.GetUser().GetUserInfo().GetHp(),
		},
	}

	appG.Response(http.StatusOK, true, "取得資料成功", nil, user)

}
