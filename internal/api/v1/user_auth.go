package v1

import (
	"api-gateway/grpc/user_auth"
	"api-gateway/pkg/app"
	"api-gateway/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LoginForm struct {
	Account  string `form:"account" valid:"Required;MaxSize(100)"`
	Password string `form:"password" valid:"Required;MinSize(8);MaxSize(100)"`
}

type LineLoginForm struct {
	Code string `form:"code" valid:"Required;MaxSize(100)"`
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
		appG.Response(httpCode, false, "登入失敗", "發生錯誤", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserAuthService.Login(ctx, &user_auth.LoginRequest{
		Account:  form.Account,
		Password: form.Password,
	})

	if err != nil {
		// 使用 status 包解析錯誤
		st, ok := status.FromError(err)
		if ok {
			code := st.Code()

			if code == codes.Unauthenticated {
				appG.Response(http.StatusUnauthorized, false, "登入失敗", st.Message(), nil)
				return
			}
			logger.Error(err.Error())
			appG.Response(http.StatusInternalServerError, false, "登入失敗", "發生錯誤", nil)
			return
		} else {
			// 不是 gRPC 錯誤，可能是其他錯誤
			logger.Error("status.FromError 發生錯誤")
			appG.Response(http.StatusInternalServerError, false, "登入失敗", "發生錯誤", nil)
			return
		}
	}
	fmt.Println(4)
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
		appG.Response(httpCode, false, "登入失敗", "發生錯誤", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := UserAuthService.LineLogin(ctx, &user_auth.LineloginRequest{
		Code: form.Code,
	})

	if err != nil {
		// 使用 status 包解析錯誤
		st, ok := status.FromError(err)
		if ok {
			appG.Response(http.StatusInternalServerError, false, "登入失敗", st.Message(), nil)
			return
		} else {
			logger.Error("status.FromError 發生錯誤")
			appG.Response(http.StatusInternalServerError, false, "登入失敗", "發生錯誤", nil)
			return
		}
	}

	appG.Response(http.StatusOK, true, "登入成功", nil, map[string]string{"token": response.GetToken()})
}

func Logout(c *gin.Context) {
	var appG = app.Gin{C: c}
	fmt.Println(1)
	jwtToken, exist := c.Get("jwtToken")
	if !exist {
		logger.Error("c.Get(jwtToken)不存在")
		appG.Response(http.StatusOK, true, "登出成功", nil, nil)
		return
	}
	fmt.Println(jwtToken.(string))

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
