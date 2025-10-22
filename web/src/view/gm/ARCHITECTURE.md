# GM管理模块架构说明

## 概述
GM管理模块已重构为独立的游戏管理模块，与系统用户管理完全分离，使用独立的API和状态管理。

## 目录结构
```
web/src/
├── api/
│   ├── gm_user.js          # 游戏用户API接口
│   └── gm_item.js          # 游戏道具API接口
├── pinia/
│   └── gm/
│       ├── user.js         # 游戏用户状态管理
│       └── item.js         # 游戏道具状态管理
└── view/
    └── gm/
        ├── index.vue       # GM管理首页
        ├── user/
        │   └── user.vue     # 游戏用户管理页面
        └── item/
            └── item.vue     # 游戏道具流水管理页面
```

## API接口设计

### 游戏用户API (gm_user.js)
- `getGMUserList(data)` - 获取游戏用户列表
- `createGMUser(data)` - 创建游戏用户
- `updateGMUser(data)` - 更新游戏用户信息
- `deleteGMUser(data)` - 删除游戏用户
- `getGMUser(id)` - 获取游戏用户详情
- `resetGMUserPassword(data)` - 重置游戏用户密码
- `setGMUserStatus(data)` - 设置游戏用户状态
- `batchOperateGMUser(data)` - 批量操作游戏用户
- `exportGMUser(data)` - 导出游戏用户数据
- `getGMUserStats()` - 获取游戏用户统计信息

### 游戏道具API (gm_item.js)
- `getGMItemList(data)` - 获取游戏道具流水列表
- `createGMItem(data)` - 创建游戏道具流水记录
- `updateGMItem(data)` - 更新游戏道具流水记录
- `deleteGMItem(data)` - 删除游戏道具流水记录
- `getGMItem(id)` - 获取游戏道具流水详情
- `batchDeleteGMItem(data)` - 批量删除游戏道具流水记录
- `exportGMItem(data)` - 导出游戏道具流水数据
- `cleanupGMItem(data)` - 清理旧数据
- `getGMItemStats(data)` - 获取游戏道具统计信息
- `getGMItemTypes()` - 获取游戏道具类型列表
- `getGMItemOperationTypes()` - 获取游戏道具操作类型列表

## 状态管理设计

### 游戏用户状态管理 (pinia/gm/user.js)
```javascript
// 状态
- userList: 用户列表
- currentUser: 当前用户
- loading: 加载状态
- total: 总数
- page: 当前页
- pageSize: 每页大小
- searchInfo: 搜索条件
- userStats: 用户统计信息

// 方法
- fetchUserList(): 获取用户列表
- createUser(): 创建用户
- updateUser(): 更新用户
- removeUser(): 删除用户
- resetPassword(): 重置密码
- setUserStatus(): 设置用户状态
- exportUsers(): 导出用户数据
- fetchUserStats(): 获取用户统计
```

### 游戏道具状态管理 (pinia/gm/item.js)
```javascript
// 状态
- itemList: 道具流水列表
- currentItem: 当前道具流水
- loading: 加载状态
- total: 总数
- page: 当前页
- pageSize: 每页大小
- searchInfo: 搜索条件
- itemStats: 道具统计信息
- itemTypes: 道具类型列表
- operationTypes: 操作类型列表

// 方法
- fetchItemList(): 获取道具流水列表
- createItem(): 创建道具流水记录
- updateItem(): 更新道具流水记录
- removeItem(): 删除道具流水记录
- exportItems(): 导出道具流水数据
- cleanupOldData(): 清理旧数据
- fetchItemStats(): 获取道具统计
```

## 页面组件设计

### 游戏用户管理页面 (user/user.vue)
- 使用 `useGMUserStore` 状态管理
- 调用 `gm_user.js` API接口
- 支持用户CRUD操作
- 支持密码重置功能
- 支持用户状态管理
- 支持角色权限设置

### 游戏道具流水管理页面 (item/item.vue)
- 使用 `useGMItemStore` 状态管理
- 调用 `gm_item.js` API接口
- 支持道具流水CRUD操作
- 支持数据导出功能
- 支持旧数据清理
- 支持多种查询条件

## 与系统用户管理的区别

| 功能 | 系统用户管理 | GM游戏用户管理 |
|------|-------------|----------------|
| 数据源 | 系统用户表 | 游戏用户表 |
| API接口 | `/api/user/*` | `/api/gm/user/*` |
| 状态管理 | `useUserStore` | `useGMUserStore` |
| 权限控制 | 系统权限 | 游戏权限 |
| 功能范围 | 系统管理 | 游戏运营 |

## 后端接口要求

### 游戏用户接口
```
POST /api/gm/user/list          # 获取游戏用户列表
POST /api/gm/user               # 创建游戏用户
PUT  /api/gm/user               # 更新游戏用户
DELETE /api/gm/user             # 删除游戏用户
GET  /api/gm/user/:id           # 获取游戏用户详情
POST /api/gm/user/resetPassword # 重置密码
POST /api/gm/user/status        # 设置用户状态
POST /api/gm/user/batch         # 批量操作
POST /api/gm/user/export        # 导出数据
GET  /api/gm/user/stats         # 获取统计信息
```

### 游戏道具接口
```
POST /api/gm/item/list          # 获取道具流水列表
POST /api/gm/item               # 创建道具流水记录
PUT  /api/gm/item               # 更新道具流水记录
DELETE /api/gm/item             # 删除道具流水记录
GET  /api/gm/item/:id           # 获取道具流水详情
POST /api/gm/item/batchDelete  # 批量删除
POST /api/gm/item/export        # 导出数据
POST /api/gm/item/cleanup       # 清理旧数据
POST /api/gm/item/stats         # 获取统计信息
GET  /api/gm/item/types         # 获取道具类型
GET  /api/gm/item/operationTypes # 获取操作类型
```

## 使用说明

1. **前端开发**：直接使用已创建的API和状态管理
2. **后端开发**：需要实现对应的API接口
3. **数据库设计**：需要创建游戏用户表和道具流水表
4. **权限配置**：需要配置GM管理相关的权限角色

## 注意事项

1. GM管理模块与系统用户管理完全独立
2. 使用不同的数据库表和API接口
3. 状态管理使用Pinia，支持响应式更新
4. 所有API调用都有错误处理机制
5. 支持分页、搜索、导出等完整功能
