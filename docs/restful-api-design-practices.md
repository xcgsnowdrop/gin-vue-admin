# RESTful API 设计最佳实践

## 一、HTTP 方法使用规范

### HTTP 方法语义

| 方法 | 语义 | 幂等性 | 用途 |
|------|------|--------|------|
| GET | 获取资源 | ✅ 是 | 查询数据，参数在 URL 中 |
| POST | 创建资源或执行操作 | ❌ 否 | 创建新资源，执行动作，数据在 Body 中 |
| PUT | 更新资源（全量） | ✅ 是 | 替换整个资源，数据在 Body 中 |
| PATCH | 更新资源（部分） | ❌ 否 | 部分更新，数据在 Body 中 |
| DELETE | 删除资源 | ✅ 是 | 删除资源，ID 在 URL 中 |

## 二、参数传递方式

### URL 参数 vs Body 参数

#### GET 方法
```
GET /api/users/123                    # 路径参数
GET /api/users?page=1&pageSize=10     # 查询参数
```
**原则**：GET 方法的参数**总是在 URL 中**（路径参数或查询参数）

#### POST 方法
```
POST /api/users
Body: { "name": "张三", "age": 25 }    # 复杂数据在 Body 中

POST /api/users/123/activate           # 简单标识符在 URL 中
Body: {} 或 { "reason": "激活" }       # 可选的操作参数在 Body 中
```

#### DELETE 方法
```
DELETE /api/users/123                  # 资源 ID 在 URL 中
```
**原则**：DELETE 方法的资源 ID**总是在 URL 中**

## 三、当前设计分析

### 撤回操作的当前实现

```go
POST /api/gm/email/audit/withdraw/:id
```

**问题点**：
1. ✅ URL 中有路径参数（符合 RESTful）
2. ⚠️ 使用了 POST 而不是 DELETE
3. ⚠️ Body 为空，参数只在 URL 中

### 为什么当前设计有问题？

1. **语义不清晰**：
   - POST 通常用于创建或执行动作
   - 撤回更像是一个状态变更操作，不是创建

2. **不符合 RESTful 规范**：
   - 如果是真正的删除 → 应该用 DELETE
   - 如果是状态变更 → 可以用 POST，但通常用 PUT/PATCH 更合适

3. **参数传递不一致**：
   - POST 方法的数据通常在 Body 中
   - 虽然 URL 参数也可以，但不是最佳实践

## 四、改进方案

### 方案1：使用 DELETE 方法（推荐）⭐

**适用场景**：如果撤回就是删除申请（软删除）

```go
DELETE /api/gm/email/audit/:id
```

**优点**：
- ✅ 语义清晰：DELETE 表示删除操作
- ✅ 符合 RESTful 规范
- ✅ ID 在 URL 中是标准做法

**缺点**：
- ⚠️ 如果撤回不是真正的删除，语义可能不准确

**实现**：
```go
func (s *EmailAuditRouter) InitEmailAuditRouter(Router *gin.RouterGroup) {
    // ...
    emailAuditRouter.DELETE(":id", gmEmailAuditApi.WithdrawApplication) // 撤回
}
```

### 方案2：使用 POST + Body 参数

**适用场景**：撤回是业务动作，不是简单的删除

```go
POST /api/gm/email/audit/withdraw
Body: { "id": 123, "reason": "撤回原因" }  // ID 在 Body 中
```

**优点**：
- ✅ 语义明确：POST 表示执行动作
- ✅ 可以传递额外参数（如撤回原因）
- ✅ 更灵活，易于扩展

**缺点**：
- ⚠️ ID 不在 URL 中，不够 RESTful

**实现**：
```go
// Request
type WithdrawEmailAuditRequest struct {
    ID     uint   `json:"id" binding:"required"`
    Reason string `json:"reason"` // 可选
}

// Router
emailAuditRouter.POST("withdraw", gmEmailAuditApi.WithdrawApplication)

// API
func (e *EmailAuditApi) WithdrawApplication(c *gin.Context) {
    var req gmReq.WithdrawEmailAuditRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // ...
    }
    // ...
}
```

