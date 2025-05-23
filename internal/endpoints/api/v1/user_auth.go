package v1

import (
	"api-gateway/grpc_proto/user_auth"
	"api-gateway/internal/dtos"
	"api-gateway/internal/endpoints/api"
	"api-gateway/pkg/app"
	"api-gateway/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserAuth struct {
	userAuthService user_auth.UserAuthServiceClient
}

func NewUserAuth(userAuthService user_auth.UserAuthServiceClient) *UserAuth {
	return &UserAuth{
		userAuthService: userAuthService,
	}
}

func (u *UserAuth) Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form dtos.LoginForm
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

	response, err := u.userAuthService.Login(ctx, &user_auth.LoginRequest{
		Account:  form.Account,
		Password: form.Password,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "登入成功", nil, map[string]string{"token": response.GetToken()})

}

func (u *UserAuth) LineLogin(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form dtos.LineLoginForm
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

	response, err := u.userAuthService.LineLogin(ctx, &user_auth.LineloginRequest{
		Code: form.Code,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "登入成功", nil, map[string]string{"token": response.GetToken()})
}

func (u *UserAuth) Logout(c *gin.Context) {
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
	_, err := u.userAuthService.Logout(ctx, empty)
	if err != nil {
		logger.Error(err.Error())
		appG.Response(http.StatusOK, true, "登出成功", nil, nil)
		return
	}

	appG.Response(http.StatusOK, true, "登出成功", nil, nil)
}
