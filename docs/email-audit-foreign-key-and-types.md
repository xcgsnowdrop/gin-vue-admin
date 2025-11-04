# 邮件审核功能 - 外键约束和字段类型说明

## 一、MySQL 外键约束详解

### 1. 什么是外键（Foreign Key）？

**外键**是数据库层面的约束，用于建立两个表之间的关联关系，确保数据的**引用完整性（Referential Integrity）**。

### 2. 外键的作用

- **数据完整性**：确保外键字段的值必须存在于被引用表的主键中
- **防止孤立数据**：避免出现"孤儿记录"（没有对应主表记录的子表记录）
- **级联操作**：可以定义删除或更新时的级联行为

### 3. 外键约束语法

```sql
CONSTRAINT `约束名` FOREIGN KEY (`外键字段`) 
  REFERENCES `被引用表` (`被引用字段`) 
  ON DELETE [RESTRICT|CASCADE|SET NULL|NO ACTION]
  ON UPDATE [RESTRICT|CASCADE|SET NULL|NO ACTION]
```

### 4. 级联操作说明

#### ON DELETE（删除时的行为）

- **RESTRICT**（默认）：阻止删除被引用的记录（如果还有子表记录引用它）
- **CASCADE**：级联删除，删除被引用记录时，同时删除所有引用它的子表记录
- **SET NULL**：删除被引用记录时，将引用它的子表记录的外键字段设置为 NULL
- **NO ACTION**：不执行任何操作（类似 RESTRICT）

#### ON UPDATE（更新时的行为）

- **RESTRICT**：阻止更新被引用的主键
- **CASCADE**：级联更新，更新主键时，同时更新所有引用它的外键
- **SET NULL**：更新主键时，将引用它的外键设置为 NULL
- **NO ACTION**：不执行任何操作

### 5. 当前实现的外键约束

```sql
-- applicant_id 外键约束
CONSTRAINT `fk_email_audit_applicant` FOREIGN KEY (`applicant_id`) 
  REFERENCES `sys_users` (`id`) 
  ON DELETE RESTRICT   -- 阻止删除申请人（如果该用户有审核申请）
  ON UPDATE CASCADE    -- 更新用户ID时，级联更新申请表中的申请人ID

-- auditor_id 外键约束
CONSTRAINT `fk_email_audit_auditor` FOREIGN KEY (`auditor_id`) 
  REFERENCES `sys_users` (`id`) 
  ON DELETE SET NULL   -- 删除审核人时，将审核人ID设为NULL（审核记录保留）
  ON UPDATE CASCADE    -- 更新用户ID时，级联更新申请表中的审核人ID
```

**设计说明：**
- `applicant_id` 使用 `RESTRICT`：因为申请必须要有申请人，如果删除申请人会导致申请记录无效
- `auditor_id` 使用 `SET NULL`：因为审核人可能被删除，但审核记录应该保留（历史记录）

## 二、GORM 字段类型选择原则

### 1. 数值类型

#### 必填字段：使用值类型
```go
ApplicantId uint `json:"applicantId" gorm:"not null"`
```
- 数据库字段：`NOT NULL`
- Go 类型：`uint`（值类型）
- 零值：`0`（但业务上应该避免，通过 `not null` 约束）

#### 可选字段：使用指针类型
```go
AuditorId *uint `json:"auditorId" gorm:"index"`
```
- 数据库字段：`DEFAULT NULL`（允许 NULL）
- Go 类型：`*uint`（指针类型）
- 零值：`nil`（表示未设置）

### 2. 时间类型

#### 必填时间：使用值类型
```go
ApplicantTime time.Time `json:"applicantTime" gorm:"not null"`
```

#### 可选时间：使用指针类型
```go
AuditTime *time.Time `json:"auditTime"`
```

### 3. 字符串类型

**建议统一使用值类型 `string`**，即使可以为空：

```go
EmailRemark string `json:"emailRemark" gorm:"type:varchar(500)"`
AuditComment string `json:"auditComment" gorm:"type:text"`
```

**原因：**
- 空字符串 `""` 可以表示"未填写"
- 避免指针类型的空指针检查
- 数据库层面可以为 `NULL`，但 Go 中统一用空字符串处理更简单

### 4. 关联对象类型

#### 必填关联：使用值类型
```go
Applicant system.SysUser `json:"applicant" gorm:"foreignKey:ApplicantId;references:ID"`
```
- 申请人总是存在的（因为 `applicant_id` 是 `NOT NULL`）
- 使用值类型，GORM 会自动填充关联数据

