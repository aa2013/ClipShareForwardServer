package web

import (
	"clipshare/types"
	"clipshare/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GlobalExceptionMiddleware handle all error on request
func GlobalExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//todo print log
				msg := fmt.Sprintf("%v", err)
				utils.LogUtil.Error(msg)
				// 返回 JSON 格式的错误响应
				errorResult(c, 0, msg, nil)
			}
		}()
		c.Next()
	}
}

// AuthMiddleware Verify if login is valid
func AuthMiddleware(debugModeIgnore bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gin.Mode() == gin.DebugMode && debugModeIgnore {
			c.Next()
			return
		}
		token := c.GetHeader(authKey)
		if token == "" || loginToken == "" || token != loginToken {
			errorResult(c, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		now := time.Now()
		if lastOperationTime != nil {
			duration := now.Sub(*lastOperationTime)
			if int(duration.Seconds()) >= types.AppConfig.Web.LoginExpiredSeconds {
				errorResult(c, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}
		}
		lastOperationTime = &now
		c.Next()
	}
}
func errorResult(c *gin.Context, code int, msg string, data any) {
	defaultCode := http.StatusInternalServerError
	if code <= 0 {
		code = defaultCode
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	c.Abort()
}
func successResult(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": data,
	})
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 可将将* 替换为指定的域名
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
