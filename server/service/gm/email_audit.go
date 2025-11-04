package gm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gmserver/global"
	"gmserver/model/gm"
	gmReq "gmserver/model/gm/request"
	gmResp "gmserver/model/gm/response"
	"gmserver/model/system"
	"gmserver/utils"

	"go.uber.org/zap"
)

type EmailAuditService struct{}

// isRootAuthority 判断是否是根角色（没有父角色）
// authority.ParentId 为 nil 或指向 0 时，表示根角色
func (s *EmailAuditService) isRootAuthority(authority system.SysAuthority) bool {
	if authority.ParentId == nil {
		return true
	}
	return *authority.ParentId == 0
}

// GetAllParentAuthorities 递归获取所有父级角色ID(包含入参角色ID)
func (s *EmailAuditService) getAllParentAuthorities(authorityId uint) []uint {
	var authorityIds []uint
	currentId := authorityId

	for {
		var authority system.SysAuthority
		if err := global.GVA_DB.First(&authority, currentId).Error; err != nil {
			break
		}

		// 添加当前角色
		authorityIds = append(authorityIds, authority.AuthorityId)

		// 如果是根角色（没有父角色或父角色ID为0），结束递归
		if s.isRootAuthority(authority) {
			break
		}

		currentId = *authority.ParentId
	}

	return authorityIds
}

// GetSameLevelAuthorities 获取同级角色ID（排除入参角色ID）
func (s *EmailAuditService) getSameLevelAuthorities(authorityId uint) []uint {
	// 获取当前角色的父角色ID
	var currentAuthority system.SysAuthority
	if err := global.GVA_DB.First(&currentAuthority, authorityId).Error; err != nil {
		return []uint{}
	}

	// 查询同级角色（相同的父角色）
	var sameLevelAuthorities []system.SysAuthority
	query := global.GVA_DB.Where("authority_id != ?", authorityId)

	// 如果是根角色，查询所有顶级角色（parent_id IS NULL 或 parent_id = 0）
	if s.isRootAuthority(currentAuthority) {
		query = query.Where("(parent_id IS NULL OR parent_id = 0)")
	} else {
		// 有父角色，查询相同父角色的同级角色
		query = query.Where("parent_id = ?", *currentAuthority.ParentId)
	}
	query.Find(&sameLevelAuthorities)

	var authorityIds []uint
	for _, auth := range sameLevelAuthorities {
		authorityIds = append(authorityIds, auth.AuthorityId)
	}

	return authorityIds
}

// GetAuditableAuthorities 获取可以审核指定申请人角色的角色ID列表
// @param applicantAuthorityId 申请人角色ID
// @return 可以审核的角色ID列表(若是根角色，则返回同级角色；否则返回所有上级角色)
func (s *EmailAuditService) GetAuditableAuthorities(applicantAuthorityId uint) []uint {
	var auditableAuthorities []uint

	// 1. 获取申请人的角色
	var applicantAuthority system.SysAuthority
	if err := global.GVA_DB.First(&applicantAuthority, applicantAuthorityId).Error; err != nil {
		return auditableAuthorities
	}

	// 2. 如果有父角色（且父角色ID不为0），父角色及其所有上级可以审核
	if !s.isRootAuthority(applicantAuthority) {
		auditableAuthorities = s.getAllParentAuthorities(*applicantAuthority.ParentId)
	} else {
		// 3. 如果是根角色（没有父角色或父角色ID为0），同级角色可以审核
		auditableAuthorities = s.getSameLevelAuthorities(applicantAuthorityId)
	}

	// 4. 排除申请人自己的角色（不能自己审核自己）
	var filtered []uint
	for _, id := range auditableAuthorities {
		if id != applicantAuthorityId {
			filtered = append(filtered, id)
		}
	}

	return filtered
}

