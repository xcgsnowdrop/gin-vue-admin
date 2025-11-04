package gm

import (
	"strconv"

	"gmserver/global"
	"gmserver/model/common/response"
	gmReq "gmserver/model/gm/request"
	"gmserver/service/gm"
	"gmserver/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmailAuditApi struct{}

// CreateApplication 创建邮件审核申请
// @Tags EmailAudit
// @Summary 创建邮件审核申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.CreateEmailAuditRequest true "创建邮件审核申请"
// @Success 200 {object} response.Response{data=gm.EmailAuditApplication,msg=string} "创建成功"
// @Router /gm/email/audit/apply [post]
func (e *EmailAuditApi) CreateApplication(c *gin.Context) {
	var req gmReq.CreateEmailAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取当前用户ID
	applicantId := utils.GetUserID(c)
	if applicantId == 0 {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	emailAuditService := gm.EmailAuditService{}
	application, err := emailAuditService.CreateApplication(req, applicantId)
	if err != nil {
		global.GVA_LOG.Error("创建邮件审核申请失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(application, c)
}

// UpdateApplication 更新邮件审核申请
// @Tags EmailAudit
// @Summary 更新邮件审核申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.UpdateEmailAuditRequest true "更新邮件审核申请"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /gm/email/audit/apply [put]
func (e *EmailAuditApi) UpdateApplication(c *gin.Context) {
	var req gmReq.UpdateEmailAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取当前用户ID
	applicantId := utils.GetUserID(c)
	if applicantId == 0 {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	emailAuditService := gm.EmailAuditService{}
	if err := emailAuditService.UpdateApplication(req, applicantId); err != nil {
		global.GVA_LOG.Error("更新邮件审核申请失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// WithdrawApplication 撤回邮件审核申请
// @Tags EmailAudit
// @Summary 撤回邮件审核申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path uint true "申请ID"
// @Success 200 {object} response.Response{msg=string} "撤回成功"
// @Router /gm/email/audit/:id [delete]
func (e *EmailAuditApi) WithdrawApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的申请ID", c)
		return
	}

	// 获取当前用户ID
	applicantId := utils.GetUserID(c)
	if applicantId == 0 {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	emailAuditService := gm.EmailAuditService{}
	if err := emailAuditService.WithdrawApplication(uint(id), applicantId); err != nil {
		global.GVA_LOG.Error("撤回邮件审核申请失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("撤回成功", c)
}

// ReviewApplication 审核邮件申请
// @Tags EmailAudit
// @Summary 审核邮件申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.ReviewEmailAuditRequest true "审核邮件申请"
// @Success 200 {object} response.Response{msg=string} "审核成功"
// @Router /gm/email/audit/review [post]
func (e *EmailAuditApi) ReviewApplication(c *gin.Context) {
	var req gmReq.ReviewEmailAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取当前用户ID和角色ID
	auditorId := utils.GetUserID(c)
	auditorAuthorityId := utils.GetUserAuthorityId(c)
	if auditorId == 0 || auditorAuthorityId == 0 {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	emailAuditService := gm.EmailAuditService{}
	result, err := emailAuditService.ReviewApplication(req, auditorId, auditorAuthorityId)
	if err != nil {
		global.GVA_LOG.Error("审核邮件申请失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 构建返回消息
	msg := result.ReviewMessage
	if result.EmailSentMsg != "" {
		msg += "，" + result.EmailSentMsg
	}

	response.OkWithDetailed(result, msg, c)
}

// GetApplicationList 获取申请列表
// @Tags EmailAudit
// @Summary 获取申请列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.SearchEmailAuditRequest true "搜索参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gm/email/audit/list [post]
func (e *EmailAuditApi) GetApplicationList(c *gin.Context) {
	var req gmReq.SearchEmailAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取当前用户ID和角色ID
	userId := utils.GetUserID(c)
	userAuthorityId := utils.GetUserAuthorityId(c)
	if userId == 0 || userAuthorityId == 0 {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	emailAuditService := gm.EmailAuditService{}
	list, total, err := emailAuditService.GetApplicationListResponse(req, userId, userAuthorityId)
	if err != nil {
		global.GVA_LOG.Error("获取申请列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetApplication 获取申请详情
// @Tags EmailAudit
// @Summary 获取申请详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path uint true "申请ID"
// @Success 200 {object} response.Response{data=gm.EmailAuditApplication,msg=string} "获取成功"
// @Router /gm/email/audit/:id [get]
func (e *EmailAuditApi) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的申请ID", c)
		return
	}

	// 获取当前用户ID和角色ID
	userId := utils.GetUserID(c)
	userAuthorityId := utils.GetUserAuthorityId(c)
	if userId == 0 || userAuthorityId == 0 {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	emailAuditService := gm.EmailAuditService{}
	application, err := emailAuditService.GetApplicationResponse(uint(id), userId, userAuthorityId)
	if err != nil {
		global.GVA_LOG.Error("获取申请详情失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(application, c)
}
