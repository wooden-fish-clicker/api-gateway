package api

import (
	"api-gateway/pkg/app"
	"api-gateway/pkg/logger"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGRPCError(err error, appG *app.Gin) {
	st, ok := status.FromError(err)
	if !ok {
		log.Printf("非 gRPC 错误: %v", err)
		appG.Response(http.StatusInternalServerError, false, "伺服器未知錯誤", map[string]string{"伺服器未知錯誤": err.Error()}, nil)
		return
	}

	code := st.Code()
	msg := st.Message()

	var statusCode int
	var errorMsg string

	switch code {
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
		errorMsg = "無效參數"
	case codes.Internal:
		statusCode = http.StatusInternalServerError
		errorMsg = "伺服器未知錯誤"
	case codes.NotFound:
		statusCode = http.StatusNotFound
		errorMsg = "查無資料"
	case codes.Unauthenticated:
		statusCode = http.StatusUnauthorized
		errorMsg = "權限不足"
	case codes.Unknown:
		statusCode = http.StatusInternalServerError
		errorMsg = "發生未知錯誤"
	case codes.AlreadyExists:
		statusCode = http.StatusConflict
		errorMsg = "資料有衝突"
	default:
		statusCode = http.StatusInternalServerError
		errorMsg = "發生未知錯誤"
	}

	log.Printf("%s: %s", errorMsg, msg)
	appG.Response(statusCode, false, errorMsg, map[string]string{errorMsg: msg}, nil)
}

func HandleValidError(httpCode int, errors map[string]string, appG *app.Gin) {
	if httpCode == http.StatusBadRequest {
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	} else if httpCode == http.StatusInternalServerError {
		logger.Error(errors)
		appG.Response(httpCode, false, "更新失敗", errors, nil)
		return
	}
}
