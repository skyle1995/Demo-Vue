package server

import (
	adminctl "Demo/controllers/admin"
	"Demo/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理员后台相关路由
func RegisterAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")

	// Admin 认证相关路由
	// 登录、刷新令牌、验证码等未鉴权入口加 IP 级限流，缓解暴力破解
	// （限流依赖 Redis，Redis 不可用时自动降级放行）
	admin.GET("/captcha", middleware.RateLimit(30, time.Minute), adminctl.CaptchaHandler)
	admin.GET("/csrf", adminctl.CSRFTokenHandler)
	admin.POST("/login", middleware.RateLimit(10, time.Minute), adminctl.LoginHandler)
	admin.POST("/refresh-token", middleware.RateLimit(30, time.Minute), adminctl.RefreshTokenHandler)

	// 公开设置API
	admin.GET("/settings/public", adminctl.SettingsPublicHandler)
	admin.GET("/portal-navigation/public", adminctl.PortalNavigationPublicListHandler)

	// 退出登录
	admin.POST("/logout", adminctl.LogoutHandler)

	// 需要认证的路由组
	authorized := admin.Group("/")
	authorized.Use(adminctl.AdminAuthRequired())
	{
		// 系统信息API
		authorized.GET("/system/info", adminctl.SystemInfoHandler)
		authorized.GET("/dashboard/login-logs", adminctl.DashboardLoginLogsHandler)

		// 个人资料API
		authorized.GET("/profile", adminctl.ProfileQueryHandler)
		authorized.POST("/profile/update", adminctl.ProfileUpdateHandler)
		authorized.POST("/profile/password", adminctl.ProfilePasswordUpdateHandler)

		// 设置API
		authorized.GET("/settings", adminctl.SettingsQueryHandler)
		authorized.POST("/settings/update", adminctl.SettingsUpdateHandler)
		authorized.POST("/settings/generate-key", adminctl.SettingsGenerateKeyHandler)
		authorized.GET("/portal-navigation", adminctl.PortalNavigationListHandler)
		authorized.POST("/portal-navigation/create", adminctl.PortalNavigationCreateHandler)
		authorized.POST("/portal-navigation/update", adminctl.PortalNavigationUpdateHandler)
		authorized.POST("/portal-navigation/delete", adminctl.PortalNavigationDeleteHandler)

		// 日志管理API
		authorized.GET("/login_logs", adminctl.LoginLogsListHandler)
		authorized.POST("/login_logs/clear", adminctl.LoginLogsClearHandler)
		authorized.GET("/logs", adminctl.LogsListHandler)
		authorized.POST("/logs/clear", adminctl.LogsClearHandler)

		// 子账号相关API (Mock)
		authorized.GET("/subaccounts/simple", adminctl.SubAccountSimpleListHandler)

		// API 密钥库管理
		apikeyGroup := authorized.Group("/apikey")
		{
			apikeyGroup.GET("/scopes", adminctl.GetApiKeyScopes)
			apikeyGroup.GET("/list", adminctl.GetApiKeyList)
			apikeyGroup.POST("/create", adminctl.CreateApiKey)
			apikeyGroup.PUT("/update", adminctl.UpdateApiKey)
			apikeyGroup.POST("/regenerate", adminctl.RegenerateApiKey)
			apikeyGroup.DELETE("/delete/:id", adminctl.DeleteApiKey)
		}
	}
}
