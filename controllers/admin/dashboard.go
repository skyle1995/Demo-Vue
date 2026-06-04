package admin

import (
	"Demo/constants"
	"Demo/controllers"
	"Demo/middleware"
	"Demo/utils/timeutil"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ============================================================================
// 全局变量
// ============================================================================

// 创建基础控制器实例
var handlersBaseController = controllers.NewBaseController()

// ============================================================================
// 辅助函数
// ============================================================================

// formatDBType 格式化数据库类型显示
// 将配置文件中的小写类型转换为友好的显示格式
func formatDBType(dbType string) string {
	switch dbType {
	case "mysql":
		return "MySQL"
	case "sqlite":
		return "SQLite"
	case "postgresql", "postgres":
		return "PostgreSQL"
	case "sqlserver":
		return "SQL Server"
	default:
		return "SQLite" // 默认显示
	}
}

// ============================================================================
// API处理器
// ============================================================================

// SystemInfoHandler 系统信息API接口
// 返回系统运行状态的JSON数据，用于前端定时刷新SystemInfoHandler 系统信息API接口
// 返回系统运行状态的JSON数据，用于前端定时刷新
func SystemInfoHandler(c *gin.Context) {
	version := constants.AppVersion
	mode := middleware.IsDevModeFromContext(c)
	dbType := viper.GetString("database.type")
	if dbType == "" {
		dbType = "sqlite"
	}
	uptime := timeutil.GetServerUptimeString()
	uptimeSeconds := int64(timeutil.GetServerUptime().Seconds())

	data := gin.H{
		"version":        version,
		"mode":           mode,
		"db_type":        formatDBType(dbType),
		"uptime":         uptime,
		"uptime_seconds": uptimeSeconds,
	}

	handlersBaseController.HandleSuccess(c, "ok", data)
}

// DashboardStatsHandler 获取数据统计信息
