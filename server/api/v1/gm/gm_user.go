package gm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	gmReq "github.com/flipped-aurora/gin-vue-admin/server/model/gm/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct{}

// GetUserList 获取游戏用户列表
// @Tags GMUser
// @Summary 获取游戏用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.SearchUserParams true "获取游戏用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/user/list [post]
func (u *UserApi) GetUserList(c *gin.Context) {
	var pageInfo gmReq.SearchUserParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := gmUserService.GetPlayerListFromMongo(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetUser 获取游戏用户详情
// @Tags GMUser
// @Summary 获取游戏用户详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "用户ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/user/{id} [get]
func (u *UserApi) GetUser(c *gin.Context) {
	id := c.Param("id")
	if user, err := gmUserService.GetUser(id); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(user, "获取成功", c)
	}
}

// CreateUser 创建游戏用户
// @Tags GMUser
// @Summary 创建游戏用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.CreateUserRequest true "创建游戏用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /gm/user [post]
func (u *UserApi) CreateUser(c *gin.Context) {
	var user gmReq.CreateUserRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmUserService.CreateUser(user); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// UpdateUser 更新游戏用户信息
// @Tags GMUser
// @Summary 更新游戏用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.UpdateUserRequest true "更新游戏用户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /gm/user [put]
func (u *UserApi) UpdateUser(c *gin.Context) {
	var user gmReq.UpdateUserRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmUserService.UpdateUser(user); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DeleteUser 删除游戏用户
// @Tags GMUser
// @Summary 删除游戏用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.DeleteUserRequest true "删除游戏用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /gm/user [delete]
func (u *UserApi) DeleteUser(c *gin.Context) {
	var req gmReq.DeleteUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmUserService.DeleteUser(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// ResetPassword 重置游戏用户密码
// @Tags GMUser
// @Summary 重置游戏用户密码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.ResetPasswordRequest true "重置游戏用户密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"重置成功"}"
// @Router /gm/user/resetPassword [post]
func (u *UserApi) ResetPassword(c *gin.Context) {
	var req gmReq.ResetPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmUserService.ResetPassword(req.ID, req.Password); err != nil {
		global.GVA_LOG.Error("重置失败!", zap.Error(err))
		response.FailWithMessage("重置失败", c)
	} else {
		response.OkWithMessage("重置成功", c)
	}
}

// SetUserStatus 设置游戏用户状态
// @Tags GMUser
// @Summary 设置游戏用户状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.SetUserStatusRequest true "设置游戏用户状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /gm/user/status [post]
func (u *UserApi) SetUserStatus(c *gin.Context) {
	var req gmReq.SetUserStatusRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmUserService.SetUserStatus(req.ID, req.Status); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
	} else {
		response.OkWithMessage("设置成功", c)
	}
}

// BatchOperate 批量操作游戏用户
// @Tags GMUser
// @Summary 批量操作游戏用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.BatchOperateRequest true "批量操作游戏用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /gm/user/batch [post]
func (u *UserApi) BatchOperate(c *gin.Context) {
	var req gmReq.BatchOperateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmUserService.BatchOperate(req.Operation, req.UserIds); err != nil {
		global.GVA_LOG.Error("操作失败!", zap.Error(err))
		response.FailWithMessage("操作失败", c)
	} else {
		response.OkWithMessage("操作成功", c)
	}
}

// ExportUsers 导出游戏用户数据
// @Tags GMUser
// @Summary 导出游戏用户数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.SearchUserParams true "导出游戏用户数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导出成功"}"
// @Router /gm/user/export [post]
func (u *UserApi) ExportUsers(c *gin.Context) {
	var pageInfo gmReq.SearchUserParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if filePath, err := gmUserService.ExportUsers(pageInfo); err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage("导出失败", c)
	} else {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=gm_users.xlsx")
		c.File(filePath)
	}
}

// GetUserStats 获取游戏用户统计信息
// @Tags GMUser
// @Summary 获取游戏用户统计信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/user/stats [get]
func (u *UserApi) GetUserStats(c *gin.Context) {
	if stats, err := gmUserService.GetUserStats(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(stats, "获取成功", c)
	}
}
