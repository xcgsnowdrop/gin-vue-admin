package gm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ItemRouter struct{}

func (s *ItemRouter) InitItemRouter(Router *gin.RouterGroup) {
	itemRouter := Router.Group("gm/item").Use(middleware.OperationRecord())
	itemRouterWithoutRecord := Router.Group("gm/item")
	{
		itemRouter.POST("", gmItemApi.CreateItem)                  // 创建游戏道具流水记录
		itemRouter.PUT("", gmItemApi.UpdateItem)                   // 更新游戏道具流水记录
		itemRouter.DELETE("", gmItemApi.DeleteItem)                // 删除游戏道具流水记录
		itemRouter.POST("batchDelete", gmItemApi.BatchDeleteItems) // 批量删除游戏道具流水记录
		itemRouter.POST("export", gmItemApi.ExportItems)           // 导出游戏道具流水数据
		itemRouter.POST("cleanup", gmItemApi.CleanupOldData)       // 清理旧数据
	}
	{
		itemRouterWithoutRecord.POST("list", gmItemApi.GetItemList)                // 获取游戏道具流水列表
		itemRouterWithoutRecord.GET(":id", gmItemApi.GetItem)                      // 获取游戏道具流水详情
		itemRouterWithoutRecord.POST("stats", gmItemApi.GetItemStats)              // 获取游戏道具统计信息
		itemRouterWithoutRecord.GET("types", gmItemApi.GetItemTypes)               // 获取游戏道具类型列表
		itemRouterWithoutRecord.GET("operationTypes", gmItemApi.GetOperationTypes) // 获取操作类型列表
	}
}
