package gm

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gmserver/model/system"

	"gorm.io/gorm"
)

// EmailAuditStatus 审核状态
type EmailAuditStatus int

const (
	EmailAuditStatusPending  EmailAuditStatus = 1 // 待审核
	EmailAuditStatusApproved EmailAuditStatus = 2 // 通过
	EmailAuditStatusRejected EmailAuditStatus = 3 // 拒绝
	EmailAuditStatusRevision EmailAuditStatus = 4 // 待修改
)

// EmailAuditApplication 邮件审核申请
type EmailAuditApplication struct {
	ID               uint             `json:"id" gorm:"primarykey"`                                                  // 主键ID
	ApplicantId      uint             `json:"applicantId" gorm:"not null;index;comment:申请人ID"`                       // 申请人ID（关联sys_users.id）
	Applicant        system.SysUser   `json:"-" gorm:"foreignKey:ApplicantId;references:ID;comment:申请人"`             // 申请人（关联SysUser，不直接暴露给前端）
	ApplicantTime    time.Time        `json:"applicantTime" gorm:"not null;comment:申请时间"`                            // 申请时间
	AuditorId        *uint            `json:"auditorId" gorm:"index;comment:审核人ID"`                                  // 审核人ID（可为空，审核时填入）
	Auditor          *system.SysUser  `json:"-" gorm:"foreignKey:AuditorId;references:ID;comment:审核人"`               // 审核人（关联SysUser，可为空，不直接暴露给前端）
	AuditTime        *time.Time       `json:"auditTime" gorm:"comment:审核时间"`                                         // 审核时间（可为空）
	AuditComment     string           `json:"auditComment" gorm:"type:text;comment:审核说明"`                            // 审核说明
	Status           EmailAuditStatus `json:"status" gorm:"type:tinyint;default:1;comment:状态：1-待审核，2-通过，3-拒绝，4-待修改"` // 状态
	EmailType        int              `json:"emailType" gorm:"type:tinyint;not null;comment:邮件类型"`                   // 邮件类型
	EmailTitle       JSONMap          `json:"emailTitle" gorm:"type:json;not null;comment:邮件标题（多语言）"`                // 邮件标题（多语言JSON）
	EmailContent     JSONMap          `json:"emailContent" gorm:"type:json;not null;comment:邮件内容（多语言）"`              // 邮件内容（多语言JSON）
	EmailAttachments JSONArray        `json:"emailAttachments" gorm:"type:json;comment:邮件附件"`                        // 邮件附件（JSON数组）
	EmailRemark      string           `json:"emailRemark" gorm:"type:varchar(500);comment:邮件备注"`                     // 邮件备注
	StartTime        *int64           `json:"startTime" gorm:"comment:开始生效时间（时间戳秒）"`                                 // 开始生效时间
	AreaIds          string           `json:"areaIds" gorm:"type:varchar(500);comment:生效区服列表（逗号分隔）"`                 // 生效区服列表
	MaxRegTime       *int64           `json:"maxRegTime" gorm:"comment:最大注册时间（时间戳秒）"`                                // 最大注册时间
	CreatedAt        time.Time        `json:"createdAt"`                                                             // 创建时间
	UpdatedAt        time.Time        `json:"updatedAt"`                                                             // 更新时间
	DeletedAt        gorm.DeletedAt   `json:"-" gorm:"index"`                                                        // 删除时间
}

// TableName 表名
func (EmailAuditApplication) TableName() string {
	return "email_audit_applications"
}

// JSONMap JSON 映射类型（用于多语言字段）
type JSONMap map[string]interface{}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// JSONArray JSON 数组类型（用于附件字段）
type JSONArray []interface{}

// Value 实现 driver.Valuer 接口
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
