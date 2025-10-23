package gm

import (
	"time"

	"gorm.io/gorm"
)

// GMUser 游戏用户模型
type GMUser struct {
	ID          uint           `json:"id" gorm:"primarykey"`                   // 主键ID
	GameUserId  string         `json:"gameUserId" gorm:"uniqueIndex;not null"` // 游戏用户ID
	UserName    string         `json:"userName" gorm:"uniqueIndex;not null"`   // 用户名
	NickName    string         `json:"nickName" gorm:"not null"`               // 昵称
	Password    string         `json:"-" gorm:"not null"`                      // 密码
	Phone       string         `json:"phone"`                                  // 手机号
	Email       string         `json:"email"`                                  // 邮箱
	HeaderImg   string         `json:"headerImg"`                              // 头像
	Enable      int            `json:"enable" gorm:"default:1"`                // 启用状态 1:启用 2:禁用
	AuthorityId string         `json:"authorityId"`                            // 权限ID
	CreatedAt   time.Time      `json:"createdAt"`                              // 创建时间
	UpdatedAt   time.Time      `json:"updatedAt"`                              // 更新时间
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`                         // 删除时间
}

// TableName 表名
func (GMUser) TableName() string {
	return "gm_users"
}
