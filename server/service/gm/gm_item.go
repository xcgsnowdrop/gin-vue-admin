package gm

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gm"
	gmReq "github.com/flipped-aurora/gin-vue-admin/server/model/gm/request"
)

type ItemService struct{}

// GetItemList 获取游戏道具流水列表
func (itemService *ItemService) GetItemList(info gmReq.SearchItemParams) (list []gm.GMItem, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&gm.GMItem{})

	// 构建查询条件
	if info.UserId != "" {
		db = db.Where("user_id LIKE ?", "%"+info.UserId+"%")
	}
	if info.ItemId != "" {
		db = db.Where("item_id LIKE ?", "%"+info.ItemId+"%")
	}
	if info.ItemName != "" {
		db = db.Where("item_name LIKE ?", "%"+info.ItemName+"%")
	}
	if info.OperationType != "" {
		db = db.Where("operation_type = ?", info.OperationType)
	}
	if len(info.DateRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.DateRange[0], info.DateRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

// GetItem 获取游戏道具流水详情
func (itemService *ItemService) GetItem(id string) (item gm.GMItem, err error) {
	var itemId uint64
	itemId, err = strconv.ParseUint(id, 10, 32)
	if err != nil {
		return item, errors.New("无效的流水ID")
	}

	err = global.GVA_DB.Where("id = ?", uint(itemId)).First(&item).Error
	return item, err
}

// CreateItem 创建游戏道具流水记录
func (itemService *ItemService) CreateItem(item gmReq.CreateItemRequest) (err error) {
	// 创建新流水记录
	newItem := gm.GMItem{
		UserId:            item.UserId,
		ItemId:            item.ItemId,
		ItemName:          item.ItemName,
		OperationType:     item.OperationType,
		QuantityChange:    item.QuantityChange,
		RemainingQuantity: item.RemainingQuantity,
		Reason:            item.Reason,
	}

	err = global.GVA_DB.Create(&newItem).Error
	return err
}

// UpdateItem 更新游戏道具流水记录
func (itemService *ItemService) UpdateItem(item gmReq.UpdateItemRequest) (err error) {
	var existingItem gm.GMItem
	err = global.GVA_DB.Where("id = ?", item.ID).First(&existingItem).Error
	if err != nil {
		return errors.New("流水记录不存在")
	}

	// 更新流水记录
	err = global.GVA_DB.Model(&existingItem).Updates(map[string]interface{}{
		"reason":     item.Reason,
		"updated_at": time.Now(),
	}).Error
	return err
}

// DeleteItem 删除游戏道具流水记录
func (itemService *ItemService) DeleteItem(id uint) (err error) {
	var item gm.GMItem
	err = global.GVA_DB.Where("id = ?", id).First(&item).Error
	if err != nil {
		return errors.New("流水记录不存在")
	}

	err = global.GVA_DB.Delete(&item).Error
	return err
}

// BatchDeleteItems 批量删除游戏道具流水记录
func (itemService *ItemService) BatchDeleteItems(ids []uint) (err error) {
	err = global.GVA_DB.Where("id IN ?", ids).Delete(&gm.GMItem{}).Error
	return err
}

// ExportItems 导出游戏道具流水数据
func (itemService *ItemService) ExportItems(info gmReq.SearchItemParams) (filePath string, err error) {
	// 这里应该实现Excel导出逻辑
	// 暂时返回模拟文件路径
	filePath = fmt.Sprintf("exports/gm_items_%d.xlsx", time.Now().Unix())
	return filePath, nil
}

// CleanupOldData 清理旧数据
func (itemService *ItemService) CleanupOldData(days int) (err error) {
	// 计算清理日期
	cleanupDate := time.Now().AddDate(0, 0, -days)

	// 删除指定天数前的数据
	err = global.GVA_DB.Where("created_at < ?", cleanupDate).Delete(&gm.GMItem{}).Error
	return err
}

// GetItemStats 获取游戏道具统计信息
func (itemService *ItemService) GetItemStats(info gmReq.SearchItemParams) (stats map[string]interface{}, err error) {
	var totalRecords, gainRecords, consumeRecords, tradeRecords, systemRecords int64

	// 构建基础查询条件
	db := global.GVA_DB.Model(&gm.GMItem{})
	if info.UserId != "" {
		db = db.Where("user_id LIKE ?", "%"+info.UserId+"%")
	}
	if info.ItemId != "" {
		db = db.Where("item_id LIKE ?", "%"+info.ItemId+"%")
	}
	if info.ItemName != "" {
		db = db.Where("item_name LIKE ?", "%"+info.ItemName+"%")
	}
	if info.OperationType != "" {
		db = db.Where("operation_type = ?", info.OperationType)
	}
	if len(info.DateRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.DateRange[0], info.DateRange[1])
	}

	// 总记录数
	err = db.Count(&totalRecords).Error
	if err != nil {
		return nil, err
	}

	// 获得记录数
	err = db.Where("operation_type = ?", "gain").Count(&gainRecords).Error
	if err != nil {
		return nil, err
	}

	// 消耗记录数
	err = db.Where("operation_type = ?", "consume").Count(&consumeRecords).Error
	if err != nil {
		return nil, err
	}

	// 交易记录数
	err = db.Where("operation_type = ?", "trade").Count(&tradeRecords).Error
	if err != nil {
		return nil, err
	}

	// 系统记录数
	err = db.Where("operation_type = ?", "system").Count(&systemRecords).Error
	if err != nil {
		return nil, err
	}

	stats = map[string]interface{}{
		"totalRecords":   totalRecords,
		"gainRecords":    gainRecords,
		"consumeRecords": consumeRecords,
		"tradeRecords":   tradeRecords,
		"systemRecords":  systemRecords,
	}

	return stats, nil
}

// GetItemTypes 获取游戏道具类型列表
func (itemService *ItemService) GetItemTypes() (types []map[string]interface{}, err error) {
	// 模拟道具类型数据，实际应该从数据库或配置中获取
	types = []map[string]interface{}{
		{"id": "1", "name": "金币", "category": "货币"},
		{"id": "2", "name": "钻石", "category": "货币"},
		{"id": "3", "name": "经验药水", "category": "消耗品"},
		{"id": "4", "name": "装备强化石", "category": "材料"},
		{"id": "5", "name": "技能书", "category": "技能"},
		{"id": "6", "name": "宠物蛋", "category": "宠物"},
	}
	return types, nil
}

// GetOperationTypes 获取操作类型列表
func (itemService *ItemService) GetOperationTypes() (types []map[string]interface{}, err error) {
	// 模拟操作类型数据
	types = []map[string]interface{}{
		{"id": "gain", "name": "获得", "description": "通过任务、活动等方式获得道具"},
		{"id": "consume", "name": "消耗", "description": "使用道具进行强化、升级等操作"},
		{"id": "trade", "name": "交易", "description": "玩家之间的道具交易"},
		{"id": "system", "name": "系统", "description": "系统自动处理的道具操作"},
	}
	return types, nil
}
