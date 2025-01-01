package v1

import (
	"api-gateway/grpc_proto/user_info"
	"api-gateway/internal/dtos"
	"api-gateway/internal/endpoints/api"
	"api-gateway/pkg/app"
	"api-gateway/pkg/logger"
	"api-gateway/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	userInfoService user_info.UserInfoServiceClient
}

func NewUserInfo(userInfoService user_info.UserInfoServiceClient) *UserInfo {
	return &UserInfo{
		userInfoService: userInfoService,
	}
}

func (u *UserInfo) Register(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form dtos.RegisterForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "建立失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "建立失敗", errors, nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := u.userInfoService.UserRegister(ctx, &user_info.CreateUserRequest{
		Account:  form.Account,
		Email:    form.Email,
		Password: form.Password,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	appG.Response(http.StatusOK, true, "建立成功", nil, map[string]string{"id": response.GetId()})
}

func (u *UserInfo) UpdateUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form dtos.UpdateUserForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	}

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "更新失敗", map[string]string{"發生錯誤": "c.Get(jwtClaims) 不存在"}, nil)
		return
	}

	claims := jwtClaims.(*utils.Claims)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := u.userInfoService.UpdateUser(ctx, &user_info.UpdateUserRequest{
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

func (u *UserInfo) UpdateUserPassword(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form dtos.UpdateUserPasswordForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	}

	jwtClaims, exist := c.Get("jwtClaims")
	if !exist {
		appG.Response(http.StatusInternalServerError, false, "更新失敗", map[string]string{"發生錯誤": "c.Get(jwtClaims) 不存在"}, nil)
		return
	}

	claims := jwtClaims.(*utils.Claims)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := u.userInfoService.UpdateUserPassword(ctx, &user_info.UpdateUserPasswordRequest{
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

func (u *UserInfo) GetCurrentUserInfo(c *gin.Context) {
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

	response, err := u.userInfoService.GetUserDeteil(ctx, &user_info.GetUserRequest{
		Id: claims.Subject,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	user := dtos.GetCurrentUserInfoResponse{
		ID:      response.GetUser().GetId(),
		Account: response.GetUser().GetAccount(),
		Email:   response.GetUser().GetEmail(),
		UserInfo: dtos.UserInfoData{
			Name:    response.GetUser().GetUserInfo().GetName(),
			Country: response.GetUser().GetUserInfo().GetCountry(),
			Points:  response.GetUser().GetUserInfo().GetPoints(),
			Hp:      response.GetUser().GetUserInfo().GetHp(),
		},
	}

	appG.Response(http.StatusOK, true, "取得資料成功", nil, user)

}

func (u *UserInfo) GetUserInfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := u.userInfoService.GetUserDeteil(ctx, &user_info.GetUserRequest{
		Id: c.Param("id"),
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	user := dtos.GetUserInfoResponse{
		ID: response.GetUser().GetId(),
		UserInfo: dtos.UserInfoData{
			Name:    response.GetUser().GetUserInfo().GetName(),
			Country: response.GetUser().GetUserInfo().GetCountry(),
			Points:  response.GetUser().GetUserInfo().GetPoints(),
			Hp:      response.GetUser().GetUserInfo().GetHp(),
		},
	}

	appG.Response(http.StatusOK, true, "取得資料成功", nil, user)

}
