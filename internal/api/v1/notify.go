package v1

import (
	"api-gateway/grpc/notification"
	"api-gateway/internal/api"
	"api-gateway/pkg/app"
	"api-gateway/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReadNotificationForm struct {
	IDs []string `json:"ids" valid:"Required"`
}

var NotificationService notification.NotificationServiceClient

func ReadNotification(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ReadNotificationForm
	)

	httpCode, errors := app.Valid(c, &form, false)
	api.HandleValidError(httpCode, errors, &appG)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := NotificationService.ReadNotification(ctx, &notification.ReadNotificationRequest{
		Ids: form.IDs,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}
	appG.Response(http.StatusOK, true, "更新成功", nil, nil)
}

func DeleteNotification(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := NotificationService.DeleteNotification(ctx, &notification.DeleteNotificationRequest{
		Id: c.Param("id"),
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}
	appG.Response(http.StatusOK, true, "刪除成功", nil, nil)
}

type GetNotificationListResponse struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	Type      int32  `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetNotificationList(c *gin.Context) {
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

	response, err := NotificationService.GetNotificationList(ctx, &notification.GetNotificationListRequest{
		UserId: claims.Subject,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	notifications := response.GetNotification()

	appG.Response(http.StatusOK, true, "取得資料成功", nil, notifications)
}
