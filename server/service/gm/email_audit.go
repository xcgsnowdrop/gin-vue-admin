package gm

import (
	"errors"
	"time"

	"gmserver/global"
	"gmserver/model/gm"
	gmReq "gmserver/model/gm/request"
	"gmserver/model/system"

	"go.uber.org/zap"
)

type EmailAuditService struct{}

// isRootAuthority 判断是否是根角色（没有父角色）
// parentId 为 nil 或指向 0 时，表示根角色
func (s *EmailAuditService) isRootAuthority(parentId *uint) bool {
	if parentId == nil {
		return true
	}
	return *parentId == 0
}

// GetAllParentAuthorities 递归获取所有父级角色ID
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
		if s.isRootAuthority(authority.ParentId) {
			break
		}

		currentId = *authority.ParentId
	}

	return authorityIds
}

// GetSameLevelAuthorities 获取同级角色ID（排除指定角色）
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
	if s.isRootAuthority(currentAuthority.ParentId) {
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

// GetAuditableAuthorities 获取可以审核指定申请的角色ID列表
func (s *EmailAuditService) GetAuditableAuthorities(applicantAuthorityId uint) []uint {
	var auditableAuthorities []uint

	// 1. 获取申请人的角色
	var applicantAuthority system.SysAuthority
	if err := global.GVA_DB.First(&applicantAuthority, applicantAuthorityId).Error; err != nil {
		return auditableAuthorities
	}

	// 2. 如果有父角色（且父角色ID不为0），父角色及其所有上级可以审核
	if !s.isRootAuthority(applicantAuthority.ParentId) {
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

// ReviewApplication 审核邮件申请
func (s *EmailAuditService) ReviewApplication(req gmReq.ReviewEmailAuditRequest, auditorId uint, auditorAuthorityId uint) error {
	var application gm.EmailAuditApplication
	if err := global.GVA_DB.First(&application, req.ID).Error; err != nil {
		return errors.New("申请不存在")
	}

	// 检查权限
	canAudit, err := s.CanUserAudit(auditorAuthorityId, req.ID)
	if err != nil {
		return err
	}
	if !canAudit {
		return errors.New("无权审核此申请")
	}

	// 只能审核待审核状态的申请
	if application.Status != gm.EmailAuditStatusPending {
		return errors.New("只能审核待审核状态的申请")
	}

	// 验证状态值
	status := gm.EmailAuditStatus(req.Status)
	if status != gm.EmailAuditStatusApproved && status != gm.EmailAuditStatusRejected && status != gm.EmailAuditStatusRevision {
		return errors.New("无效的审核状态")
	}

	// 更新审核信息
	now := time.Now()
	application.AuditorId = &auditorId
	application.AuditTime = &now
	application.AuditComment = req.AuditComment
	application.Status = status

	if err := global.GVA_DB.Save(&application).Error; err != nil {
		global.GVA_LOG.Error("审核邮件申请失败", zap.Error(err))
		return err
	}

	return nil
}

// GetApplicationList 获取申请列表（根据用户角色返回不同数据）
func (s *EmailAuditService) GetApplicationList(req gmReq.SearchEmailAuditRequest, userId uint, userAuthorityId uint) ([]gm.EmailAuditApplication, int64, error) {
	var applications []gm.EmailAuditApplication
	var total int64

	query := global.GVA_DB.Model(&gm.EmailAuditApplication{})

	// 判断用户角色，返回不同的数据
	// 这里简化处理：普通用户只能看自己的，管理员可以看到所有
	// 实际可以根据角色配置更细粒度的权限
	var isAdmin bool
	var superAdmin system.SysAuthority
	if err := global.GVA_DB.Where("authority_id = ?", userAuthorityId).First(&superAdmin).Error; err == nil {
		// 检查是否是超级管理员（假设 authority_id = 1 或根角色）
		if s.isRootAuthority(superAdmin.ParentId) || userAuthorityId == 1 {
			isAdmin = true
		}
	}

	if !isAdmin {
		// 普通用户：只能看到自己的申请 + 需要自己审核的申请
		// 获取所有可以审核的角色列表（反向：哪些角色的申请可以被当前用户审核）
		// 思路：遍历所有角色，检查当前用户是否可以审核该角色的申请
		var allAuditableApplicantIds []uint

		// 获取所有用户
		var allUsers []system.SysUser
		global.GVA_DB.Select("id", "authority_id").Find(&allUsers)

		for _, user := range allUsers {
			if user.ID == userId {
				continue // 跳过自己
			}
			auditableAuthorities := s.GetAuditableAuthorities(user.AuthorityId)
			for _, authId := range auditableAuthorities {
				if authId == userAuthorityId {
					allAuditableApplicantIds = append(allAuditableApplicantIds, user.ID)
					break
				}
			}
		}

		if len(allAuditableApplicantIds) > 0 {
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

	// 查询列表
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		global.GVA_LOG.Error("获取申请列表失败", zap.Error(err))
		return nil, 0, err
	}

	return applications, total, nil
}

// GetApplication 获取申请详情
func (s *EmailAuditService) GetApplication(id uint, userId uint, userAuthorityId uint) (*gm.EmailAuditApplication, error) {
	var application gm.EmailAuditApplication
	if err := global.GVA_DB.First(&application, id).Error; err != nil {
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
				if s.isRootAuthority(authority.ParentId) || authority.AuthorityId == 1 {
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

// convertMapToJSONMap 转换 map[string]string 为 JSONMap
func convertMapToJSONMap(m map[string]string) gm.JSONMap {
	result := make(gm.JSONMap)
	for k, v := range m {
		result[k] = v
	}
	return result
}
