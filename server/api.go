package server

import (
	"net/http"
	"time"

	"Demo/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterApiRoutes 注册对外开放的 /api/v1 接口（以 API 密钥库鉴权）。
// 脚手架内置一个示例端点，演示 RequireAPIKey 的用法；
// 业务项目在此按需新增接口，并用对应 scope 保护。
func RegisterApiRoutes(apiGroup *gin.RouterGroup) {
	v1 := apiGroup.Group("/v1")
	{
		// 示例：需要携带具备 api:read 能力的密钥才能访问。
		// 调用方式：Header `X-API-Key: sk_xxx` 或 `Authorization: Bearer sk_xxx`。
		v1.GET("/ping",
			middleware.RateLimit(600, time.Minute),
			middleware.RequireAPIKey("api:read"),
			func(c *gin.Context) {
				key := middleware.CurrentAPIKey(c)
				c.JSON(http.StatusOK, gin.H{
					"code": 0,
					"msg":  "pong",
					"data": gin.H{"key_name": key.Name, "ts": time.Now().Unix()},
				})
			},
		)
	}
}
