# Postman 测试邮件审核接口指南

## 一、路由路径说明

### 完整URL路径组成

```
DELETE [RouterPrefix]/gm/email/audit/:id
```

**路径解析**：
- `RouterPrefix`: 系统路由前缀（通常在 `config.yaml` 中配置，常见值为 `/api` 或空）
- `gm/email/audit`: 邮件审核路由组
- `:id`: 路径参数，申请ID（要撤回的申请ID）

**为什么使用 DELETE 方法？**
- ✅ DELETE 方法语义清晰：表示删除/撤回资源
- ✅ 符合 RESTful 规范：DELETE 方法的ID在 URL 中是标准做法
- ✅ GET、DELETE 等方法使用 URL 参数是常见做法

### 实际URL示例

假设 `RouterPrefix = "/api"`，申请ID为 `123`：

```
DELETE http://localhost:8888/api/gm/email/audit/123
```

如果 `RouterPrefix = ""`（空）：

```
DELETE http://localhost:8888/gm/email/audit/123
```

## 二、Postman 测试步骤

### 1. 基本设置

1. **请求方法**：选择 `DELETE`
2. **URL**：输入完整URL，将 `:id` 替换为实际申请ID

   ```
   DELETE http://localhost:8888/api/gm/email/audit/123
   ```
   其中 `123` 是要撤回的申请ID

### 2. Headers 设置

在 Headers 标签页中添加：

| Key | Value | 说明 |
|-----|-------|------|
| `x-token` | `你的JWT Token` | 必填，用于身份认证 |
| `Content-Type` | `application/json` | 可选，根据后端要求 |

**获取 Token**：
1. 先调用登录接口获取 token
2. 将 token 复制到 `x-token` header 中

### 3. Body 设置

**重要**：DELETE 方法使用**路径参数**传递申请ID（在URL中），**不需要在 Body 中传递数据**。这是 DELETE 方法的标准做法。

Body 设置：
- **Body 标签页选择 "none"**
- **或直接留空**

### 4. 完整示例

```
Method: DELETE
URL: http://localhost:8888/api/gm/email/audit/123

Headers:
  x-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

Body: (none 或留空，DELETE 方法不需要 Body)
```

## 三、路径参数说明

### 什么是路径参数？

路径参数（Path Parameter）是直接在 URL 路径中传递的参数，格式为 `:id` 或 `{id}`。

### 在 Postman 中的使用

**方式1：直接替换**
```
URL: http://localhost:8888/api/gm/email/audit/withdraw/123
```
直接将 `:id` 替换为实际值 `123`

**方式2：使用 Params 标签（Postman 高级用法）**
1. 在 Params 标签页添加：
   - Key: `id`
   - Value: `123`
2. Postman 会自动将 URL 更新为：
   ```
   http://localhost:8888/api/gm/email/audit/withdraw/123
   ```

## 四、测试其他接口

### 创建申请
```
POST http://localhost:8888/api/gm/email/audit/apply

Body (JSON):
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
  "startTime": 1761926400,
  "areaIds": "1,2,3"
}
```

### 更新申请
```
PUT http://localhost:8888/api/gm/email/audit/apply

Body (JSON):
{
  "id": 1,
  "emailType": 1,
  "emailTitle": {
    "en": "Updated Email"
  },
  ...
}
```

### 审核申请
```
POST http://localhost:8888/api/gm/email/audit/review

Body (JSON):
{
  "id": 1,
  "status": 2,
  "auditComment": "审核通过"
}
```

### 获取申请列表
```
POST http://localhost:8888/api/gm/email/audit/list

Body (JSON):
{
  "page": 1,
  "pageSize": 10,
  "status": 1
}
```

### 获取申请详情
```
GET http://localhost:8888/api/gm/email/audit/123
```

## 五、常见问题

### Q1: 为什么 DELETE 请求不需要 Body？

A: DELETE 方法是 RESTful 规范中用于删除资源的标准方法，资源ID通过**路径参数**传递（在URL中），这是 DELETE 方法的标准做法。后端通过 `c.ShouldBindUri(&id)` 从URL路径中读取ID。

**为什么使用 DELETE 而不是 POST？**
- ✅ DELETE 方法语义清晰：表示删除/撤回资源
- ✅ 符合 RESTful 规范：DELETE 方法的ID在 URL 中是标准做法
- ✅ 与其他删除接口保持一致（如 `DELETE /gm/user`）
- ✅ GET、DELETE 等方法使用 URL 参数是常见做法

### Q2: 如何知道路由前缀是什么？

A: 查看 `server/config.yaml` 中的 `router-prefix` 配置，或查看其他API接口的URL格式。

### Q3: 如果ID不存在会怎样？

A: 后端会返回错误："申请不存在"。

### Q4: 只能撤回自己的申请吗？

A: 是的，后端会验证 `application.ApplicantId == applicantId`，只能撤回自己创建的申请。

### Q5: 只能撤回"待审核"状态的申请吗？

A: 是的，后端会验证 `application.Status == EmailAuditStatusPending`（值为1）。

## 六、响应示例

### 成功响应
```json
{
  "code": 0,
  "data": null,
  "msg": "撤回成功",
  "success": true
}
```

### 错误响应
```json
{
  "code": 7,
  "data": null,
  "msg": "只能撤回自己的申请",
  "success": false
}
```

## 七、Postman Collection 建议

建议创建一个 Postman Collection，包含以下请求：

1. **登录** - 获取 Token
2. **创建邮件审核申请**
3. **获取申请列表**
4. **获取申请详情** - 使用路径参数
5. **更新申请**
6. **撤回申请** - 使用路径参数 `withdraw/:id`
7. **审核申请**

每个请求都应包含：
- 正确的 Headers（x-token）
- 正确的 URL
- 正确的 Body（如果需要）