#### 可选关联：使用指针类型
```go
Auditor *system.SysUser `json:"auditor" gorm:"foreignKey:AuditorId;references:ID"`
```
- 审核人可能不存在（因为 `auditor_id` 可以为 `NULL`）
- 使用指针类型，可以明确区分"未设置"（`nil`）和"已设置但关联对象为空"

### 5. 字段类型选择总结表

| 字段类型 | 数据库约束 | Go 类型 | 说明 |
|---------|-----------|---------|------|
| 必填数值 | `NOT NULL` | `uint`, `int` | 值类型 |
| 可选数值 | `DEFAULT NULL` | `*uint`, `*int` | 指针类型，`nil` 表示未设置 |
| 必填时间 | `NOT NULL` | `time.Time` | 值类型 |
| 可选时间 | `DEFAULT NULL` | `*time.Time` | 指针类型，`nil` 表示未设置 |
| 字符串 | 可为 `NULL` | `string` | 值类型，空字符串表示空 |
| 必填关联 | `NOT NULL` | `Struct` | 值类型，GORM 自动填充 |
| 可选关联 | `DEFAULT NULL` | `*Struct` | 指针类型，`nil` 表示未关联 |

## 三、当前 EmailAuditApplication 字段类型说明

```go
type EmailAuditApplication struct {
    // 主键 - 必填
    ID uint `json:"id" gorm:"primarykey"`
    
    // 外键字段 - 必填
    ApplicantId uint `json:"applicantId" gorm:"not null;index"`
    
    // 关联对象 - 必填（因为 ApplicantId 必填）
    Applicant system.SysUser `json:"applicant" gorm:"foreignKey:ApplicantId;references:ID"`
    
    // 外键字段 - 可选
    AuditorId *uint `json:"auditorId" gorm:"index"`
    
    // 关联对象 - 可选（因为 AuditorId 可选）
    Auditor *system.SysUser `json:"auditor" gorm:"foreignKey:AuditorId;references:ID"`
    
    // 时间字段 - 必填
    ApplicantTime time.Time `json:"applicantTime" gorm:"not null"`
    
    // 时间字段 - 可选
    AuditTime *time.Time `json:"auditTime"`
    
    // 字符串字段 - 可为空（使用空字符串表示）
    AuditComment string `json:"auditComment" gorm:"type:text"`
    EmailRemark string `json:"emailRemark" gorm:"type:varchar(500)"`
    
    // 数值字段 - 可选
    StartTime *int64 `json:"startTime"`
    MaxRegTime *int64 `json:"maxRegTime"`
}
```

## 四、GORM 外键关联说明

### 1. GORM 是否需要明确指定外键？

**不是必须的，但明确指定更好**：

- **自动推断**：GORM 可以根据命名约定自动推断
  - `Applicant` 关联 → 自动推断 `foreignKey:ApplicantID`
  - `SysUser.ID` → 自动推断 `references:ID`
  
- **明确指定**：更清晰，避免歧义
  ```go
  Applicant system.SysUser `gorm:"foreignKey:ApplicantId;references:ID"`
  ```

### 2. foreignKey 和 references 说明

- `foreignKey:ApplicantId`：当前表（`email_audit_applications`）中的外键字段
- `references:ID`：关联表（`sys_users`）中的主键字段（`SysUser.ID`）

### 3. 使用关联字段

```go
// 查询时预加载关联数据
var application EmailAuditApplication
db.Preload("Applicant").Preload("Auditor").First(&application, id)

// 访问关联数据
fmt.Println(application.Applicant.Username)  // 申请人用户名
if application.Auditor != nil {
    fmt.Println(application.Auditor.Username)  // 审核人用户名
}
```

## 五、注意事项

1. **数据库外键约束 vs GORM 关联**
   - GORM 的 `foreignKey` 和 `references` 只是告诉 GORM 如何关联数据
   - **不会自动创建数据库外键约束**
   - 数据库外键约束需要在 SQL 迁移文件中手动创建

2. **循环导入问题**
   - `EmailAuditApplication` 引用 `system.SysUser`
   - 避免 `SysUser` 引用 `EmailAuditApplication`，否则会产生循环导入

3. **关联字段不存储到数据库**
   - `Applicant` 和 `Auditor` 是关联对象，不会在数据库中创建列
   - 数据库中只有 `applicant_id` 和 `auditor_id` 字段

4. **性能考虑**
   - 使用 `Preload` 预加载关联数据，避免 N+1 查询问题
   - 不需要关联数据时，不要使用 `Preload`

