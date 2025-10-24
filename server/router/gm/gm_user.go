package gm

import (
	"gmserver/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("gm/user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("gm/user")
	{
		userRouter.POST("", gmUserApi.CreateUser)                 // 创建游戏用户
		userRouter.PUT("", gmUserApi.UpdateUser)                  // 更新游戏用户信息
		userRouter.DELETE("", gmUserApi.DeleteUser)               // 删除游戏用户
		userRouter.POST("resetPassword", gmUserApi.ResetPassword) // 重置游戏用户密码
		userRouter.POST("status", gmUserApi.SetUserStatus)        // 设置游戏用户状态
		userRouter.POST("batch", gmUserApi.BatchOperate)          // 批量操作游戏用户
		userRouter.POST("export", gmUserApi.ExportUsers)          // 导出游戏用户数据
	}
	{
		userRouterWithoutRecord.POST("list", gmUserApi.GetUserList)  // 获取游戏用户列表
		userRouterWithoutRecord.GET(":id", gmUserApi.GetUser)        // 获取游戏用户详情
		userRouterWithoutRecord.GET("stats", gmUserApi.GetUserStats) // 获取游戏用户统计信息
	}
}
