package middleware

import (
	// "api-gateway/internal/models"
	"api-gateway/pkg/app"
	"api-gateway/pkg/redis"
	"api-gateway/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CheckJwtToken struct {
	Redis *redis.Redis
}

func NewCheckJwtToken(redis *redis.Redis) *CheckJwtToken {
	return &CheckJwtToken{Redis: redis}
}

func (cjt *CheckJwtToken) CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		appG := app.Gin{C: c}

		//取得Authorization裡面的token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appG.Response(http.StatusUnauthorized, false, "未登入", "未登入", nil)
			c.Abort()
			return
		}

		//分割Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appG.Response(http.StatusUnauthorized, false, "未登入", "Invalid or missing Bearer token", nil)
			c.Abort()
			return
		}

		//檢查redis裡面的黑名單的token
		err := cjt.checkRedisJwtBlackList(parts[1])
		if err != nil && err.Error() != "redis: nil" {
			appG.Response(http.StatusInternalServerError, false, "未登入", "發生錯誤", nil)
			c.Abort()
			return
		} else if err == nil {
			appG.Response(http.StatusUnauthorized, false, "未登入", "此token已經失效", nil)
			c.Abort()
			return
		}

		//驗證JWT
		claims, err1 := utils.ParseJwtToken(parts[1])
		if err1 != nil {
			appG.Response(http.StatusInternalServerError, false, "未登入", err1.Error(), nil)
			c.Abort()
			return
		}

		// //檢查用戶權限
		// if err2 := checkUserAndRole(claims); err2 != nil {
		// 	appG.Response(http.StatusUnauthorized, false, "沒有權限", err2.Error(), nil)
		// 	c.Abort()
		// 	return
		// }

		c.Set("jwtToken", parts[1])
		c.Set("jwtClaims", claims)

		c.Next()
	}
}

// func checkUserAndRole(c *utils.Claims) error {

// 	// userIdStr := c.Subject
// 	// userId, _ := strconv.ParseUint(userIdStr, 10, 64)

// 	switch c.Role {
// 	case "admin":
// 		// admin := models.Admin{}
// 		// if _, error := db.CheckExist(uint(userId), admin); error != nil {
// 		// 	return errors.New("用戶不存在")
// 		// }
// 		// return nil
// 	case "user":
//
// 	default:
// 		return errors.New("沒有權限")
// 	}

// 	return errors.New("沒有權限")
// }

func (cjt *CheckJwtToken) checkRedisJwtBlackList(token string) error {
	_, err := cjt.Redis.Client.Get(context.Background(), "jwt:blacklist:"+token).Result()

	if err != nil {
		return err
	}

	return nil
}