// CanUserAudit 判断用户是否可以审核某个申请
// 注意：基于单角色设计，每个用户只能有一个角色，避免权限绕过风险
func (s *EmailAuditService) CanUserAudit(userAuthorityId uint, applicationId uint) (bool, error) {
	// 1. 获取申请
	var application gm.EmailAuditApplication
	if err := global.GVA_DB.First(&application, applicationId).Error; err != nil {
		return false, err
	}

	// 2. 获取申请人的角色ID
	var applicant system.SysUser
	if err := global.GVA_DB.First(&applicant, application.ApplicantId).Error; err != nil {
		return false, err
	}

	// 3. 获取可以审核的角色列表
	auditableAuthorities := s.GetAuditableAuthorities(applicant.AuthorityId)

	// 4. 判断当前用户角色是否在可审核列表中
	for _, id := range auditableAuthorities {
		if id == userAuthorityId {
			return true, nil
		}
	}

	return false, nil
}

// CreateApplication 创建邮件审核申请
func (s *EmailAuditService) CreateApplication(req gmReq.CreateEmailAuditRequest, applicantId uint) (*gm.EmailAuditApplication, error) {
	// 转换附件数据
	var attachments gm.JSONArray
	if req.EmailAttachments != nil {
		attachments = make(gm.JSONArray, len(req.EmailAttachments))
		for i, att := range req.EmailAttachments {
			attachments[i] = att
		}
	}

	application := &gm.EmailAuditApplication{
		ApplicantId:      applicantId,
		ApplicantTime:    time.Now(),
		Status:           gm.EmailAuditStatusPending,
		EmailType:        req.EmailType,
		EmailTitle:       convertMapToJSONMap(req.EmailTitle),
		EmailContent:     convertMapToJSONMap(req.EmailContent),
		EmailAttachments: attachments,
		EmailRemark:      req.EmailRemark,
		StartTime:        req.StartTime,
		AreaIds:          req.AreaIds,
		MaxRegTime:       req.MaxRegTime,
	}

	if err := global.GVA_DB.Create(application).Error; err != nil {
		global.GVA_LOG.Error("创建邮件审核申请失败", zap.Error(err))
		return nil, err
	}

	return application, nil
}

// UpdateApplication 更新邮件审核申请（只能更新自己的申请，且状态为待审核或待修改）
func (s *EmailAuditService) UpdateApplication(req gmReq.UpdateEmailAuditRequest, applicantId uint) error {
	var application gm.EmailAuditApplication
	if err := global.GVA_DB.First(&application, req.ID).Error; err != nil {
		return errors.New("申请不存在")
	}

	// 只能更新自己的申请
	if application.ApplicantId != applicantId {
		return errors.New("只能更新自己的申请")
	}

	// 只能更新待审核或待修改状态的申请
	if application.Status != gm.EmailAuditStatusPending && application.Status != gm.EmailAuditStatusRevision {
		return errors.New("只能更新待审核或待修改状态的申请")
	}

	// 更新字段
	if req.EmailType > 0 {
		application.EmailType = req.EmailType
	}
	if req.EmailTitle != nil {
		application.EmailTitle = convertMapToJSONMap(req.EmailTitle)
	}
	if req.EmailContent != nil {
		application.EmailContent = convertMapToJSONMap(req.EmailContent)
	}
	if req.EmailAttachments != nil {
		attachments := make(gm.JSONArray, len(req.EmailAttachments))
		for i, att := range req.EmailAttachments {
			attachments[i] = att
		}
		application.EmailAttachments = attachments
	}
	if req.EmailRemark != "" {
		application.EmailRemark = req.EmailRemark
	}
	if req.StartTime != nil {
		application.StartTime = req.StartTime
	}
	if req.AreaIds != "" {
		application.AreaIds = req.AreaIds
	}
	if req.MaxRegTime != nil {
		application.MaxRegTime = req.MaxRegTime
	}

	// 如果状态是待修改，更新后改为待审核
	if application.Status == gm.EmailAuditStatusRevision {
		application.Status = gm.EmailAuditStatusPending
		application.AuditorId = nil
		application.AuditTime = nil
		application.AuditComment = ""
	}

	if err := global.GVA_DB.Save(&application).Error; err != nil {
		global.GVA_LOG.Error("更新邮件审核申请失败", zap.Error(err))
		return err
	}

	return nil
}

