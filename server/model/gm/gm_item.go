package gm

import (
	"time"

	"gorm.io/gorm"
)

// GMItem 游戏道具流水模型
type GMItem struct {
	ID                uint           `json:"id" gorm:"primarykey"`              // 主键ID
	UserId            string         `json:"userId" gorm:"not null"`            // 游戏用户ID
	ItemId            string         `json:"itemId" gorm:"not null"`            // 道具ID
	ItemName          string         `json:"itemName" gorm:"not null"`          // 道具名称
	OperationType     string         `json:"operationType" gorm:"not null"`     // 操作类型 gain:获得 consume:消耗 trade:交易 system:系统
	QuantityChange    int            `json:"quantityChange" gorm:"not null"`    // 数量变化
	RemainingQuantity int            `json:"remainingQuantity" gorm:"not null"` // 剩余数量
	Reason            string         `json:"reason"`                            // 操作原因
	CreatedAt         time.Time      `json:"createdAt"`                         // 创建时间
	UpdatedAt         time.Time      `json:"updatedAt"`                         // 更新时间
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`                    // 删除时间
}

// TableName 表名
func (GMItem) TableName() string {
	return "gm_items"
}
