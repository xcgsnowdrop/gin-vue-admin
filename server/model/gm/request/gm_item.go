package request

// SearchItemParams 搜索游戏道具流水参数
type SearchItemParams struct {
	Page          int      `json:"page" form:"page"`                   // 页码
	PageSize      int      `json:"pageSize" form:"pageSize"`           // 每页数量
	UserId        string   `json:"userId" form:"userId"`               // 游戏用户ID
	ItemId        string   `json:"itemId" form:"itemId"`               // 道具ID
	ItemName      string   `json:"itemName" form:"itemName"`           // 道具名称
	OperationType string   `json:"operationType" form:"operationType"` // 操作类型
	DateRange     []string `json:"dateRange" form:"dateRange"`         // 时间范围
}

// CreateItemRequest 创建游戏道具流水请求
type CreateItemRequest struct {
	UserId            string `json:"userId" binding:"required"`            // 游戏用户ID
	ItemId            string `json:"itemId" binding:"required"`            // 道具ID
	ItemName          string `json:"itemName" binding:"required"`          // 道具名称
	OperationType     string `json:"operationType" binding:"required"`     // 操作类型
	QuantityChange    int    `json:"quantityChange" binding:"required"`    // 数量变化
	RemainingQuantity int    `json:"remainingQuantity" binding:"required"` // 剩余数量
	Reason            string `json:"reason"`                               // 操作原因
}

// UpdateItemRequest 更新游戏道具流水请求
type UpdateItemRequest struct {
	ID     uint   `json:"id" binding:"required"`     // 流水ID
	Reason string `json:"reason" binding:"required"` // 操作原因
}

// DeleteItemRequest 删除游戏道具流水请求
type DeleteItemRequest struct {
	ID uint `json:"id" binding:"required"` // 流水ID
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"` // 流水ID列表
}

// CleanupRequest 清理旧数据请求
type CleanupRequest struct {
	Days int `json:"days" binding:"required"` // 清理天数
}