### 方案3：使用 PUT/PATCH 更新状态

**适用场景**：撤回是状态变更，不是删除

```go
PUT /api/gm/email/audit/:id/withdraw
Body: { "reason": "撤回原因" }  // 可选参数
```

**优点**：
- ✅ 语义清晰：PUT 表示更新操作
- ✅ 符合 RESTful 规范
- ✅ ID 在 URL 中

**缺点**：
- ⚠️ 路径较长

### 方案4：保持 POST，但改进路径设计（当前方案的优化）

**适用场景**：保持当前设计，但改进路径

```go
POST /api/gm/email/audit/:id/withdraw  // 更清晰的路径
```

**优点**：
- ✅ 路径更清晰，表示对某个资源执行撤回动作
- ✅ 符合 RESTful 的动作命名规范
- ✅ 不需要修改太多代码

## 五、推荐方案对比

| 方案 | HTTP方法 | URL格式 | Body | 语义 | 推荐度 |
|------|----------|---------|------|------|--------|
| 当前 | POST | `/withdraw/:id` | 空 | ⚠️ 不清晰 | ⭐⭐ |
| 方案1 | DELETE | `/:id` | 空 | ✅ 清晰 | ⭐⭐⭐⭐⭐ |
| 方案2 | POST | `/withdraw` | `{id, reason}` | ✅ 清晰 | ⭐⭐⭐⭐ |
| 方案3 | PUT | `/:id/withdraw` | `{reason}` | ✅ 清晰 | ⭐⭐⭐⭐ |
| 方案4 | POST | `/:id/withdraw` | 空/可选 | ✅ 清晰 | ⭐⭐⭐ |

## 六、行业最佳实践

### GitHub API 示例
```
DELETE /repos/:owner/:repo/issues/:id         # 删除 issue
POST /repos/:owner/:repo/issues/:id/lock     # 锁定 issue（动作）
PUT /repos/:owner/:repo/issues/:id/lock       # 更新锁定状态
```

### 常见模式

1. **删除资源**：
   ```
   DELETE /api/users/:id
   ```

2. **执行动作**（带资源ID）：
   ```
   POST /api/users/:id/activate
   POST /api/users/:id/deactivate
   POST /api/users/:id/reset-password
   ```

3. **执行动作**（不带资源ID，ID在Body）：
   ```
   POST /api/users/batch-delete
   Body: { "ids": [1, 2, 3] }
   ```

## 七、最终建议

### 推荐：使用 DELETE 方法（方案1）

**理由**：
1. **语义最清晰**：撤回申请实际上是删除申请记录
2. **符合 RESTful 规范**：DELETE 用于删除资源
3. **ID 在 URL 中**：这是 DELETE 方法的标准做法
4. **简单直观**：不需要额外的路径

**实现示例**：
```go
// Router
emailAuditRouter.DELETE(":id", gmEmailAuditApi.WithdrawApplication)

// API
func (e *EmailAuditApi) WithdrawApplication(c *gin.Context) {
    var id uint
    if err := c.ShouldBindUri(&id); err != nil {
        response.FailWithMessage("无效的申请ID", c)
        return
    }
    // ... 其他逻辑
}
```

### 备选：POST + 动作路径（方案4）

如果撤回不是真正的删除（如只是状态变更），可以使用：
```go
POST /api/gm/email/audit/:id/withdraw
```

## 八、总结

**当前设计的问题**：
- POST 方法在 URL 中传参虽然可以，但不是最佳实践
- DELETE 方法更适合删除/撤回操作

**改进方向**：
- ✅ 推荐使用 DELETE 方法
- ✅ 或者使用 POST + 动作路径 `/:id/withdraw`
- ✅ 避免 POST `/withdraw/:id` 这种设计

