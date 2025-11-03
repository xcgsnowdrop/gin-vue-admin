package gm

import (
	"gmserver/middleware"

	"github.com/gin-gonic/gin"
)

type EmailAuditRouter struct{}

func (s *EmailAuditRouter) InitEmailAuditRouter(Router *gin.RouterGroup) {
	emailAuditRouter := Router.Group("email/audit").Use(middleware.OperationRecord())
	emailAuditRouterWithoutRecord := Router.Group("email/audit")
	{
		emailAuditRouter.POST("apply", gmEmailAuditApi.CreateApplication)   // 创建邮件审核申请
		emailAuditRouter.PUT("apply", gmEmailAuditApi.UpdateApplication)    // 更新邮件审核申请
		emailAuditRouter.DELETE(":id", gmEmailAuditApi.WithdrawApplication) // 撤回邮件审核申请（使用DELETE更符合RESTful规范）
		emailAuditRouter.POST("review", gmEmailAuditApi.ReviewApplication)  // 审核邮件申请
	}
	{
		emailAuditRouterWithoutRecord.POST("list", gmEmailAuditApi.GetApplicationList) // 获取邮件审核申请列表
		emailAuditRouterWithoutRecord.GET(":id", gmEmailAuditApi.GetApplication)       // 获取邮件审核申请详情
	}
}
