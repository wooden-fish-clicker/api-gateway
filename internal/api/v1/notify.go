package v1

import (
	"api-gateway/grpc/notification"
	"api-gateway/internal/api"
	"api-gateway/internal/dtos"
	"api-gateway/pkg/app"
	"api-gateway/pkg/logger"
	"api-gateway/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Notify struct {
	notificationService notification.NotificationServiceClient
}

func NewNotify(notificationService notification.NotificationServiceClient) *Notify {
	return &Notify{
		notificationService: notificationService,
	}
}

func (n *Notify) ReadNotification(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form dtos.ReadNotificationForm
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := n.notificationService.ReadNotification(ctx, &notification.ReadNotificationRequest{
		Ids: form.IDs,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}
	appG.Response(http.StatusOK, true, "更新成功", nil, nil)
}

func (n *Notify) DeleteNotification(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := n.notificationService.DeleteNotification(ctx, &notification.DeleteNotificationRequest{
		Id: c.Param("id"),
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}
	appG.Response(http.StatusOK, true, "刪除成功", nil, nil)
}

func (n *Notify) GetNotificationList(c *gin.Context) {
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

	response, err := n.notificationService.GetNotificationList(ctx, &notification.GetNotificationListRequest{
		UserId: claims.Subject,
	})

	if err != nil {
		api.HandleGRPCError(err, &appG)
		return
	}

	notifications := response.GetNotification()

	appG.Response(http.StatusOK, true, "取得資料成功", nil, notifications)
}
