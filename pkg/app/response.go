package app

import (
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Status struct {
	IsSuccess bool        `json:"is_success"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details"`
}

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, isSuccess bool, message string, Details interface{}, data interface{}) {
	status := Status{
		IsSuccess: isSuccess,
		Message:   message,
		Details:   Details,
	}

	g.C.JSON(httpCode, Response{
		Status: status,
		Data:   data,
	})
	return
}