// WithdrawApplication 撤回邮件审核申请（只能撤回自己的申请，且状态为待审核）
func (s *EmailAuditService) WithdrawApplication(applicationId uint, applicantId uint) error {
	var application gm.EmailAuditApplication
	if err := global.GVA_DB.First(&application, applicationId).Error; err != nil {
		return errors.New("申请不存在")
	}

	// 只能撤回自己的申请
	if application.ApplicantId != applicantId {
		return errors.New("只能撤回自己的申请")
	}

	// 只能撤回待审核状态的申请
	if application.Status != gm.EmailAuditStatusPending {
		return errors.New("只能撤回待审核状态的申请")
	}

	// 软删除
	if err := global.GVA_DB.Delete(&application).Error; err != nil {
		global.GVA_LOG.Error("撤回邮件审核申请失败", zap.Error(err))
		return err
	}

	return nil
}

// ReviewApplicationResult 审核结果
type ReviewApplicationResult struct {
	Success       bool   `json:"success"`       // 审核是否成功
	EmailSent     bool   `json:"emailSent"`     // 邮件是否发送成功（仅当审核通过时有效）
	EmailSentMsg  string `json:"emailSentMsg"`  // 邮件发送结果消息（仅当审核通过时有效）
	ReviewMessage string `json:"reviewMessage"` // 审核结果消息
}

// ReviewApplication 审核邮件申请
func (s *EmailAuditService) ReviewApplication(req gmReq.ReviewEmailAuditRequest, auditorId uint, auditorAuthorityId uint) (*ReviewApplicationResult, error) {
	result := &ReviewApplicationResult{
		Success: true,
	}
	var application gm.EmailAuditApplication
	if err := global.GVA_DB.First(&application, req.ID).Error; err != nil {
		return nil, errors.New("申请不存在")
	}

	// 检查权限
	canAudit, err := s.CanUserAudit(auditorAuthorityId, req.ID)
	if err != nil {
		return nil, errors.New("无权审核此申请")
	}
	if !canAudit {
		return nil, errors.New("无权审核此申请")
	}

	// 只能审核待审核状态的申请
	if application.Status != gm.EmailAuditStatusPending {
		return nil, errors.New("只能审核待审核状态的申请")
	}

	// 验证状态值
	status := gm.EmailAuditStatus(req.Status)
	if status != gm.EmailAuditStatusApproved && status != gm.EmailAuditStatusRejected && status != gm.EmailAuditStatusRevision {
		return nil, errors.New("无效的审核状态")
	}

	// 更新审核信息
	now := time.Now()
	application.AuditorId = &auditorId
	application.AuditTime = &now
	application.AuditComment = req.AuditComment
	application.Status = status

	if err := global.GVA_DB.Save(&application).Error; err != nil {
		global.GVA_LOG.Error("审核邮件申请失败", zap.Error(err))
		return nil, err
	}

	// 设置审核结果消息
	switch status {
	case gm.EmailAuditStatusApproved:
		result.ReviewMessage = "审核通过"
	case gm.EmailAuditStatusRejected:
		result.ReviewMessage = "审核拒绝"
	case gm.EmailAuditStatusRevision:
		result.ReviewMessage = "需要修改"
	}

	// 如果审核通过，调用游戏API发送邮件
	if status == gm.EmailAuditStatusApproved {
		if err := s.sendEmailToGame(application); err != nil {
			result.EmailSent = false
			result.EmailSentMsg = fmt.Sprintf("邮件发送失败: %v", err)
			global.GVA_LOG.Error("发送邮件到游戏服务器失败", zap.Error(err))
		} else {
			result.EmailSent = true
			result.EmailSentMsg = "邮件发送成功"
		}
	}

	return result, nil
}

// sendEmailToGame 发送邮件到游戏服务器
func (s *EmailAuditService) sendEmailToGame(application gm.EmailAuditApplication) error {
	// 构建游戏API请求数据
	gameEmailReq := s.buildGameEmailRequest(application)

	// 调用游戏API
	apiResp, err := utils.CallGameAPIPOST("/email/system/send", gameEmailReq)
	if err != nil {
		return fmt.Errorf("调用游戏API失败: %v", err)
	}

	// 检查返回的code是否为0
	if apiResp.Code != 0 {
		return fmt.Errorf("游戏API返回错误: code=%d, msg=%s", apiResp.Code, apiResp.Msg)
	}

	global.GVA_LOG.Info("邮件发送到游戏服务器成功",
		zap.Uint("application_id", application.ID),
		zap.Int("game_api_code", apiResp.Code))

	return nil
}

