# 邮件审核权限设计文档

## 一、审核权限设计方案分析

### 方案对比

#### 方案1：基于角色层级审核（推荐）⭐

**原理**：
- 下级角色的申请由上级角色审核
- 利用现有的角色层次结构（ParentId）
- 审核人自己申请时，由同级或上级审核（交叉审核）

**示例**：
```
角色层级结构：
超级管理员 (AuthorityId: 1)
├── 运营管理员 (AuthorityId: 2)
│   ├── 运营专员 (AuthorityId: 3)
│   └── 客服专员 (AuthorityId: 4)
└── 内容管理员 (AuthorityId: 5)
    └── 内容编辑 (AuthorityId: 6)
```

**审核规则**：
1. 运营专员申请 → 运营管理员审核
2. 运营管理员申请 → 超级管理员审核
3. 超级管理员申请 → 由其他超级管理员或同级审核

**优点**：
- ✅ 符合组织层级结构
- ✅ 实现简单，利用现有数据结构
- ✅ 权限清晰，职责明确
- ✅ 可扩展性强

**缺点**：
- ⚠️ 需要维护角色层级关系
- ⚠️ 同级之间需要额外处理

#### 方案2：基于角色配置审核

**原理**：
- 配置表：哪些角色可以审核哪些角色
- 需要额外的配置表和配置接口

**缺点**：
- ❌ 需要额外的配置管理
- ❌ 配置复杂度高
- ❌ 不利用现有数据结构

#### 方案3：独立审核组

**原理**：
- 创建独立的审核角色组
- 所有申请都由审核组审核

**缺点**：
- ❌ 不符合实际业务场景
- ❌ 无法体现层级关系

---

## 二、推荐方案详细设计（方案1：角色层级 + 交叉审核）

### 2.1 审核权限确定规则

```
函数：GetAuditAuthorities(applicantAuthorityId)
1. 获取申请人角色
2. 查找父角色（ParentId）
   - 如果有父角色 → 父角色及其所有上级角色都可以审核
   - 如果没有父角色 → 同级角色可以审核
3. 特殊情况：超级管理员申请
   - 由其他超级管理员审核（避免自审）
```

### 2.2 具体实现逻辑

```go
// 获取可以审核该申请的角色列表
func GetAuditableAuthorities(applicantAuthorityId uint) []uint {
    var authorities []uint
    
    // 1. 获取申请人角色
    var applicantAuthority system.SysAuthority
    db.First(&applicantAuthority, applicantAuthorityId)
    
    // 2. 如果有父角色，父角色及其上级可以审核
    if applicantAuthority.ParentId != nil {
        // 递归获取所有上级角色
        authorities = getAllParentAuthorities(*applicantAuthority.ParentId)
    } else {
        // 3. 如果没有父角色，同级角色可以审核
        var sameLevelAuthorities []system.SysAuthority
        db.Where("parent_id IS NULL").Find(&sameLevelAuthorities)
        for _, auth := range sameLevelAuthorities {
            if auth.AuthorityId != applicantAuthorityId {
                authorities = append(authorities, auth.AuthorityId)
            }
        }
    }
    
    // 4. 特殊情况：不能自己审核自己
    authorities = removeElement(authorities, applicantAuthorityId)
    
    return authorities
}
```

### 2.3 交叉审核处理

**场景**：审核人A自己发送申请

**处理方式**：
1. 正常情况：A的申请由A的上级审核
2. 同级情况：A的申请由A的同级审核（排除A自己）
3. 顶级情况：超级管理员的申请由其他超级管理员审核

**实现**：
```go
// 判断用户是否可以审核某个申请
func CanUserAudit(userAuthorityId uint, applicationId uint) bool {
    // 1. 获取申请
    var application EmailAuditApplication
    db.First(&application, applicationId)
    
    // 2. 获取申请人的角色
    applicantAuthorityId := getApplicantAuthorityId(application.ApplicantId)
    
    // 3. 获取可以审核的角色列表
    auditableAuthorities := GetAuditableAuthorities(applicantAuthorityId)
    
    // 4. 判断当前用户角色是否在可审核列表中
    return contains(auditableAuthorities, userAuthorityId)
}
```

---

## 三、数据库设计

### 3.1 邮件审核申请表

```sql
CREATE TABLE `email_audit_applications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `applicant_id` bigint unsigned NOT NULL COMMENT '申请人ID（关联sys_users.id）',
  `applicant_time` datetime NOT NULL COMMENT '申请时间',
  `auditor_id` bigint unsigned DEFAULT NULL COMMENT '审核人ID（关联sys_users.id）',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `audit_comment` text COMMENT '审核说明',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1-待审核，2-通过，3-拒绝，4-待修改',
  `email_type` tinyint NOT NULL COMMENT '邮件类型',
  `email_title` json NOT NULL COMMENT '邮件标题（多语言）',
  `email_content` json NOT NULL COMMENT '邮件内容（多语言）',
  `email_attachments` json DEFAULT NULL COMMENT '邮件附件',
  `email_remark` varchar(500) DEFAULT NULL COMMENT '邮件备注',
  `start_time` bigint DEFAULT NULL COMMENT '开始生效时间（时间戳）',
  `area_ids` varchar(500) DEFAULT NULL COMMENT '生效区服列表（逗号分隔）',
  `max_reg_time` bigint DEFAULT NULL COMMENT '最大注册时间（时间戳）',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_applicant_id` (`applicant_id`),
  KEY `idx_auditor_id` (`auditor_id`),
  KEY `idx_status` (`status`),
  KEY `idx_applicant_time` (`applicant_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件审核申请表';
