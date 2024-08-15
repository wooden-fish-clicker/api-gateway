package v1

import (
	"api-gateway/grpc/user_info"
	"api-gateway/pkg/app"
	"api-gateway/pkg/logger"
	"api-gateway/pkg/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterForm struct {
	Account  string `form:"account" valid:"Required;MaxSize(100)"`
	Email    string `form:"email" valid:"Required;Email;MaxSize(255)"`
	Password string `form:"password" valid:"Required;MinSize(8);MaxSize(100)"`
}

type UpdateUserForm struct {
	Account string `form:"account" valid:"Required;MaxSize(100)"`
	Email   string `form:"email" valid:"Required;Email;MaxSize(255)"`
	Name    string `form:"name" valid:"Required;MaxSize(100)"`
	Country string `form:"country" valid:"Required;MaxSize(100)"`
}

type UpdateUserPasswordForm struct {
	OldPassword string `form:"old_password" valid:"Required;MinSize(8);MaxSize(100)"`
	NewPassword string `form:"new_password" valid:"Required;MinSize(8);MaxSize(100)"`
}

var UserInfoService user_info.UserInfoServiceClient

func Register(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form RegisterForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "建立失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "建立失敗", "發生錯誤", nil)
		return
	}

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
		errMsg := strings.Split(err.Error(), " = ")[2]
		if errMsg == "帳號或信箱已存在" {
			appG.Response(http.StatusConflict, false, "建立失敗", err.Error(), nil)
			return
		}
		appG.Response(http.StatusInternalServerError, false, "建立失敗", err.Error(), nil)
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
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "更新失敗", "發生錯誤", nil)
		return
	}

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "更新失敗", "c.Get(jwtClaims) 不存在", nil)
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
		errMsg := strings.Split(err.Error(), " = ")[2]
		if errMsg == "帳號或信箱已存在" {
			appG.Response(http.StatusConflict, false, "更新失敗", err.Error(), nil)
			return
		} else if errMsg == "找不到此id" {
			appG.Response(http.StatusNotFound, false, "更新失敗", err.Error(), nil)
			return
		}
		appG.Response(http.StatusInternalServerError, false, "更新失敗", err.Error(), nil)
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
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "更新失敗", "發生錯誤", nil)
		return
	}

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "更新失敗", "c.Get(jwtClaims) 不存在", nil)
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
		errMsg := strings.Split(err.Error(), " = ")[2]
		if errMsg == "舊密碼錯誤" {
			appG.Response(http.StatusForbidden, false, "更新失敗", err.Error(), nil)
			return
		} else if errMsg == "找不到此id" {
			appG.Response(http.StatusNotFound, false, "更新失敗", err.Error(), nil)
			return
		}
		appG.Response(http.StatusInternalServerError, false, "更新失敗", err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, true, "更新成功", nil, nil)

}

func GetCurrentUserInfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "更新失敗", "c.Get(jwtClaims) 不存在", nil)
		return
	}

	claims := jwtClaims.(*utils.Claims)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserInfoService.GetUserDeteil(ctx, &user_info.GetUserRequest{
		Id: claims.Subject,
	})

	if err != nil {
		errMsg := strings.Split(err.Error(), " = ")[2]
		if errMsg == "找不到此id" {
			appG.Response(http.StatusNotFound, false, "更新失敗", err.Error(), nil)
			return
		}
		appG.Response(http.StatusInternalServerError, false, "更新失敗", err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, true, "更新成功", nil, response)

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
		errMsg := strings.Split(err.Error(), " = ")[2]
		if errMsg == "找不到此id" {
			appG.Response(http.StatusNotFound, false, "更新失敗", err.Error(), nil)
			return
		}
		appG.Response(http.StatusInternalServerError, false, "更新失敗", err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, true, "更新成功", nil, response)

}
