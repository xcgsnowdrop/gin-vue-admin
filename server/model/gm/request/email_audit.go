package request

// CreateEmailAuditRequest 创建邮件审核申请请求（统一处理系统邮件和私人邮件）
// 通过 PlayerId 字段区分：PlayerId 为空为系统邮件，PlayerId 有值为私人邮件
type CreateEmailAuditRequest struct {
	EmailType        int                      `json:"emailType" binding:"required"`    // 邮件类型
	EmailTitle       map[string]string        `json:"emailTitle" binding:"required"`   // 邮件标题（多语言）
	EmailContent     map[string]string        `json:"emailContent" binding:"required"` // 邮件内容（多语言）
	EmailAttachments []map[string]interface{} `json:"emailAttachments"`                // 邮件附件
	EmailRemark      string                   `json:"emailRemark"`                     // 邮件备注
	StartTime        *int64                   `json:"startTime" binding:"required"`    // 开始生效时间（时间戳秒）
	AreaIds          string                   `json:"areaIds"`                         // 生效区服列表（逗号分隔，系统邮件专用，PlayerId为空时有效）
	MaxRegTime       *int64                   `json:"maxRegTime"`                      // 最大注册时间（时间戳秒，系统邮件专用，PlayerId为空时有效）
	PlayerId         string                   `json:"playerId"`                        // 目标玩家ID（私人邮件专用，PlayerId有值时表示私人邮件）
}

// UpdateEmailAuditRequest 更新邮件审核申请请求（统一处理系统邮件和私人邮件）
// 通过 PlayerId 字段区分：PlayerId 为空为系统邮件，PlayerId 有值为私人邮件
type UpdateEmailAuditRequest struct {
	ID               uint                     `json:"id" binding:"required"` // 申请ID
	EmailType        int                      `json:"emailType"`             // 邮件类型
	EmailTitle       map[string]string        `json:"emailTitle"`            // 邮件标题（多语言）
	EmailContent     map[string]string        `json:"emailContent"`          // 邮件内容（多语言）
	EmailAttachments []map[string]interface{} `json:"emailAttachments"`      // 邮件附件
	EmailRemark      string                   `json:"emailRemark"`           // 邮件备注
	StartTime        *int64                   `json:"startTime"`             // 开始生效时间（时间戳秒）
	AreaIds          string                   `json:"areaIds"`               // 生效区服列表（逗号分隔，系统邮件专用，PlayerId为空时有效）
	MaxRegTime       *int64                   `json:"maxRegTime"`            // 最大注册时间（时间戳秒，系统邮件专用，PlayerId为空时有效）
	PlayerId         string                   `json:"playerId"`              // 目标玩家ID（私人邮件专用，PlayerId有值时表示私人邮件）
}

// ReviewEmailAuditRequest 审核邮件申请请求
type ReviewEmailAuditRequest struct {
	ID           uint   `json:"id" binding:"required"`     // 申请ID
	Status       int    `json:"status" binding:"required"` // 审核状态：2-通过，3-拒绝，4-待修改
	AuditComment string `json:"auditComment"`              // 审核说明
}

// SearchEmailAuditRequest 搜索邮件审核申请请求
type SearchEmailAuditRequest struct {
	Page        int    `json:"page" form:"page"`               // 页码
	PageSize    int    `json:"pageSize" form:"pageSize"`       // 每页数量
	Status      *int   `json:"status" form:"status"`           // 状态筛选
	ApplicantId *uint  `json:"applicantId" form:"applicantId"` // 申请人ID筛选
	AuditorId   *uint  `json:"auditorId" form:"auditorId"`     // 审核人ID筛选
	StartTime   *int64 `json:"startTime" form:"startTime"`     // 申请开始时间
	EndTime     *int64 `json:"endTime" form:"endTime"`         // 申请结束时间
	PlayerId    string `json:"playerId" form:"playerId"`       // 玩家ID筛选（私人邮件专用）
	IsSystem    *bool  `json:"isSystem" form:"isSystem"`       // 是否为系统邮件：true-系统邮件，false-私人邮件，nil-全部
}
