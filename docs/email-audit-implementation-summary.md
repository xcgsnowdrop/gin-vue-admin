# 邮件审核功能实现总结

## 一、功能概述

已实现完整的邮件审核功能，包括：
- ✅ 新增邮件申请
- ✅ 修改邮件申请
- ✅ 撤回邮件申请
- ✅ 邮件审核（通过/拒绝/待修改）
- ✅ 申请列表查询
- ✅ 申请详情查询

## 二、文件结构

### Model 层
- `server/model/gm/email_audit.go` - 数据模型定义
- `server/model/gm/request/email_audit.go` - 请求参数定义

### Service 层
- `server/service/gm/email_audit.go` - 业务逻辑实现（包含权限判断）

### API 层
- `server/api/v1/gm/email_audit.go` - HTTP 接口处理

### Router 层
- `server/router/gm/email_audit.go` - 路由配置

### 注册文件
- `server/api/v1/gm/enter.go` - API 组注册
- `server/service/gm/enter.go` - Service 组注册
- `server/router/gm/enter.go` - Router 组注册
- `server/initialize/router.go` - 路由初始化

## 三、数据库设计

### 表名：`email_audit_applications`

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键ID |
| applicant_id | bigint | 申请人ID（关联sys_users.id） |
| applicant_time | datetime | 申请时间 |
| auditor_id | bigint | 审核人ID（可为空） |
| audit_time | datetime | 审核时间（可为空） |
| audit_comment | text | 审核说明 |
| status | tinyint | 状态：1-待审核，2-通过，3-拒绝，4-待修改 |
| email_type | tinyint | 邮件类型 |
| email_title | json | 邮件标题（多语言JSON） |
| email_content | json | 邮件内容（多语言JSON） |
| email_attachments | json | 邮件附件（JSON数组） |
| email_remark | varchar(500) | 邮件备注 |
| start_time | bigint | 开始生效时间（秒级时间戳） |
| area_ids | varchar(500) | 生效区服列表（逗号分隔） |
| max_reg_time | bigint | 最大注册时间（秒级时间戳） |

**创建SQL**：见 `server/docs/email_audit_migration.sql`

## 四、审核权限设计（核心）

### 权限确定规则

采用**角色层级 + 交叉审核**方案：

1. **下级申请由上级审核**：
   - 如果申请人有父角色 → 父角色及其所有上级角色可以审核

2. **同级交叉审核**：
   - 如果申请人没有父角色（顶级角色）→ 同级角色可以审核（排除自己）

3. **不能自己审核自己**：
   - 申请人不能审核自己的申请

### 实现逻辑

```go
// 获取可以审核指定申请的角色列表
func GetAuditableAuthorities(applicantAuthorityId uint) []uint {
    // 1. 获取申请人的角色
    // 2. 如果有父角色 → 获取所有父级角色
    // 3. 如果没有父角色 → 获取同级角色（排除自己）
    // 4. 排除申请人自己的角色
}
```

### 交叉审核处理

**场景**：审核人A自己发送申请

**处理方式**：
- A的申请由A的上级审核（正常流程）
- 如果A是顶级角色，由其他同级角色审核
- 确保不能自己审核自己

## 五、API 接口

### 1. 创建邮件审核申请
- **路径**：`POST /gm/email/audit/apply`
- **权限**：所有登录用户
- **参数**：`CreateEmailAuditRequest`

### 2. 更新邮件审核申请
- **路径**：`PUT /gm/email/audit/apply`
- **权限**：只能更新自己的申请，且状态为"待审核"或"待修改"

### 3. 撤回邮件审核申请
- **路径**：`POST /gm/email/audit/withdraw/:id`
- **权限**：只能撤回自己的申请，且状态为"待审核"

### 4. 审核邮件申请
- **路径**：`POST /gm/email/audit/review`
- **权限**：需要有审核权限（通过权限检查）
- **参数**：`ReviewEmailAuditRequest`

### 5. 获取申请列表
- **路径**：`POST /gm/email/audit/list`
- **权限**：根据角色返回不同数据
  - 申请人：看到自己的申请
  - 审核人：看到待审核的申请
  - 管理员：看到所有申请

### 6. 获取申请详情
- **路径**：`GET /gm/email/audit/:id`
- **权限**：申请人或审核人可以查看

## 六、状态流转

```
待审核 (1)
  ↓ 审核通过
通过 (2) → 可以发送邮件
  ↓ 审核拒绝
拒绝 (3) → 申请结束
  ↓ 需要修改
待修改 (4) → 可以重新编辑申请
  ↓ 修改后重新提交
待审核 (1)
```

## 七、权限检查流程

### 申请时
1. 获取当前登录用户ID
2. 创建申请，保存 `applicant_id` 和 `applicant_time`
3. 状态初始化为"待审核" (1)

### 审核时
1. 获取当前用户角色ID
2. 获取申请人的角色ID
3. 调用 `CanUserAudit` 检查权限
4. 验证通过 → 允许审核
5. 验证失败 → 返回"无权审核"

### 查看列表时
1. 获取当前用户角色ID
2. 判断是否是管理员
3. 普通用户：只看到自己的申请 + 需要自己审核的申请
4. 管理员：看到所有申请

## 八、使用示例

### 创建申请
```json
POST /gm/email/audit/apply
{
  "emailType": 1,
  "emailTitle": {
    "en": "Test Email",
    "ja": "テストメール"
  },
  "emailContent": {
    "en": "Content",
    "ja": "内容"
  },
  "emailAttachments": [
    {"id": 1, "type": 1, "num": 100}
  ],
  "startTime": 1761926400,
  "areaIds": "1,2,3"
}
```

### 审核申请
```json
POST /gm/email/audit/review
{
  "id": 1,
  "status": 2,  // 2-通过，3-拒绝，4-待修改
  "auditComment": "审核通过"
}
```

## 九、注意事项

1. **数据库迁移**：需要先执行 `server/docs/email_audit_migration.sql` 创建表

2. **权限配置**：确保角色层级关系（ParentId）配置正确

3. **状态管理**：
   - 只有"待审核"状态的申请可以被审核
   - 只有"待审核"或"待修改"状态的申请可以被修改
   - 只有"待审核"状态的申请可以被撤回

4. **数据安全**：
   - 所有接口都需要 JWT 认证
   - 权限检查在 Service 层实现
   - 操作记录通过中间件记录

## 十、后续优化建议

1. **多级审核**：支持多级审核流程（一级、二级审核）
2. **审核通知**：审核时发送通知给申请人
3. **审核统计**：统计审核效率和通过率
4. **审核模板**：支持审核意见模板
5. **批量审核**：支持批量审核功能

