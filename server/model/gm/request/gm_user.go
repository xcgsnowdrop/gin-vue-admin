package request

// SearchUserParams 搜索游戏用户参数
type SearchUserParams struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页数量
	UserId   string `json:"userId" form:"userId"`     // 游戏用户ID
	UserName string `json:"userName" form:"userName"` // 用户名
	NickName string `json:"nickName" form:"nickName"` // 昵称
	Phone    string `json:"phone" form:"phone"`       // 手机号
	Email    string `json:"email" form:"email"`       // 邮箱
}

// CreateUserRequest 创建游戏用户请求
type CreateUserRequest struct {
	GameUserId  string `json:"gameUserId" binding:"required"` // 游戏用户ID
	UserName    string `json:"userName" binding:"required"`   // 用户名
	NickName    string `json:"nickName" binding:"required"`   // 昵称
	Password    string `json:"password" binding:"required"`   // 密码
	Phone       string `json:"phone"`                         // 手机号
	Email       string `json:"email"`                         // 邮箱
	HeaderImg   string `json:"headerImg"`                     // 头像
	Enable      int    `json:"enable"`                        // 启用状态
	AuthorityId string `json:"authorityId"`                   // 权限ID
}

// UpdateUserRequest 更新游戏用户请求
type UpdateUserRequest struct {
	ID          uint   `json:"id" binding:"required"` // 用户ID
	GameUserId  string `json:"gameUserId"`            // 游戏用户ID
	UserName    string `json:"userName"`              // 用户名
	NickName    string `json:"nickName"`              // 昵称
	Phone       string `json:"phone"`                 // 手机号
	Email       string `json:"email"`                 // 邮箱
	HeaderImg   string `json:"headerImg"`             // 头像
	Enable      int    `json:"enable"`                // 启用状态
	AuthorityId string `json:"authorityId"`           // 权限ID
}

// DeleteUserRequest 删除游戏用户请求
type DeleteUserRequest struct {
	ID uint `json:"id" binding:"required"` // 用户ID
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	ID       uint   `json:"id" binding:"required"`       // 用户ID
	Password string `json:"password" binding:"required"` // 新密码
}

// SetUserStatusRequest 设置用户状态请求
type SetUserStatusRequest struct {
	ID     uint `json:"id" binding:"required"`     // 用户ID
	Status int  `json:"status" binding:"required"` // 状态
}

// BatchOperateRequest 批量操作请求
type BatchOperateRequest struct {
	Operation string `json:"operation" binding:"required"` // 操作类型
	UserIds   []uint `json:"userIds" binding:"required"`   // 用户ID列表
}