// buildGameEmailRequest 构建游戏API请求数据
func (s *EmailAuditService) buildGameEmailRequest(application gm.EmailAuditApplication) map[string]interface{} {
	// 转换附件格式
	var attachments []map[string]interface{}
	if application.EmailAttachments != nil {
		for _, att := range application.EmailAttachments {
			if attMap, ok := att.(map[string]interface{}); ok {
				// 确保类型正确（JSON解析时数字可能是float64）
				attachment := make(map[string]interface{})
				if id, ok := attMap["id"].(float64); ok {
					attachment["id"] = int(id)
				} else if id, ok := attMap["id"].(int); ok {
					attachment["id"] = id
				}
				if typ, ok := attMap["type"].(float64); ok {
					attachment["type"] = int(typ)
				} else if typ, ok := attMap["type"].(int); ok {
					attachment["type"] = typ
				}
				if num, ok := attMap["num"].(float64); ok {
					attachment["num"] = int(num)
				} else if num, ok := attMap["num"].(int); ok {
					attachment["num"] = num
				}
				attachments = append(attachments, attachment)
			}
		}
	}

	// 转换区服ID列表（从逗号分隔的字符串转为数组）
	var areaIds []int
	if application.AreaIds != "" {
		areaIdsStr := strings.Split(application.AreaIds, ",")
		for _, idStr := range areaIdsStr {
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				if id, err := strconv.Atoi(idStr); err == nil {
					areaIds = append(areaIds, id)
				}
			}
		}
	}

	// 构建请求数据
	req := map[string]interface{}{
		"type": application.EmailType,
		"senderI18n": map[string]string{
			"en":    "GM System Administrator",
			"zh-TW": "GM系統管理員",
			"ja":    "GMシステム管理者",
			"ko":    "GM 시스템 관리자",
		},
		"titleI18n":   application.EmailTitle,
		"contentI18n": application.EmailContent,
		"attachments": attachments,
		"remark":      application.EmailRemark,
		"startTime":   nil,
		"areaIds":     areaIds,
		"maxRegTime":  nil,
	}

	// 设置开始时间
	if application.StartTime != nil {
		req["startTime"] = *application.StartTime
	}

	// 设置最大注册时间
	if application.MaxRegTime != nil {
		req["maxRegTime"] = *application.MaxRegTime
	}

	return req
}