```

### 3.2 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| applicant_id | bigint | 申请人ID（sys_users.id） |
| applicant_time | datetime | 申请时间 |
| auditor_id | bigint | 审核人ID（可为空，审核时填入） |
| audit_time | datetime | 审核时间（可为空） |
| audit_comment | text | 审核说明 |
| status | tinyint | 状态：1-待审核，2-通过，3-拒绝，4-待修改 |
| email_type | tinyint | 邮件类型 |
| email_title | json | 邮件标题（多语言对象） |
| email_content | json | 邮件内容（多语言对象） |
| email_attachments | json | 邮件附件数组 |
| email_remark | varchar | 邮件备注 |
| start_time | bigint | 开始生效时间（秒级时间戳） |
| area_ids | varchar | 生效区服列表（逗号分隔） |
| max_reg_time | bigint | 最大注册时间（秒级时间戳） |

---

## 四、状态流转

```
待审核 (1)
  ↓ 审核通过
通过 (2) → 直接发送邮件
  ↓ 审核拒绝
拒绝 (3) → 申请结束
  ↓ 需要修改
待修改 (4) → 可以重新编辑申请
  ↓ 修改后重新提交
待审核 (1)
```

---

## 五、API 接口设计

### 5.1 新增邮件申请
- **路径**：`POST /gm/email/audit/apply`
- **权限**：所有登录用户
- **参数**：邮件详情（标题、内容、附件等）

### 5.2 修改邮件申请
- **路径**：`PUT /gm/email/audit/apply/:id`
- **权限**：只能修改自己的申请，且状态为"待审核"或"待修改"
- **参数**：邮件详情

### 5.3 撤回邮件申请
- **路径**：`POST /gm/email/audit/withdraw/:id`
- **权限**：只能撤回自己的申请，且状态为"待审核"
- **逻辑**：软删除或状态改为"已撤回"

### 5.4 审核邮件申请
- **路径**：`POST /gm/email/audit/review/:id`
- **权限**：需要审核权限（通过 CanUserAudit 检查）
- **参数**：审核结果（通过/拒绝/待修改）、审核说明

### 5.5 获取申请列表
- **路径**：`POST /gm/email/audit/list`
- **权限**：根据角色返回不同数据
  - 申请人：只能看到自己的申请
  - 审核人：可以看到待审核的申请
  - 管理员：可以看到所有申请

### 5.6 获取申请详情
- **路径**：`GET /gm/email/audit/:id`
- **权限**：申请人或审核人可以查看

---

## 六、权限检查流程

### 6.1 申请时
```
1. 用户登录 → 获取用户ID和角色ID
2. 创建申请 → 保存 applicant_id 和 applicant_time
3. 状态初始化为"待审核"
```

### 6.2 审核时
```
1. 获取当前用户角色ID
2. 获取申请的申请人角色ID
3. 调用 CanUserAudit 检查权限
4. 如果有权限 → 允许审核
5. 如果没有权限 → 返回"无权审核"
```

### 6.3 查看列表时
```
1. 获取当前用户角色ID
2. 如果是申请人 → 查询 applicant_id = 当前用户ID
3. 如果是审核人 → 查询状态=待审核 AND 当前用户角色在可审核列表中
4. 如果是超级管理员 → 查询所有
```

---

## 七、最佳实践建议

### 7.1 权限设计原则
1. **最小权限原则**：只授予必要的审核权限
2. **职责分离**：申请人和审核人分离
3. **不可自审**：自己不能审核自己的申请
4. **层级清晰**：利用现有角色层级结构

### 7.2 安全建议
1. **审核操作记录**：记录所有审核操作日志
2. **状态变更记录**：记录状态变更历史
3. **权限验证**：每个接口都要验证权限
4. **数据校验**：严格校验输入数据

### 7.3 扩展建议
1. **多级审核**：支持多级审核流程
2. **审核通知**：审核时发送通知
3. **审核统计**：统计审核效率和通过率
4. **审核模板**：支持审核意见模板

---

## 八、总结

**推荐方案**：角色层级 + 交叉审核（方案1）

**核心优势**：
- ✅ 利用现有数据结构（角色层级）
- ✅ 实现简单，维护成本低
- ✅ 符合实际业务场景
- ✅ 可扩展性强

**关键实现点**：
1. 获取父角色列表（递归）
2. 判断同级角色（ParentId 相同）
3. 排除自己审核自己
4. 在 API 层做权限检查

