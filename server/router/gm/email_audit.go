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
		// 邮件审核申请（统一处理系统邮件和私人邮件，通过 PlayerId 字段区分）
		emailAuditRouter.POST("apply", gmEmailAuditApi.CreateApplication) // 创建邮件审核申请（PlayerId为空为系统邮件，有值为私人邮件）
		emailAuditRouter.PUT("apply", gmEmailAuditApi.UpdateApplication)  // 更新邮件审核申请（PlayerId为空为系统邮件，有值为私人邮件）

		// 通用操作
		emailAuditRouter.DELETE(":id", gmEmailAuditApi.WithdrawApplication) // 撤回邮件审核申请（使用DELETE更符合RESTful规范）
		emailAuditRouter.POST("review", gmEmailAuditApi.ReviewApplication)  // 审核邮件申请（统一处理系统邮件和私人邮件）
	}
	{
		emailAuditRouterWithoutRecord.POST("list", gmEmailAuditApi.GetApplicationList) // 获取邮件审核申请列表（支持系统/私人邮件筛选）
		emailAuditRouterWithoutRecord.GET(":id", gmEmailAuditApi.GetApplication)       // 获取邮件审核申请详情
	}
}