// GetApplicationList 获取申请列表（根据用户角色返回不同数据）
// 注意：基于单角色设计，每个用户只能有一个角色
func (s *EmailAuditService) GetApplicationList(req gmReq.SearchEmailAuditRequest, userId uint, userAuthorityId uint) ([]gm.EmailAuditApplication, int64, error) {
	var applications []gm.EmailAuditApplication
	var total int64

	query := global.GVA_DB.Model(&gm.EmailAuditApplication{})

	var userAuthority system.SysAuthority
	if err := global.GVA_DB.Where("authority_id = ?", userAuthorityId).First(&userAuthority).Error; err != nil {
		return nil, 0, errors.New("获取用户角色失败")
	}

	// 检查是否是根角色, 根据用户角色，返回不同的数据
	if !s.isRootAuthority(userAuthority) {

		// 获取所有当前用户可以审核的用户ID列表（反向：哪些用户的申请可以被当前用户审核）
		var allAuditableApplicantIds []uint

		// 获取所有用户
		var allUsers []system.SysUser
		global.GVA_DB.Select("id", "authority_id").Find(&allUsers)

		for _, user := range allUsers {
			if user.ID == userId {
				continue // 跳过自己
			}
			// 获取当前用户可以审核的申请人角色ID列表，此处会返回所有上级角色ID列表
			auditableAuthorities := s.GetAuditableAuthorities(user.AuthorityId)
			for _, authId := range auditableAuthorities {
				// 如果当前遍历的user的所有上级角色中，有一个角色ID与当前用户角色ID相同，则将该user.ID添加到可审核用户ID列表中
				if authId == userAuthorityId {
					allAuditableApplicantIds = append(allAuditableApplicantIds, user.ID)
					break
				}
			}
		}

		if len(allAuditableApplicantIds) > 0 {
			// 如果可审核用户ID列表不为空，则查询条件为：申请人ID为当前用户ID，或者状态为待审核且申请人ID在当前用户可审核的用户ID列表中
			query = query.Where("applicant_id = ? OR (status = ? AND applicant_id IN ?)",
				userId, gm.EmailAuditStatusPending, allAuditableApplicantIds)
		} else {
			// 如果没有可审核的申请人，只能看自己的
			query = query.Where("applicant_id = ?", userId)
		}
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 申请人筛选
	if req.ApplicantId != nil {
		query = query.Where("applicant_id = ?", *req.ApplicantId)
	}

	// 审核人筛选
	if req.AuditorId != nil {
		query = query.Where("auditor_id = ?", *req.AuditorId)
	}

	// 时间范围筛选
	if req.StartTime != nil {
		query = query.Where("applicant_time >= ?", time.Unix(*req.StartTime, 0))
	}
	if req.EndTime != nil {
		query = query.Where("applicant_time <= ?", time.Unix(*req.EndTime, 0))
	}

	// 获取总数
	query.Count(&total)

	// 分页
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 查询列表（预加载关联数据）
	if err := query.Preload("Applicant").Preload("Auditor").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		global.GVA_LOG.Error("获取申请列表失败", zap.Error(err))
		return nil, 0, err
	}

	return applications, total, nil
}

// GetApplicationListResponse 获取申请列表（返回响应 DTO）
func (s *EmailAuditService) GetApplicationListResponse(req gmReq.SearchEmailAuditRequest, userId uint, userAuthorityId uint) ([]gmResp.EmailAuditApplicationResponse, int64, error) {
	applications, total, err := s.GetApplicationList(req, userId, userAuthorityId)
	if err != nil {
		return nil, 0, err
	}
	return gmResp.ToEmailAuditApplicationResponseList(applications), total, nil
}

// GetApplication 获取申请详情
func (s *EmailAuditService) GetApplication(id uint, userId uint, userAuthorityId uint) (*gm.EmailAuditApplication, error) {
	var application gm.EmailAuditApplication
	if err := global.GVA_DB.Preload("Applicant").Preload("Auditor").First(&application, id).Error; err != nil {
		return nil, errors.New("申请不存在")
	}

	// 权限检查：申请人或审核人可以查看
	canView := application.ApplicantId == userId

	// 检查是否有审核权限
	if !canView {
		canAudit, err := s.CanUserAudit(userAuthorityId, id)
		if err == nil && canAudit {
			canView = true
		}
	}

	// 超级管理员可以查看所有
	if !canView {
		var user system.SysUser
		if err := global.GVA_DB.First(&user, userId).Error; err == nil {
			var authority system.SysAuthority
			if err := global.GVA_DB.First(&authority, user.AuthorityId).Error; err == nil {
				if s.isRootAuthority(authority) {
					canView = true
				}
			}
		}
	}

	if !canView {
		return nil, errors.New("无权查看此申请")
	}

	return &application, nil
}

// GetApplicationResponse 获取申请详情（返回响应 DTO）
func (s *EmailAuditService) GetApplicationResponse(id uint, userId uint, userAuthorityId uint) (*gmResp.EmailAuditApplicationResponse, error) {
	application, err := s.GetApplication(id, userId, userAuthorityId)
	if err != nil {
		return nil, err
	}
	return gmResp.ToEmailAuditApplicationResponse(application), nil
}

// convertMapToJSONMap 转换 map[string]string 为 JSONMap
func convertMapToJSONMap(m map[string]string) gm.JSONMap {
	result := make(gm.JSONMap)
	for k, v := range m {
		result[k] = v
	}
	return result
}
