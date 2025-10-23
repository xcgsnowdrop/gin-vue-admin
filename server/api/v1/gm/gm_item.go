package gm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	gmReq "github.com/flipped-aurora/gin-vue-admin/server/model/gm/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ItemApi struct{}

// GetItemList 获取游戏道具流水列表
// @Tags GMItem
// @Summary 获取游戏道具流水列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.SearchItemParams true "获取游戏道具流水列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/item/list [post]
func (i *ItemApi) GetItemList(c *gin.Context) {
	var pageInfo gmReq.SearchItemParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := gmItemService.GetItemList(pageInfo); err != nil {
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

// GetItem 获取游戏道具流水详情
// @Tags GMItem
// @Summary 获取游戏道具流水详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "流水ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/item/{id} [get]
func (i *ItemApi) GetItem(c *gin.Context) {
	id := c.Param("id")
	if item, err := gmItemService.GetItem(id); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(item, "获取成功", c)
	}
}

// CreateItem 创建游戏道具流水记录
// @Tags GMItem
// @Summary 创建游戏道具流水记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.CreateItemRequest true "创建游戏道具流水记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /gm/item [post]
func (i *ItemApi) CreateItem(c *gin.Context) {
	var item gmReq.CreateItemRequest
	err := c.ShouldBindJSON(&item)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmItemService.CreateItem(item); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// UpdateItem 更新游戏道具流水记录
// @Tags GMItem
// @Summary 更新游戏道具流水记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.UpdateItemRequest true "更新游戏道具流水记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /gm/item [put]
func (i *ItemApi) UpdateItem(c *gin.Context) {
	var item gmReq.UpdateItemRequest
	err := c.ShouldBindJSON(&item)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmItemService.UpdateItem(item); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DeleteItem 删除游戏道具流水记录
// @Tags GMItem
// @Summary 删除游戏道具流水记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.DeleteItemRequest true "删除游戏道具流水记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /gm/item [delete]
func (i *ItemApi) DeleteItem(c *gin.Context) {
	var req gmReq.DeleteItemRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmItemService.DeleteItem(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// BatchDeleteItems 批量删除游戏道具流水记录
// @Tags GMItem
// @Summary 批量删除游戏道具流水记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.BatchDeleteRequest true "批量删除游戏道具流水记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /gm/item/batchDelete [post]
func (i *ItemApi) BatchDeleteItems(c *gin.Context) {
	var req gmReq.BatchDeleteRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmItemService.BatchDeleteItems(req.IDs); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// ExportItems 导出游戏道具流水数据
// @Tags GMItem
// @Summary 导出游戏道具流水数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.SearchItemParams true "导出游戏道具流水数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导出成功"}"
// @Router /gm/item/export [post]
func (i *ItemApi) ExportItems(c *gin.Context) {
	var pageInfo gmReq.SearchItemParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if filePath, err := gmItemService.ExportItems(pageInfo); err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage("导出失败", c)
	} else {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=gm_items.xlsx")
		c.File(filePath)
	}
}

// CleanupOldData 清理旧数据
// @Tags GMItem
// @Summary 清理旧数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.CleanupRequest true "清理旧数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"清理成功"}"
// @Router /gm/item/cleanup [post]
func (i *ItemApi) CleanupOldData(c *gin.Context) {
	var req gmReq.CleanupRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := gmItemService.CleanupOldData(req.Days); err != nil {
		global.GVA_LOG.Error("清理失败!", zap.Error(err))
		response.FailWithMessage("清理失败", c)
	} else {
		response.OkWithMessage("清理成功", c)
	}
}

// GetItemStats 获取游戏道具统计信息
// @Tags GMItem
// @Summary 获取游戏道具统计信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gmReq.SearchItemParams true "获取游戏道具统计信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/item/stats [post]
func (i *ItemApi) GetItemStats(c *gin.Context) {
	var pageInfo gmReq.SearchItemParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if stats, err := gmItemService.GetItemStats(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(stats, "获取成功", c)
	}
}

// GetItemTypes 获取游戏道具类型列表
// @Tags GMItem
// @Summary 获取游戏道具类型列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/item/types [get]
func (i *ItemApi) GetItemTypes(c *gin.Context) {
	if types, err := gmItemService.GetItemTypes(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(types, "获取成功", c)
	}
}

// GetOperationTypes 获取操作类型列表
// @Tags GMItem
// @Summary 获取操作类型列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gm/item/operationTypes [get]
func (i *ItemApi) GetOperationTypes(c *gin.Context) {
	if types, err := gmItemService.GetOperationTypes(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(types, "获取成功", c)
	}
}
