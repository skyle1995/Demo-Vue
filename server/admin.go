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
	// （Redis 可用时多实例共享计数；不可用时自动降级为进程内限流，不会放行）
	admin.GET("/captcha", middleware.RateLimit(30, time.Minute), adminctl.CaptchaHandler)
	// 验证码类型 + 滑动拼图验证码（未鉴权入口，同样加 IP 级限流）
	admin.GET("/captcha/type", middleware.RateLimit(60, time.Minute), adminctl.CaptchaTypeHandler)
	admin.GET("/captcha/slide", middleware.RateLimit(30, time.Minute), adminctl.SlideCaptchaHandler)
	admin.POST("/captcha/slide/verify", middleware.RateLimit(60, time.Minute), adminctl.SlideCaptchaVerifyHandler)
	admin.GET("/captcha/click", middleware.RateLimit(30, time.Minute), adminctl.ClickCaptchaHandler)
	admin.POST("/captcha/click/verify", middleware.RateLimit(60, time.Minute), adminctl.ClickCaptchaVerifyHandler)
	// CSRF 令牌接口同样限流（旧实现未限流，可被脚本无限批量取令牌）
	admin.GET("/csrf", middleware.RateLimit(30, time.Minute), adminctl.CSRFTokenHandler)
	// 登录：IP 级限流（10 次/分钟）之外，LoginHandler 内还有账号级失败锁定（services.LoginGuard）
	admin.POST("/login", middleware.RateLimit(10, time.Minute), adminctl.LoginHandler)
	admin.POST("/refresh-token", middleware.RateLimit(30, time.Minute), adminctl.RefreshTokenHandler)

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

		// 设置API（系统级，仅超级管理员）
		authorized.GET("/settings", adminctl.SuperAdminRequired(), adminctl.SettingsQueryHandler)
		authorized.POST("/settings/update", adminctl.SuperAdminRequired(), adminctl.SettingsUpdateHandler)
		authorized.POST("/settings/generate-key", adminctl.SuperAdminRequired(), adminctl.SettingsGenerateKeyHandler)
		authorized.GET("/navigation", adminctl.PortalNavigationListHandler)
		authorized.POST("/navigation/create", adminctl.PortalNavigationCreateHandler)
		authorized.POST("/navigation/update", adminctl.PortalNavigationUpdateHandler)
		authorized.POST("/navigation/delete", adminctl.PortalNavigationDeleteHandler)

		// 日志管理API
		authorized.GET("/login_logs", adminctl.LoginLogsListHandler)
		authorized.POST("/login_logs/clear", adminctl.LoginLogsClearHandler)
		authorized.POST("/login_logs/batch-delete", adminctl.LoginLogsBatchDeleteHandler)
		authorized.GET("/logs", adminctl.LogsListHandler)
		authorized.POST("/logs/clear", adminctl.LogsClearHandler)
		authorized.POST("/logs/batch-delete", adminctl.LogsBatchDeleteHandler)

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
