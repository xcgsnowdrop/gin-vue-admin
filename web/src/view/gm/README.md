# GM管理模块

## 概述
GM管理模块是gin-vue-admin项目中的游戏管理功能模块，主要用于游戏运营管理。

## 目录结构
```
web/src/view/gm/
├── index.vue          # GM管理首页（路由容器）
├── user/
│   └── user.vue       # 游戏用户管理
├── item/
│   └── item.vue       # 游戏道具流水管理
└── README.md          # 说明文档
```

## 功能模块

### 1. GM管理首页 (index.vue)
- 作为路由容器，承载子页面
- 组件名：GM

### 2. 游戏用户管理 (user/user.vue)
- 管理游戏用户信息
- 支持用户查询、新增、编辑、删除
- 包含游戏用户ID、用户名、昵称等字段
- 组件名：GMUser

### 3. 游戏道具流水管理 (item/item.vue)
- 管理游戏道具的获得、消耗、交易记录
- 支持按用户ID、道具ID、操作类型、时间范围查询
- 支持数据导出和旧数据清理
- 组件名：GMItem

## 菜单配置

### 在后台管理系统中添加菜单：

1. **根菜单 - GM管理**
   - 名称：GM管理
   - 路径：/gm
   - 组件：view/gm/index.vue
   - 图标：game-controller（或其他游戏相关图标）

2. **子菜单 - 游戏用户管理**
   - 名称：游戏用户管理
   - 路径：/gm/user
   - 组件：view/gm/user/user.vue
   - 父菜单：GM管理

3. **子菜单 - 道具流水管理**
   - 名称：道具流水管理
   - 路径：/gm/item
   - 组件：view/gm/item/item.vue
   - 父菜单：GM管理

## 权限配置

建议为GM管理模块配置专门的权限角色：
- GM管理员：可以访问所有GM管理功能
- GM操作员：只能查看和操作部分功能

## 数据接口

### 游戏用户管理接口
- 获取用户列表：GET /api/gm/users
- 创建用户：POST /api/gm/users
- 更新用户：PUT /api/gm/users/:id
- 删除用户：DELETE /api/gm/users/:id

### 道具流水管理接口
- 获取流水列表：GET /api/gm/items
- 创建流水记录：POST /api/gm/items
- 更新流水记录：PUT /api/gm/items/:id
- 删除流水记录：DELETE /api/gm/items/:id
- 导出数据：GET /api/gm/items/export
- 清理旧数据：DELETE /api/gm/items/cleanup

## 注意事项

1. 确保在pathInfo.json中正确配置了组件路径
2. 菜单配置需要在后台管理系统中手动添加
3. 实际使用时需要实现对应的后端API接口
4. 建议为GM管理功能配置专门的权限控制
