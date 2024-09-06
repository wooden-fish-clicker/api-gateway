package v1

import (
	"api-gateway/grpc/user_auth"
	"api-gateway/internal/api"
	"api-gateway/pkg/app"
	"api-gateway/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LoginForm struct {
	Account  string `json:"account" valid:"Required;MaxSize(100)"`
	Password string `json:"password" valid:"Required;MinSize(8);MaxSize(100)"`
}

type LineLoginForm struct {
	Code string `json:"code" valid:"Required;MaxSize(100)"`
}

var UserAuthService user_auth.UserAuthServiceClient

func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LoginForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "登入失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "登入失敗", errors, nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserAuthService.Login(ctx, &user_auth.LoginRequest{
		Account:  form.Account,
		Password: form.Password,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "登入成功", nil, map[string]string{"token": response.GetToken()})

}

func LineLogin(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LineLoginForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "登入失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "登入失敗", errors, nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserAuthService.LineLogin(ctx, &user_auth.LineloginRequest{
		Code: form.Code,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "登入成功", nil, map[string]string{"token": response.GetToken()})
}

func Logout(c *gin.Context) {
	var appG = app.Gin{C: c}

	jwtToken, exist := c.Get("jwtToken")
	if !exist {
		logger.Error("c.Get(jwtToken)不存在")
		appG.Response(http.StatusOK, true, "登出成功", nil, nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("authorization", jwtToken.(string))
	ctx = metadata.NewOutgoingContext(ctx, md)
	empty := new(emptypb.Empty)
	_, err := UserAuthService.Logout(ctx, empty)
	if err != nil {
		logger.Error(err.Error())
		appG.Response(http.StatusOK, true, "登出成功", nil, nil)
		return
	}

	appG.Response(http.StatusOK, true, "登出成功", nil, nil)
}
