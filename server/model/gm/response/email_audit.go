package response

import (
	"time"

	"gmserver/model/gm"
)

// SimpleUserInfo 简化的用户信息（只包含前端需要的字段）
type SimpleUserInfo struct {
	ID       uint   `json:"id"`       // 用户ID
	NickName string `json:"nickName"` // 用户昵称
}

// EmailAuditApplicationResponse 邮件审核申请响应（用于前端展示）
type EmailAuditApplicationResponse struct {
	ID               uint                `json:"id"`               // 主键ID
	ApplicantId      uint                `json:"applicantId"`      // 申请人ID
	Applicant        SimpleUserInfo      `json:"applicant"`        // 申请人信息（简化）
	ApplicantTime    time.Time           `json:"applicantTime"`    // 申请时间
	AuditorId        *uint               `json:"auditorId"`        // 审核人ID
	Auditor          *SimpleUserInfo     `json:"auditor"`          // 审核人信息（简化，可为空）
	AuditTime        *time.Time          `json:"auditTime"`        // 审核时间
	AuditComment     string              `json:"auditComment"`     // 审核说明
	Status           gm.EmailAuditStatus `json:"status"`           // 状态
	EmailType        int                 `json:"emailType"`        // 邮件类型
	EmailTitle       gm.JSONMap          `json:"emailTitle"`       // 邮件标题（多语言JSON）
	EmailContent     gm.JSONMap          `json:"emailContent"`     // 邮件内容（多语言JSON）
	EmailAttachments gm.JSONArray        `json:"emailAttachments"` // 邮件附件（JSON数组）
	EmailRemark      string              `json:"emailRemark"`      // 邮件备注
	StartTime        *int64              `json:"startTime"`        // 开始生效时间（时间戳秒）
	AreaIds          string              `json:"areaIds"`          // 生效区服列表（逗号分隔）
	MaxRegTime       *int64              `json:"maxRegTime"`       // 最大注册时间（时间戳秒）
	PlayerId         string              `json:"playerId"`         // 目标玩家ID（私人邮件专用）
	CreatedAt        time.Time           `json:"createdAt"`        // 创建时间
	UpdatedAt        time.Time           `json:"updatedAt"`        // 更新时间
}

// ToEmailAuditApplicationResponse 将 model 转换为 response DTO
func ToEmailAuditApplicationResponse(app *gm.EmailAuditApplication) *EmailAuditApplicationResponse {
	if app == nil {
		return nil
	}

	response := &EmailAuditApplicationResponse{
		ID:               app.ID,
		ApplicantId:      app.ApplicantId,
		ApplicantTime:    app.ApplicantTime,
		AuditorId:        app.AuditorId,
		AuditTime:        app.AuditTime,
		AuditComment:     app.AuditComment,
		Status:           app.Status,
		EmailType:        app.EmailType,
		EmailTitle:       app.EmailTitle,
		EmailContent:     app.EmailContent,
		EmailAttachments: app.EmailAttachments,
		EmailRemark:      app.EmailRemark,
		StartTime:        app.StartTime,
		AreaIds:          app.AreaIds,
		MaxRegTime:       app.MaxRegTime,
		PlayerId:         app.PlayerId,
		CreatedAt:        app.CreatedAt,
		UpdatedAt:        app.UpdatedAt,
	}

	// 转换申请人信息（只包含必要字段）
	response.Applicant = SimpleUserInfo{
		ID:       app.Applicant.ID,
		NickName: app.Applicant.NickName,
	}

	// 转换审核人信息（如果存在）
	if app.Auditor != nil {
		response.Auditor = &SimpleUserInfo{
			ID:       app.Auditor.ID,
			NickName: app.Auditor.NickName,
		}
	}

	return response
}

// ToEmailAuditApplicationResponseList 将 model 列表转换为 response DTO 列表
func ToEmailAuditApplicationResponseList(apps []gm.EmailAuditApplication) []EmailAuditApplicationResponse {
	if apps == nil {
		return nil
	}

	responses := make([]EmailAuditApplicationResponse, len(apps))
	for i := range apps {
		responses[i] = *ToEmailAuditApplicationResponse(&apps[i])
	}

	return responses
}
